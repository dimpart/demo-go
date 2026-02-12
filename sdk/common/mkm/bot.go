/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
package mkm

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/mkm"
)

/**
 *  Bot User
 *  ~~~~~~~~
 */
type Bot interface {
	User

	/**
	 *  Bot Document
	 */
	Profile() Document

	/**
	 *  Get provider ID
	 *
	 * @return ICP ID, bot group
	 */
	Provider() ID
}

func NewBot(did ID) Bot {
	bot := &BotUser{}
	return bot.Init(did)
}

type BotUser struct {
	BaseUser
}

func (user *BotUser) Init(uid ID) Bot {
	if user.BaseUser.Init(uid) != nil {
	}
	return user
}

// Override
func (user *BotUser) Profile() Document {
	docs := user.Documents()
	return GetLastVisa(docs)
}

// Override
func (user *BotUser) Provider() ID {
	doc := user.Profile()
	if doc == nil {
		return nil
	}
	pid := doc.GetProperty("provider")
	return ParseID(pid)
}
