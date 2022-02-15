package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/list_proposals", list_proposals)
	http.ListenAndServe(os.ExpandEnv(":${PORT}"), nil)
}
