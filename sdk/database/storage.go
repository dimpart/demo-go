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

	root string

	password SymmetricKey

	//
	//  memory caches
	//

	identityKeyTable      map[string]PrivateKey   // meta keys: ID -> SK
	communicationKeyTable map[string][]PrivateKey // visa keys: ID -> []SK
	decryptionKeyTable    map[string][]DecryptKey // visa keys: ID -> []SK

	metaTable     map[string]Meta       // meta: ID -> meta
	documentTable map[string][]Document // document: type -> ID -> doc

	ansTable map[string]ID // ANS: string -> ID

	loginCommandTable map[string]LoginCommand    // ID -> Login Command
	loginMessageTable map[string]ReliableMessage // ID -> Login Message

	users        []ID
	contactTable map[string][]ID // user contacts: ID -> []ID

	memberTable map[string][]ID // group members: ID -> []ID
}

func NewStorage(root string) *Storage {
	db := &Storage{

		root: root,

		password: NewPlainKey(),

		// private keys
		identityKeyTable:      make(map[string]PrivateKey, 8),
		communicationKeyTable: make(map[string][]PrivateKey, 8),
		decryptionKeyTable:    make(map[string][]DecryptKey, 8),

		// meta
		metaTable: make(map[string]Meta, 1024),

		// documents
		documentTable: make(map[string][]Document, 1024),

		// ANS
		ansTable: make(map[string]ID, 1024),

		// login info
		loginCommandTable: make(map[string]LoginCommand, 1024),
		loginMessageTable: make(map[string]ReliableMessage, 1024),

		// local users
		users:        make([]ID, 0, 1),
		contactTable: make(map[string][]ID, 1),

		// group info
		memberTable: make(map[string][]ID, 1024),
	}
	// load ANS
	db.ansTable = loadANS(db)
	// OK
	return db
}

/**
 *  Password for private key encryption
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 */
func (db *Storage) Password() SymmetricKey {
	return db.password
}
func (db *Storage) SetPassword(password string) {
	db.password = GeneratePassword(password)
}

/**
 *  Root Directory
 *  ~~~~~~~~~~~~~~
 *
 *  File directory for database
 */
func (db *Storage) Root() string {
	return db.root
}
func (db *Storage) SetRoot(root string) {
	if PathIsExist(root) {
		db.root = root
	} else {
		panic(root)
	}
}

// Directory for MKM entity: '.dim/mkm/{zzz}/{ADDRESS}'
func (db *Storage) mkmDir(identifier ID) string {
	address := identifier.Address().String()
	//pos := len(address)
	//z := address[pos-1 : pos]
	//y := address[pos-2 : pos-1]
	//x := address[pos-3 : pos-2]
	//w := address[pos-4 : pos-3]
	//return PathJoin(db.Root(), "mkm", z, y, x, w, address)
	return PathJoin(db.Root(), "public", address)
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
func (db *Storage) readList(path string) []interface{} {
	info := ReadJSONFile(path)
	if info == nil {
		return nil
	}
	array, ok := info.([]interface{})
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
	password := db.password
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
	password := db.password
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
