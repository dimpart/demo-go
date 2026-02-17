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
package cpu

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/client"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

type HandshakeCommandProcessor struct {
	BaseCommandProcessor
}

// Override
func (cpu *HandshakeCommandProcessor) GetMessenger() IClientMessenger {
	messenger := cpu.BaseCommandProcessor.GetMessenger()
	cm, ok := messenger.(IClientMessenger)
	if ok {
		return cm
	}
	return nil
}

// Override
func (cpu *HandshakeCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(HandshakeCommand)
	if !ok {
		return nil
	}
	messenger := cpu.GetMessenger()
	session := messenger.GetClientSession()
	// update station's default ID ('station@anywhere') to sender (real ID)
	station := session.GetStation()
	oid := station.ID()
	sender := rMsg.Sender()
	if oid == nil || oid.IsBroadcast() {
		station.SetID(sender)
	} else if !oid.Equal(sender) {
		panic("station ID not match: " + oid.String() + ", " + sender.String())
	}
	// handle handshake command with title & session key
	title := command.Title()
	newKey := command.SessionKey()
	oldKey := session.GetSessionKey()
	if title == "DIM?" {
		// S -> C: station ask client to handshake again
		if oldKey == "" {
			// first handshake response with new session key
			messenger.Handshake(newKey)
		} else if oldKey == newKey {
			// duplicated handshake response?
			// or session expired and the station ask to handshake again?
			messenger.Handshake(newKey)
		} else {
			// connection changed?
			// erase session key to handshake again
			session.SetSessionKey("")
		}
	} else if title == "DIM!" {
		// S -> C: handshake accepted by station
		if oldKey == "" {
			// normal handshake response,
			// update session key to change state to 'running'
			session.SetSessionKey(newKey)
		} else if oldKey == newKey {
			// duplicated handshake response?
			// set it again here to invoke the flutter channel
			session.SetSessionKey(newKey)
		} else {
			// FIXME: handshake error
			// erase session key to handshake again
			session.SetSessionKey("")
		}
	} else {
		// C -> S: Hello world!
		panic("handshake from other user?")
	}
	return nil
}
