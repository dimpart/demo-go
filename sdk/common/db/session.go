package db

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

// LoginDBI defines the interface for user login history persistence operations
//
// Manages storage and retrieval of the most recent login command/message for user authentication
type LoginDBI interface {

	// GetLoginCommandMessage retrieves the last login command and associated reliable message for a user
	//
	// Used to verify recent login activity and session state
	//
	// Parameters:
	//   - user - ID of the user to retrieve login history for
	// Returns: Pair[LoginCommand, ReliableMessage] containing last login data (zero-value Pair if no login found)
	GetLoginCommandMessage(user ID) Pair[LoginCommand, ReliableMessage]

	// SaveLoginCommandMessage persists the last login command and associated reliable message for a user
	//
	// Overwrites existing login data (only stores the most recent login)
	//
	// Parameters:
	//   - user    - ID of the user associated with the login
	//   - content - LoginCommand containing login details (device, agent, station info)
	//   - rMsg    - ReliableMessage associated with the login command (for traceability)
	// Returns: true if login data saved successfully, false on database error
	SaveLoginCommandMessage(user ID, content LoginCommand, rMsg ReliableMessage) bool
}

type SessionDBI interface {
	LoginDBI

	ProviderDBI
	StationDBI
}
