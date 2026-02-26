package db

import (
	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/common/mkm"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- MetaTable

// Override
func (db *Storage) SaveMeta(meta Meta, entity ID) bool {
	// 1. verify meta with ID
	if !MetaMatchID(entity, meta) {
		return false
	}
	// 2. cache it
	db.metaTable[entity.String()] = meta
	// 3. save into local storage
	return saveMeta(db, meta, entity)
}

// Override
func (db *Storage) GetMeta(entity ID) Meta {
	// 1. try from memory cache
	meta := db.metaTable[entity.String()]
	if meta == nil {
		// 2. try from local storage
		meta = loadMeta(db, entity)
		if meta == nil {
			db.metaTable[entity.String()] = emptyMeta // placeholder
		} else {
			db.metaTable[entity.String()] = meta
		}
	} else if meta == emptyMeta {
		meta = nil
	}
	return meta
}

/**
 *  Meta file for Entities (User/Group)
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/meta.js'
 */

func metaPath(db *Storage, did ID) string {
	return PathJoin(db.mkmDir(did), "meta.js")
}

func loadMeta(db *Storage, did ID) Meta {
	path := metaPath(db, did)
	db.log("Loading meta: " + path)
	return ParseMeta(db.readMap(path))
}

func saveMeta(db *Storage, meta Meta, did ID) bool {
	info := meta.Map()
	path := metaPath(db, did)
	db.log("Saving meta: " + path)
	return db.writeMap(path, info)
}

// place holder
var emptyMeta Meta = &BaseMeta{}
