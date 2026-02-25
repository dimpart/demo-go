package digest

import (
	"crypto/sha1"

	. "github.com/dimchat/mkm-go/digest"
)

func NewSHA1Digester() MessageDigester {
	return &SHA1Digester{}
}

type SHA1Digester struct {
	//MessageDigester
}

// Override
func (SHA1Digester) Digest(data []byte) []byte {
	hash := sha1.Sum(data)
	return hash[:]
}

//
//  SHA-1
//

var sha1Digester MessageDigester = nil

func SetSHA1Digester(digester MessageDigester) {
	sha1Digester = digester
}

func SHA1(bytes []byte) []byte {
	return sha1Digester.Digest(bytes)
}
