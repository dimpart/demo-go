package dimp

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type EntityChecker struct {
	//IEntityChecker

	// query checkers (string = ID)
	metaQueries    IFrequencyChecker[string]
	docsQueries    IFrequencyChecker[string]
	membersQueries IFrequencyChecker[string]

	// response checker (string = ID)
	documentResponses IFrequencyChecker[string]

	// recent time checkers (string = ID)
	lastDocumentTimes IRecentTimeChecker[string]
	lastHistoryTimes  IRecentTimeChecker[string]

	// group => member (string = GID)
	lastActiveMembers map[string]ID

	// protected
	database AccountDBI

	// delegates
	Request IEntityRequest
	Respond IEntityRespond
}

func (checker *EntityChecker) Init(db AccountDBI) IEntityChecker {
	checker.metaQueries = NewFrequencyChecker[string](QUERY_EXPIRES)
	checker.docsQueries = NewFrequencyChecker[string](QUERY_EXPIRES)
	checker.membersQueries = NewFrequencyChecker[string](QUERY_EXPIRES)

	checker.documentResponses = NewFrequencyChecker[string](QUERY_EXPIRES)

	checker.lastDocumentTimes = NewRecentTimeChecker[string]()
	checker.lastHistoryTimes = NewRecentTimeChecker[string]()

	checker.lastActiveMembers = make(map[string]ID, 128)

	checker.database = db
	return checker
}

// protected
func (checker *EntityChecker) IsMetaQueryExpired(did ID) bool {
	return checker.metaQueries.IsExpired(did.String(), nil, false)
}

// protected
func (checker *EntityChecker) IsDocumentQueryExpired(did ID) bool {
	return checker.docsQueries.IsExpired(did.String(), nil, false)
}

// protected
func (checker *EntityChecker) IsMembersQueryExpired(did ID) bool {
	return checker.membersQueries.IsExpired(did.String(), nil, false)
}

// protected
func (checker *EntityChecker) IsDocumentResponseExpired(did ID, force bool) bool {
	return checker.documentResponses.IsExpired(did.String(), nil, force)
}

// Override
func (checker *EntityChecker) SetLastActiveMember(group, member ID) {
	checker.lastActiveMembers[group.String()] = member
}

// Override
func (checker *EntityChecker) GetLastActiveMember(group ID) ID {
	return checker.lastActiveMembers[group.String()]
}

// Override
func (checker *EntityChecker) SetLastDocumentTime(did ID, current Time) bool {
	return checker.lastDocumentTimes.SetLastTime(did.String(), current)
}

// Override
func (checker *EntityChecker) SetLastGroupHistoryTime(gid ID, current Time) bool {
	return checker.lastHistoryTimes.SetLastTime(gid.String(), current)
}

//
//  Meta
//

// Override
func (checker *EntityChecker) CheckMeta(did ID, meta Meta) bool {
	if checker.NeedsQueryMeta(did, meta) {
		//if !checker.IsMetaQueryExpired(did) {
		//	// query not expired yet
		//	return false
		//}
		return checker.Request.QueryMeta(did)
	}
	// no need to query meta again
	return false
}

// protected
func (checker *EntityChecker) NeedsQueryMeta(did ID, meta Meta) bool {
	if did.IsBroadcast() {
		// broadcast entity has no meta to query
		return false
	} else if meta == nil {
		// meta not found, sure to query
		return true
	}
	return false
}

//
//  Documents
//

// Override
func (checker *EntityChecker) CheckDocuments(did ID, documents []Document) bool {
	if checker.NeedsQueryDocuments(did, documents) {
		//if !checker.IsDocumentQueryExpired(did) {
		//	// query not expired yet
		//	return false
		//}
		return checker.Request.QueryDocuments(did, documents)
	}
	// no need to update documents now
	return false
}

// protected
func (checker *EntityChecker) NeedsQueryDocuments(did ID, documents []Document) bool {
	if did.IsBroadcast() {
		// broadcast entity has no document to query
		return false
	} else if len(documents) == 0 {
		// documents not found, sure to query
		return true
	}
	current := checker.GetLastDocumentTime(did, documents)
	return checker.lastDocumentTimes.IsExpired(did.String(), current)
}

func (checker *EntityChecker) GetLastDocumentTime(did ID, documents []Document) Time {
	if len(documents) == 0 {
		return nil
	}
	var lastTime Time
	var docTime Time
	for _, doc := range documents {
		docTime = doc.Time()
		if TimeIsNil(docTime) {
			//panic("document error")
		} else if TimeIsNil(lastTime) || TimeIsBefore(docTime, lastTime) {
			lastTime = docTime
		}
	}
	return lastTime
}

//
//  Group Members
//

// Override
func (checker *EntityChecker) CheckMembers(group ID, members []ID) bool {
	if checker.NeedsQueryMembers(group, members) {
		//if !checker.IsMembersQueryExpired(group) {
		//	// query not expired yet
		//	return false
		//}
		return checker.Request.QueryMembers(group, members)
	}
	// no need to update group members now
	return false
}

// protected
func (checker *EntityChecker) NeedsQueryMembers(group ID, members []ID) bool {
	if group.IsBroadcast() {
		// broadcast group has no members to query
		return false
	} else if len(members) == 0 {
		// members not found, sure to query
		return true
	}
	current := checker.GetLastGroupHistoryTime(group)
	return checker.lastHistoryTimes.IsExpired(group.String(), current)
}

func (checker *EntityChecker) GetLastGroupHistoryTime(group ID) Time {
	array := checker.database.GetGroupHistories(group)
	if len(array) == 0 {
		return nil
	}
	var lastTime Time
	var hisTime Time
	var his GroupCommand
	for _, pair := range array {
		his = pair.First()
		if his == nil {
			//panic("group command error")
			continue
		}
		hisTime = his.Time()
		if TimeIsNil(hisTime) {
			//panic("group command error")
		} else if TimeIsNil(lastTime) || TimeIsBefore(hisTime, lastTime) {
			lastTime = hisTime
		}
	}
	return lastTime
}
