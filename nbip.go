package main

import (
	"context"
	"github.com/google/go-github/github"
	"log"
	"strings"
	"sync"
	"time"
)

var nbips []struct {
	README string
	NBIP   struct{} // TODO: FIX
	RESULT struct{} // TODO: FIX
}

type getContentResponse struct {
	branchName *string
	result     *github.RepositoryContent
	resp       *github.Response
}

var wg *sync.WaitGroup

func sync_branches() {
	ctx := context.Background()
	branches, _, err := client.Repositories.ListBranches(ctx, config.github_owner, config.github_repository, nil)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(branches)
	activeBranchUrls := []string{}
	wg = &sync.WaitGroup{}
	resultChannel := make(chan *getContentResponse, len(branches))
	for _, branch := range branches {
		if strings.HasPrefix(*branch.Name, "NBIP-") {
			wg.Add(1)
			response := getContentResponse{branch.Name, nil, nil}
			resultChannel <- &response
			go func(branchName *string) {
				resultJson, _, resp, _ := client.Repositories.GetContents(
					ctx, config.github_owner, config.github_repository, "result.json",
					&github.RepositoryContentGetOptions{Ref: *branchName})
				response.result, response.resp = resultJson, resp
				wg.Done()
			}(branch.Name)
		}
	}
	wg.Wait()
	//close(resultChannel)
	for len(resultChannel) > 0 {
		// If "result.json" does not exist, there is always an error != nil
		response := <-resultChannel
		resultJson, resp := response.result, response.resp
		if resultJson == nil && resp.StatusCode == 404 && resp.Header.Get("X-Github-Request-Id") != "" {
			// Active NBIP branch
			activeBranchUrls = append(activeBranchUrls, *response.branchName)
		}
	}
	log.Println(activeBranchUrls)
	time.Sleep(6_00_000_000_000)
	go sync_branches()
}
