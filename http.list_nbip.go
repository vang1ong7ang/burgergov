package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/list_nbip", func(w http.ResponseWriter, r *http.Request) {
		// get all branch names starts with `NBIP-`
		resp, err := http.Get("https://api.github.com/repos/" + config.github_repo + "/branches")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var branches []struct{ Name string }
		json.Unmarshal(body, &branches)
		result := []string{}
		for _, branch := range branches {
			result = append(result, branch.Name)
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
