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
	"fmt"

	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/crypto"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
	. "github.com/dimpart/demo-go/sdk/extensions"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type Database interface {
	PrivateKeyDBI
	MetaDBI
	DocumentDBI

	AddressNameDBI
	LoginDBI

	UserDBI
	ContactDBI
	GroupDBI

	GroupHistoryDBI

	// root directory for database
	SetRoot(root string)
}

/**
 *  Local Storage
 *  ~~~~~~~~~~~~~
 */
type Storage struct {
	//Database

	_root string

	_password SymmetricKey

	//
	//  memory caches
	//

	_identityKeyTable      map[string]PrivateKey   // meta keys: ID -> SK
	_communicationKeyTable map[string][]PrivateKey // visa keys: ID -> []SK
	_decryptionKeyTable    map[string][]DecryptKey // visa keys: ID -> []SK

	_metaTable     map[string]Meta       // meta: ID -> meta
	_documentTable map[string][]Document // document: type -> ID -> doc

	_ansTable map[string]ID // ANS: string -> ID

	_loginCommandTable map[string]LoginCommand    // ID -> Login Command
	_loginMessageTable map[string]ReliableMessage // ID -> Login Message

	_users        []ID
	_contactTable map[string][]ID // user contacts: ID -> []ID

	_memberTable map[string][]ID // group members: ID -> []ID
}

func (db *Storage) Init() Database {

	db._root = "/var/dim"

	db._password = NewPlainKey()

	// private keys
	db._identityKeyTable = make(map[string]PrivateKey, 8)
	db._communicationKeyTable = make(map[string][]PrivateKey, 8)
	db._decryptionKeyTable = make(map[string][]DecryptKey, 8)

	// meta
	db._metaTable = make(map[string]Meta, 1024)

	// documents
	db._documentTable = make(map[string][]Document, 1024)

	// ANS
	db._ansTable = loadANS(db) // make(map[string]ID)

	// login info
	db._loginCommandTable = make(map[string]LoginCommand, 1024)
	db._loginMessageTable = make(map[string]ReliableMessage, 1024)

	// local users
	db._users = make([]ID, 0, 1)
	db._contactTable = make(map[string][]ID, 1)

	// group info
	db._memberTable = make(map[string][]ID, 1024)

	return db
}

/**
 *  Password for private key encryption
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 */
func (db *Storage) Password() SymmetricKey {
	return db._password
}
func (db *Storage) SetPassword(password string) {
	db._password = GeneratePassword(password)
}

/**
 *  Root Directory
 *  ~~~~~~~~~~~~~~
 *
 *  File directory for database
 */
func (db *Storage) Root() string {
	return db._root
}
func (db *Storage) SetRoot(root string) {
	if PathIsExist(root) {
		db._root = root
	} else {
		panic(root)
	}
}

// Directory for MKM entity: '.dim/mkm/{zzz}/{ADDRESS}'
func (db *Storage) mkmDir(identifier ID) string {
	address := identifier.Address().String()
	pos := len(address)
	z := address[pos-1 : pos]
	y := address[pos-2 : pos-1]
	x := address[pos-3 : pos-2]
	w := address[pos-4 : pos-3]
	return PathJoin(db.Root(), "mkm", z, y, x, w, address)
}

func (db *Storage) prepareDir(filepath string) bool {
	dir := PathDir(filepath)
	return MakeDirs(dir)
}

//
//  DOS
//

//func (db *Storage) isExist(path string) bool {
//	return PathIsExist(path)
//}
//func (db *Storage) remove(path string) bool {
//	return PathRemove(path)
//}

func (db *Storage) readText(path string) string {
	return ReadTextFile(path)
}
func (db *Storage) readMap(path string) StringKeyMap {
	info := ReadJSONFile(path)
	if info == nil {
		return nil
	}
	dict, ok := info.(StringKeyMap)
	if !ok {
		return nil
	}
	return dict
}
func (db *Storage) readList(path string) []StringKeyMap {
	info := ReadJSONFile(path)
	if info == nil {
		return nil
	}
	array, ok := info.([]StringKeyMap)
	if !ok {
		return nil
	}
	return array
}
func (db *Storage) readSecret(path string) []byte {
	data := ReadBinaryFile(path)
	if data == nil {
		return nil
	}
	password := db._password
	return password.Decrypt(data, password.Map())
}

func (db *Storage) writeText(path string, text string) bool {
	if !db.prepareDir(path) {
		panic(path)
	}
	return WriteTextFile(path, text)
}
func (db *Storage) writeMap(path string, container StringKeyMap) bool {
	if !db.prepareDir(path) {
		panic(path)
	}
	return WriteJSONFile(path, container)
}
func (db *Storage) writeList(path string, container []StringKeyMap) bool {
	if !db.prepareDir(path) {
		panic(path)
	}
	return WriteJSONFile(path, container)
}
func (db *Storage) writeSecret(path string, data []byte) bool {
	if !db.prepareDir(path) {
		panic(path)
	}
	password := db._password
	binary := password.Encrypt(data, password.Map())
	return WriteBinaryFile(path, binary)
}

//
//  Log
//

func (db *Storage) debug(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogDebug(msg)
}

func (db *Storage) log(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogInfo(msg)
}

func (db *Storage) warning(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogWarning(msg)
}

func (db *Storage) error(msg string) {
	msg = fmt.Sprintf("Storage > %s", msg)
	LogError(msg)
}

// Singleton
var sharedDatabase Database

func SharedDatabase() Database {
	return sharedDatabase
}

func init() {
	db := &Storage{}
	sharedDatabase = db.Init()
}
