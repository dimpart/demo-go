package sdk

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// each query will be expired after 10 minutes
	QUERY_EXPIRES = DurationOfMinutes(10)

	// each respond will be expired after 10 minutes
	RESPOND_EXPIRES = DurationOfMinutes(10)
)

// IEntityChecker defines the interface for validating entity data freshness and query eligibility
//
// Core responsibilities:
//   - Track active group members and document/group history timestamps
//   - Check if entity data (Meta/Documents/Group Members) needs to be queried/updated
type IEntityChecker interface {

	// SetLastActiveMember records the most recently active member of a specific group
	//
	// Used to track group activity and prioritize member status checks
	//
	// Parameters:
	//   - group  - Group ID to update active member for
	//   - member - User ID of the most recently active group member
	SetLastActiveMember(group, member ID)

	// GetLastActiveMember retrieves the most recently active member of a specific group
	//
	// Parameters:
	//   - group - Group ID to retrieve active member for
	// Returns: User ID of the last active member (nil/zero value if no active member recorded)
	GetLastActiveMember(group ID) ID

	// SetLastDocumentTime updates the 'SDT' (Sender Document Time) for an entity
	//
	// Tracks the timestamp of the last document update to determine freshness
	//
	// Parameters:
	//   - did     - Entity ID (user/device) to update document time for
	//   - current - New timestamp to set as the last document time
	// Returns: true if timestamp updated successfully, false on validation/error
	SetLastDocumentTime(did ID, current Time) bool

	// SetLastGroupHistoryTime updates the 'GHT' (Group History Time) for a group
	//
	// Tracks the timestamp of the last group history update (e.g., member changes)
	//
	// Parameters:
	//   - gid     - Group ID to update history time for
	//   - current - New timestamp to set as the last group history time
	// Returns: true if timestamp updated successfully, false on validation/error
	SetLastGroupHistoryTime(gid ID, current Time) bool

	// -------------------------------------------------------------------------
	//  Meta Validation
	// -------------------------------------------------------------------------

	// CheckMeta determines if Meta data for an entity needs to be queried from the network
	//
	// Evaluates freshness of existing Meta to decide if a network query is required
	//
	// Parameters:
	//   - did  - Entity ID (user/device) to check Meta for
	//   - meta - Existing Meta data for the entity (may be nil/empty)
	// Returns: true if Meta needs to be queried (stale/missing), false if current
	CheckMeta(did ID, meta Meta) bool

	// -------------------------------------------------------------------------
	//  Document Validation
	// -------------------------------------------------------------------------

	// CheckDocuments determines if documents for an entity need to be queried/updated
	//
	// Evaluates freshness of existing documents to decide if a network query/update is required
	//
	// Parameters:
	//   - did       - Entity ID (user/device) to check documents for
	//   - documents - Existing document list for the entity (may be empty)
	// Returns: true if documents need to be queried/updated (stale/missing), false if current
	CheckDocuments(did ID, documents []Document) bool

	// -------------------------------------------------------------------------
	//  Group Member Validation
	// -------------------------------------------------------------------------

	// CheckMembers determines if group member list needs to be queried from the network
	//
	// Evaluates freshness of existing member list to decide if a network query is required
	//
	// Parameters:
	//   - group   - Group ID to check members for
	//   - members - Existing group member list (may be empty)
	// Returns: true if members need to be queried (stale/missing), false if current
	CheckMembers(group ID, members []ID) bool
}

// IEntityRequest defines the interface for requesting entity data from the network
//
// Handles network requests for Meta/Documents/Group Members with duplicate request prevention
//
// All methods should check expiration status (via isXXXQueryExpired()) before sending requests
type IEntityRequest interface {

	// QueryMeta sends a network request for Metadata of a specific entity
	//
	// Precondition: Call isMetaQueryExpired() first to check if request is allowed;
	// Prevents duplicate requests for the same Metadata
	//
	// Parameters:
	//   - did - Entity ID (user/device) to request Meta for
	// Returns: true if request sent successfully, false if duplicate/invalid
	QueryMeta(did ID) bool

	// QueryDocuments sends a network request for documents of a specific entity
	//
	// Precondition: Call isDocumentQueryExpired() first to check if request is allowed;
	// Evaluates existing documents and sends request only if needed (prevents duplicates)
	//
	// Parameters:
	//   - did  - Entity ID (user/device) to request documents for
	//   - docs - Existing document list (used to evaluate need for request)
	// Returns: true if request sent successfully, false if duplicate/invalid
	QueryDocuments(did ID, docs []Document) bool

	// QueryMembers sends a network request for group member list of a specific group
	//
	// Precondition: Call isMembersQueryExpired() first to check if request is allowed;
	// Evaluates existing members and sends request only if needed (prevents duplicates)
	//
	// Parameters:
	//   - gid     - Group ID to request members for
	//   - members - Existing group member list (used to evaluate need for request)
	// Returns: true if request sent successfully, false if duplicate/invalid
	QueryMembers(gid ID, members []ID) bool
}

// IEntityRespond defines the interface for responding to entity data requests
//
// Manages outbound entity data (e.g., Visa documents) with rate limiting and update checks
type IEntityRespond interface {

	// SendVisa sends a Visa document to a contact with rate limiting and update checks
	//
	// Sending Rules:
	//   - If updated=true: Force send (ignores rate limit, used for document updates)
	//   - If updated=false: Send only once every 10 minutes (rate limit)
	// Parameters:
	//   - visa     - Visa document to send to the contact
	//   - receiver - Contact ID to send the Visa document to
	//   - updated  - Flag indicating if the Visa document has been updated
	// Returns: true if Visa sent successfully, false if rate limited/failed
	SendVisa(visa Visa, receiver ID, updated bool) bool
}
