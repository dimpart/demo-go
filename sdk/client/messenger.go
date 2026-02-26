package sdk

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimpart/demo-go/sdk/client/network"
	. "github.com/dimpart/demo-go/sdk/common"
)

// IClientMessenger defines the interface for client-side message communication
//
// Extends ICommonMessenger with client-specific operations: handshake, broadcast, login, and presence reporting
type IClientMessenger interface {
	ICommonMessenger

	// GetClientSession retrieves the client-specific session interface
	//
	// Returns: IClientSession instance for client session management
	GetClientSession() IClientSession

	// Handshake sends a handshake command to the current network station
	//
	// Establishes a secure session with the station using the provided session key
	//
	// Parameters:
	//   - sessionKey - Session key to use for the handshake response (authentication)
	Handshake(sessionKey string)

	// HandshakeSuccess is the callback invoked when handshake completes successfully
	//
	// Triggers post-handshake operations (e.g., session initialization, document broadcast)
	HandshakeSuccess()

	// BroadcastDocuments sends Meta and Visa documents to all network stations
	//
	// Controls broadcast behavior based on update status:
	//   - updated=true: Force broadcast (document changed, send immediately)
	//   - updated=false: Conditional broadcast (only send if needed/fresh)
	// Parameters:
	//   - updated - Flag indicating if documents have been updated
	BroadcastDocuments(updated bool)

	// BroadcastLogin sends a login command to maintain roaming state across stations
	//
	// Keeps the client's login state active across network stations (roaming support)
	//
	// Parameters:
	//   - sender    - User ID of the login sender (current client user)
	//   - userAgent - Client user agent string (e.g., "DIM-Client/1.0 (iOS)")
	BroadcastLogin(sender ID, userAgent string)

	// ReportOnline sends a presence report command to mark the user as online
	//
	// Updates the user's online status across the network
	//
	// Parameters:
	//   - sender - User ID to mark as online
	ReportOnline(sender ID)

	// ReportOffline sends a presence report command to mark the user as offline
	//
	// Updates the user's offline status across the network
	//
	// Parameters:
	//   - sender - User ID to mark as offline
	ReportOffline(sender ID)
}

// ClientMessenger implements the IClientMessenger interface
//
// # Wraps CommonMessenger to provide client-specific message communication capabilities
//
// Core responsibilities: Handshake management, document broadcasting, and presence reporting
type ClientMessenger struct {
	//IClientMessenger
	*CommonMessenger
}

func NewClientMessenger(session Session, facebook ICommonFacebook, database CipherKeyDelegate) *ClientMessenger {
	return &ClientMessenger{
		CommonMessenger: NewCommonMessenger(session, facebook, database),
	}
}
