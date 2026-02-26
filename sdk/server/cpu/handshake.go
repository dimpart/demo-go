package cpu

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

type HandshakeCommandProcessor struct {
	*BaseCommandProcessor
}

func (cpu *HandshakeCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	hsCmd, _ := cmd.(HandshakeCommand)
	title := hsCmd.Title()
	if title == "DIM?" || title == "DIM!" {
		// S -> C
		return NewTextContent("Handshake command error: " + title)
	}
	// C -> S: Hello world!
	//sessionKey := hsCmd.Session()
	return nil
}
