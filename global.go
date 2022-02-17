package main

import (
	"context"
	"encoding/json"
	"io"
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
				defer func() {
					data.lock.Lock()
					defer data.lock.Unlock()
					data.nbips = nbips
				}()
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

			}()
			func() {
				data.lock.RLock()
				defer data.lock.RUnlock()
				for k := range data.nbips {
					log.Println("[MAINTAINED]: ", k)
				}
			}()
			func() {
				// load nobug
			}()
			func() {
				// count
				for _, v := range data.nbips {
					if v.RESULT.TIMESTAMP != 0 {
						continue
					}
					// count
					v.RESULT.YES = 0 // TODO
					v.RESULT.NO = 0  // TODO
				}
			}()
		}
	}()
}
