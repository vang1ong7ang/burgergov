package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func init() {
	http.HandleFunc("/nbips.json", func(w http.ResponseWriter, r *http.Request) {
		ghctx, ghcancel := context.WithTimeout(context.Background(), time.Second*5)
		defer ghcancel()
		// TODO: list all branches (100 max now)
		branches, rsp, err := client.Repositories.ListBranches(ghctx, config.github_owner, config.github_repository,
			&github.ListOptions{
				Page: 1,
				PerPage: 100,
			})
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			log.Println("[ERROR]: ", err, branches, rsp)
			return
		}
		result := []string{}
		for _, branch := range branches {
			if strings.HasPrefix(branch.GetName(), "NBIP-") {
				result = append(result, branch.GetName())
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println("[ERROR]: ", err)
		}
	})
}
