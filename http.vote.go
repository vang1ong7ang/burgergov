package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nspcc-dev/neo-go/pkg/crypto/hash"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

var (
	DATA_PROPOSAL_ROOT = "proposals"
)

func init() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		voter := r.Form.Get("voter")
		nbip := r.Form.Get("nbip")
		yes := r.Form.Get("yes")
		tp := r.Form.Get("type")
		signature := r.Form.Get("signature")
		extra := r.Form.Get("extra")

		vote := &Vote{}
		if v, err := util.Uint160DecodeStringBE(voter); err != nil {
			http.Error(w, "invalid voter which can not convert to uint160", http.StatusBadRequest)
			return
		} else {
			vote.Voter = v
		}
		if v, err := strconv.ParseUint(nbip, 10, 64); err != nil {
			http.Error(w, "invalid nbip which can not convert to uint64", http.StatusBadRequest)
			return
		} else {
			vote.Nbip = v
		}
		if v, err := strconv.ParseBool(yes); err != nil {
			http.Error(w, "invalid yes which can not convert to bool", http.StatusBadRequest)
			return
		} else {
			vote.Yes = v
		}
		vote.Type = tp
		vote.Signature = signature
		vote.Extra = extra
		if !vote.CheckSig() {
			http.Error(w, "signature check failed", http.StatusUnauthorized)
			return
		}

		log.Println("vote check parameter passed")

		dirUrl := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s/%d",
			config.github_repo, DATA_PROPOSAL_ROOT, vote.Nbip)
		dirRsp, err := http.Get(dirUrl)
		if err != nil {
			http.Error(w, "failed to contact to github", http.StatusInternalServerError)
			return
		}
		if dirRsp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("invalid nbip [%d], proposal not found", vote.Nbip), http.StatusBadRequest)
			return
		}
		// TODO: check the proposal path is a real directory
		log.Println("directory check passed")

		storeUrl := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s/%d/%s.json",
			config.github_repo, DATA_PROPOSAL_ROOT, vote.Nbip, vote.Voter.String())
		getRsp, err := http.Get(storeUrl)
		if err != nil {
			log.Printf("failed to contact to github, err: %s\n", err.Error())
			http.Error(w, "failed to contact to github", http.StatusInternalServerError)
			return
		}
		switch getRsp.StatusCode {
		case http.StatusOK:
			fallthrough
		case http.StatusFound:
			http.Error(w, "repeat vote, previous vote had been recorded", http.StatusBadRequest)
			return
		case http.StatusForbidden:
			http.Error(w, "failed to contact to github", http.StatusInternalServerError)
			return
		case http.StatusNotFound:
		default:
			log.Printf("Unknown status code [%d] when checking repeat submitting\n", getRsp.StatusCode)
		}

		log.Println("repeat submitting check passed")

		body, _ := json.Marshal(vote)
		// TODO: add auth token
		putReq, err := http.NewRequest("PUT", storeUrl, bytes.NewReader(body))
		if err != nil {
			log.Printf("failed to construct the put request, err: %s\n", err.Error())
			http.Error(w, "failed to construct the put request", http.StatusInternalServerError)
			return
		}
		putResp, err := http.DefaultClient.Do(putReq)
		if err != nil {
			log.Printf("failed to contact to github, err: %s\n", err.Error())
			http.Error(w, "failed to contact to github", http.StatusInternalServerError)
			return
		}
		if putResp.StatusCode != http.StatusOK {
			log.Printf("failed to submit vote to github, err: %s\n", err.Error())
			http.Error(w, "failed to submit vote to github", http.StatusInternalServerError)
			return
		}
	})
}

type Vote struct {
	Voter     util.Uint160
	Nbip      uint64
	Yes       bool
	Type      string
	Signature string
	Extra     string
}

func (v *Vote) Message() string {
	return fmt.Sprintf("[%d,%d,%t]", 860833102, v.Nbip, v.Yes)
}

func (v *Vote) CheckSig() bool {
	switch v.Type {
	case "NEOLINE":
		var extra struct {
			PublicKey string
			Salt      string
		}
		if err := json.Unmarshal([]byte(v.Extra), extra); err != nil {
			log.Printf("invalid extra info for NEOLINE typed vote, err: %s\n", err.Error())
			return false
		}
		pubKey, err := keys.NewPublicKeyFromString(extra.PublicKey)
		if err := json.Unmarshal([]byte(v.Extra), extra); err != nil {
			log.Printf("invalid public key info for NEOLINE typed vote, err: %s\n", err.Error())
			return false
		}
		if pubKey.Address() != v.Voter.String() {
			log.Printf("public key info for NEOLINE typed vote not match the voter, pubKey address: %s, voter address: %s\n", pubKey.Address(), v.Voter)
			return false
		}

		message := hex.EncodeToString([]byte(extra.Salt + v.Message()))
		message = fmt.Sprintf("%s%s%s%s", "010001f0", hex.EncodeToString([]byte{byte(len(message) / 2)}), message, "0000")
		msg, err := hex.DecodeString(message)
		if err != nil {
			panic(err)
		}
		hashedData := hash.Sha256(msg)
		sig, err := hex.DecodeString(v.Signature)
		if err != nil {
			log.Printf("invalid signature: %s, err: %s\n", v.Signature, err.Error())
			return false
		}
		return pubKey.Verify(sig, hashedData.BytesBE())
	default:
		log.Printf("unsupported vote type: %s\n", v.Type)
		return false
	}
}
