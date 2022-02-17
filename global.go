package main

import (
	"context"
	"log"
	"os"
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
	nbips []struct {
		readme []byte
		nbip   struct {
			SCRIPTHASH util.Uint160
			METHOD     string
			ARGS       []interface{}
			TIMESTAMP  int64
		}
		result struct {
			// TODO
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
			// load nbips
			// load nobug
		}
	}()
}
