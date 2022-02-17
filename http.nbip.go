package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/google/go-github/github"
)

func get_text_content(pathname string, ref string) ([]byte, error) {
	ghctx, ghcancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ghcancel()
	reader, err := client.Repositories.DownloadContents(ghctx,
		config.github_owner, config.github_repository, pathname,
		&github.RepositoryContentGetOptions{
			Ref: ref,
		})
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
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
			result, err := get_text_content(target, fmt.Sprintf("NBIP-%s", id))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Write(result)
		case "all.json":
			readme, err := get_text_content("README.md", fmt.Sprintf("NBIP-%s", id))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			nbip, err := get_text_content("nbip.json", fmt.Sprintf("NBIP-%s", id))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			result := struct {
				README string
				NBIP   json.RawMessage
			}{
				string(readme),
				nbip,
			}
			if err := json.NewEncoder(w).Encode(result); err != nil {
				log.Println("[ERROR]: ", err)
			}
		default:
			http.NotFound(w, r)
			return
		}
	})
}
