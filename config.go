package main

import "os"

var config struct {
	github_token string
	github_repo  string
}

func init() {
	config.github_token = os.ExpandEnv("${GITHUBTOKEN}")
	config.github_repo = os.ExpandEnv("${GITHUBREPO}")
}
