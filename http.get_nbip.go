package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/get_nbip", func(w http.ResponseWriter, r *http.Request) {
		// get the `readme.md` and `nbip.json`
		result := struct {
			README string
			NBIP   interface{}
		}{}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
