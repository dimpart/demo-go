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
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/**
 *  ANS
 */
type AddressNameServer struct {
	//AddressNameService

	_reserved   map[string]bool
	_caches     map[string]ID
	_namesTable map[ID][]string
}

func (ans *AddressNameServer) Init() AddressNameService {
	// reserved names
	ans._reserved = make(map[string]bool, len(KEYWORDS))
	for _, item := range KEYWORDS {
		ans._reserved[item] = true
	}

	// constant ANS records
	ans._caches = make(map[string]ID, 1024)
	ans.setID("all", EVERYONE)
	ans.setID("everyone", EVERYONE)
	ans.setID("anyone", ANYONE)
	ans.setID("owner", ANYONE)
	ans.setID("founder", FOUNDER)

	// temp
	ans._namesTable = make(map[ID][]string, 128)

	return ans
}

func (ans *AddressNameServer) setID(name string, identifier ID) {
	if ValueIsNil(identifier) {
		delete(ans._caches, name)
	} else {
		ans._caches[name] = identifier
	}
}

// Override
func (ans *AddressNameServer) IsReserved(name string) bool {
	return ans._reserved[name]
}

// Override
func (ans *AddressNameServer) GetID(name string) ID {
	return ans._caches[name]
}

// Override
func (ans *AddressNameServer) GetNames(identifier ID) []string {
	array := ans._namesTable[identifier]
	if array == nil {
		array = make([]string, 0, 1)
		// TODO: update all tables?
		for key, value := range ans._caches {
			if identifier.Equal(value) {
				array = append(array, key)
			}
		}
	}
	return array
}
