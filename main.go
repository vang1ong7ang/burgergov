package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/list_proposals", list_proposals)
	http.HandleFunc("/can_propose", CanPropose)
	http.HandleFunc("/can_vote", CanVote)
	http.HandleFunc("/propose", Propose)
	http.HandleFunc("/vote", Vote)
	http.HandleFunc("/vote_status", VoteStatus)
	http.ListenAndServe(os.ExpandEnv(":${PORT}"), nil)
}

