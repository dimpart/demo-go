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

type BaseMuteCommand struct {
	//MuteCommand
	BaseCommand

	// mute-list
	_list []ID
}

func (content *BaseMuteCommand) InitWithMap(dict StringKeyMap) MuteCommand {
	if content.BaseCommand.InitWithMap(dict) != nil {
		// lazy load
		content._list = nil
	}
	return content
}

func (content *BaseMuteCommand) InitWithList(list []ID) MuteCommand {
	if content.BaseCommand.Init(MUTE) != nil {
		if !ValueIsNil(list) {
			content.SetMuteList(list)
		}
	}
	return content
}

// Override
func (content *BaseMuteCommand) MuteList() []ID {
	if content._list == nil {
		list := content.Get("list")
		if list != nil {
			content._list = IDConvert(list)
		}
	}
	return content._list
}

// Override
func (content *BaseMuteCommand) SetMuteList(list []ID) {
	if ValueIsNil(list) {
		content.Remove("list")
	} else {
		content.Set("list", IDRevert(list))
	}
	content._list = list
}
