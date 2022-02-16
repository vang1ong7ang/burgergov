package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/list_nbip", func(w http.ResponseWriter, r *http.Request) {
		// get all branch names starts with `NBIP-`
		result := []string{}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
