/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package db

import (
	"strings"

	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

//-------- GroupTable

// Override
func (db *Storage) GetFounder(group ID) ID {
	// TODO: get group founder
	return nil
}

// Override
func (db *Storage) GetOwner(group ID) ID {
	// TODO: get group owner
	return nil
}

// Override
func (db *Storage) GetAdministrators(group ID) []ID {
	// TODO: get administrators
	return nil
}

// Override
func (db *Storage) SaveAdministrators(admins []ID, group ID) bool {
	// TODO: save administrators
	return true
}

// Override
func (db *Storage) GetMembers(group ID) []ID {
	arr := db._memberTable[group.String()]
	if arr == nil {
		arr = loadMembers(db, group)
		db._memberTable[group.String()] = arr
	}
	return arr
}

// Override
func (db *Storage) SaveMembers(members []ID, group ID) bool {
	db._memberTable[group.String()] = members
	return saveMembers(db, group, members)
}

/**
 *  Members file for Group
 *  ~~~~~~~~~~~~~~~~~~~~~~
 *
 *  file path: '.dim/mkm/{zzz}/{ADDRESS}/members.txt'
 */

func membersPath(db *Storage, group ID) string {
	return PathJoin(db.mkmDir(group), "members.txt")
}

func loadMembers(db *Storage, group ID) []ID {
	path := membersPath(db, group)
	db.log("Loading members for group: " + group.String())
	text := db.readText(path)
	lines := strings.Split(text, "\n")
	members := make([]ID, 0, len(lines))
	for _, rec := range lines {
		id := ParseID(rec)
		if id != nil {
			members = append(members, id)
		}
	}
	return members
}

func saveMembers(db *Storage, group ID, members []ID) bool {
	text := ""
	lines := IDRevert(members)
	for _, rec := range lines {
		text = text + rec + "\n"
	}
	path := membersPath(db, group)
	db.log("Saving members for group: " + group.String())
	return db.writeText(path, text)
}

//-------- GroupKeysTable

// Override
func (db *Storage) SaveGroupHistory(content GroupCommand, rMsg ReliableMessage, group ID) bool {
	// TODO:
	return false
}

// Override
func (db *Storage) GetGroupHistories(group ID) []Pair[GroupCommand, ReliableMessage] {
	// TODO:
	return nil
}

// Override
func (db *Storage) GetResetCommandMessage(group ID) Pair[ResetCommand, ReliableMessage] {
	// TODO:
	return nil
}

// Override
func (db *Storage) ClearGroupMemberHistories(group ID) bool {
	// TODO:
	return false
}

// Override
func (db *Storage) ClearGroupAdminHistories(group ID) bool {
	// TODO:
	return false
}
