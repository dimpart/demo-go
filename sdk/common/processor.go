/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package dimp

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type ICommonMessageProcessor interface {
	Processor

	GetEntityChecker() IEntityChecker
}

/**
 *  Common Processor
 *  ~~~~~~~~~~~~~~~~
 */
type CommonMessageProcessor struct {
	//ICommonMessageProcessor
	MessageProcessor
}

func (processor *CommonMessageProcessor) Init(facebook Facebook, messenger Messenger) ICommonMessageProcessor {
	if processor.MessageProcessor.Init(facebook, messenger) != nil {
	}
	return processor
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
