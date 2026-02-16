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
	db._metaTable[entity.String()] = meta
	// 3. save into local storage
	return saveMeta(db, meta, entity)
}

// Override
func (db *Storage) LoadMeta(entity ID) Meta {
	// 1. try from memory cache
	meta := db._metaTable[entity.String()]
	if meta == nil {
		// 2. try from local storage
		meta = loadMeta(db, entity)
		if meta == nil {
			db._metaTable[entity.String()] = emptyMeta // placeholder
		} else {
			db._metaTable[entity.String()] = meta
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
