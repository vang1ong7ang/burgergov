package main

import (
	"github.com/nspcc-dev/neo-go/pkg/util"
)

type scripthash struct {
	util.Uint160
}

func (me scripthash) MarshalText() (text []byte, err error) {
	return []byte(`0x` + me.StringLE()), nil
}

func ScripthashDecodeStringLE(s string) (scripthash, error) {
	sh, e := util.Uint160DecodeStringLE(s)
	return scripthash{sh}, e
}
