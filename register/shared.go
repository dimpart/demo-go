package main

import (
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/client"
	. "github.com/dimpart/demo-go/sdk/client/ext"
	. "github.com/dimpart/demo-go/sdk/common"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/database"
)

var clientFacebook IClientFacebook = nil

func createArchivist(facebook Facebook, adb AccountDBI) ICommonArchivist {
	archivist := &CommonArchivist{}
	return archivist.Init(facebook, adb)
}

func createEntityChecker(adb AccountDBI) IEntityChecker {
	emitter := &CheckEmitter{}
	checker := &EntityChecker{}
	checker.Init(adb)
	checker.Request = emitter
	checker.Respond = emitter
	return checker
}

func createDatabase() Database {
	db := &Storage{}
	return db.Init()
}

func createFacebook() IClientFacebook {
	// create database
	database := createDatabase()
	// create facebook with database
	facebook := &ClientFacebook{}
	facebook.Init(database)
	// set archivist & barrack
	archivist := createArchivist(facebook, database)
	facebook.Archivist = archivist
	facebook.Barrack = archivist
	// set entity checker
	facebook.Checker = createEntityChecker(database)
	return facebook
}

func SharedFacebook() IClientFacebook {
	facebook := clientFacebook
	if facebook == nil {
		facebook = createFacebook()
		clientFacebook = facebook
	}
	return facebook
}

func init() {
	lib := &LibraryLoader{}
	lib.Init(nil, nil)
	lib.Run()
}
