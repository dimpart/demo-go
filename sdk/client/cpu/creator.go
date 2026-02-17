/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
 *
 *  Delegate for CPU factory
 */
type ClientContentProcessorCreator struct {
	BaseContentProcessorCreator
}

//-------- IProcessorCreator

func (creator *ClientContentProcessorCreator) CreateContentProcessor(msgType MessageType) ContentProcessor {
	switch msgType {
	// application customized
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

func NewReceiptCommandProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	cpu := &ReceiptCommandProcessor{}
	cpu.Init(facebook, messenger)
	return cpu
}

func NewHandshakeCommandProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	cpu := &HandshakeCommandProcessor{}
	cpu.Init(facebook, messenger)
	return cpu
}

func NewLoginCommandProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	cpu := &LoginCommandProcessor{}
	cpu.Init(facebook, messenger)
	return cpu
}
