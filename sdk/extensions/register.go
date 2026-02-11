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
package dimp

import (
	"math/rand"

	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type UserInfo struct {
	ID ID

	Meta Meta
	Visa Document

	IdentityKey      SignKey
	CommunicationKey DecryptKey
}

type GroupInfo struct {
	ID ID

	Meta     Meta
	Bulletin Document
}

/**
 *  Generate user account info
 *
 * @param nickname - nickname
 * @param avatar   - photo URL
 * @return user info
 */
func GenerateUserInfo(nickname string, avatar TransportableFile) *UserInfo {
	//
	//  Step 1. generate private key (with asymmetric algorithm)
	//
	idKey := GeneratePrivateKey(ECC)
	//
	//  Step 2. generate meta with private key
	//
	meta := GenerateMeta(ETH, idKey, "")
	//
	//  Step 3. generate ID with meta
	//
	identifier := GenerateID(meta, USER, "")
	//
	//  Step 4. generate visa document and sign with private key
	//
	msgKey := GeneratePrivateKey(RSA)
	visaKey := msgKey.PublicKey().(EncryptKey)
	visa := createVisa(identifier, visaKey, idKey, nickname, avatar)
	//
	//  OK
	//
	return &UserInfo{
		ID:               identifier,
		Meta:             meta,
		Visa:             visa,
		IdentityKey:      idKey,
		CommunicationKey: msgKey.(DecryptKey),
	}
}

/**
 *  Generate bot user account info
 *
 * @param seed   - meta seed for ID.name
 * @param name   - bot name
 * @param avatar - photo URL
 * @return bot info
 */
func GenerateBotInfo(seed string, name string, avatar TransportableFile) *UserInfo {
	if seed == "" {
		seed = "bot"
	}
	//
	//  Step 1. generate private key (with asymmetric algorithm)
	//
	privateKey := GeneratePrivateKey(RSA)
	//
	//  Step 2. generate meta with private key
	//
	meta := GenerateMeta(MKM, privateKey, seed)
	//
	//  Step 3. generate ID with meta
	//
	identifier := GenerateID(meta, BOT, "")
	//
	//  Step 4. generate visa document and sign with private key
	//
	visaKey := privateKey.PublicKey().(EncryptKey)
	visa := createVisa(identifier, visaKey, privateKey, name, avatar)
	//
	//  OK
	//
	return &UserInfo{
		ID:               identifier,
		Meta:             meta,
		Visa:             visa,
		IdentityKey:      privateKey,
		CommunicationKey: privateKey.(DecryptKey),
	}
}

/**
 *  Generate station account info
 *
 * @param seed - meta seed for ID.name
 * @param name - station name
 * @param logo - service provider logo URL
 * @param host - station IP
 * @param port - station port
 * @return station info
 */
func GenerateStationInfo(seed string, name string, logo TransportableFile, host string, port uint16) *UserInfo {
	if seed == "" {
		seed = "station"
	}
	//
	//  Step 1. generate private key (with asymmetric algorithm)
	//
	privateKey := GeneratePrivateKey(RSA)
	//
	//  Step 2. generate meta with private key
	//
	meta := GenerateMeta(MKM, privateKey, seed)
	//
	//  Step 3. generate ID with meta
	//
	identifier := GenerateID(meta, STATION, "")
	//
	//  Step 4. generate visa document and sign with private key
	//
	visaKey := privateKey.PublicKey().(EncryptKey)
	profile := createVisa(identifier, visaKey, nil, name, logo)
	profile.SetProperty("host", host)
	profile.SetProperty("port", port)
	profile.Sign(privateKey)
	//
	//  OK
	//
	return &UserInfo{
		ID:               identifier,
		Meta:             meta,
		Visa:             profile,
		IdentityKey:      privateKey,
		CommunicationKey: privateKey.(DecryptKey),
	}
}

/**
 *  Generate group account info
 *
 * @param founder - group founder
 * @param title   - group name
 * @return Group object
 */
func GenerateGroupInfo(founder *UserInfo, title string, seed string) *GroupInfo {
	//
	//  Step 0. prepare seed for group ID
	//
	if seed == "" {
		r := rand.Int31n(999990000) + 10000 // 10,000 ~ 999,999,999
		seed = "Group-" + string(r)
	}
	//
	//  Step 1. get private key
	//
	privateKey := founder.IdentityKey
	//
	//  Step 2. generate meta with founder's private key
	//
	meta := GenerateMeta(MKM, privateKey, seed)
	//
	//  Step 3. generate ID with meta
	//
	identifier := GenerateID(meta, GROUP, "")
	//
	//  Step 4. generate bulletin document and sign with founder's private key
	//
	bulletin := createBulletin(identifier, privateKey, title, founder.ID)
	//
	//  OK
	//
	return &GroupInfo{
		ID:       identifier,
		Meta:     meta,
		Bulletin: bulletin,
	}
}

//
//  Document Creation
//

func createVisa(uid ID, visaKey EncryptKey, idKey SignKey, nickname string, avatar TransportableFile) Visa {
	doc := &BaseVisa{}
	if doc.Init() != nil {
		doc.SetStringer("did", uid)
		// App ID
		doc.SetProperty("app_id", "chat.dim.tarsier")
		// nickname
		doc.SetName(nickname)
		// avatar
		if avatar != nil {
			doc.SetAvatar(avatar)
		}
		// public key
		doc.SetPublicKey(visaKey)
	}
	if idKey != nil {
		// sign it
		sig := doc.Sign(idKey)
		if sig == nil {
			panic("failed to sign visa: " + uid.String())
			return nil
		}
	}
	return doc
}

func createBulletin(gid ID, privateKey SignKey, title string, founder ID) Bulletin {
	doc := &BaseBulletin{}
	if doc.Init() != nil {
		doc.SetStringer("did", gid)
		// App ID
		doc.SetProperty("app_id", "chat.dim.tarsier")
		// group founder
		doc.SetProperty("founder", founder.String())
		// group name
		doc.SetName(title)
	}
	// sign it
	sig := doc.Sign(privateKey)
	if sig == nil {
		panic("failed to sign bulletin: " + gid.String())
		return nil
	}
	return doc
}
