package network

import . "github.com/dimchat/dkd-go/protocol"

type Departure interface {
}

// IMessageQueue defines the interface for managing outgoing message queues with departure tracking
//
// Core functionality: Append outgoing messages to queue while preventing duplicates
type IMessageQueue interface {

	// Append adds an outgoing reliable message to the queue with departure metadata
	//
	// Associates the message with a Departure (ship) for delivery tracking
	// Prevents duplicate messages from being added to the queue
	//
	// Parameters:
	//   - rMsg - Outgoing reliable message to queue (encrypted/signed)
	//   - ship - Departure metadata for delivery tracking (e.g., target station, priority)
	// Returns: true if message added successfully, false if duplicate message
	Append(rMsg ReliableMessage, ship Departure) bool
}

type MessageQueue struct {
	//IMessageQueue
}
