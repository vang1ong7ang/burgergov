package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.ExpandEnv(":${PORT}")
	log.Println("[LISTEN]: ", addr)

	http.ListenAndServe(addr, nil)
}
