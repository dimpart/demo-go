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
	"fmt"

	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/plugins-go/mkm"
)

func getName(network EntityType) string {
	switch network {
	case BOT:
		return "Robot"
	case STATION:
		return "Station"
	case ICP:
		return "CP"
	case ISP:
		return "SP"
	}
	if EntityTypeIsUser(network) {
		return "User"
	}
	if EntityTypeIsGroup(network) {
		return "Group"
	}
	return "Unknown"
}

func AnonymousGetName(identifier ID) string {
	name := identifier.Name()
	if name == "" {
		name = getName(identifier.Type())
	}
	return name + " (" + AnonymousGetNumberString(identifier.Address()) + ")"
}

func AnonymousGetNumber(address Address) uint32 {
	btc, ok := address.(*BTCAddress)
	if ok {
		return btcNumber(btc.String())
	}
	eth, ok := address.(*ETHAddress)
	if ok {
		return ethNumber(eth.String())
	}
	panic(address)
}
func AnonymousGetNumberString(address Address) string {
	number := AnonymousGetNumber(address)
	str := fmt.Sprintf("%010d", number)
	return str[0:3] + "-" + str[3:6] + "-" + str[6:]
}

func btcNumber(address string) uint32 {
	data := Base58Decode(address)
	return userNumber(data[len(data)-4:])
}
func ethNumber(address string) uint32 {
	data := HexDecode(address[2:])
	return userNumber(data[len(data)-4:])
}

func userNumber(cc []byte) uint32 {
	length := len(cc)
	return (uint32(cc[length-4]) << 24) |
		(uint32(cc[length-3]) << 16) |
		(uint32(cc[length-2]) << 8) |
		uint32(cc[length-1])
}
