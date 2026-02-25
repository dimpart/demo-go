package digest

import (
	"crypto/md5"

	. "github.com/dimchat/mkm-go/digest"
)

func NewMD5Digester() MessageDigester {
	return &MD5Digester{}
}

type MD5Digester struct {
	//MessageDigester
}

// Override
func (MD5Digester) Digest(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}

//
//  MD-5
//

var md5Digester MessageDigester = nil

func SetMD5Digester(digester MessageDigester) {
	md5Digester = digester
}

func MD5(bytes []byte) []byte {
	return md5Digester.Digest(bytes)
}
