package main

import (
	"context"
	"log"
	"os"
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
		synctime time.Time
		readme   []byte
		nbip     struct {
			TIMESTAMP  int64
			SCRIPTHASH util.Uint160
			METHOD     string
			ARGS       []interface{}
		}
		result struct {
			TIMESTAMP int64
			PASSED    bool
			YES       uint64
			NO        uint64
		}
	}
	nobug []struct {
	}
}

var client *github.Client

func init() {
	config.github_token = os.ExpandEnv("${GITHUBTOKEN}")
	config.github_owner = os.ExpandEnv("${GITHUBOWNER}")
	config.github_repository = os.ExpandEnv("${GITHUBREPOSITORY}")
	config.listen_address = os.ExpandEnv(":${PORT}")
	log.Println("[CONFIG]: ", config)

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.github_token})
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)

	go func() {
		for ; ; time.Sleep(time.Hour) {
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				gb, gr, err := client.Repositories.ListBranches(ctx, config.github_owner, config.github_repository, &github.ListOptions{
					PerPage: 100,
				})
				if err != nil {
					log.Println("[ERROR]: ", "[SYNC]:", err, gr)
				}
				synctime := time.Now()
				for _, branch := range gb {
					name := branch.GetName()
					if strings.HasPrefix(name, "NBIP-") == false {
						continue
					}
					// get readme
					// get nbip.json
					// get result.json
					func() {
						data.lock.Lock()
						defer data.lock.Unlock()
						nbip := data.nbips[name]
						nbip.synctime = synctime
						nbip.readme = nil // TODO
						// nbip.nbip
						// nbip.result
						data.nbips[name] = nbip
					}()
				}
			}()
		}
	}()
}
