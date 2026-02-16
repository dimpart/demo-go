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
	db._documentTable[entity.String()] = array
	// 3. save into local storage
	return saveDocuments(db, entity, array)
}

// Override
func (db *Storage) LoadDocuments(entity ID) []Document {
	// 1. try from memory cache
	docs := db._documentTable[entity.String()]
	if docs == nil {
		// 2. try from local storage
		docs = loadDocuments(db, entity)
		if docs == nil {
			docs = []Document{} // placeholder
		}
		db._documentTable[entity.String()] = docs
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
