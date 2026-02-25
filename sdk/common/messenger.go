package sdk

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
	Messenger
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
 *  Message Transceiver
 */

type CommonMessenger struct {
	//ICommonMessenger
	*BaseMessenger

	// protected
	Session     Session
	Facebook    ICommonFacebook
	Transmitter Transmitter
}

func NewCommonMessenger(session Session, facebook ICommonFacebook, database CipherKeyDelegate) *CommonMessenger {
	return &CommonMessenger{
		BaseMessenger: NewBaseMessenger(facebook, database),
		Session:       session,
		Facebook:      facebook,
		Transmitter:   nil, // NewMessageTransmitter(facebook, messenger)
	}
}

func (messenger *CommonMessenger) GetSession() Session {
	return messenger.Session
}

func (messenger *CommonMessenger) GetFacebook() ICommonFacebook {
	return messenger.Facebook
}

func (messenger *CommonMessenger) GetCipherKeyDelegate() CipherKeyDelegate {
	return messenger.CipherKeyDelegate
}

func (messenger *CommonMessenger) GetMessagePacker() Packer {
	return messenger.Packer
}

func (messenger *CommonMessenger) SetMessagePacker(packer Packer) {
	messenger.Packer = packer
}

func (messenger *CommonMessenger) GetMessageProcessor() Processor {
	return messenger.Processor
}

func (messenger *CommonMessenger) SetMessageProcessor(processor Processor) {
	messenger.Processor = processor
}

//-------- ITransceiver

//// Override
//func (messenger *CommonMessenger) SerializeMessage(rMsg ReliableMessage) []byte {
//	// fix meta attachment
//	// fix visa attachment
//	return messenger.BaseMessenger.SerializeMessage(rMsg)
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
	rMsg := messenger.BaseMessenger.DeserializeMessage(data)
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
	return messenger.BaseMessenger.EncryptKey(data, receiver, iMsg)
}

// Override
func (messenger *CommonMessenger) EncodeKey(bundle EncryptedBundle, receiver ID, iMsg InstantMessage) StringKeyMap {
	keys := messenger.BaseMessenger.EncodeKey(bundle, receiver, iMsg)
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
		return messenger.BaseMessenger.SerializeKey(password, iMsg)
	}
	// 1. remove before serializing key
	password.Remove("reused")
	password.Remove("digest")
	// 2. serialize key without flags
	data := messenger.BaseMessenger.SerializeKey(password, iMsg)
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
//	return messenger.BaseMessenger.SerializeContent(content, password, iMsg)
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
