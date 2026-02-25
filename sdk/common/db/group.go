package db

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

// GroupDBI defines the interface for group core information persistence operations
//
// Manages storage and retrieval of group ownership, membership, and administrator data
type GroupDBI interface {

	// GetFounder retrieves the founder (creator) ID of a specific group
	//
	// The founder is the original creator of the group (immutable)
	//
	// Parameters:
	//   - group - Group ID to retrieve founder for
	// Returns: Founder's user ID (nil/zero value if group not found)
	GetFounder(group ID) ID

	// GetOwner retrieves the current owner ID of a specific group
	//
	// The owner may change (e.g., transfer ownership) unlike the founder
	//
	// Parameters:
	//   - group - Group ID to retrieve owner for
	// Returns: Current group owner's user ID (nil/zero value if group not found)
	GetOwner(group ID) ID

	// GetMembers retrieves the full list of members for a specific group
	//
	// Returns all regular members (excludes administrators unless they are also regular members)
	//
	// Parameters:
	//   - group - Group ID to retrieve members for
	// Returns: Slice of member user IDs (empty slice if no members or group not found)
	GetMembers(group ID) []ID

	// SaveMembers overwrites the full member list for a specific group
	//
	// Replaces the entire member list (not incremental update/addition)
	//
	// Parameters:
	//   - members - Slice of user IDs to set as group members
	//   - group   - Group ID to associate with the member list
	// Returns: true if member list saved successfully, false on database error
	SaveMembers(members []ID, group ID) bool

	// GetAdministrators retrieves the list of administrators for a specific group
	//
	// Administrators have elevated permissions for group management
	//
	// Parameters:
	//   - group - Group ID to retrieve administrators for
	// Returns: Slice of admin user IDs (empty slice if no admins or group not found)
	GetAdministrators(group ID) []ID

	// SaveAdministrators overwrites the administrator list for a specific group
	//
	// Replaces the entire admin list (not incremental update/addition)
	//
	// Parameters:
	//   - members - Slice of user IDs to set as group administrators
	//   - group   - Group ID to associate with the admin list
	// Returns: true if admin list saved successfully, false on database error
	SaveAdministrators(members []ID, group ID) bool
}

// GroupHistoryDBI defines the interface for group command history persistence operations
//
// Manages storage, retrieval, and cleanup of group operation commands (invite, join, reset, etc.)
type GroupHistoryDBI interface {

	// SaveGroupHistory persists a group operation command and its associated reliable message
	//
	// Supported command types:
	//   1. invite   - Invite users to group
	//   2. expel    - Remove users from group (deprecated)
	//   3. join     - User joins group
	//   4. quit     - User leaves group
	//   5. reset    - Reset group membership
	//   6. resign   - Admin resigns from position
	// Parameters:
	//   - content - GroupCommand instance representing the operation
	//   - rMsg    - ReliableMessage associated with the command (for traceability)
	//   - group   - Group ID the command applies to
	// Returns: true if history saved successfully, false on database error
	SaveGroupHistory(content GroupCommand, rMsg ReliableMessage, group ID) bool

	// GetGroupHistories retrieves all command history for a specific group
	//
	// Returns all supported command types (invite, expel, join, quit, reset, resign)
	//
	// Parameters:
	//   - group - Group ID to retrieve command history for
	// Returns: Slice of Pair[GroupCommand, ReliableMessage] (empty slice if no history)
	GetGroupHistories(group ID) []Pair[GroupCommand, ReliableMessage]

	// GetResetCommandMessage retrieves the most recent 'reset' command for a specific group
	//
	// The reset command typically resets group configuration or membership
	//
	// Parameters:
	//   - group - Group ID to retrieve reset command for
	// Returns: Pair[ResetCommand, ReliableMessage] (zero value Pair if no reset command found)
	GetResetCommandMessage(group ID) Pair[ResetCommand, ReliableMessage]

	// ClearGroupMemberHistories deletes member-related command history for a specific group
	//
	// Removes history for: invite, expel (deprecated), join, quit, reset commands
	//
	// Parameters:
	//   - group - Group ID to clear member history for
	// Returns: true if history cleared successfully, false on database error
	ClearGroupMemberHistories(group ID) bool

	// ClearGroupAdminHistories deletes administrator-related command history for a specific group
	//
	// Removes history for: resign commands (admin resignation)
	//
	// Parameters:
	//   - group - Group ID to clear admin history for
	// Returns: true if history cleared successfully, false on database error
	ClearGroupAdminHistories(group ID) bool
}
