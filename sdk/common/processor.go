package sdk

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/utils"
)

// ICommonMessageProcessor defines the interface for common message processing
//
// Extends Processor with access to entity validation utilities
type ICommonMessageProcessor interface {
	Processor

	// GetEntityChecker retrieves the entity validation helper
	//
	// Used to check entity data freshness during message processing
	//
	// Returns: IEntityChecker instance for entity validation
	GetEntityChecker() IEntityChecker
}

// CommonMessageProcessor implements the ICommonMessageProcessor interface
//
// Wraps MessageProcessor to provide common message processing capabilities
type CommonMessageProcessor struct {
	//ICommonMessageProcessor
	*MessageProcessor
}

func NewCommonMessageProcessor(facebook Facebook, messenger Messenger) *CommonMessageProcessor {
	return &CommonMessageProcessor{
		MessageProcessor: NewMessageProcessor(facebook, messenger),
	}
}

func (processor *CommonMessageProcessor) GetFacebook() ICommonFacebook {
	facebook := processor.MessageProcessor.Facebook
	return facebook.(ICommonFacebook)
}

func (processor *CommonMessageProcessor) GetMessenger() ICommonMessenger {
	messenger := processor.MessageProcessor.Messenger
	return messenger.(ICommonMessenger)
}

// protected
func (processor *CommonMessageProcessor) GetEntityChecker() IEntityChecker {
	facebook := processor.GetFacebook()
	return facebook.GetEntityChecker()
}

func (processor *CommonMessageProcessor) checkVisaTime(content Content, rMsg ReliableMessage) bool {
	checker := processor.GetEntityChecker()
	if checker == nil {
		panic("should not happen")
		return false
	}
	docUpdate := false
	// check sender document time
	lastDocTime := rMsg.GetTime("SDT", nil)
	if lastDocTime != nil {
		now := TimeNow()
		if TimeIsAfter(now, lastDocTime) {
			// calibrate the clock
			lastDocTime = now
		}
		sender := rMsg.Sender()
		docUpdate = checker.SetLastDocumentTime(sender, lastDocTime)
		// check whether it needs update now
		if docUpdate {
			// checking for new visa
			facebook := processor.GetFacebook()
			facebook.GetDocuments(sender)
		}
	}
	return docUpdate
}

// Override
func (processor *CommonMessageProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	responses := processor.MessageProcessor.ProcessContent(content, rMsg)

	// check sender's document times from the message
	// to make sure the user info synchronized
	processor.checkVisaTime(content, rMsg)

	return responses
}
