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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type GroupDBI interface {

	/**
	 *  Get group founder/owner from DB
	 */
	LoadFounder(group ID) ID
	LoadOwner(group ID) ID

	/**
	 *  Get group members from DB
	 */
	LoadMembers(group ID) []ID
	SaveMembers(members []ID, group ID) bool

	/**
	 *  Get group admins from DB
	 */
	LoadAdministrators(group ID) []ID
	SaveAdministrators(members []ID, group ID) bool
}

type GroupHistoryDBI interface {

	/**
	 *  Save group commands
	 *      1. invite
	 *      2. expel (deprecated)
	 *      3. join
	 *      4. quit
	 *      5. reset
	 *      6. resign
	 *
	 * @param content - group command
	 * @param rMsg    - group command message
	 * @param group   - group ID
	 * @return false on failed
	 */
	SaveGroupHistory(content GroupCommand, rMsg ReliableMessage, group ID) bool

	/**
	 *  Load group commands
	 *      1. invite
	 *      2. expel (deprecated)
	 *      3. join
	 *      4. quit
	 *      5. reset
	 *      6. resign
	 *
	 * @param group - group ID
	 * @return history list
	 */
	LoadGroupHistories(group ID) []*Pair[GroupCommand, ReliableMessage]

	/**
	 *  Load last 'reset' group command
	 *
	 * @param group - group ID
	 * @return reset command message
	 */
	LoadResetCommandMessage(group ID) *Pair[ResetCommand, ReliableMessage]

	/**
	 *  Clean group commands for members:
	 *      1. invite
	 *      2. expel (deprecated)
	 *      3. join
	 *      4. quit
	 *      5. reset
	 *
	 * @param group - group ID
	 * @return false on failed
	 */
	ClearGroupMemberHistories(group ID) bool

	/**
	 *  Clean group commands for administrators
	 *      1. resign
	 *
	 * @param group - group ID
	 * @return false on failed
	 */
	ClearGroupAdminHistories(group ID)
}
