package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := struct {
			Name      string
			Timestamp int64
		}{
			os.Args[0],
			time.Now().Unix(),
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Println("[ERROR]: ", err)
		}
	})
}
