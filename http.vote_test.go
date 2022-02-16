package main

import (
	"testing"

	"github.com/nspcc-dev/neo-go/pkg/util"
)

func Test_Vote_CheckSig(t *testing.T) {
	voter, err := util.Uint160DecodeStringLE("13a192c56738900f9918d7f1ec07d9d8c278b804")
	if err != nil {
		t.Fail()
	}
	vote := &Vote{
		Voter:     voter,
		Nbip:      0,
		Yes:       true,
		Type:      "NEOLINE",
		Signature: "c0bf05182fc055f8017b4da1bde0361a552ea9d434014122c78625c701f98d3b31ec222bdfef0182c56698967de9dee771e9e153d0fafe355790dfbe0f13f5fe",
		Extra:     "{\"Salt\": \"4baa4b4257cfc7971907de0cd68aecf5\", \"PublicKey\":\"021319f4f1ece7181760b876dcc08f91dbffc4b55743c78a9fe827a2cc7a8e2600\"}",
	}
	if vote.CheckSig() != true {
		t.Fail()
	}
}
