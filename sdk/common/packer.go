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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common/dkd"
	. "github.com/dimpart/demo-go/sdk/common/mkm"
)

type ICommonMessagePacker interface {
	Packer

	GetArchivist() Archivist

	// for checking whether user's ready
	GetVisaKey(user ID) EncryptKey

	/**
	 *  Check sender before verifying received message
	 *
	 * @param rMsg - network message
	 * @return false on verify key not found
	 */
	CheckSender(rMsg ReliableMessage) bool

	/**
	 *  Check receiver before encrypting message
	 *
	 * @param iMsg - plain message
	 * @return false on encrypt key not found
	 */
	CheckReceiver(iMsg InstantMessage) bool

	/**
	 *  Check meta &amp; visa
	 *
	 * @param rMsg - received message
	 * @return false on error
	 */
	CheckAttachments(rMsg ReliableMessage) bool
}

type IMessageWaitingQueue interface {

	/**
	 *  Add income message in a queue for waiting sender's visa
	 *
	 * @param rMsg - incoming message
	 * @param info - error info
	 */
	SuspendReliableMessage(rMsg ReliableMessage, info StringKeyMap)

	/**
	 *  Add outgo message in a queue for waiting receiver's visa
	 *
	 * @param iMsg - outgo message
	 * @param info - error info
	 */
	SuspendInstantMessage(iMsg InstantMessage, info StringKeyMap)
}

type CommonMessagePacker struct {
	//ICommonMessagePacker
	MessagePacker

	// protected
	Queue IMessageWaitingQueue
}

func (packer *CommonMessagePacker) Init(facebook IFacebook, messenger IMessenger) ICommonMessagePacker {
	if packer.MessagePacker.Init(facebook, messenger) != nil {
		// TODO:
		packer.Queue = nil
	}
	return packer
}

// protected
func (packer *CommonMessagePacker) GetArchivist() Archivist {
	return packer.Facebook.GetArchivist()
}

// protected
func (packer *CommonMessagePacker) GetVisaKey(user ID) EncryptKey {
	facebook := packer.Facebook
	documents := facebook.GetDocuments(user)
	doc := GetLastVisa(documents)
	if doc != nil /*&& doc.IsValid()*/ {
		return doc.PublicKey()
	}
	meta := facebook.GetMeta(user)
	if meta != nil /*&& meta.IsValid()*/ {
		metaKey := meta.PublicKey()
		if encKey, ok := metaKey.(EncryptKey); ok {
			return encKey
		}
	}
	return nil
}

// protected
func (packer *CommonMessagePacker) CheckSender(rMsg ReliableMessage) bool {
	sender := rMsg.Sender()
	// check sender's meta & document
	visa := GetVisaAttachment(rMsg)
	if visa != nil {
		// first handshake?
		did := ParseID(visa.Get("did"))
		if did == nil {
			//panic("visa error")
		} else if sender.Equal(did) {
			return true
		}
		//panic("visa ID not matched")
		return false
	} else if packer.GetVisaKey(sender) != nil {
		// sender is OK
		return true
	}
	// sender not ready, suspend message for waiting document
	info := StringKeyMap{
		"message": "verify key not found",
		"user":    sender.String(),
	}
	packer.Queue.SuspendReliableMessage(rMsg, info)
	return false
}

// protected
func (packer *CommonMessagePacker) CheckReceiver(iMsg InstantMessage) bool {
	receiver := iMsg.Receiver()
	if receiver.IsBroadcast() {
		// broadcast message
		return true
	} else if receiver.IsGroup() {
		// NOTICE: station will never send group message, so
		//         we don't need to check group info here; and
		//         if a client wants to send group message,
		//         that should be sent to a group bot first,
		//         and the bot will split it for all members.
		return false
	} else if packer.GetVisaKey(receiver) != nil {
		// receiver is OK
		return true
	}
	// receiver not ready, suspend message for waiting document
	info := StringKeyMap{
		"message": "encrypt key not found",
		"user":    receiver.String(),
	}
	packer.Queue.SuspendInstantMessage(iMsg, info)
	return false
}

// Override
func (packer *CommonMessagePacker) EncryptMessage(iMsg InstantMessage) SecureMessage {
	// 1. check contact info
	// 2. check group members info
	if !packer.CheckReceiver(iMsg) {
		// receiver not ready
		return nil
	}
	return packer.MessagePacker.EncryptMessage(iMsg)
}

// protected
func (packer *CommonMessagePacker) CheckAttachments(rMsg ReliableMessage) bool {
	archivist := packer.GetArchivist()
	if archivist == nil {
		panic("archivist not ready")
		return false
	}
	sender := rMsg.Sender()
	// [Meta Protocol]
	meta := GetMetaAttachment(rMsg)
	if meta != nil {
		archivist.SaveMeta(meta, sender)
	}
	// [Visa Protocol]
	visa := GetVisaAttachment(rMsg)
	if visa != nil {
		archivist.SaveDocument(visa, sender)
	}
	//
	//  TODO: check [Visa Protocol] before calling this
	//        make sure the sender's meta(visa) exists
	//        (do in by application)
	//
	return true
}

// Override
func (packer *CommonMessagePacker) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	// 1. check receiver/group with local user
	// 2. check sender's visa info
	if !packer.CheckSender(rMsg) {
		// sender not ready
		return nil
	}
	// make sure meta exists before verifying message
	if !packer.CheckAttachments(rMsg) {
		return nil
	}
	return packer.MessagePacker.VerifyMessage(rMsg)
}

// Override
func (packer *CommonMessagePacker) SignMessage(sMsg SecureMessage) ReliableMessage {
	if rMsg, ok := sMsg.(ReliableMessage); ok {
		// already signed
		return rMsg
	}
	return packer.MessagePacker.SignMessage(sMsg)
}
