package proposal

import (
)

type Proposal struct {
	ProposalID string // biginteger
	Scripthash string // uint160
	Method string
	Args string // ByteString[]
	VotingDeadline int64 // BigInteger
}
