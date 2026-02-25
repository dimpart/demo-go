package digest

import (
	. "github.com/dimchat/mkm-go/digest"
	"golang.org/x/crypto/sha3"
)

func NewKECCAK256Digester() MessageDigester {
	return &KECCAK256Digester{}
}

type KECCAK256Digester struct {
	//MessageDigester
}

// Override
func (KECCAK256Digester) Digest(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}
