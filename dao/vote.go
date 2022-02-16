package dao

import (
	"fmt"
	"math/big"
)

type Vote struct {
	Voter      Voter
	ProposalID big.Int
	Accept     bool
	Salt       string
	Signature  string
}

func (v *Vote) Message() string {
	return fmt.Sprintf("[%d,%d,%t]", 860833102, v.ProposalID, v.Accept)
}
