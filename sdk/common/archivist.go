package dimp

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/mem"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/mkm"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/common/mkm"
)

type ICommonArchivist interface {
	Archivist
	Barrack

	//GetFacebook() Facebook
}

type CommonArchivist struct {
	//ICommonArchivist

	Facebook Facebook
	Database AccountDBI
}

func (archivist *CommonArchivist) Init(facebook Facebook, database AccountDBI) ICommonArchivist {
	archivist.Facebook = facebook
	archivist.Database = database
	return archivist
}

//
//  Barrack
//

// Override
func (archivist *CommonArchivist) CacheUser(user User) {
	if user.DataSource() == nil {
		user.SetDataSource(archivist.Facebook)
	}
	did := user.ID()
	sharedUserCache.Put(did.String(), user)
}

// Override
func (archivist *CommonArchivist) CacheGroup(group Group) {
	if group.DataSource() == nil {
		group.SetDataSource(archivist.Facebook)
	}
	did := group.ID()
	sharedGroupCache.Put(did.String(), group)
}

// Override
func (archivist *CommonArchivist) GetUser(uid ID) User {
	return sharedUserCache.Get(uid.String())
}

// Override
func (archivist *CommonArchivist) GetGroup(group ID) Group {
	return sharedGroupCache.Get(group.String())
}

// Override
func (archivist *CommonArchivist) CreateUser(uid ID) User {
	network := uid.Type()
	// check user type
	switch network {
	case STATION:
		return NewStation(uid)
	case BOT:
		return NewBot(uid)
	}
	// general user, or 'anyone@anywhere'
	user := &BaseUser{}
	return user.Init(uid)
}

// Override
func (archivist *CommonArchivist) CreateGroup(gid ID) Group {
	network := gid.Type()
	// check group type
	switch network {
	case ISP:
		return NewProvider(gid)
	}
	// general group, or 'everyone@everywhere'
	group := &BaseGroup{}
	return group.Init(gid)
}

//
//  Archivist
//

// Override
func (archivist *CommonArchivist) SaveMeta(meta Meta, did ID) bool {
	//
	//  1. check valid
	//
	valid := archivist.checkMeta(meta, did)
	if !valid {
		//panic("meta not valid")
		return false
	}
	//
	//  2. check duplicated
	facebook := archivist.Facebook
	old := facebook.GetMeta(did)
	if old != nil {
		return true
	}
	//
	//  3. save into database
	db := archivist.Database
	return db.SaveMeta(meta, did)
}

// protected
func (archivist *CommonArchivist) checkMeta(meta Meta, did ID) bool {
	return meta.IsValid() && MetaMatchID(did, meta)
}

// Override
func (archivist *CommonArchivist) SaveDocument(doc Document, did ID) bool {
	//
	//  1. check valid
	//
	valid := archivist.checkDocumentValid(doc, did)
	if !valid {
		// document not valid
		return false
	}
	//
	//  2. check expired
	//
	if archivist.checkDocumentExpired(doc, did) {
		// drop expired document
		return false
	}
	//
	//  3. save into database
	db := archivist.Database
	return db.SaveDocument(doc, did)
}

// protected
func (archivist *CommonArchivist) checkDocumentValid(doc Document, did ID) bool {
	docTime := doc.Time()
	// check document time
	if docTime == nil {
		//panic("document error")
	} else {
		// calibrate the clock
		// make sure the document time is not in the far future
		nearFuture := TimeToFloat64(TimeNow()) + 1800
		if TimeToFloat64(docTime) > nearFuture {
			//panic("document time error")
			return false
		}
	}
	// check valid
	return archivist.verifyDocument(doc, did)
}

// protected
func (archivist *CommonArchivist) verifyDocument(doc Document, did ID) bool {
	/*/
	if doc.IsValid() {
		return true
	}
	// check ID
	docID := ParseID(doc.Get("did"))
	if docID == nil {
		//panic("document ID not found")
		return false
	} else if !docID.Address().Equal(did.Address()) {
		// ID not matched
		return false
	}
	/*/
	facebook := archivist.Facebook
	// verify with meta.key
	meta := facebook.GetMeta(did)
	if meta == nil {
		//panic("failed to get meta: " + did.String())
		return false
	}
	metaKey := meta.PublicKey()
	return doc.Verify(metaKey)
}

// protected
func (archivist *CommonArchivist) checkDocumentExpired(doc Document, did ID) bool {
	docType := GetDocumentType(doc)
	if docType == "" {
		docType = "*"
	}
	facebook := archivist.Facebook
	// check old documents with type
	documents := facebook.GetDocuments(did)
	old := GetLastDocument(documents, docType)
	return old != nil && DocumentIsExpired(doc, old)
}

// Override
func (archivist *CommonArchivist) LocalUsers() []ID {
	db := archivist.Database
	return db.GetLocalUsers()
}

/**
 * Call it when received 'UIApplicationDidReceiveMemoryWarningNotification',
 * this will remove 50% of cached objects
 *
 * @return number of survivors
 */
func (archivist *CommonArchivist) ReduceMemory() int {
	addressCache := GetAddressCache()
	idCache := GetIDCache()
	cnt1 := addressCache.ReduceMemory()
	cnt2 := idCache.ReduceMemory()
	cnt3 := sharedUserCache.ReduceMemory()
	cnt4 := sharedGroupCache.ReduceMemory()
	return cnt1 + cnt2 + cnt3 + cnt4
}

var sharedUserCache MemoryCache[string, User] = NewThanosCache[string, User]()
var sharedGroupCache MemoryCache[string, Group] = NewThanosCache[string, Group]()
