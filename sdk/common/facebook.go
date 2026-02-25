package sdk

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
	Facebook

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
	*BaseFacebook

	Database AccountDBI

	Checker IEntityChecker

	currentUser User
}

func NewCommonFacebook(db EntityDataSource) *CommonFacebook {
	return &CommonFacebook{
		BaseFacebook: NewBaseFacebook(db),
		Database:     nil,
		Checker:      nil,
		currentUser:  nil,
	}
}

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
	users := facebook.Database.GetLocalUsers()
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
	return facebook.BaseFacebook.SelectUser(receiver)
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
	return facebook.BaseFacebook.SelectMember(members)
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
	meta := facebook.Database.GetMeta(did)
	facebook.Checker.CheckMeta(did, meta)
	return meta
}

// Override
func (facebook *CommonFacebook) GetDocuments(did ID) []Document {
	docs := facebook.Database.GetDocuments(did)
	facebook.Checker.CheckDocuments(did, docs)
	return docs
}

//
//  User DataSource
//

// Override
func (facebook *CommonFacebook) GetContacts(user ID) []ID {
	return facebook.Database.GetContacts(user)
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
