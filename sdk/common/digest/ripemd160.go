package digest

import (
	"crypto"

	. "github.com/dimchat/mkm-go/digest"
	_ "golang.org/x/crypto/ripemd160"
)

func NewRIPEMD160Digester() MessageDigester {
	return &RIPEMD160Digester{}
}

type RIPEMD160Digester struct {
	//MessageDigester
}

// Override
func (RIPEMD160Digester) Digest(data []byte) []byte {
	hash := crypto.RIPEMD160.New()
	hash.Write(data)
	return hash.Sum(nil)
}
