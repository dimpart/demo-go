package sdk

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
	visaKey, ok := msgKey.PublicKey().(EncryptKey)
	if !ok {
		panic("message key error")
	}
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
	visaKey, ok := privateKey.PublicKey().(EncryptKey)
	if !ok {
		panic("private key error")
	}
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
	doc := NewBaseVisa(nil, "", nil)
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
	doc := NewBaseBulletin(nil, "", nil)
	doc.SetStringer("did", gid)
	// App ID
	doc.SetProperty("app_id", "chat.dim.tarsier")
	// group founder
	doc.SetProperty("founder", founder.String())
	// group name
	doc.SetName(title)
	// sign it
	sig := doc.Sign(privateKey)
	if sig == nil {
		panic("failed to sign bulletin: " + gid.String())
		return nil
	}
	return doc
}
