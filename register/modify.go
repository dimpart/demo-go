package main

import (
	"fmt"

	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

func showUser(did ID) bool {
	info := getUserInfo(did)
	if info == nil {
		println("user not found: " + did.String())
		return false
	}
	visa := info.Visa
	if visa == nil {
		println("visa not found: " + did.String())
		return false
	}
	name := visa.GetProperty("name")
	avatar := visa.GetProperty("avatar")
	println(fmt.Sprintf("name: \"%s\", avatar: %s", name, avatar))
	return true
}

func doModify(path string, args []string) bool {
	if len(args) > 0 {
		did := ParseID(args[0])
		if did != nil {
			showUser(did)
		}
		// arguments
		name := getOptionString(args, "--name")
		var avatar TransportableFile
		url := getOptionString(args, "--avatar")
		if url != "" {
			avatar = CreateTransportableFile(nil, "", ParseURL(url), nil)
		}
		println("modify: ", name, avatar)
		return true
	}
	doHelp(path, []string{"modify"})
	return false
}
