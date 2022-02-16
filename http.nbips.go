package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func init() {
	http.HandleFunc("/nbips.json", func(w http.ResponseWriter, r *http.Request) {
		req := url.URL{Scheme: "https", Host: "api.github.com", Path: path.Join("/", "repos", config.github_repo, "branches")}
		resp, err := http.Get(req.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: log
			return
		}
		defer resp.Body.Close()
		var branches []struct{ Name string }
		if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: log
			return
		}
		result := []string{}
		for _, branch := range branches {
			if strings.HasPrefix(branch.Name, "NBIP-") {
				result = append(result, branch.Name)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
