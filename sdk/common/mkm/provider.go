/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
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
package mkm

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/mkm"
)

/**
 *  DIM Station Owner
 *  ~~~~~~~~~~~~~~~~~
 */
type Provider interface {
	Group

	/**
	 *  Provider Document
	 */
	Profile() Document

	Stations() []StringKeyMap
}

func NewProvider(pid ID) Provider {
	return &ServiceProvider{
		BaseGroup: NewBaseGroup(pid),
	}
}

/**
 *  DIM Station Owner
 *  ~~~~~~~~~~~~~~~~~
 */
type ServiceProvider struct {
	*BaseGroup
}

// Override
func (sp *ServiceProvider) Profile() Document {
	docs := sp.Documents()
	return GetLastVisa(docs)
}

// Override
func (sp *ServiceProvider) Stations() []StringKeyMap {
	doc := sp.Profile()
	if doc != nil {
		stations := doc.GetProperty("stations")
		if array, ok := stations.([]StringKeyMap); ok {
			return array
		}
	}
	// TODO: load from local storage
	return nil
}

//
//  Comparison
//

func sameStation(a, b Station) bool {
	if a == b {
		// same object
		return true
	}
	return checkIdentifiers(a.ID(), b.ID()) &&
		checkHosts(a.Host(), b.Host()) &&
		checkPorts(a.Port(), b.Port())
}

func checkIdentifiers(a, b ID) bool {
	if a == b {
		// same object
		return true
	} else if a.IsBroadcast() || b.IsBroadcast() {
		return true
	}
	return a.Equal(b)
}

func checkHosts(a, b string) bool {
	if a == "" || b == "" {
		return true
	}
	return a == b
}

func checkPorts(a, b uint16) bool {
	if a == 0 || b == 0 {
		return true
	}
	return a == b
}
