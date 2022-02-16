package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/nbip/", func(w http.ResponseWriter, r *http.Request) {
		// /nbip/{NUM}/nbip.json
		// /nbip/{NUM}/README.md
		// /nbip/{NUM}/all.json
		
		// log.Println(r.URL.Path)
		
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
