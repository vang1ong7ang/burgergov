package main

import (
	"net/http"
)

func list_proposals(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
