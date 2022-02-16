package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/list_nbip", func(w http.ResponseWriter, r *http.Request) {
		// get all branch names starts with `NBIP-`
		client := http.Client{}
		req, _ := http.NewRequest("GET", "https://api.github.com/repos/" + config.github_repo + "/branches", nil)
		req.Header.Add("Authorization", "token "+config.github_token)
		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var branches []struct{ Name string }
		json.Unmarshal(body, &branches)
		result := []string{}
		for _, branch := range branches {
			if strings.HasPrefix(branch.Name, "NBIP-") {
				result = append(result, branch.Name)
			}
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
