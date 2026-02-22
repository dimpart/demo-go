/* license: https://mit-license.org
 *
 *  Ming-Ke-Ming : Decentralized User Identity Authentication
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
	"strings"

	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/mem"
	. "github.com/dimchat/plugins-go/mkm"
)

func NewCompatibleAddressFactory() AddressFactory {
	return &compatibleAddressFactory{}
}

/**
 *  Compatible Address Factory
 */
type compatibleAddressFactory struct {
	//AddressFactory
}

// Override
func (compatibleAddressFactory) GenerateAddress(meta Meta, network EntityType) Address {
	address := meta.GenerateAddress(network)
	if address != nil {
		cache := GetAddressCache()
		cache.Put(address.String(), address)
	}
	return address
}

// Override
func (compatibleAddressFactory) ParseAddress(str string) Address {
	cache := GetAddressCache()
	address := cache.Get(str)
	if address == nil {
		address = parseAddress(str)
		if address != nil {
			cache.Put(str, address)
		}
	}
	return address
}

// protected
func parseAddress(str string) Address {
	if str == "" {
		//panic("address is empty")
		return nil
	}
	size := len(str)
	if size == 0 {
		//panic("address is empty")
		return nil
	} else if size == 8 {
		// "anywhere"
		lower := strings.ToLower(str)
		if ANYWHERE.Equal(lower) {
			return ANYWHERE
		}
	} else if size == 10 {
		// "everywhere"
		lower := strings.ToLower(str)
		if EVERYWHERE.Equal(lower) {
			return EVERYWHERE
		}
	}
	var result Address
	if 26 <= size && size <= 35 {
		// BTC
		result = ParseBTCAddress(str)
	} else if size == 42 {
		// ETH
		result = ParseETHAddress(str)
	} else {
		//panic("invalid address")
	}
	//
	//  TODO: parse for other types of address
	//
	if result == nil && 4 <= size && size <= 64 {
		return NewUnknownAddress(str)
	}
	panic("invalid address: " + str)
	return result
}

/**
 *  Unsupported Address
 *  ~~~~~~~~~~~~~~~~~~~
 */
type UnknownAddress struct {
	//Address
	ConstantString
}

func NewUnknownAddress(address string) Address {
	return &UnknownAddress{
		ConstantString: *NewConstantString(address),
	}
}

// Override
func (UnknownAddress) Network() EntityType {
	return USER
}
