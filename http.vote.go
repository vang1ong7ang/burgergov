package main

import (
	"net/http"
)

func init() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
	})
}
