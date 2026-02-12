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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

type StationInfo struct {
	ID ID

	Host string
	Port uint16

	SP ID // provider

	Chosen int
}

func ConvertStationInfo(array []StringKeyMap) []*StationInfo {
	stations := make([]*StationInfo, 0, len(array))
	var info *StationInfo
	var did ID
	var host string
	var port uint16
	var provider ID
	var chosen int
	for _, item := range array {
		did = ParseID(item["did"])
		if did == nil {
			did = ParseID(item["ID"])
		}
		host = ConvertString(item["host"], "")
		port = ConvertUInt16(item["port"], 0)
		provider = ParseID(item["provider"])
		chosen = ConvertInt(item["chosen"], 0)
		if host == "" || port == 0 /* || provider == nil*/ {
			// station socket error
			continue
		}
		info = &StationInfo{
			ID:     did,
			Host:   host,
			Port:   port,
			SP:     provider,
			Chosen: chosen,
		}
		stations = append(stations, info)
	}
	return stations
}

func RevertStationInfo(stations []*StationInfo) []StringKeyMap {
	array := make([]StringKeyMap, len(stations))
	for index, item := range stations {
		array[index] = StringKeyMap{
			"ID":       item.ID.String(),
			"did":      item.ID.String(),
			"host":     item.Host,
			"port":     item.Port,
			"provider": item.SP.String(),
			"chosen":   item.Chosen,
		}
	}
	return array
}

type StationDBI interface {

	// get list of (SP_ID, chosen)
	AllStations(provider ID) []*StationInfo

	AddStation(sid ID, host string, port uint16, provider ID, chosen int) bool

	UpdateStation(sid ID, host string, port uint16, provider ID, chosen int) bool

	RemoveStation(host string, port uint16, provider ID) bool

	RemoveStations(provider ID) bool
}
