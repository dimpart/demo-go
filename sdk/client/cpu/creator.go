package cpu

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimchat/sdk-go/dkd"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/**
 *  CPU Creator
 *  ~~~~~~~~~~~
 *  Delegate for CPU factory
 */

type ClientContentProcessorCreator struct {
	*BaseContentProcessorCreator
}

//-------- IProcessorCreator

func (creator *ClientContentProcessorCreator) CreateContentProcessor(msgType MessageType) ContentProcessor {
	switch msgType {
	// application customized
	case ContentType.APPLICATION, ContentType.CUSTOMIZED:
		return NewCustomizedContentProcessor(creator.Facebook, creator.Messenger)
	case "application", "customized":
		return NewCustomizedContentProcessor(creator.Facebook, creator.Messenger)
	// forward content
	case "forward":
		return NewForwardContentProcessor(creator.Facebook, creator.Messenger)
	// array content
	case "array":
		return NewArrayContentProcessor(creator.Facebook, creator.Messenger)
	// default commands
	case "command":
		return NewBaseCommandProcessor(creator.Facebook, creator.Messenger)
	//// history command
	//case "history":
	//	return NewHistoryCommandProcessor(creator.Facebook, creator.Messenger)
	// default contents
	case "*":
		// must return a default processor for unknown type
		return NewBaseContentProcessor(creator.Facebook, creator.Messenger)
	}

	// others
	return creator.BaseContentProcessorCreator.CreateContentProcessor(msgType)
}

func (creator *ClientContentProcessorCreator) CreateCommandProcessor(msgType MessageType, cmdName string) ContentProcessor {
	switch cmdName {
	case RECEIPT:
		return NewReceiptCommandProcessor(creator.Facebook, creator.Messenger)
	case HANDSHAKE:
		return NewHandshakeCommandProcessor(creator.Facebook, creator.Messenger)
	case LOGIN:
		return NewLoginCommandProcessor(creator.Facebook, creator.Messenger)
	}
	// others
	return creator.BaseContentProcessorCreator.CreateCommandProcessor(msgType, cmdName)
}

//
//  Factories
//

func NewReceiptCommandProcessor(facebook Facebook, messenger Messenger) ContentProcessor {
	return &ReceiptCommandProcessor{
		BaseCommandProcessor: NewBaseCommandProcessor(facebook, messenger),
	}
}

func NewHandshakeCommandProcessor(facebook Facebook, messenger Messenger) ContentProcessor {
	return &HandshakeCommandProcessor{
		BaseCommandProcessor: NewBaseCommandProcessor(facebook, messenger),
	}
}

func NewLoginCommandProcessor(facebook Facebook, messenger Messenger) ContentProcessor {
	return &LoginCommandProcessor{
		BaseCommandProcessor: NewBaseCommandProcessor(facebook, messenger),
	}
}

func NewCustomizedContentProcessor(facebook Facebook, messenger Messenger) *CustomizedContentProcessor {
	return &CustomizedContentProcessor{
		BaseContentProcessor: NewBaseContentProcessor(facebook, messenger),
	}
}
