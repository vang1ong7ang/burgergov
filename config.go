package main

import (
	"log"
	"os"
)

var config struct {
	github_token string
	github_repo  string
	listen_address string
}

func init() {
	config.github_token = os.ExpandEnv("${GITHUBTOKEN}")
	config.github_repo = os.ExpandEnv("${GITHUBREPO}")
	config.listen_address = os.ExpandEnv(":${PORT}")
	log.Println("[CONFIG]: ", config)
}
