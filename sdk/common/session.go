package sdk

import (
	"fmt"

	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/db"
)

// SocketAddress defines the interface for network socket address (host + port)
type SocketAddress interface {
	fmt.Stringer

	// Host returns the IP/domain name part of the socket address
	//
	// Returns: Host string (e.g., "192.168.1.1", "example.com")
	Host() string

	// Port returns the numeric port part of the socket address
	//
	// Returns: Port number as uint16 (e.g., 8080, 9394)
	Port() uint16
}

// Session defines the interface for a network communication session
//
// Extends Transmitter with session management, database access, and message queuing
type Session interface {
	Transmitter

	// GetDatabase retrieves the session-specific database interface
	//
	// Returns: SessionDBI instance for session-related data operations
	GetDatabase() SessionDBI

	// GetRemoteAddress retrieves the remote socket address of the session
	//
	// Returns: SocketAddress containing remote host and port
	GetRemoteAddress() SocketAddress

	// GetSessionKey returns the cryptographic session key for secure communication
	//
	// Returns: Session key string (empty string if not established)
	GetSessionKey() string

	/**
	 *  Update user ID
	 *
	 * @param uid - login user ID
	 * @return true on changed
	 */

	// GetID returns the current user ID associated with the session
	//
	// Returns: User ID (nil/zero value if no user logged in)
	GetID() ID
	SetID(uid ID) bool

	// IsActive checks if the session is currently active
	//
	// Returns: true if session is active, false if inactive
	IsActive() bool

	// SetActive updates the session's active status and timestamp
	//
	// Parameters:
	//   - active - Active status flag (true = session active, false = inactive)
	//   - when   - Timestamp (typically current time) of status change
	// Returns: true if active status was changed, false if no change
	SetActive(active bool, when Time) bool

	// QueueMessagePackage adds a serialized message to the delivery queue
	//
	// Enqueues messages for later transmission (handles rate limiting/backpressure)
	//
	// Parameters:
	//   - rMsg     - ReliableMessage associated with the serialized data
	//   - data     - Serialized byte array of the message
	//   - priority - Delivery priority (smaller = faster)
	// Returns: true if message queued successfully, false on queue error
	QueueMessagePackage(rMsg ReliableMessage, data []byte, priority int) bool
}
