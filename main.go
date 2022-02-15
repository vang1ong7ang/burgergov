package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(os.ExpandEnv(":${PORT}"), nil)
}
