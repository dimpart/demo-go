/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
package mkm

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

//
//  Meta Utils
//

/**
 *  Check whether meta matches with entity ID
 *  <p>
 *      (must call this when received a new meta from network)
 *  </p>
 *
 * @param did - entity ID
 * @return true on matched
 */
func MetaMatchID(did ID, meta Meta) bool {
	// check ID.name
	seed := meta.Seed()
	name := did.Name()
	if name != seed {
		return false
	}
	// check ID.address
	old := did.Address()
	gen := GenerateAddress(meta, old.Network())
	return old.Equal(gen)
}

/**
 *  Check whether meta matches with public key
 *
 * @param pKey - public key
 * @return true on matched
 */
func MetaMatchPublicKey(pKey VerifyKey, meta Meta) bool {
	// check whether the public key equals to meta.key
	if meta.PublicKey().Equal(pKey) {
		return true
	}
	// check with seed & fingerprint
	seed := meta.Seed()
	if seed == "" {
		// NOTICE: ID with BTC/ETH address has no name, so
		//         just compare the key.data to check matching
		return false
	}
	fingerprint := meta.Fingerprint()
	if fingerprint == nil || fingerprint.IsEmpty() {
		// fingerprint should not be empty here
		return false
	}
	// check whether keys equal by verifying signature
	data := UTF8Encode(seed)
	sig := fingerprint.Bytes()
	return pKey.Verify(data, sig)
}

//
//  Document Utils
//

func GetDocumentType(doc Document) string {
	helper := GetGeneralAccountHelper()
	return helper.GetDocumentType(doc.Map(), "")
}

/**
 *  Check whether this time is before old time
 */
func DocumentTimeIsBefore(oldTime, thisTime Time) bool {
	if TimeIsNil(oldTime) || TimeIsNil(thisTime) {
		return false
	}
	return TimeToInt64(thisTime) < TimeToInt64(oldTime)
}

/**
 *  Check whether this document's time is before old document's time
 */
func DocumentIsExpired(thisDoc, oldDoc Document) bool {
	oldTime := oldDoc.Time()
	thisTime := thisDoc.Time()
	return DocumentTimeIsBefore(oldTime, thisTime)
}

/**
 *  Select last document matched the type
 */
func GetLastDocument(documents []Document, docType string) Document {
	if documents == nil {
		return nil
	} else if docType == "*" {
		docType = ""
	}
	checkType := docType != ""

	var lastDoc Document
	var thisType string
	var matched bool
	for _, doc := range documents {
		// 1. check type
		if checkType {
			thisType = GetDocumentType(doc)
			matched = thisType == "" || docType == thisType
			if !matched {
				// type not matched, ignore it
				continue
			}
		}
		// 2. check time
		if lastDoc != nil && DocumentIsExpired(doc, lastDoc) {
			// skip old document
			continue
		}
		// got it
		lastDoc = doc
	}
	return lastDoc
}

/**
 *  Select last visa document
 */
func GetLastVisa(documents []Document) Visa {
	if documents == nil {
		return nil
	}
	var lastVisa Visa
	var thisVisa Visa
	var matched bool
	for _, doc := range documents {
		// 1. check type
		thisVisa, matched = doc.(Visa)
		if !matched {
			// type not matched, ignore it
			continue
		}
		// 2. check time
		if lastVisa != nil && DocumentIsExpired(doc, lastVisa) {
			// skip old document
			continue
		}
		// got it
		lastVisa = thisVisa
	}
	return lastVisa
}

/**
 *  Select last bulletin document
 */
func GetLastBulletin(documents []Document) Bulletin {
	if documents == nil {
		return nil
	}
	var lastBulletin Bulletin
	var thisBulletin Bulletin
	var matched bool
	for _, doc := range documents {
		// 1. check type
		thisBulletin, matched = doc.(Bulletin)
		if !matched {
			// type not matched, ignore it
			continue
		}
		// 2. check time
		if lastBulletin != nil && DocumentIsExpired(doc, lastBulletin) {
			// skip old document
			continue
		}
		// got it
		lastBulletin = thisBulletin
	}
	return lastBulletin
}

//
//  Document Command Builder
//

func RespondDocument(did ID, meta Meta, document Document) DocumentCommand {
	docs := []Document{document}
	return RespondDocuments(did, meta, docs)
}

func RespondDocuments(did ID, meta Meta, documents []Document) DocumentCommand {
	// check document ID
	address := did.Address()
	var docID ID
	for _, doc := range documents {
		docID = ParseID(doc.Get("did"))
		if docID == nil {
			//panic("document ID not found")
			continue
		} else if docID.Equal(address) {
			// OK
			continue
		}
		if docID.Address().Equal(address) {
			// TODO: check ID.name
			continue
		}
		//panic("document ID not matched")
		return nil
	}
	return NewCommandForRespondDocuments(did, meta, documents)
}
