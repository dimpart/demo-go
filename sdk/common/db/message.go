package db

import (
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// CipherKeyDBI defines the interface for symmetric cipher key persistence operations
//
// # Manages storage and retrieval of directional symmetric keys used for encrypting messages between entities
//
// Supports 1:1 communication (user ↔ user/contact) and 1:N communication (user → group)
type CipherKeyDBI interface {

	// GetCipherKey retrieves the symmetric cipher key for encrypting messages from sender to receiver
	//
	// Keys are directional: sender→receiver key may differ from receiver→sender key
	//
	// Parameters:
	//   - sender   - ID of the message sender (user/contact ID)
	//   - receiver - ID of the message receiver (user/contact/group ID)
	// Returns: SymmetricKey for message encryption (nil if no key found for the direction)
	GetCipherKey(sender, receiver ID) SymmetricKey

	// SaveCipherKey caches a symmetric cipher key for reuse in directional message encryption
	//
	// Stores the key associated with the sender→receiver direction to avoid re-negotiation
	//
	// Parameters:
	//   - sender   - ID of the message sender (user/contact ID)
	//   - receiver - ID of the message receiver (user/contact/group ID)
	//   - key      - SymmetricKey to cache for future encryption use
	// Returns: true if key saved successfully, false on database error
	SaveCipherKey(sender, receiver ID, key SymmetricKey) bool
}

// GroupKeysDBI defines the interface for group message key persistence operations
//
// # Manages storage and retrieval of encrypted symmetric keys for group message communication
//
// Each group member may have unique encoded keys for secure group messaging
type GroupKeysDBI interface {

	// GetGroupKeys retrieves the encoded symmetric message keys for a specific group and sender
	//
	// Returns a StringKeyMap containing key-value pairs of encoded group message keys
	//
	// Parameters:
	//   - group  - ID of the group the keys apply to
	//   - sender - ID of the group member (sender) associated with the keys
	// Returns: StringKeyMap with encoded group keys (nil/empty map if no keys found)
	GetGroupKeys(group, sender ID) StringKeyMap

	// SaveGroupKeys persists encoded symmetric message keys for a specific group and sender
	//
	// Stores keys in StringKeyMap format (key-value pairs of encoded key data)
	//
	// Parameters:
	//   - group  - ID of the group the keys apply to
	//   - sender - ID of the group member (sender) associated with the keys
	//   - keys   - StringKeyMap containing encoded group message keys
	// Returns: true if keys saved successfully, false on database error
	SaveGroupKeys(group, sender ID, keys StringKeyMap) bool
}

type MessageDBI interface {
	CipherKeyDBI

	GroupKeysDBI
}
