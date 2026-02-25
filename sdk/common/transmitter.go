package sdk

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/utils"
)

// Transmitter defines the interface for sending messages with priority-based delivery
//
// Core capabilities: Send message content, instant messages, and reliable messages
//
// Priority Rule: Smaller integer values indicate higher delivery priority (faster transmission)
type Transmitter interface {

	// SendContent constructs and sends a message from sender to receiver with specified priority
	//
	// Creates both InstantMessage (plain content) and ReliableMessage (encrypted/signed)
	//
	// Parameters:
	//   - content  - Message content to send (core payload)
	//   - sender   - ID of the message sender (nil/zero value for current user)
	//   - receiver - ID of the message receiver (user/group/bot)
	//   - priority - Delivery priority (smaller = faster)
	// Returns: Pair[InstantMessage, ReliableMessage] (zero-value Pair if sending fails)
	SendContent(content Content, sender, receiver ID, priority int) Pair[InstantMessage, ReliableMessage]

	// SendInstantMessage sends a plaintext instant message with specified priority
	//
	// Converts the plain InstantMessage to a ReliableMessage for network transmission
	//
	// Parameters:
	//   - iMsg     - Plaintext instant message to send
	//   - priority - Delivery priority (smaller = faster)
	// Returns: ReliableMessage (nil if sending fails)
	SendInstantMessage(iMsg InstantMessage, priority int) ReliableMessage

	// SendReliableMessage sends an encrypted/signed reliable message with specified priority
	//
	// ReliableMessage includes encryption, signature, and delivery guarantees
	//
	// Parameters:
	//   - rMsg     - Encrypted and signed reliable message to send
	//   - priority - Delivery priority (smaller = faster)
	// Returns: true if message sent successfully, false on transmission error
	SendReliableMessage(rMsg ReliableMessage, priority int) bool
}

// MessageTransmitter implements the Transmitter interface
//
// Wraps TwinsHelper to provide message transmission capabilities
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
