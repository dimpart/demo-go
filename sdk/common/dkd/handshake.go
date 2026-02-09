/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
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
package dkd

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

func getState(title string, session string) HandshakeState {
	// check message text
	if title == "" {
		return HandshakeInit
	}
	if title == "DIM!" /*|| message == "OK!"*/ {
		return HandshakeSuccess
	}
	if title == "DIM?" {
		return HandshakeAgain
	}
	// check session key
	if session == "" {
		return HandshakeStart
	}
	return HandshakeRestart
}

type BaseHandshakeCommand struct {
	//HandshakeCommand
	BaseCommand
}

func (content *BaseHandshakeCommand) InitWithTitle(title string, sessionKey string) HandshakeCommand {
	if content.BaseCommand.Init(HANDSHAKE) != nil {
		// text message
		content.Set("title", title)
		// session key
		content.Set("session", sessionKey)
	}
	return content
}

// Override
func (content *BaseHandshakeCommand) Title() string {
	return content.GetString("title", "")
}

// Override
func (content *BaseHandshakeCommand) SessionKey() string {
	return content.GetString("session", "")
}

// Override
func (content *BaseHandshakeCommand) State() HandshakeState {
	title := content.Title()
	session := content.SessionKey()
	return getState(title, session)
}
