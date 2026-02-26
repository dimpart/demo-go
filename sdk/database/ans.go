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
	return db.ansTable[alias]
}

// Override
func (db *Storage) AddRecord(identifier ID, alias string) bool {
	if len(alias) == 0 || ValueIsNil(identifier) {
		return false
	}
	if len(db.ansTable) == 0 {
		panic("ANS not initialized")
	}
	// cache it
	db.ansTable[alias] = identifier
	// save them
	return saveANS(db, db.ansTable)
}

// Override
func (db *Storage) RemoveRecord(alias string) bool {
	if len(alias) == 0 || db.ansTable[alias] == nil {
		return false
	}
	// remove it
	delete(db.ansTable, alias)
	// save them
	return saveANS(db, db.ansTable)
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
