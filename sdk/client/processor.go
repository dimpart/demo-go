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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/common"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type ClientMessageProcessor struct {
	CommonMessageProcessor
}

// Override
func (processor *ClientMessageProcessor) GetMessenger() ICommonMessenger {
	messenger := processor.CommonMessageProcessor.GetMessenger()
	cm, ok := messenger.(ICommonMessenger)
	if ok {
		return cm
	}
	return nil
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
