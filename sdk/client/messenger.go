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
package dimp

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimpart/demo-go/sdk/client/network"
	. "github.com/dimpart/demo-go/sdk/common"
)

type IClientMessenger interface {
	ICommonMessenger

	GetClientSession() IClientSession

	/**
	 *  Send handshake command to current station
	 *
	 * @param sessionKey - respond session key
	 */
	Handshake(sessionKey string)

	/**
	 *  Callback for handshake success
	 */
	HandshakeSuccess()

	/**
	 *  Broadcast meta &amp; visa document to all stations
	 */
	BroadcastDocuments(updated bool)

	/**
	 *  Send login command to keep roaming
	 */
	BroadcastLogin(sender ID, userAgent string)

	/**
	 *  Send report command to keep user online
	 */
	ReportOnline(sender ID)

	/**
	 *  Send report command to let user offline
	 */
	ReportOffline(sender ID)
}

/**
 *  Client Messenger for Handshake &amp; Broadcast Report
 */
type ClientMessenger struct {
	CommonMessenger
}

func (messenger *ClientMessenger) Init(session Session, facebook ICommonFacebook, database CipherKeyDelegate) *ClientMessenger {
	if messenger.CommonMessenger.Init(session, facebook, database) != nil {
	}
	return messenger
}
