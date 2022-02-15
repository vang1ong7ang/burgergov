package dao

import (
	"github.com/nspcc-dev/neo-go/pkg/util"
)

type Proposal struct {
	ProposalID     int64
	Scripthash     util.Uint160
	Method         string
	Args           []string
	VotingDeadline int64
}


type ProposalStatus struct {
	ProposalID     int64
	AcceptNum      int64
	DeclineNum     int64
}
