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
package protocol

import (
	"strings"

	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
)

//
//  Broadcast Utils
//

func BroadcastGroupSeed(group ID) string {
	name := group.Name()
	if name != "" {
		size := len(name)
		if size == 8 && strings.ToLower(name) == "everyone" {
			name = ""
		}
	}
	return name
}

func BroadcastGroupFounder(group ID) ID {
	name := BroadcastGroupSeed(group)
	if name == "" {
		// Consensus: the founder of group 'everyone@everywhere'
		//            'Albert Moky'
		return FOUNDER
	}
	// DISCUSS: who should be the founder of group 'xxx@everywhere'?
	//          'anyone@anywhere', or 'xxx.founder@anywhere'
	return ParseID(name + ".founder@anywhere")
}

func BroadcastGroupOwner(group ID) ID {
	name := BroadcastGroupSeed(group)
	if name == "" {
		// Consensus: the owner of group 'everyone@everywhere'
		//            'anyone@anywhere'
		return ANYONE
	}
	// DISCUSS: who should be the owner of group 'xxx@everywhere'?
	//          'anyone@anywhere', or 'xxx.owner@anywhere'
	return ParseID(name + ".owner@anywhere")
}

func BroadcastGroupMembers(group ID) []ID {
	name := BroadcastGroupSeed(group)
	if name == "" {
		// Consensus: the member of group 'everyone@everywhere'
		//            'anyone@anywhere'
		return []ID{
			ANYONE,
		}
	}
	// DISCUSS: who should be the member of group 'xxx@everywhere'?
	//          'anyone@anywhere', or 'xxx.member@anywhere'
	owner := ParseID(name + ".owner@anywhere")
	member := ParseID(name + ".member@anywhere")
	return []ID{
		owner,
		member,
	}
}
