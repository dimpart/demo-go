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
	. "github.com/dimpart/demo-go/sdk/utils"
)

// -------------------------------------------------------------------------
//  Meta Utility Functions
// -------------------------------------------------------------------------

// MetaMatchID verifies if a Meta instance matches a given entity ID (DID)
//
// # Critical validation step: MUST be called when receiving new Meta from the network
//
// Validation logic:
//  1. Compares Meta.Seed() (name) with ID.Name()
//  2. Generates address from Meta (using ID's network) and compares with ID.Address()
//
// Parameters:
//   - did  - Entity ID (DID) to validate against the Meta
//   - meta - Meta instance to validate (typically received from network)
//
// Returns: true if Meta matches the ID (name and address are consistent), false otherwise
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

// MetaMatchPublicKey verifies if a Meta instance matches a given public key
//
// Performs two-level validation for maximum compatibility:
//  1. Direct comparison of Meta.PublicKey() with the provided VerifyKey
//  2. Signature verification (Meta.Fingerprint() signs Meta.Seed()) for legacy/compatibility support
//
// Special Case: BTC/ETH addresses have no seed (name) - only direct key comparison is performed
//
// Parameters:
//   - pKey - Public key (VerifyKey) to validate against the Meta
//   - meta - Meta instance to validate
//
// Returns: true if Meta matches the public key (either direct match or valid signature), false otherwise
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

// -------------------------------------------------------------------------
//  Document Utility Functions
// -------------------------------------------------------------------------

// GetDocumentID retrieves the unique identifier for a Document
//
// # Delegates to GeneralAccountHelper for consistent ID generation across document types
//
// Parameters:
//   - doc - Document instance to get ID for
//
// Returns: Unique ID associated with the document (zero value if doc is nil)
func GetDocumentID(doc Document) ID {
	helper := GetGeneralAccountHelper()
	return helper.GetDocumentID(doc.Map())
}

// GetDocumentType retrieves the type identifier for a Document
//
// # Delegates to GeneralAccountHelper with empty default type fallback
//
// Parameters:
//   - doc - Document instance to get type for
//
// Returns: Document type string (empty string if type not found or doc is nil)
func GetDocumentType(doc Document) string {
	helper := GetGeneralAccountHelper()
	return helper.GetDocumentType(doc.Map(), "")
}

// DocumentTimeIsBefore checks if a timestamp is chronologically before another
//
// # Handles nil/empty Time values safely (returns false if either time is nil)
//
// Parameters:
//   - oldTime  - Reference timestamp to compare against
//   - thisTime - Timestamp to check (if it is before oldTime)
//
// Returns: true if thisTime is before oldTime (and both are non-nil), false otherwise
func DocumentTimeIsBefore(oldTime, thisTime Time) bool {
	if TimeIsNil(oldTime) || TimeIsNil(thisTime) {
		return false
	}
	return TimeIsBefore(oldTime, thisTime)
}

// DocumentIsExpired checks if a document is expired relative to another (by timestamp)
//
// # A document is considered expired if its timestamp is before the reference document's timestamp
//
// Parameters:
//   - thisDoc - Document to check for expiration
//   - oldDoc  - Reference document (typically the latest valid document)
//
// Returns: true if thisDoc is expired (older than oldDoc), false otherwise
func DocumentIsExpired(thisDoc, oldDoc Document) bool {
	oldTime := oldDoc.Time()
	thisTime := thisDoc.Time()
	return DocumentTimeIsBefore(oldTime, thisTime)
}

// GetLastDocument retrieves the most recent document of a specified type from a list
//
// # Filters documents by type (supports "*" for all types) and selects the latest by timestamp
//
// Parameters:
//   - documents - Slice of Document instances to filter and sort
//   - docType   - Target document type (use "*" or empty string for all types)
//
// Returns: Most recent Document matching the type (nil if no matching documents)
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

// GetLastVisa retrieves the most recent Visa document from a list
//
// # Filters documents by type (Visa interface) and selects the latest by timestamp
//
// Parameters:
//   - documents - Slice of Document instances to filter for Visa documents
//
// Returns: Most recent Visa document (nil if no Visa documents found)
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

// GetLastBulletin retrieves the most recent Bulletin document from a list
//
// # Filters documents by type (Bulletin interface) and selects the latest by timestamp
//
// Parameters:
//   - documents - Slice of Document instances to filter for Bulletin documents
//
// Returns: Most recent Bulletin document (nil if no Bulletin documents found)
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
