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
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/mem"
	. "github.com/dimchat/plugins-go/mkm"
)

func NewCompatibleMetaFactory(version MetaType) MetaFactory {
	return &compatibleMetaFactory{
		BaseMetaFactory{
			Type: version,
		},
	}
}

/**
 *  Compatible Meta Factory
 */
type compatibleMetaFactory struct {
	BaseMetaFactory
}

// Override
func (factory compatibleMetaFactory) ParseMeta(info StringKeyMap) Meta {
	// check 'type', 'key', 'seed', 'fingerprint'
	if !ContainsKey(info, "type") || !ContainsKey(info, "key") {
		// meta.type should not be empty
		// meta.key should not be empty
		return nil
	} else if !ContainsKey(info, "seed") {
		if ContainsKey(info, "fingerprint") {
			//panic("meta error")
			return nil
		}
	} else if !ContainsKey(info, "fingerprint") {
		//panic("meta error")
		return nil
	}
	// create meta for type
	var out Meta
	helper := GetGeneralAccountHelper()
	version := helper.GetMetaType(info, "")
	switch version {
	case "1", "mkm", "MKM":
		out = NewMetaWithMap(info)
		break
	case "2", "btc", "BTC":
		out = NewBTCMetaWithMap(info)
		break
	case "4", "eth", "ETH":
		out = NewETHMetaWithMap(info)
		break
	default:
		break
	}
	if out.IsValid() {
		return out
	}
	//panic("meta error")
	return nil
}
