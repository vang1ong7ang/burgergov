package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

func init() {
	word := map[bool]string{true: "FOR", false: "AGAINST"}
	strptr := func(v string) *string { return &v }

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}
		voter := r.Form.Get("voter")
		if strings.HasPrefix(voter, "0x") == false {
			http.Error(w, "invalid voter", http.StatusBadRequest)
			return
		}
		sh, err := util.Uint160DecodeStringLE(voter[2:])
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
		nbipStr := fmt.Sprintf("NBIP-%d", nbip)
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
			if pk.GetScriptHash() != sh {
				http.Error(w, "invalid public key", http.StatusBadRequest)
				return
			}
			message := fmt.Sprintf("I VOTE %s %s.", word[yes], nbipStr)
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
		content, err := json.Marshal(struct {
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
		})
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			log.Println("[ERROR]: ", err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		grcr, gr, err := client.Repositories.CreateFile(ctx, config.github_owner, config.github_repository, fmt.Sprintf("0x%s.json", sh.StringLE()), &github.RepositoryContentFileOptions{
			Branch:  &nbipStr,
			Content: content,
			Message: strptr(fmt.Sprintf("0x%s VOTE %s %s.", sh.StringLE(), word[yes], nbipStr)),
			Author: &github.CommitAuthor{
				Name:  strptr(address.Uint160ToString(sh)),
				Email: strptr(fmt.Sprintf("%s@NOREPLY", address.Uint160ToString(sh))),
			},
		})
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			log.Println("[ERROR]: ", err, grcr, gr)
			return
		}
		data.append_votes(nbipStr, scripthash{sh}, yes)
	})
}
