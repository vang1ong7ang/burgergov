package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

func get_text_content(target string, id string) ([]byte, error) {
	client := &http.Client{}
	values := url.Values{}
	values.Add("ref", "NBIP-"+id)
	urlStringObj := url.URL{Scheme: "https", Host: "api.github.com", RawQuery: values.Encode(), Path: path.Join("/", "repos", config.github_repo, "contents", target)}
	req, err := http.NewRequest("GET", urlStringObj.String(), nil)
	req.Header.Add("Authorization", "token "+config.github_token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var downloadUrl struct{ Download_url string }
	json.Unmarshal(body, &downloadUrl)
	if downloadUrl.Download_url == "" {
		return body, err
	}

	resp, err = http.Get(downloadUrl.Download_url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func init() {
	http.HandleFunc("/nbip/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		dir := path.Dir(urlPath)
		target := path.Base(urlPath)
		head := path.Dir(dir)
		id := path.Base(dir)
		if head != "/nbip" {
			http.NotFound(w, r)
			return
		}
		switch target {
		case "README.md":
			fallthrough
		case "nbip.json":
			result, err := get_text_content(target, id)
			if err != nil {
				http.NotFound(w, r)
			}
			w.Write(result)
		case "all.json":
			readmeResult, readmeErr := get_text_content("README.md", id)
			nbipResult, nbipErr := get_text_content("nbip.json", id)
			if readmeErr != nil || nbipErr != nil {
				http.NotFound(w, r)
				return
			}
			result := struct {
				README string
				NBIP   json.RawMessage
			}{string(readmeResult), nbipResult}
			if err := json.NewEncoder(w).Encode(result); err != nil {
				log.Println("[ERROR]: ", err)
			}
		default:
			http.NotFound(w, r)
			return
		}
	})
}
