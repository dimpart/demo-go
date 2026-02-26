package db

import (
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/plugins-go/crypto"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- PrivateKeyTable

// Override
func (db *Storage) SavePrivateKey(key PrivateKey, keyType string, user ID) bool {
	if keyType == META_KEY {
		if !cacheIdentityKey(db, user, key) {
			return false
		}
		return saveIdentityKey(db, user, key)
	}
	if !cacheCommunicationKey(db, user, key) {
		return false
	}
	keys := getCommunicationKeys(db, user)
	return saveCommunicationKeys(db, user, keys)
}

// Override
func (db *Storage) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return getDecryptionKeys(db, user)
}

// Override
func (db *Storage) GetPrivateKeyForSignature(user ID) SignKey {
	keys := getCommunicationKeys(db, user)
	if len(keys) > 0 {
		// sign message with communication key
		return keys[0]
	}
	// if communication keys not exists, use identity key to sign message
	return getIdentityKey(db, user)
}

// Override
func (db *Storage) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return getIdentityKey(db, user)
}

/**
 *  Private Key file for Local Users
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  1. Identity Key      - paired to meta.key, CONSTANT
 *     file path: '.dim/private/{ADDRESS}/secret'
 *
 *  2. Communication Key - paired to visa.key, VOLATILE
 *     file path: '.dim/private/{ADDRESS}/secret_keys'
 */

func identityKeyPath(db *Storage, user ID) string {
	return PathJoin(db.Root(), "private", user.Address().String(), "secret.js")
}

func communicationKeysPath(db *Storage, user ID) string {
	return PathJoin(db.Root(), "private", user.Address().String(), "secret_keys.js")
}

func loadIdentityKey(db *Storage, user ID) PrivateKey {
	path := identityKeyPath(db, user)
	db.log("Loading identity key: " + path)
	data := db.readSecret(path)
	if data == nil {
		return nil
	}
	json := UTF8Decode(data)
	dict := JSONDecodeMap(json)
	return ParsePrivateKey(dict)
}
func loadCommunicationKeys(db *Storage, user ID) []PrivateKey {
	keys := make([]PrivateKey, 0, 1)
	path := communicationKeysPath(db, user)
	db.log("Loading communication keys: " + path)
	data := db.readSecret(path)
	if data != nil {
		json := UTF8Decode(data)
		arr, _ := JSONDecode(json).([]interface{})
		for _, item := range arr {
			k := ParsePrivateKey(item)
			if k == nil {
				panic(item)
			}
			keys = append(keys, k)
		}
	}
	return keys
}

func saveIdentityKey(db *Storage, user ID, key PrivateKey) bool {
	info := key.Map()
	path := identityKeyPath(db, user)
	db.log("Saving identity key: " + path)
	json := JSONEncodeMap(info)
	data := UTF8Encode(json)
	return db.writeSecret(path, data)
}
func saveCommunicationKeys(db *Storage, user ID, keys []PrivateKey) bool {
	arr := make([]interface{}, 0, len(keys))
	for _, item := range keys {
		arr = append(arr, item.Map())
	}
	path := communicationKeysPath(db, user)
	db.log("Saving communication keys: " + path)
	json := JSONEncode(arr)
	data := UTF8Encode(json)
	return db.writeSecret(path, data)
}

// place holder
var emptyPrivateKey PrivateKey = &ECCPrivateKey{}

func getIdentityKey(db *Storage, user ID) PrivateKey {
	// 1. try from memory cache
	key := db.identityKeyTable[user.String()]
	if key == nil {
		// 2. try from local storage
		key = loadIdentityKey(db, user)
		if key == nil {
			// place an empty key for cache
			db.identityKeyTable[user.String()] = emptyPrivateKey
		} else {
			// cache it
			db.identityKeyTable[user.String()] = key
		}
	} else if key == emptyPrivateKey {
		db.error("Private key not found: " + user.String())
		key = nil
	}
	return key
}

func getCommunicationKeys(db *Storage, user ID) []PrivateKey {
	// 1. try from memory cache
	keys := db.communicationKeyTable[user.String()]
	if keys == nil {
		// 2. try from local storage
		keys = loadCommunicationKeys(db, user)
		// 3. cache them
		db.communicationKeyTable[user.String()] = keys
	}
	return keys
}
func getDecryptionKeys(db *Storage, user ID) []DecryptKey {
	// 1. try from memory cache
	keys := db.decryptionKeyTable[user.String()]
	if len(keys) == 0 {
		var decKey DecryptKey
		var ok bool
		// 2. get communication keys
		msgKeys := getCommunicationKeys(db, user)
		keys = make([]DecryptKey, 0, len(msgKeys)+1)
		for _, item := range msgKeys {
			decKey, ok = item.(DecryptKey)
			if ok && decKey != nil {
				keys = append(keys, decKey)
			}
		}
		// 3. check identity key
		idKey := getIdentityKey(db, user)
		decKey, ok = idKey.(DecryptKey)
		if ok && decKey != nil && findKey(msgKeys, idKey) < 0 {
			keys = append(keys, decKey)
		}
		// 4. cache them
		db.decryptionKeyTable[user.String()] = keys
	}
	return keys
}

func cacheIdentityKey(db *Storage, user ID, key PrivateKey) bool {
	old := getIdentityKey(db, user)
	if old != nil {
		// identity key won't change
		return false
	}
	db.identityKeyTable[user.String()] = key
	return true
}

func cacheCommunicationKey(db *Storage, user ID, key PrivateKey) bool {
	keys := getCommunicationKeys(db, user)
	index := findKey(keys, key)
	if index == 0 {
		return false // nothing changed
	} else if index > 0 {
		keys = removeKey(keys, index) // move to the front
	} else if len(keys) > 2 {
		keys = keys[:2] // keep only last three records
	}
	keys = insertKey(keys, key)
	db.communicationKeyTable[user.String()] = keys
	// reset decryption keys
	delete(db.decryptionKeyTable, user.String())
	return true
}

func findKey(keys []PrivateKey, key PrivateKey) int {
	for index, item := range keys {
		if key.Equal(item) {
			return index
		}
	}
	return -1
}
func removeKey(keys []PrivateKey, index int) []PrivateKey {
	return append(keys[:index], keys[index+1:]...)
}
func insertKey(keys []PrivateKey, key PrivateKey) []PrivateKey {
	arr := make([]PrivateKey, 0, len(keys)+1)
	arr = append(arr, key)
	return append(arr, keys...)
}
