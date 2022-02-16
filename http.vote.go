package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

func init() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}
		voter, err := util.Uint160DecodeStringBE(r.Form.Get("voter"))
		if err != nil {
			http.Error(w, "invalid voter", http.StatusBadRequest)
			return
		}
		// TODO: check voter
		// balance
		nbip, err := strconv.Atoi(r.Form.Get("nbip"))
		if err != nil {
			http.Error(w, "invalid nbip", http.StatusBadRequest)
			return
		}
		if nbip <= 0 {
			http.Error(w, "invalid nbip", http.StatusBadRequest)
			return
		}
		// TODO: check nbip
		// active
		// voted
		yes, err := strconv.ParseBool(r.Form.Get("yes"))
		if err != nil {
			http.Error(w, "invalid yes", http.StatusBadRequest)
			return
		}
		profile := r.Form.Get("profile")
		signature := r.Form.Get("signature")
		extra := r.Form.Get("extra")
		switch profile {
		case "NEOLINE":
			if len(extra) < 32 {
				http.Error(w, "invalid extra", http.StatusBadRequest)
				return
			}
			pk, err := keys.NewPublicKeyFromString(extra[32:])
			if err != nil {
				http.Error(w, "invalid extra", http.StatusBadRequest)
				return
			}
			if pk.GetScriptHash() != voter {
				http.Error(w, "invalid public key", http.StatusBadRequest)
				return
			}
			word := map[bool]string{true: "FOR", false: "AGAINST"}
			message := fmt.Sprintf("I VOTE %s NBIP-%d.", word[yes], nbip)
			if len(message) > 64 {
				http.Error(w, "invalid nbip", http.StatusBadRequest)
				return
			}
			msg := []byte{0x01, 0x00, 0x01, 0xf0, byte(len(message)) + 32}
			msg = append(msg, extra[:32]...)
			msg = append(msg, message...)
			msg = append(msg, 0, 0)
			digest := sha256.Sum256(msg)
			sig, err := hex.DecodeString(signature)
			if err != nil {
				http.Error(w, "invalid signautre", http.StatusBadRequest)
				return
			}
			if pk.Verify(sig, digest[:]) == false {
				http.Error(w, "invalid signature", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "invalid profile", http.StatusBadRequest)
			return
		}

		content := struct {
			NBIP      int
			YES       bool
			PROFILE   string
			SIGNATURE string
			EXTRA     string
		}{
			nbip,
			yes,
			profile,
			signature,
			extra,
		}

		// TODO: upload
		_ = content
	})
}
