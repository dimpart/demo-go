package main

import (
	. "github.com/dimpart/demo-go/sdk/client"
	. "github.com/dimpart/demo-go/sdk/client/ext"
	. "github.com/dimpart/demo-go/sdk/common"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/database"
)

var clientFacebook IClientFacebook = nil

func createEntityChecker(adb AccountDBI) IEntityChecker {
	emitter := &CheckEmitter{}
	checker := NewEntityChecker(adb)
	checker.Request = emitter
	checker.Respond = emitter
	return checker
}

func createFacebook() IClientFacebook {
	// create database for facebook
	database := NewStorage("/var/dim")
	facebook := NewClientFacebook(database)
	// set archivist & barrack
	archivist := NewCommonArchivist(facebook, database)
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
	lib := NewLibraryLoader()
	lib.Run()
}
