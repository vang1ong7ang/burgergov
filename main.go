package main

import (
	"net/http"

	"github.com/neoburger/burgergov/dao"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {

	})
	http.HandleFunc("/can_propose", CanPropose)
	http.HandleFunc("/can_vote", CanVote)
	http.HandleFunc("/propose", Propose)
	http.HandleFunc("/vote", Vote)
	http.HandleFunc("/vote_status", VoteStatus)
	http.ListenAndServe("0.0.0.0:443", nil)
}

// TODO
func CanPropose(rw http.ResponseWriter, r *http.Request) {
	// 1. proposer has NoBug token
	// 2. proposer has no active proposal yet
}
func CanVote(rw http.ResponseWriter, r *http.Request) {
	// 1. voter has NoBug token
	// 2. voter has not vote this existing proposal yet
}
func Propose(rw http.ResponseWriter, r *http.Request) {
	// 1. check can propose
	// 2. assign id and timestamp
	// 3. make the propose (to github)
}
func Vote(rw http.ResponseWriter, r *http.Request) {
	// 1. check can vote
	// 2. make the vote (to github)
}
func VoteStatus(rw http.ResponseWriter, r *http.Request) {
	// return data calculated by cache
}

type CanProposeReq struct {
	Proposer dao.Voter
}
type CanProposeResp struct {
	can bool
}
type CanVoteReq struct {
	Voter      dao.Voter
	ProposalID int64
}
type CanVoteResp struct {
	can bool
}
type ProposeReq struct {
	Proposal *dao.Proposal
	Readme   string
}
type ProposeResp struct {
}
type VoteReq struct {
	Vote *dao.Vote
}
type VoteResp struct {
}
type VoteStatusReq struct {
	ProposalID int64
}
type VoteStatusResp struct {
	proposal *dao.Proposal
	status   *dao.ProposalStatus
}
