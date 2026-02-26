package db

import (
	. "github.com/dimchat/core-go/msg"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- LoginTable

// Override
func (db *Storage) GetLoginCommandMessage(user ID) Pair[LoginCommand, ReliableMessage] {
	cmd, msg := getLoginInfo(db, user)
	return NewPair[LoginCommand, ReliableMessage](cmd, msg)
}

// Override
func (db *Storage) SaveLoginCommandMessage(user ID, cmd LoginCommand, msg ReliableMessage) bool {
	if !cacheLoginInfo(db, user, cmd, msg) {
		return false
	}
	return saveLoginInfo(db, user, cmd, msg)
}

/**
 *  Login info for Users
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/login.js'
 */

func loginInfoPath(db *Storage, identifier ID) string {
	return PathJoin(db.mkmDir(identifier), "login.js")
}

func loadLoginInfo(db *Storage, identifier ID) (cmd LoginCommand, msg ReliableMessage) {
	path := loginInfoPath(db, identifier)
	db.log("Loading login info: " + path)
	info := db.readMap(path)
	if info == nil {
		return nil, nil
	}
	cmd, _ = ParseContent(info["cmd"]).(LoginCommand)
	msg = ParseReliableMessage(info["msg"])
	return cmd, msg
}

func saveLoginInfo(db *Storage, user ID, cmd LoginCommand, msg ReliableMessage) bool {
	info := NewMap()
	info["cmd"] = cmd.Map()
	info["msg"] = msg.Map()
	path := loginInfoPath(db, user)
	db.log("Saving login info: " + path)
	return db.writeMap(path, info)
}

// place holder
var emptyMessage = &NetworkMessage{}

func getLoginInfo(db *Storage, user ID) (cmd LoginCommand, msg ReliableMessage) {
	// 1. try from memory cache
	msg = db.loginMessageTable[user.String()]
	if msg == nil {
		// 2. try from local storage
		cmd, msg = loadLoginInfo(db, user)
		if msg == nil {
			// place an empty message for cache
			db.loginMessageTable[user.String()] = emptyMessage
		} else {
			// cache them
			db.loginCommandTable[user.String()] = cmd
			db.loginMessageTable[user.String()] = msg
		}
	} else if msg == emptyMessage {
		cmd = nil
		msg = nil
	} else {
		cmd = db.loginCommandTable[user.String()]
	}
	return cmd, msg
}

func cacheLoginInfo(db *Storage, user ID, cmd LoginCommand, msg ReliableMessage) bool {
	// 1. verify sender ID
	if msg.Sender().Equal(user) == false {
		return false
	}
	// 2. check last login time
	old, _ := getLoginInfo(db, user)
	if old != nil {
		oldTime := old.Time().Unix()
		newTime := cmd.Time().Unix()
		if newTime <= oldTime {
			// expired command, drop it
			return false
		}
	}
	// 3. cache them
	db.loginCommandTable[user.String()] = cmd
	db.loginMessageTable[user.String()] = msg
	return true
}
