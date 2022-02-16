package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var config struct {
	github_token   string
	github_owner   string
	github_repo    string
	listen_address string
}
var client *github.Client

func init() {
	config.github_token = os.ExpandEnv("${GITHUBTOKEN}")
	config.github_owner = os.ExpandEnv("${GITHUBOWNER}")
	config.github_repo = os.ExpandEnv("${GITHUBREPO}")
	config.listen_address = os.ExpandEnv(":${PORT}")
	log.Println("[CONFIG]: ", config)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.github_token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)
}
