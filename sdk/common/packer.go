package sdk

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

// ICommonMessagePacker defines the interface for common message packing/validation
//
// Extends Packer with sender/receiver validation and attachment checks
//
// Ensures messages are properly encrypted/signed before transmission
type ICommonMessagePacker interface {
	Packer

	// GetVisaKey retrieves the encryption key from a user's Visa document
	//
	// Used to check if a user is ready for encrypted communication
	//
	// Parameters:
	//   - user - User ID to retrieve Visa encryption key for
	// Returns: EncryptKey (nil if user's Visa key not found)
	GetVisaKey(user ID) EncryptKey

	// CheckSender validates the sender before verifying a received message
	//
	// Verifies sender identity and retrieves necessary verification keys
	//
	// Parameters:
	//   - rMsg - Received reliable message to validate sender for
	// Returns: true if sender is valid (verify key found), false if verify key missing/invalid
	CheckSender(rMsg ReliableMessage) bool

	// CheckReceiver validates the receiver before encrypting an outbound message
	//
	// Verifies receiver identity and retrieves necessary encryption keys
	//
	// Parameters:
	//   - iMsg - Outbound instant message to validate receiver for
	// Returns: true if receiver is valid (encrypt key found), false if encrypt key missing/invalid
	CheckReceiver(iMsg InstantMessage) bool

	// CheckAttachments validates Meta and Visa attachments in a received message
	//
	// Ensures attached metadata and identity documents are valid and unmodified
	//
	// Parameters:
	//   - rMsg - Received reliable message to check attachments for
	// Returns: true if attachments are valid, false on validation error
	CheckAttachments(rMsg ReliableMessage) bool
}

// IMessageWaitingQueue defines the interface for message queuing during validation delays
//
// Manages suspended messages waiting for sender/receiver Visa documents
type IMessageWaitingQueue interface {

	// SuspendReliableMessage adds an incoming message to queue waiting for sender's Visa
	//
	// Used when sender's Visa is missing/invalid (message held until Visa is available)
	//
	// Parameters:
	//   - rMsg - Incoming reliable message to suspend
	//   - info - StringKeyMap containing error/queue metadata (why message was suspended)
	SuspendReliableMessage(rMsg ReliableMessage, info StringKeyMap)

	// SuspendInstantMessage adds an outbound message to queue waiting for receiver's Visa
	//
	// Used when receiver's Visa is missing/invalid (message held until Visa is available)
	//
	// Parameters:
	//   - iMsg - Outbound instant message to suspend
	//   - info - StringKeyMap containing error/queue metadata (why message was suspended)
	SuspendInstantMessage(iMsg InstantMessage, info StringKeyMap)
}

// CommonMessagePacker implements the ICommonMessagePacker interface
//
// Wraps MessagePacker and integrates with IMessageWaitingQueue for suspended messages
type CommonMessagePacker struct {
	//ICommonMessagePacker
	*MessagePacker

	// Queue manages suspended messages waiting for Visa validation
	//
	// Protected field: Should not be accessed directly by external packages
	Queue IMessageWaitingQueue
}

func NewCommonMessagePacker(facebook Facebook, messenger Messenger) *CommonMessagePacker {
	return &CommonMessagePacker{
		MessagePacker: NewMessagePacker(facebook, messenger),
		// TODO:
		Queue: nil,
	}
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
	facebook := packer.Facebook
	sender := rMsg.Sender()
	// [Meta Protocol]
	meta := GetMetaAttachment(rMsg)
	if meta != nil {
		facebook.SaveMeta(meta, sender)
	}
	// [Visa Protocol]
	visa := GetVisaAttachment(rMsg)
	if visa != nil {
		facebook.SaveDocument(visa, sender)
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
