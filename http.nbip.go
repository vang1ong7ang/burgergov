package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
)

func init() {
	http.HandleFunc("/nbip/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		dir := filepath.Dir(path)
		target := filepath.Base(path)
		head := filepath.Dir(dir)
		id := filepath.Base(dir)
		if head != "/nbip" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		println(id)
		switch target {
		case "README.md":
			// todo
		case "nbip.json":
			// todo
		case "all.json":
			// todo
		}

		// to be removed
		
		// /nbip/{NUM}/nbip.json
		// /nbip/{NUM}/README.md
		// /nbip/{NUM}/all.json
		
		
		// all.json
		result := struct {
			README string
			NBIP   interface{}
		}{}
		
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
