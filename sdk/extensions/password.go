package sdk

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/types"
)

var KeySize = 32
var BlockSize = 16

/**
 *  This is for generating symmetric key with a text string
 */
func GeneratePassword(password string) SymmetricKey {
	data := UTF8Encode(password)
	digest := SHA256(data)
	// AES key data
	dataLen := len(data)
	length := KeySize - dataLen
	if length > 0 {
		// format: {digest_prefix}+{pwd_data}
		merged := make([]byte, KeySize)
		BytesCopy(digest, 0, merged, 0, uint(length))
		BytesCopy(data, 0, merged, uint(length), uint(dataLen))
		data = merged
	} else if length < 0 {
		data = digest
	}
	// AES iv
	iv := digest[KeySize-BlockSize:]
	// generate AES key
	key := NewMap()
	key["algorithm"] = AES
	key["data"] = Base64Encode(data)
	key["iv"] = Base64Encode(iv)
	return ParseSymmetricKey(key)
}
