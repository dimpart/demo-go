package db

import . "github.com/dimchat/mkm-go/protocol"

// MetaDBI defines the interface for Metadata persistence operations
//
// Manages storage and retrieval of Meta information associated with entity IDs
type MetaDBI interface {

	// GetMeta retrieves the Meta instance associated with a specific entity ID
	//
	// Parameters:
	//   - entity - Entity ID (user/device ID) to retrieve Meta for
	// Returns: Meta instance (nil if no Meta found for the entity)
	GetMeta(entity ID) Meta

	// SaveMeta persists a Meta instance associated with a specific entity ID
	//
	// Overwrites existing Meta if the entity ID already has associated Meta
	//
	// Parameters:
	//   - meta   - Meta instance to persist (must be valid Meta implementation)
	//   - entity - Entity ID (user/device ID) to associate with the Meta
	// Returns: true if Meta saved successfully, false on database error
	SaveMeta(meta Meta, entity ID) bool
}

// DocumentDBI defines the interface for document persistence operations
//
// Manages storage and retrieval of documents (Visa/Bulletin/etc.) associated with entity IDs
type DocumentDBI interface {

	// GetDocuments retrieves all documents associated with a specific entity ID
	//
	// Returns all document types (Visa, Bulletin, etc.) for the entity
	//
	// Parameters:
	//   - entity - Entity ID (user/device ID) to retrieve documents for
	// Returns: Slice of Document instances (empty slice if no documents found)
	GetDocuments(entity ID) []Document

	// SaveDocument persists a single document associated with a specific entity ID
	//
	// Appends to existing documents (does not overwrite existing documents)
	//
	// Parameters:
	//   - doc    - Document instance to persist (Visa/Bulletin/etc.)
	//   - entity - Entity ID (user/device ID) to associate with the document
	// Returns: true if document saved successfully, false on database error
	SaveDocument(doc Document, entity ID) bool
}

// UserDBI defines the interface for local user management operations
//
// Manages storage and retrieval of locally registered user IDs
type UserDBI interface {

	// GetLocalUsers retrieves the list of locally registered user IDs
	//
	// Returns all users configured on the local device/client
	//
	// Returns: Slice of ID instances (empty slice if no local users registered)
	GetLocalUsers() []ID

	// SaveLocalUsers overwrites the list of locally registered user IDs
	//
	// Replaces the entire local user list (not incremental update)
	//
	// Parameters:
	//   - users - Slice of user IDs to set as local users
	// Returns: true if user list saved successfully, false on database error
	SaveLocalUsers(users []ID) bool
}

// ContactDBI defines the interface for contact list management operations
//
// Manages storage and retrieval of contact lists associated with specific users
type ContactDBI interface {

	// GetContacts retrieves the contact list for a specific user
	//
	// Returns all contact IDs associated with the user's address book
	//
	// Parameters:
	//   - user - User ID to retrieve contact list for
	// Returns: Slice of contact ID instances (empty slice if no contacts found)
	GetContacts(user ID) []ID

	// SaveContacts overwrites the contact list for a specific user
	//
	// Replaces the entire contact list (not incremental update)
	//
	// Parameters:
	//   - contacts - Slice of contact IDs to set for the user
	//   - user     - User ID to associate with the contact list
	// Returns: true if contact list saved successfully, false on database error
	SaveContacts(contacts []ID, user ID) bool
}

type AccountDBI interface {
	PrivateKeyDBI

	MetaDBI
	DocumentDBI

	UserDBI
	ContactDBI

	GroupDBI
	GroupHistoryDBI
}
