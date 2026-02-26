package db

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- DocumentTable

// Override
func (db *Storage) SaveDocument(doc Document, entity ID) bool {
	// 1. check valid
	if !doc.IsValid() {
		return false
	}
	// TODO: check old documents
	array := make([]Document, 1)
	array[0] = doc
	// 2. cache it
	db.documentTable[entity.String()] = array
	// 3. save into local storage
	return saveDocuments(db, entity, array)
}

// Override
func (db *Storage) GetDocuments(entity ID) []Document {
	// 1. try from memory cache
	docs := db.documentTable[entity.String()]
	if docs == nil {
		// 2. try from local storage
		docs = loadDocuments(db, entity)
		if docs == nil {
			docs = []Document{} // placeholder
		}
		db.documentTable[entity.String()] = docs
	}
	return docs
}

/**
 *  Document for Entities (User/Group)
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/documents.js'
 */

func documentPath(db *Storage, did ID) string {
	return PathJoin(db.mkmDir(did), "documents.js")
}

func loadDocuments(db *Storage, did ID) []Document {
	path := documentPath(db, did)
	db.log("Loading document: " + path)
	array := db.readList(path)
	return DocumentConvert(array)
}

func saveDocuments(db *Storage, did ID, docs []Document) bool {
	path := documentPath(db, did)
	db.log("Saving document: " + path)
	array := DocumentRevert(docs)
	return db.writeList(path, array)
}
