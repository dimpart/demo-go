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
	"strings"

	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- AddressNameTable

// Override
func (db *Storage) GetIdentifier(alias string) ID {
	return db._ansTable[alias]
}

// Override
func (db *Storage) AddRecord(identifier ID, alias string) bool {
	if len(alias) == 0 || ValueIsNil(identifier) {
		return false
	}
	if len(db._ansTable) == 0 {
		panic("ANS not initialized")
	}
	// cache it
	db._ansTable[alias] = identifier
	// save them
	return saveANS(db, db._ansTable)
}

// Override
func (db *Storage) RemoveRecord(alias string) bool {
	if len(alias) == 0 || db._ansTable[alias] == nil {
		return false
	}
	// remove it
	delete(db._ansTable, alias)
	// save them
	return saveANS(db, db._ansTable)
}

/**
 *  Address Name Service
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/ans.txt'
 */

func ansPath(db *Storage) string {
	return PathJoin(db.Root(), "ans.txt")
}

func loadANS(db *Storage) map[string]ID {
	table := make(map[string]ID)

	path := ansPath(db)
	db.log("Loading ANS records: " + path)
	text := db.readText(path)
	lines := strings.Split(text, "\n")
	for _, rec := range lines {
		if len(rec) == 0 {
			// skip empty line
			continue
		}
		pair := strings.Split(rec, "\t")
		if len(pair) != 2 {
			db.error("Invalid ANS record: " + rec)
			continue
		}
		table[strings.TrimSpace(pair[0])] = ParseID(strings.TrimSpace(pair[1]))
	}
	//
	//  Reserved names
	//
	table["all"] = EVERYONE
	table[EVERYONE.Name()] = EVERYONE
	table[ANYONE.Name()] = ANYONE
	table["owner"] = ANYONE
	table["founder"] = FOUNDER

	return table
}

func saveANS(db *Storage, records map[string]ID) bool {
	text := ""
	for key, value := range records {
		if value != nil {
			text += key + "\t" + value.String() + "\n"
		}
	}
	path := ansPath(db)
	db.log("Saving ANS records: " + path)
	return db.writeText(path, text)
}
