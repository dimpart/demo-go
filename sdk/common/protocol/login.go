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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

const LOGIN = "login"

// LoginCommand defines the interface for user login commands in the DIM network
//
// # Implements the Command interface for authentication and session setup
//
//	Data Format: {
//	    "type": 0x88,
//	    "sn": 123,
//
//	    "command": "login",
//	    "time": 0,
//	    //---- client info ----
//	    "did": "{UserID}",
//	    "device": "DeviceID",      // Device identifier (optional)
//	    "agent": "UserAgent",      // Client user agent (optional, e.g., "DIM-Client/1.0")
//	    //---- server info ----
//	    "station": {
//	        "did": "{StationID}",
//	        "host": "{IP}",
//	        "port": 9394
//	    },
//	    "provider": {
//	        "did": "{SP_ID}"
//	    }
//	}
type LoginCommand interface {
	Command

	// ID returns the user's unique identifier (did field)
	//
	// Returns: User ID associated with the login request
	ID() ID

	// Device returns the client device identifier (optional field)
	//
	// Returns: Device ID string (empty string if not set)
	Device() string
	SetDevice(device string)

	// Agent returns the client user agent string (optional field)
	//
	// Returns: User agent string (empty string if not set)
	Agent() string
	SetAgent(agent string)

	// StationInfo returns the DIM network station information
	//
	// Returns: StringKeyMap containing station did, host, and port
	StationInfo() StringKeyMap
	SetStationInfo(station StringKeyMap)

	// ProviderInfo returns the DIM network service provider information
	//
	// Returns: StringKeyMap containing provider did
	ProviderInfo() StringKeyMap
	SetProviderInfo(sp StringKeyMap)
}
