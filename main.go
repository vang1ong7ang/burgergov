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
}
func CanVote(rw http.ResponseWriter, r *http.Request) {
}
func Propose(rw http.ResponseWriter, r *http.Request) {
}
func Vote(rw http.ResponseWriter, r *http.Request) {
}
func VoteStatus(rw http.ResponseWriter, r *http.Request) {
}

type CanProposeReq struct {
	Proposer dao.Voter
}
type CanProposeResp struct {
	can bool
}
type CanVoteReq struct {
	Voter dao.Voter
}
type CanVoteResp struct {
	can bool
}
type ProposeReq struct {
	Proposal *dao.Proposal
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
