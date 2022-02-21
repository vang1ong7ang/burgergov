package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/counts.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data.get_counts()); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println("[ERROR]: ", err)
		}
	})
}
