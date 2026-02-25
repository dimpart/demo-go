/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
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
package protocol

import (
	"fmt"

	. "github.com/dimchat/core-go/protocol"
)

type HandshakeState uint8

const (
	HandshakeInit    HandshakeState = iota
	HandshakeStart                  // C -> S, without session key(or session expired)
	HandshakeAgain                  // S -> C, with new session key
	HandshakeRestart                // C -> S, with new session key
	HandshakeSuccess                // S -> C, handshake accepted
)

func (state HandshakeState) String() string {
	switch state {
	case HandshakeInit:
		return "HandshakeInit"
	case HandshakeStart:
		return "HandshakeStart"
	case HandshakeAgain:
		return "HandshakeAgain"
	case HandshakeRestart:
		return "HandshakeRestart"
	case HandshakeSuccess:
		return "HandshakeSuccess"
	default:
		return fmt.Sprintf("HandshakeState(%d)", state)
	}
}

const HANDSHAKE = "handshake"

// HandshakeCommand defines the interface for handshake commands (session initialization)
//
// # Implements the Command interface for DIM network session establishment
//
//	Data Format: {
//	    "type": 0x88,
//	    "sn": 123,
//
//	    "command": "handshake",
//	    "title": "Hello world!",   // Handshake state indicator ("DIM?", "DIM!")
//	    "session": "{SESSION_KEY}" // Session key for authenticated communication
//	}
type HandshakeCommand interface {
	Command

	// Title returns the handshake state indicator (e.g., "DIM?", "DIM!")
	//
	// Returns: String representing the current handshake state
	Title() string

	// SessionKey returns the session key for authenticated communication
	//
	// Returns: Session key string (empty string if not established)
	SessionKey() string

	// State returns the structured HandshakeState derived from the title
	//
	// Returns: Enumerated HandshakeState value
	State() HandshakeState
}
