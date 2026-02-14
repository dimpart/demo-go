/* license: https://mit-license.org
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
package dimp

import (
	"fmt"

	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/db"
)

type SocketAddress interface {
	fmt.Stringer

	Host() string
	Port() uint16
}

type Session interface {
	Transmitter

	GetDatabase() SessionDBI

	/**
	 *  Get remote socket address
	 *
	 * @return host & port
	 */
	GetRemoteAddress() SocketAddress

	// session key
	GetSessionKey() string

	/**
	 *  Update user ID
	 *
	 * @param uid - login user ID
	 * @return true on changed
	 */
	SetID(uid ID) bool
	GetID() ID

	/**
	 *  Update active flag
	 *
	 * @param active - flag
	 * @param when   - now
	 * @return true on changed
	 */
	SetActive(active bool, when Time) bool
	IsActive() bool

	/**
	 *  Pack message into a waiting queue
	 *
	 * @param rMsg     - network message
	 * @param data     - serialized message
	 * @param priority - smaller is faster
	 * @return false on error
	 */
	QueueMessagePackage(rMsg ReliableMessage, data []byte, priority int) bool
}
