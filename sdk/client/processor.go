package sdk

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type ClientMessageProcessor struct {
	*CommonMessageProcessor
}

func NewClientMessageProcessor(facebook Facebook, messenger Messenger) *ClientMessageProcessor {
	return &ClientMessageProcessor{
		CommonMessageProcessor: NewCommonMessageProcessor(facebook, messenger),
	}
}

func (processor *ClientMessageProcessor) checkGroupTimes(content Content, rMsg ReliableMessage) {
	// TODO: check 'GDT' & 'GHT' in rMsg
}

func (processor *ClientMessageProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	responses := processor.CommonMessageProcessor.ProcessContent(content, rMsg)

	// check group document & history times from the message
	// to make sure the group info synchronized
	processor.checkGroupTimes(content, rMsg)

	if len(responses) == 0 {
		// respond nothing
		return nil
	} else if _, ok := responses[0].(HandshakeCommand); ok {
		// urgent command
		return responses
	}
	messenger := processor.GetMessenger()
	if messenger == nil {
		panic("messenger not ready")
	}
	receiver := rMsg.Receiver()
	user := processor.SelectLocalUser(receiver)
	if user == nil {
		panic(receiver)
	}
	sender := rMsg.Sender()
	// check responses
	for _, res := range responses {
		if res == nil {
			// should not happen
			continue
		} else if _, ok := res.(ReceiptCommand); ok {
			if sender.Type() == STATION {
				// no need to respond receipt to station
				LogInfo("drop receipt responding to station: " + sender.String())
				continue
			}
		} else if _, ok := res.(TextContent); ok {
			if sender.Type() == STATION {
				// no need to respond text message to station
				LogInfo("drop text msg responding to station: " + sender.String())
				continue
			}
		}
		// normal response
		messenger.SendContent(res, user.ID(), sender, 1)
	}
	// DON'T respond to station directly
	return nil
}
