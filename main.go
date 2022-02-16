package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.ExpandEnv(":${PORT}")
	log.Println("[LISTEN]: ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println("[END]: ", err)
	}
}
