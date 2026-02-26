package cpu

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type LoginCommandProcessor struct {
	*BaseCommandProcessor
}

func (cpu *LoginCommandProcessor) Execute(cmd Command, rMsg ReliableMessage) Content {
	sender := rMsg.Sender()
	info := NewMap()
	info["ID"] = sender.String()
	info["cmd"] = cmd.Map()
	// post notification: USER_ONLINE
	NotificationPost("user_online", cpu, info)
	return NewReceiptCommand("Login received", rMsg.Envelope(), cmd)
}
