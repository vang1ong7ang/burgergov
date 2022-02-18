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
	"github.com/nspcc-dev/neo-go/pkg/util"
	"golang.org/x/oauth2"
)

var config struct {
	github_token      string
	github_owner      string
	github_repository string
	listen_address    string
}

var data struct {
	lock  sync.RWMutex
	nbips map[string]struct {
		SYNCTIME time.Time
		README   string
		NBIP     struct {
			TIMESTAMP  int64
			SCRIPTHASH util.Uint160
			METHOD     string
			ARGS       []interface{}
		}
		RESULT struct {
			TIMESTAMP int64
			PASSED    bool
			YES       uint64
			NO        uint64
		}
	}
	nobug map[util.Uint160]uint64
	votes map[string]map[util.Uint160]bool
}

var client *github.Client

func init() {
	config.github_token = os.ExpandEnv("${GITHUBTOKEN}")
	config.github_owner = os.ExpandEnv("${GITHUBOWNER}")
	config.github_repository = os.ExpandEnv("${GITHUBREPOSITORY}")
	config.listen_address = os.ExpandEnv(":${PORT}")
	log.Println("[LISTEN]:", config.listen_address)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.github_token})
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)

	http.DefaultClient.Timeout = time.Second * 10

	data.nbips = make(map[string]struct {
		SYNCTIME time.Time
		README   string
		NBIP     struct {
			TIMESTAMP  int64
			SCRIPTHASH util.Uint160
			METHOD     string
			ARGS       []interface{}
		}
		RESULT struct {
			TIMESTAMP int64
			PASSED    bool
			YES       uint64
			NO        uint64
		}
	})
	data.nobug = make(map[util.Uint160]uint64)
	data.votes = make(map[string]map[util.Uint160]bool)

	go func() {
		for ; ; time.Sleep(time.Hour) {
			func() {
				nbips := make(map[string]struct {
					SYNCTIME time.Time
					README   string
					NBIP     struct {
						TIMESTAMP  int64
						SCRIPTHASH util.Uint160
						METHOD     string
						ARGS       []interface{}
					}
					RESULT struct {
						TIMESTAMP int64
						PASSED    bool
						YES       uint64
						NO        uint64
					}
				})
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
					var readme []byte
					var nbipjson []byte
					var resultjson []byte
					func() {
						reader, err := client.Repositories.DownloadContents(ctx, config.github_owner, config.github_repository, "README.md", &github.RepositoryContentGetOptions{
							Ref: name,
						})
						if err != nil {
							// TODO: log
							return
						}
						defer reader.Close()
						readme, err = io.ReadAll(reader)
						if err != nil {
							// TODO: log
						}
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
						nbipjson, err = io.ReadAll(reader)
						if err != nil {
							// TODO: log
						}
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
						resultjson, err = io.ReadAll(reader)
						if err != nil {
							// TODO: log
						}
					}()
					func() {
						var nbip struct {
							SYNCTIME time.Time
							README   string
							NBIP     struct {
								TIMESTAMP  int64
								SCRIPTHASH util.Uint160
								METHOD     string
								ARGS       []interface{}
							}
							RESULT struct {
								TIMESTAMP int64
								PASSED    bool
								YES       uint64
								NO        uint64
							}
						}
						nbip.SYNCTIME = synctime
						nbip.README = string(readme)
						defer func() { nbips[name] = nbip }()
						if len(nbipjson) == 0 {
							return
						}
						if err := json.Unmarshal(nbipjson, &nbip.NBIP); err != nil {
							// TODO: log
						}
						if len(resultjson) == 0 {
							return
						}
						if err := json.Unmarshal(resultjson, &nbip.RESULT); err != nil {
							// TODO: log
						}
					}()
				}
				data.lock.Lock()
				defer data.lock.Unlock()
				data.nbips = nbips

			}()
			func() {
				// TODO: load more addresses
				nobug := make(map[util.Uint160]uint64)
				req := strings.NewReader(`{"jsonrpc": "2.0","method": "GetAssetHoldersByContractHash","params": {"ContractHash":"0x54806765d451e2b0425072730d527d05fbfa9817","Limit":4096,"Skip":0},"id": 1}`)
				rsp, err := http.Post("https://neofura.ngd.network", "application/json", req)
				if err != nil {
					log.Println("[ERROR]: ", "[SYNC]:", err, rsp)
					return
				}
				defer rsp.Body.Close()
				var output struct {
					Result struct {
						Result []struct {
							Address util.Uint160
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
				data.lock.Lock()
				defer data.lock.Unlock()
				data.nobug = nobug
			}()
			func() {
				data.lock.RLock()
				defer data.lock.RUnlock()
				for k := range data.nbips {
					log.Println("[MAINTAINED]: ", k)
				}
				log.Println("[MAINTAINED]:", "[HOLDER]:", len(data.nobug))
			}()
			func() {
				for name, v := range data.nbips {
					if v.RESULT.TIMESTAMP != 0 {
						continue
					}
					// count
					currentPage := 1
					perPage := 100
					for{
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						commits, _, err := client.Repositories.ListCommits(ctx, config.github_owner, config.github_repository,
							&github.CommitsListOptions{SHA: name, ListOptions: github.ListOptions{Page: currentPage, PerPage: perPage}})
						cancel()
						if err != nil{break}
						if commits == nil{break}

						if data.votes[name] == nil {
							data.votes[name] = make(map[util.Uint160]bool)
						}

						for _, commit := range commits {
							message := *commit.Commit.Message
							re := regexp.MustCompile("(0x\\w{40}) VOTE (FOR|AGAINST) " + name)
							match := re.FindStringSubmatch(message)
							if match == nil {
								continue
							}
							voter, err := util.Uint160DecodeStringBE(match[1][2:])
							if err != nil {
								continue
							}
							if _, ok := data.votes[name][voter]; !ok {
								data.votes[name][voter] = match[2] == "FOR"
							}
						}
						if len(commits) < perPage {break} else {currentPage += 1}
					}
					//log.Println(data.votes[name])
				}
			}()
		}
	}()
}
