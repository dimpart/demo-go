/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
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
package dkd

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/**
 *  Handshake command message
 *
 *  data format: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "handshake",    // command name
 *      title   : "Hello world!", // "DIM?", "DIM!"
 *      session : "{SESSION_KEY}" // session key
 *  }
 */

func NewHandshakeCommand(title, sessionKey string) HandshakeCommand {
	return NewBaseHandshakeCommand(nil, title, sessionKey)
}

func NewHandshakeCommandWithMap(dict StringKeyMap) Command {
	return NewBaseHandshakeCommand(dict, "", "")
}

/**
 *  Login Command
 *
 *  data format: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command  : "login",
 *      time     : 0,
 *      //---- client info ----
 *      did      : "{UserID}",
 *      device   : "DeviceID",  // (optional)
 *      agent    : "UserAgent", // (optional)
 *      //---- server info ----
 *      station  : {
 *          did  : "{StationID}",
 *          host : "{IP}",
 *          port : 9394
 *      },
 *      provider : {
 *          did  : "{SP_ID}"
 *      }
 *  }
 */

func NewLoginCommand(did ID) LoginCommand {
	return NewBaseLoginCommand(nil, did)
}

func NewLoginCommandWithMap(dict StringKeyMap) Command {
	return NewBaseLoginCommand(dict, nil)
}

/**
 *  Mute Command
 *
 *  data format: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "mute",
 *      list    : []      // mute-list
 *  }
 */

func NewMuteCommand(list []ID) MuteCommand {
	return NewBaseMuteCommand(nil, list)
}

func NewMuteCommandWithMap(dict StringKeyMap) Command {
	return NewBaseMuteCommand(dict, nil)
}

/**
 *  Block Command
 *
 *  data format: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "block",
 *      list    : []      // block-list
 *  }
 */

func NewBlockCommand(list []ID) BlockCommand {
	return NewBaseBlockCommand(nil, list)
}

func NewBlockCommandWithMap(dict StringKeyMap) Command {
	return NewBaseBlockCommand(dict, nil)
}

/**
 *  Report Command
 *
 *  data format: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "report",
 *      title   : "online",      // or "offline"
 *      //---- extra info
 *      time    : 1234567890,    // timestamp
 *  }
 */

func NewReportCommand(title string) ReportCommand {
	return NewBaseReportCommand(nil, title)
}

func NewReportCommandWithMap(dict StringKeyMap) Command {
	return NewBaseReportCommand(dict, "")
}

/**
 *  Application Customized Content
 */

func NewCustomizedContent(app, mod, act string) CustomizedContent {
	return NewAppCustomizedContent(nil, "", app, mod, act)
}

func NewCustomizedContentWithMap(dict StringKeyMap) Content {
	return NewAppCustomizedContent(dict, "", "", "", "")
}
