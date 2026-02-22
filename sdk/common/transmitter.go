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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type Transmitter interface {

	/**
	 *  Send content from sender to receiver with priority
	 *
	 * @param sender   - from where, null for current user
	 * @param receiver - to where
	 * @param content  - message content
	 * @param priority - smaller is faster
	 * @return (iMsg, None) on error
	 */
	SendContent(content Content, sender, receiver ID, priority int) Pair[InstantMessage, ReliableMessage]

	/**
	 *  Send instant message with priority
	 *
	 * @param iMsg     - plain message
	 * @param priority - smaller is faster
	 * @return null on error
	 */
	SendInstantMessage(iMsg InstantMessage, priority int) ReliableMessage

	/**
	 *  Send reliable message with priority
	 *
	 * @param rMsg     - encrypted &amp; signed message
	 * @param priority - smaller is faster
	 * @return false on error
	 */
	SendReliableMessage(rMsg ReliableMessage, priority int) bool
}

type MessageTransmitter struct {
	//Transmitter
	*TwinsHelper
}

func NewMessageTransmitter(facebook Facebook, messenger Messenger) *MessageTransmitter {
	return &MessageTransmitter{
		TwinsHelper: NewTwinsHelper(facebook, messenger),
	}
}

// protected
func (transmitter *MessageTransmitter) GetFacebook() ICommonFacebook {
	facebook := transmitter.TwinsHelper.Facebook
	return facebook.(ICommonFacebook)
}

// protected
func (transmitter *MessageTransmitter) GetMessenger() ICommonMessenger {
	messenger := transmitter.TwinsHelper.Messenger
	return messenger.(ICommonMessenger)
}

// Override
func (transmitter *MessageTransmitter) SendContent(content Content, sender, receiver ID, priority int) Pair[InstantMessage, ReliableMessage] {
	if sender == nil {
		facebook := transmitter.GetFacebook()
		current := facebook.GetCurrentUser()
		sender = current.ID()
	}
	env := CreateEnvelope(sender, receiver, nil)
	iMsg := CreateInstantMessage(env, content)
	messenger := transmitter.GetMessenger()
	rMsg := messenger.SendInstantMessage(iMsg, priority)
	return NewPair[InstantMessage, ReliableMessage](iMsg, rMsg)
}

// private
func (transmitter *MessageTransmitter) attachVisaTime(sender ID, iMsg InstantMessage) bool {
	content := iMsg.Content()
	if _, ok := content.(Command); ok {
		// no need to attach times for command
		return false
	}
	facebook := transmitter.GetFacebook()
	doc := facebook.GetVisa(sender)
	if doc == nil {
		panic("failed to get visa document for sender: " + sender.String())
		return false
	}
	// attach sender document time
	lastDocTime := doc.Time()
	if TimeIsNil(lastDocTime) {
		//panic("document error")
	} else {
		iMsg.SetTime("SDT", lastDocTime)
	}
	return true
}

// Override
func (transmitter *MessageTransmitter) SendInstantMessage(iMsg InstantMessage, priority int) ReliableMessage {
	sender := iMsg.Sender()
	//
	//  0. check cycled message
	//
	if iMsg.Receiver().Equal(sender) {
		//panic("drop cycled message")
		return nil
	}
	// attach sender's document times
	// for the receiver to check whether user info synchronized
	transmitter.attachVisaTime(sender, iMsg)
	//
	//  1. encrypt message
	//
	messenger := transmitter.GetMessenger()
	sMsg := messenger.EncryptMessage(iMsg)
	if sMsg == nil {
		panic("public key not found")
		return nil
	}
	//
	//  2. sign message
	//
	rMsg := messenger.SignMessage(sMsg)
	if rMsg == nil {
		// TODO: set msg.state = error
		panic("failed to sign message")
	}
	//
	//  3. send message
	//
	ok := messenger.SendReliableMessage(rMsg, priority)
	if !ok {
		// failed
		return nil
	}
	return rMsg
}

// Override
func (transmitter *MessageTransmitter) SendReliableMessage(rMsg ReliableMessage, priority int) bool {
	sender := rMsg.Sender()
	// 0. check cycled message
	if rMsg.Receiver().Equal(sender) {
		//panic("drop cycled message")
		return false
	}
	// 1. serialize message
	messenger := transmitter.GetMessenger()
	data := messenger.SerializeMessage(rMsg)
	if data == nil {
		//panic("failed to serialize message")
		return false
	}
	// 2. call gatekeeper to send the message data package
	//    put message package into the waiting queue of current session
	session := messenger.GetSession()
	return session.QueueMessagePackage(rMsg, data, priority)
}
