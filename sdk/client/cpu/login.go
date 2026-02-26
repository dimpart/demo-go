package cpu

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/client"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

type LoginCommandProcessor struct {
	*BaseCommandProcessor
}

func (cpu *LoginCommandProcessor) GetMessenger() IClientMessenger {
	messenger := cpu.BaseCommandProcessor.Messenger
	return messenger.(IClientMessenger)
}

// private
func (cpu *LoginCommandProcessor) getDatabase() SessionDBI {
	messenger := cpu.GetMessenger()
	session := messenger.GetSession()
	return session.GetDatabase()
}

// Override
func (cpu *LoginCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(LoginCommand)
	if !ok {
		return nil
	}
	sender := command.ID()
	// save login command to session db
	db := cpu.getDatabase()
	if db.SaveLoginCommandMessage(sender, command, rMsg) {
		// OK
	} else {
		panic("failed to save login command")
	}
	// no need to response login command
	return nil
}
