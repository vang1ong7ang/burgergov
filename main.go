package main

import (
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(config.listen_address, nil); err != nil {
		log.Println("[END]: ", err)
	}
}
