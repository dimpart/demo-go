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
package main

import (
	"fmt"
	"strings"

	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/extensions"
)

func getUserInfo(identifier ID) *UserInfo {
	facebook := SharedFacebook()
	return &UserInfo{
		ID:               identifier,
		Meta:             facebook.GetMeta(identifier),
		Visa:             facebook.GetDocument(identifier, VISA),
		IdentityKey:      facebook.GetPrivateKeyForVisaSignature(identifier),
		CommunicationKey: facebook.GetPrivateKeysForDecryption(identifier)[0],
	}
}

func saveInfo(identifier ID, meta Meta, doc Document, idKey SignKey, msgKey DecryptKey) bool {
	fmt.Println("******** ID:", identifier)
	//fmt.Println("******** meta:", meta.Map())
	//fmt.Println("******** doc:", doc.Map())
	//fmt.Println("******** id key:", idKey.Map())
	//fmt.Println("******** msg key:", msgKey.Map())
	facebook := SharedFacebook()
	database := facebook.GetDatabase()
	// id key
	identityKey, ok := idKey.(PrivateKey)
	if ok && identityKey != nil {
		if database.SavePrivateKey(identityKey, "M", identifier) == false {
			return false
		}
	}
	// msg key
	communicationKey, ok := msgKey.(PrivateKey)
	if ok && communicationKey != nil {
		if database.SavePrivateKey(communicationKey, "V", identifier) == false {
			return false
		}
	}
	// meta
	if facebook.SaveMeta(meta, identifier) == false {
		return false
	}
	// document
	if facebook.SaveDocument(doc, identifier) == false {
		return false
	}
	// OK
	return true
}

func doGenerate(path string, args []string) bool {
	if len(args) > 0 {
		// arguments
		seed := getOptionString(args, "--seed")
		name := getOptionString(args, "--name")
		var avatar TransportableFile
		url := getOptionString(args, "--avatar")
		if url != "" {
			avatar = CreateTransportableFile(nil, "", ParseURL(url), nil)
		}
		// check account type
		aType := strings.ToLower(args[0])
		if aType == "user" {
			info := GenerateUserInfo(name, avatar)
			return saveInfo(info.ID, info.Meta, info.Visa, info.IdentityKey, info.CommunicationKey)
		} else if aType == "group" {
			founder := ParseID(getOptionString(args, "--founder"))
			if founder != nil {
				info := GenerateGroupInfo(getUserInfo(founder), name, seed)
				return saveInfo(info.ID, info.Meta, info.Bulletin, nil, nil)
			}
		} else if aType == "station" {
			var logo TransportableFile
			url = getOptionString(args, "--logo")
			if url != "" {
				logo = CreateTransportableFile(nil, "", ParseURL(url), nil)
			}
			host := getOptionString(args, "--host")
			port := getOptionInteger(args, "--port")
			info := GenerateStationInfo(seed, name, logo, host, uint16(port))
			return saveInfo(info.ID, info.Meta, info.Visa, info.IdentityKey, info.CommunicationKey)
		} else if aType == "robot" {
			info := GenerateBotInfo(seed, name, avatar)
			return saveInfo(info.ID, info.Meta, info.Visa, info.IdentityKey, info.CommunicationKey)
		}
	}
	doHelp(path, []string{"generate"})
	return false
}
