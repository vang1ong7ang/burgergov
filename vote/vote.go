package vote

import (
)

type Vote struct {
	Voter string // uint160
	ProposalID string // biginteger
	accept bool
}
