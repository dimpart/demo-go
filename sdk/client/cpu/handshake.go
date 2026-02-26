package cpu

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/client"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

type HandshakeCommandProcessor struct {
	*BaseCommandProcessor
}

func (cpu *HandshakeCommandProcessor) GetMessenger() IClientMessenger {
	messenger := cpu.BaseCommandProcessor.Messenger
	return messenger.(IClientMessenger)
}

// Override
func (cpu *HandshakeCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(HandshakeCommand)
	if !ok {
		return nil
	}
	messenger := cpu.GetMessenger()
	session := messenger.GetClientSession()
	// update station's default ID ('station@anywhere') to sender (real ID)
	station := session.GetStation()
	oid := station.ID()
	sender := rMsg.Sender()
	if oid == nil || oid.IsBroadcast() {
		station.SetID(sender)
	} else if !oid.Equal(sender) {
		panic("station ID not match: " + oid.String() + ", " + sender.String())
	}
	// handle handshake command with title & session key
	title := command.Title()
	newKey := command.SessionKey()
	oldKey := session.GetSessionKey()
	if title == "DIM?" {
		// S -> C: station ask client to handshake again
		if oldKey == "" {
			// first handshake response with new session key
			messenger.Handshake(newKey)
		} else if oldKey == newKey {
			// duplicated handshake response?
			// or session expired and the station ask to handshake again?
			messenger.Handshake(newKey)
		} else {
			// connection changed?
			// erase session key to handshake again
			session.SetSessionKey("")
		}
	} else if title == "DIM!" {
		// S -> C: handshake accepted by station
		if oldKey == "" {
			// normal handshake response,
			// update session key to change state to 'running'
			session.SetSessionKey(newKey)
		} else if oldKey == newKey {
			// duplicated handshake response?
			// set it again here to invoke the flutter channel
			session.SetSessionKey(newKey)
		} else {
			// FIXME: handshake error
			// erase session key to handshake again
			session.SetSessionKey("")
		}
	} else {
		// C -> S: Hello world!
		panic("handshake from other user?")
	}
	return nil
}
