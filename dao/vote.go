package dao

import (
	"math/big"

	"github.com/nspcc-dev/neo-go/pkg/util"
)

type Vote struct {
	Proposer Voter
	ProposalID big.Int
	accept bool
}
