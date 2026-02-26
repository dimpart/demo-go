package db

import (
	"strings"

	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- ContactTable

// Override
func (db *Storage) GetContacts(user ID) []ID {
	arr := db.contactTable[user.String()]
	if arr == nil {
		arr = loadContacts(db, user)
		db.contactTable[user.String()] = arr
	}
	return arr
}

func (db *Storage) AddContact(contact ID, user ID) bool {
	arr := db.GetContacts(user)
	for _, item := range arr {
		if contact.Equal(item) {
			// duplicated
			return false
		}
	}
	arr = append(arr, contact)
	return db.SaveContacts(arr, user)
}

func (db *Storage) RemoveContact(contact ID, user ID) bool {
	arr := db.GetContacts(user)
	var pos = -1
	for index, id := range arr {
		if contact.Equal(id) {
			pos = index
			break
		}
	}
	if pos == -1 {
		// contact ID not found
		return false
	}
	arr = append(arr[:pos], arr[pos+1:]...)
	return db.SaveContacts(arr, user)
}

// Override
func (db *Storage) SaveContacts(contacts []ID, user ID) bool {
	db.contactTable[user.String()] = contacts
	return saveContacts(db, user, contacts)
}

/**
 *  Contacts file for User
 *  ~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/protected/{ADDRESS}/contacts.txt'
 */

func contactsPath(db *Storage, user ID) string {
	return PathJoin(db.Root(), "protected", user.Address().String(), "contacts.txt")
}

func loadContacts(db *Storage, user ID) []ID {
	path := contactsPath(db, user)
	db.log("Loading contacts for user: " + user.String())
	text := db.readText(path)
	lines := strings.Split(text, "\n")
	contacts := make([]ID, 0, len(lines))
	for _, rec := range lines {
		id := ParseID(rec)
		if id != nil {
			contacts = append(contacts, id)
		}
	}
	return contacts
}

func saveContacts(db *Storage, user ID, contacts []ID) bool {
	text := ""
	lines := IDRevert(contacts)
	for _, rec := range lines {
		text = text + rec + "\n"
	}
	path := contactsPath(db, user)
	db.log("Saving contacts for user: " + user.String())
	return db.writeText(path, text)
}
