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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

type BaseLoginCommand struct {
	//LoginCommand
	BaseCommand
}

func (content *BaseLoginCommand) InitWithID(did ID) LoginCommand {
	if content.BaseCommand.Init(LOGIN) != nil {
		content.Set("did", did.String())
	}
	return content
}

//
//  Client Info
//

// Override
func (content *BaseLoginCommand) ID() ID {
	did := content.Get("did")
	return ParseID(did)
}

// Override
func (content *BaseLoginCommand) Device() string {
	return content.GetString("device", "")
}

// Override
func (content *BaseLoginCommand) SetDevice(device string) {
	content.Set("device", device)
}

// Override
func (content *BaseLoginCommand) Agent() string {
	return content.GetString("agent", "")
}

// Override
func (content *BaseLoginCommand) SetAgent(agent string) {
	content.Set("agent", agent)
}

//
//  Server Info
//

// Override
func (content *BaseLoginCommand) StationInfo() StringKeyMap {
	info := content.Get("station")
	if dict, ok := info.(StringKeyMap); ok {
		return dict
	}
	return nil
}

// Override
func (content *BaseLoginCommand) SetStationInfo(station StringKeyMap) {
	content.Set("station", station)
}

// Override
func (content *BaseLoginCommand) ProviderInfo() StringKeyMap {
	info := content.Get("provider")
	if dict, ok := info.(StringKeyMap); ok {
		return dict
	}
	return nil
}

// Override
func (content *BaseLoginCommand) SetProviderInfo(sp StringKeyMap) {
	content.Set("provider", sp)
}
