package dao

import (
	"math/big"
)

type Vote struct {
	Proposer   Voter
	ProposalID big.Int
	accept     bool
}
