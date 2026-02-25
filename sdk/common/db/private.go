package db

import (
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

//goland:noinspection GoSnakeCaseUsage
const (
	META_KEY = "M"
	VISA_KEY = "V"
)

// PrivateKeyDBI defines the interface for private key persistence operations
//
// Manages storage and retrieval of user private keys with different usage purposes
type PrivateKeyDBI interface {

	// SavePrivateKey stores a private key associated with a user and usage type
	//
	// Key Type Identifiers:
	//   'M' - Key matches meta.key (for Meta signature/validation)
	//   'V' - Key matches visa.key (for Visa document signature)
	// Parameters:
	//   - key     - Private key to persist (must be valid PrivateKey implementation)
	//   - keyType - Usage type identifier ('M' or 'V')
	//   - user    - User ID to associate with the private key
	// Returns: true if key saved successfully, false on database error
	SavePrivateKey(key PrivateKey, keyType string, user ID) bool

	// GetPrivateKeysForDecryption retrieves all private keys for a user marked for decryption
	//
	// Returns keys implementing the DecryptKey interface (for message/data decryption)
	//
	// Parameters:
	//   - user - User ID to retrieve decryption keys for
	// Returns: Slice of DecryptKey instances (empty slice if no keys found or error)
	GetPrivateKeysForDecryption(user ID) []DecryptKey

	// GetPrivateKeyForSignature retrieves the primary private key for a user marked for general signature
	//
	// Returns the first valid key implementing the SignKey interface (for general-purpose signing)
	//
	// Parameters:
	//   - user - User ID to retrieve signature key for
	// Returns: Primary SignKey instance (nil if no signature key found)
	GetPrivateKeyForSignature(user ID) SignKey

	// GetPrivateKeyForVisaSignature retrieves the private key for a user matched with meta.key (Visa signing)
	//
	// Returns key implementing SignKey interface specifically for Visa document signature
	//
	// Parameters:
	//   - user - User ID to retrieve Visa signature key for
	// Returns: Visa-specific SignKey instance (nil if no matching key found)
	GetPrivateKeyForVisaSignature(user ID) SignKey
}

//
//  Conveniences
//

func ConvertDecryptKeys(privateKeys []PrivateKey) []DecryptKey {
	decryptKeys := make([]DecryptKey, 0, len(privateKeys))
	for _, item := range privateKeys {
		if key, ok := item.(DecryptKey); ok {
			decryptKeys = append(decryptKeys, key)
		}
	}
	return decryptKeys
}

func ConvertPrivateKeys(decryptKeys []DecryptKey) []PrivateKey {
	privateKeys := make([]PrivateKey, 0, len(decryptKeys))
	for _, item := range decryptKeys {
		if key, ok := item.(PrivateKey); ok {
			privateKeys = append(privateKeys, key)
		}
	}
	return privateKeys
}

func RevertPrivateKeys(privateKeys []PrivateKey) []StringKeyMap {
	array := make([]StringKeyMap, len(privateKeys))
	for index, item := range privateKeys {
		array[index] = item.Map()
	}
	return array
}

func InsertPrivateKey(key PrivateKey, privateKeys []PrivateKey) []PrivateKey {
	index := FindPrivateKey(key, privateKeys)
	if index == 0 {
		// nothing change
		return nil
	} else if index > 0 {
		// move to the front
		privateKeys = append(privateKeys[:index], privateKeys[index+1:]...)
	} else {
		// private key not found,
		// prepare to insert
		size := len(privateKeys)
		if size > 2 {
			// keep only last three records
			privateKeys = privateKeys[:size]
		}
	}
	// insert to the front
	return append([]PrivateKey{key}, privateKeys...)
}

func FindPrivateKey(key PrivateKey, privateKeys []PrivateKey) int {
	data := key.GetString("data", "")
	if data == "" {
		panic("private key data not found")
	}
	for index, item := range privateKeys {
		if item.Get("data") == data {
			return index
		}
	}
	return -1
}
