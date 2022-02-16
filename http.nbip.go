package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func get_text_content(w http.ResponseWriter, r *http.Request, target string, id string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/" + config.github_repo + "/contents/" + target + "?ref=NBIP-" + id, nil)
	req.Header.Add("Authorization", "token "+config.github_token)
	resp, err := client.Do(req)
	if err != nil {log.Println(err)}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var downloadUrl struct{ Download_url string }
	json.Unmarshal(body, &downloadUrl)
	if downloadUrl.Download_url == ""{
		http.NotFound(w, r)
		return nil
	}

	resp, err = http.Get(downloadUrl.Download_url)
	if err != nil {log.Println(err)}
	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	return result
}

func init() {
	http.HandleFunc("/nbip/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		dir := filepath.Dir(path)
		target := filepath.Base(path)
		head := filepath.Dir(dir)
		id := filepath.Base(dir)
		if head != "/nbip" && head != "\\nbip" {
			http.NotFound(w, r)
			return
		}
		switch target {
		case "README.md":
			fallthrough
		case "nbip.json":
			result := get_text_content(w, r, target, id)
			w.Write(result)
		case "all.json":
			result := struct {
				README string
				NBIP   json.RawMessage
			}{string(get_text_content(w, r, "README.md", id)), get_text_content(w, r, "nbip.json", id)}
			if err := json.NewEncoder(w).Encode(result); err != nil {
				log.Println("[ERROR]: ", err)
			}
		default:
			http.NotFound(w, r)
			return
		}
	})
}
