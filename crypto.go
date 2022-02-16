package main

import (
	"encoding/hex"

	"github.com/nspcc-dev/neo-go/pkg/crypto/hash"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
)

func init() {
	pubKey, _ := keys.NewPublicKeyFromString("021319f4f1ece7181760b876dcc08f91dbffc4b55743c78a9fe827a2cc7a8e2600")
	salt := "1e42534301a51fd977f919850f1dcbc7"
	message := hex.EncodeToString([]byte(salt + "[860833102,0,true]"))
	message = hex.EncodeToString([]byte{byte(len(message) / 2)}) + message
	message = "010001f0" + message + "0000"
	msg, _ := hex.DecodeString(message)
	hashedData := hash.Sha256(msg)
	signature := "00cdb322937e6dbe233c019a9583fd27f4cba2b6055291279eaaeb70fc93d28f1621de6d01a17c2a894232d7b7d0cf54382018e2e2630f0f171f495d8200e1f4"
	sig, _ := hex.DecodeString(signature)
	println(pubKey.Verify(sig, hashedData.BytesBE()))
}

/*
export function num2VarInt(num: number): string {
  if (num < 0xfd) {
    return num2hexstring(num);
  } else if (num <= 0xffff) {
    // uint16
    return "fd" + num2hexstring(num, 2, true);
  } else if (num <= 0xffffffff) {
    // uint32
    return "fe" + num2hexstring(num, 4, true);
  } else {
    // uint64
    return "ff" + num2hexstring(num, 8, true);
  }
}
*/
