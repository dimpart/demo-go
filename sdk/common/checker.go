package dimp

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

type IEntityChecker interface {

	// active member for group
	SetLastActiveMember(group, member ID)
	GetLastActiveMember(group ID) ID

	// update 'SDT' - Sender Document Time
	SetLastDocumentTime(did ID, current Time) bool

	// update 'GHT' - Group History Time
	SetLastGroupHistoryTime(gid ID, current Time) bool

	//
	//  Meta
	//

	/**
	 *  Check meta for querying
	 *
	 * @param did  - entity ID
	 * @param meta - exists meta
	 * @return true on querying
	 */
	CheckMeta(did ID, meta Meta) bool

	//
	//  Documents
	//

	/**
	 *  Check documents for querying/updating
	 *
	 * @param identifier - entity ID
	 * @param documents  - exist document
	 * @return true on querying
	 */
	CheckDocuments(did ID, documents []Document) bool

	//
	//  Group Members
	//

	/**
	 *  Check group members for querying
	 *
	 * @param group   - group ID
	 * @param members - exist members
	 * @return true on querying
	 */
	CheckMembers(group ID, members []ID) bool
}

type IEntityRequest interface {

	/**
	 *  Request for meta with entity ID
	 *  (call 'isMetaQueryExpired()' before sending command)
	 *
	 * @param did - entity ID
	 * @return false on duplicated
	 */
	QueryMeta(did ID) bool

	/**
	 *  Request for documents with entity ID
	 *  (call 'isDocumentQueryExpired()' before sending command)
	 *
	 * @param did  - entity ID
	 * @param docs - exist documents
	 * @return false on duplicated
	 */
	QueryDocuments(did ID, docs []Document) bool

	/**
	 *  Request for group members with group ID
	 *  (call 'isMembersQueryExpired()' before sending command)
	 *
	 * @param gid     - group ID
	 * @param members - exist members
	 * @return false on duplicated
	 */
	QueryMembers(gid ID, members []ID) bool
}

type IEntityRespond interface {

	///  Send my visa document to contact
	///  if document is updated, force to send it again.
	///  else only send once every 10 minutes.
	SendVisa(visa Visa, receiver ID, updated bool) bool
}
