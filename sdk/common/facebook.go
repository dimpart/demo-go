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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/mkm"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/common/mkm"
	. "github.com/dimpart/demo-go/sdk/extensions"
)

type ICommonFacebook interface {
	IFacebook

	GetDatabase() AccountDBI
	GetEntityChecker() IEntityChecker

	//
	//  Current User
	//
	GetCurrentUser() User
	SetCurrentUser(user User)

	//
	//  Documents
	//
	GetDocument(did ID, docType string) Document
	GetVisa(uid ID) Visa
	GetBulletin(gid ID) Bulletin

	GetName(did ID) string
}

/**
 *  Common Facebook
 *  ~~~~~~~~~~~~~~~
 *  Barrack for Server/Client
 */
type CommonFacebook struct {
	//ICommonFacebook
	Facebook

	Database AccountDBI

	Checker IEntityChecker

	currentUser User
}

//func (facebook *CommonFacebook) Init() ICommonFacebook {
//	if facebook.Facebook.Init() != nil {
//		facebook.Database = nil
//		facebook.Checker = nil
//		facebook.currentUser = nil
//	}
//	return facebook
//}

func (facebook *CommonFacebook) GetDatabase() AccountDBI {
	return facebook.Database
}

func (facebook *CommonFacebook) GetEntityChecker() IEntityChecker {
	return facebook.Checker
}

//
//  Current User
//

func (facebook *CommonFacebook) GetCurrentUser() User {
	currentUser := facebook.currentUser
	if currentUser != nil {
		return currentUser
	}
	users := facebook.Database.LoadLocalUsers()
	if len(users) == 0 {
		return nil
	}
	currentUser = facebook.GetUser(users[0])
	facebook.currentUser = currentUser
	return currentUser
}

func (facebook *CommonFacebook) SetCurrentUser(user User) {
	//if user.DataSource() == nil {
	//	user.SetDataSource(facebook)
	//}
	facebook.currentUser = user
}

// Override
func (facebook *CommonFacebook) SelectUser(receiver ID) ID {
	currentUser := facebook.GetCurrentUser()
	if currentUser != nil {
		current := currentUser.ID()
		if receiver.IsBroadcast() {
			// broadcast message can be decrypted by anyone, so
			// just return current user here
			return current
		} else if receiver.Equal(current) {
			return current
		}
	}
	// check local users
	return facebook.Facebook.SelectUser(receiver)
}

// Override
func (facebook *CommonFacebook) SelectMember(members []ID) ID {
	currentUser := facebook.GetCurrentUser()
	if currentUser != nil {
		// group message (recipient not designated)
		current := currentUser.ID()
		// the messenger will check group info before decrypting message,
		// so we can trust that the group's meta & members MUST exist here.
		for _, member := range members {
			if member.Equal(current) {
				return current
			}
		}
	}
	// check local users
	return facebook.Facebook.SelectMember(members)
}

//
//  Documents
//

func (facebook *CommonFacebook) GetDocument(did ID, docType string) Document {
	documents := facebook.GetDocuments(did)
	doc := GetLastDocument(documents, docType)
	// compatible for document type
	if doc == nil && docType == VISA {
		doc = GetLastDocument(documents, PROFILE)
	}
	return doc
}

func (facebook *CommonFacebook) GetVisa(uid ID) Visa {
	documents := facebook.GetDocuments(uid)
	return GetLastVisa(documents)
}

func (facebook *CommonFacebook) GetBulletin(gid ID) Bulletin {
	documents := facebook.GetDocuments(gid)
	return GetLastBulletin(documents)
}

func (facebook *CommonFacebook) GetName(did ID) string {
	var docType string
	if did.IsUser() {
		docType = VISA
	} else if did.IsGroup() {
		docType = BULLETIN
	} else {
		docType = "*"
	}
	// get name from document
	doc := facebook.GetDocument(did, docType)
	if doc != nil {
		name := ConvertString(doc.GetProperty("name"), "")
		if name != "" {
			return name
		}
	}
	// get name from ID
	return AnonymousGetName(did)
}

//
//  Entity DataSource
//

// Override
func (facebook *CommonFacebook) GetMeta(did ID) Meta {
	meta := facebook.Database.LoadMeta(did)
	facebook.Checker.CheckMeta(did, meta)
	return meta
}

// Override
func (facebook *CommonFacebook) GetDocuments(did ID) []Document {
	docs := facebook.Database.LoadDocuments(did)
	facebook.Checker.CheckDocuments(did, docs)
	return docs
}

//
//  User DataSource
//

// Override
func (facebook *CommonFacebook) GetContacts(user ID) []ID {
	return facebook.Database.LoadContacts(user)
}

// Override
func (facebook *CommonFacebook) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return facebook.Database.GetPrivateKeysForDecryption(user)
}

// Override
func (facebook *CommonFacebook) GetPrivateKeyForSignature(user ID) SignKey {
	return facebook.Database.GetPrivateKeyForSignature(user)
}

// Override
func (facebook *CommonFacebook) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return facebook.Database.GetPrivateKeyForVisaSignature(user)
}
