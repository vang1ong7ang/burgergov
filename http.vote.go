package main

import (
	"net/http"
	"strconv"
)

func init() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}
		voter := r.Form.Get("voter")
		// TODO: check voter
		nbip, err := strconv.Atoi(r.Form.Get("nbip"))
		if err != nil {
			http.Error(w, "invalid nbip", http.StatusBadRequest)
			return
		}
		// TODO: check nbip
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
			_ = voter
			_ = signature
			_ = extra
			_ = nbip
			_ = yes
			// TODO
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
