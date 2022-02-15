package dao

import (
	"github.com/nspcc-dev/neo-go/pkg/util"
)

type Proposal struct {
	ProposalID      int64
	Scripthash      util.Uint160
	Method          string
	Args            []string
	VotingTimestamp int64
	Readme			string        // must contain when propose new proposal, will not show when list
	ReadmeURL       string        // assigned by github
}

type ProposalStatus struct {
	ProposalID int64
	AcceptNum  int64
	DeclineNum int64
}
