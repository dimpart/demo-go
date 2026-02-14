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
	"fmt"

	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/crypto"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type ICommonMessenger interface {
	IMessenger
	Transmitter

	GetSession() Session

	GetFacebook() ICommonFacebook

	GetCipherKeyDelegate() CipherKeyDelegate

	GetMessagePacker() Packer
	SetMessagePacker(packer Packer)

	GetMessageProcessor() Processor
	SetMessageProcessor(processor Processor)
}

/**
 *  Common Messenger
 *  ~~~~~~~~~~~~~~~~
 */
type CommonMessenger struct {
	//ICommonMessenger
	Messenger

	// protected
	Session     Session
	Facebook    ICommonFacebook
	Transmitter Transmitter
}

func (messenger *CommonMessenger) Init(session Session, facebook ICommonFacebook, database CipherKeyDelegate) ICommonMessenger {
	if messenger.Messenger.Init(facebook) != nil {
		messenger.Session = session
		messenger.Facebook = facebook
		//messenger.Transmitter = (&MessageTransmitter{}).Init(facebook, messenger)
		messenger.Messenger.CipherKeyDelegate = database
		messenger.Messenger.Packer = nil
		messenger.Messenger.Processor = nil
	}
	return messenger
}

func (messenger *CommonMessenger) GetSession() Session {
	return messenger.Session
}

func (messenger *CommonMessenger) GetFacebook() ICommonFacebook {
	return messenger.Facebook
}

func (messenger *CommonMessenger) GetCipherKeyDelegate() CipherKeyDelegate {
	return messenger.Messenger.CipherKeyDelegate
}

func (messenger *CommonMessenger) GetMessagePacker() Packer {
	return messenger.Messenger.Packer
}

func (messenger *CommonMessenger) SetMessagePacker(packer Packer) {
	messenger.Messenger.Packer = packer
}

func (messenger *CommonMessenger) GetMessageProcessor() Processor {
	return messenger.Messenger.Processor
}

func (messenger *CommonMessenger) SetMessageProcessor(processor Processor) {
	messenger.Messenger.Processor = processor
}

//-------- ITransceiver

//// Override
//func (messenger *CommonMessenger) SerializeMessage(rMsg ReliableMessage) []byte {
//	// fix meta attachment
//	// fix visa attachment
//	return messenger.Messenger.SerializeMessage(rMsg)
//}

// Override
func (messenger *CommonMessenger) DeserializeMessage(data []byte) ReliableMessage {
	size := len(data)
	if size <= 8 {
		// message data error
		return nil
		//} else if data[0] != '{' || data[size-1] != '}' {
		//	// only support JSON format now
		//	return nil
	}
	rMsg := messenger.Messenger.DeserializeMessage(data)
	if rMsg != nil {
		// fix meta attachment
		// fix visa attachment
	}
	return rMsg
}

//-------- IInstantMessageDelegate

// Override
func (messenger *CommonMessenger) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) EncryptedBundle {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed to encrypt key: ", r)
		}
	}()
	return messenger.Messenger.EncryptKey(data, receiver, iMsg)
}

// Override
func (messenger *CommonMessenger) EncodeKey(bundle EncryptedBundle, receiver ID, iMsg InstantMessage) StringKeyMap {
	keys := messenger.Messenger.EncodeKey(bundle, receiver, iMsg)
	if len(keys) > 0 {
		// TODO: fixEncodeKeys
	}
	return keys
}

// Override
func (messenger *CommonMessenger) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	// TODO: reuse message key

	// 0. check message key
	reused := password.Get("reused")
	digest := password.Get("digest")
	if reused == nil && digest == nil {
		// flags not exists, serialize it directly
		return messenger.Messenger.SerializeKey(password, iMsg)
	}
	// 1. remove before serializing key
	password.Remove("reused")
	password.Remove("digest")
	// 2. serialize key without flags
	data := messenger.Messenger.SerializeKey(password, iMsg)
	// 3. put it back after serialized
	if ConvertBool(reused, false) {
		password.Set("reused", true)
	}
	if digest != nil {
		password.Set("digest", digest)
	}
	return data
}

//// Override
//func (messenger *CommonMessenger) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
//	// fix content
//	return messenger.Messenger.SerializeContent(content, password, iMsg)
//}

//
//  Interfaces for Transmitting Message
//

// Override
func (messenger *CommonMessenger) SendContent(content Content, sender, receiver ID, priority int) Pair[InstantMessage, ReliableMessage] {
	transmitter := messenger.Transmitter
	return transmitter.SendContent(content, sender, receiver, priority)
}

// Override
func (messenger *CommonMessenger) SendInstantMessage(iMsg InstantMessage, priority int) ReliableMessage {
	transmitter := messenger.Transmitter
	return transmitter.SendInstantMessage(iMsg, priority)
}

// Override
func (messenger *CommonMessenger) SendReliableMessage(rMsg ReliableMessage, priority int) bool {
	transmitter := messenger.Transmitter
	return transmitter.SendReliableMessage(rMsg, priority)
}
