package dao

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/nspcc-dev/neo-go/pkg/crypto/hash"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
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

func (v *Vote) CheckSig() bool {
	pubKey := keys.PublicKey(v.Voter)
	message := hex.EncodeToString([]byte(v.Salt + v.Message()))
	message = fmt.Sprintf("%s%s%s%s", "010001f0", hex.EncodeToString([]byte{byte(len(message) / 2)}), message, "0000")
	msg, err := hex.DecodeString(message)
	if err != nil {
		panic(err)
	}
	hashedData := hash.Sha256(msg)
	sig, err := hex.DecodeString(v.Signature)
	if err != nil {
		return false
	}
	return pubKey.Verify(sig, hashedData.BytesBE())
}
