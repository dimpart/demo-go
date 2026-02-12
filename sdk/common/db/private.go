/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
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

type PrivateKeyDBI interface {

	/**
	 *  Save private key for user
	 *
	 * @param user - user ID
	 * @param key - private key
	 * @param type - 'M' for matching meta.key; or 'V' for matching visa.key
	 * @return false on error
	 */
	SavePrivateKey(key PrivateKey, keyType string, user ID) bool

	/**
	 *  Get private keys for user
	 *
	 * @param user - user ID
	 * @return all keys marked for decryption
	 */
	GetPrivateKeysForDecryption(user ID) []DecryptKey

	/**
	 *  Get private key for user
	 *
	 * @param user - user ID
	 * @return first key marked for signature
	 */
	GetPrivateKeyForSignature(user ID) PrivateKey

	/**
	 *  Get private key for user
	 *
	 * @param user - user ID
	 * @return the private key matched with meta.key
	 */
	GetPrivateKeyForVisaSignature(user ID) PrivateKey
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
