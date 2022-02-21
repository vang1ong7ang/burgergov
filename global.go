package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var config struct {
	github_token      string
	github_owner      string
	github_repository string
	listen_address    string
	neofura_request   string
}

type state_nbip struct {
	SYNCTIME time.Time
	README   string
	NBIP     *struct {
		TIMESTAMP  int64
		SCRIPTHASH scripthash
		METHOD     string
		ARGS       []interface{}
	}
	RESULT *struct {
		TIMESTAMP  int64
		BLOCKINDEX uint64
		PASSED     bool
		YES        uint64
		NO         uint64
	}
}

type state_count struct {
	YES uint64
	NO  uint64
}

type state_vote struct {
	TIMESTAMP time.Time
	YES       bool
}

type state struct {
	lock sync.RWMutex

	nbips map[string]state_nbip
	nobug map[scripthash]uint64
	votes map[string]map[scripthash]state_vote

	counts map[string]state_count
}

func (me *state) set_nbips(v map[string]state_nbip) {
	me.lock.Lock()
	defer me.lock.Unlock()
	me.nbips = v
}

func (me *state) get_nbips() map[string]state_nbip {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.nbips
}

func (me *state) set_nobug(v map[scripthash]uint64) {
	me.lock.Lock()
	defer me.lock.Unlock()
	me.nobug = v
}

func (me *state) get_nobug() map[scripthash]uint64 {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.nobug
}

func (me *state) update_votes(v map[string]map[scripthash]state_vote, t time.Time) {
	me.lock.Lock()
	defer me.lock.Unlock()
	for kpv, pv := range me.votes {
		for ksv, sv := range pv {
			if sv.TIMESTAMP.After(t) {
				if item, ok := v[kpv]; ok {
					if _, ok := item[ksv]; ok == false {
						item[ksv] = sv
					}
				}
			}
		}
	}
	me.votes = v
}

func (me *state) set_votes(v map[string]map[scripthash]state_vote) {
	me.lock.Lock()
	defer me.lock.Unlock()
	me.votes = v
}

func (me *state) append_votes(nbip string, voter scripthash, v bool) {
	me.lock.Lock()
	defer me.lock.Unlock()
	if nbip, ok := me.votes[nbip]; ok {
		nbip[voter] = state_vote{
			TIMESTAMP: time.Now(),
			YES:       v,
		}
	}
}

func (me *state) get_votes() map[string]map[scripthash]state_vote {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.votes
}

func (me *state) set_counts(v map[string]state_count) {
	me.lock.Lock()
	defer me.lock.Unlock()
	me.counts = v
}

func (me *state) get_counts() map[string]state_count {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.counts
}

func (me *state) biz_refresh_nbips() {
	nbips := make(map[string]state_nbip)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	gb, gr, err := client.Repositories.ListBranches(ctx, config.github_owner, config.github_repository, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		log.Println("[ERROR]: ", "[SYNC]:", err, gr)
		return
	}
	synctime := time.Now()
	for _, branch := range gb {
		name := branch.GetName()
		if strings.HasPrefix(name, "NBIP-") == false {
			continue
		}
		item := struct {
			SYNCTIME time.Time
			README   string
			NBIP     *struct {
				TIMESTAMP  int64
				SCRIPTHASH scripthash
				METHOD     string
				ARGS       []interface{}
			}
			RESULT *struct {
				TIMESTAMP  int64
				BLOCKINDEX uint64
				PASSED     bool
				YES        uint64
				NO         uint64
			}
		}{
			synctime,
			"",
			nil,
			nil,
		}
		func() {
			// TODO: RETRY
			reader, err := client.Repositories.DownloadContents(ctx, config.github_owner, config.github_repository, "README.md", &github.RepositoryContentGetOptions{
				Ref: name,
			})
			if err != nil {
				// TODO: log
				return
			}
			defer reader.Close()
			readme, err := io.ReadAll(reader)
			if err != nil {
				// TODO: log
				return
			}
			item.README = string(readme)
		}()
		func() {
			reader, err := client.Repositories.DownloadContents(ctx, config.github_owner, config.github_repository, "nbip.json", &github.RepositoryContentGetOptions{
				Ref: name,
			})
			if err != nil {
				// TODO: log
				return
			}
			defer reader.Close()
			nbip := new(struct {
				TIMESTAMP  int64
				SCRIPTHASH scripthash
				METHOD     string
				ARGS       []interface{}
			})
			if err := json.NewDecoder(reader).Decode(nbip); err != nil {
				// TODO: log
				return
			}
			item.NBIP = nbip
		}()
		func() {
			reader, err := client.Repositories.DownloadContents(ctx, config.github_owner, config.github_repository, "result.json", &github.RepositoryContentGetOptions{
				Ref: name,
			})
			if err != nil {
				// TODO: log
				return
			}
			defer reader.Close()
			result := new(struct {
				TIMESTAMP  int64
				BLOCKINDEX uint64
				PASSED     bool
				YES        uint64
				NO         uint64
			})
			if err := json.NewDecoder(reader).Decode(result); err != nil {
				// TODO: log
				return
			}
			item.RESULT = result
		}()
		nbips[name] = item
	}
	me.set_nbips(nbips)
}

func (me *state) biz_refresh_nobug() {
	// TODO: load more addresses
	nobug := make(map[scripthash]uint64)
	req := strings.NewReader(config.neofura_request)
	rsp, err := http.Post("https://neofura.ngd.network", "application/json", req)
	if err != nil {
		log.Println("[ERROR]: ", "[SYNC]:", err, rsp)
		return
	}
	defer rsp.Body.Close()
	var output struct {
		Result struct {
			Result []struct {
				Address scripthash
				Balance string
			}
		}
	}
	if err := json.NewDecoder(rsp.Body).Decode(&output); err != nil {
		log.Println("[ERROR]: ", "[SYNC]:", err)
		return
	}
	for _, r := range output.Result.Result {
		balance, err := strconv.ParseUint(r.Balance, 10, 64)
		if err != nil {
			log.Println("[ERROR]: ", "[SYNC]: ", err, "balance: ", r.Balance, "addr: ", r.Address)
			return
		}
		nobug[r.Address] = balance
	}
	me.set_nobug(nobug)
}

func (me *state) biz_refresh_votes() {
	start_time := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	reg := regexp.MustCompile(`^0x(\w{40}) VOTE (FOR|AGAINST) NBIP-\d+\.$`)
	votes := make(map[string]map[scripthash]state_vote)
	for k, v := range me.get_nbips() {
		if v.RESULT != nil {
			continue
		}
		for item, p := make(map[scripthash]state_vote), 1; p < 64; p++ {
			grc, gr, err := client.Repositories.ListCommits(ctx, config.github_owner, config.github_repository, &github.CommitsListOptions{SHA: k, ListOptions: github.ListOptions{Page: p, PerPage: 100}})
			if err != nil {
				log.Println("ERROR", err, gr)
				// TODO: see if EOF is here
				break
			}
			for _, commit := range grc {
				message := commit.GetCommit().GetMessage()
				match := reg.FindStringSubmatch(message)
				if len(match) != 3 {
					continue
				}
				voter, err := ScripthashDecodeStringLE(match[1])
				if err != nil {
					continue
				}
				item[voter] = state_vote{TIMESTAMP: commit.Author.CreatedAt.Time, YES: match[2] == "FOR"}
			}
			if len(grc) < 100 {
				votes[k] = item
				break
			}
		}
	}

	me.update_votes(votes, start_time)
}

func (me *state) biz_refresh_counts() {
	nobug := me.get_nobug()
	votes := me.get_votes()
	counts := make(map[string]state_count)
	for k, v := range votes {
		count := state_count{}
		for k, v := range v {
			if v.YES {
				count.YES += nobug[k]
			} else {
				count.NO += nobug[k]
			}
		}
		counts[k] = count
	}
	me.set_counts(counts)
}

func (me *state) biz_log() {
	for k := range me.get_nbips() {
		log.Println("[NBIP]:", k)
	}
	log.Println("[ADDRESSES]:", len(me.get_nobug()))
	log.Println("[ACTIVE]:", len(me.get_votes()))
}

func (me *state) biz_refresh() {
	me.biz_refresh_nbips()
	me.biz_refresh_nobug()
	me.biz_refresh_votes()
	me.biz_refresh_counts()
	me.biz_log()
}

var data state

var client *github.Client

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.github_token = os.ExpandEnv("${GITHUBTOKEN}")
	config.github_owner = os.ExpandEnv("${GITHUBOWNER}")
	config.github_repository = os.ExpandEnv("${GITHUBREPOSITORY}")
	config.listen_address = os.ExpandEnv(":${PORT}")
	config.neofura_request = os.ExpandEnv(`{"jsonrpc": "2.0","method": "GetAssetHoldersByContractHash","params": {"ContractHash":"${NOBUG}","Limit":4096,"Skip":0},"id": 1}`)

	log.Println("[LISTEN]:", config.listen_address)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.github_token})
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)

	http.DefaultClient.Timeout = time.Second * 10

	data.set_nbips(make(map[string]state_nbip))
	data.set_nobug(make(map[scripthash]uint64))
	data.set_votes(make(map[string]map[scripthash]state_vote))
	data.set_counts(make(map[string]state_count))

	go func() {
		for ; ; time.Sleep(time.Hour) {
			data.biz_refresh()
		}
	}()
}
