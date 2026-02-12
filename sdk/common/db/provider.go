/* license: https://mit-license.org
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
package db

import (
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

var GSP = NewID("gsp", EVERYWHERE, "")

type ProviderInfo struct {
	ID     ID
	Chosen int
}

func ConvertProviderInfo(array []StringKeyMap) []*ProviderInfo {
	providers := make([]*ProviderInfo, 0, len(array))
	var info *ProviderInfo
	var did ID
	var chosen int
	for _, item := range array {
		did = ParseID(item["did"])
		if did == nil {
			did = ParseID(item["ID"])
			if did == nil {
				// SP ID error
				continue
			}
		}
		chosen = ConvertInt(item["chosen"], 0)
		info = &ProviderInfo{
			ID:     did,
			Chosen: chosen,
		}
		providers = append(providers, info)
	}
	return providers
}

func RevertProviderInfo(providers []*ProviderInfo) []StringKeyMap {
	array := make([]StringKeyMap, len(providers))
	for index, item := range providers {
		array[index] = StringKeyMap{
			"ID":     item.ID.String(),
			"did":    item.ID.String(),
			"chosen": item.Chosen,
		}
	}
	return array
}

type ProviderDBI interface {

	// get list of (SP_ID, chosen)
	AllProviders() []*ProviderInfo

	AddProvider(pid ID, chosen bool) bool

	UpdateProvider(pid ID, chosen bool) bool

	RemoveProvider(pid ID) bool
}
