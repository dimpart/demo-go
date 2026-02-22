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
package db

import . "github.com/dimchat/mkm-go/protocol"

type ProviderInfo struct {
	ID     ID
	Name   string
	URL    string
	Chosen bool
}

func NewProviderInfo(pid ID, name string, url string, chosen bool) *ProviderInfo {
	return &ProviderInfo{
		ID:     pid,
		Name:   name,
		URL:    url,
		Chosen: chosen,
	}
}

type StationInfo struct {
	ID     ID
	Name   string
	Host   string
	Port   uint16
	Chosen bool
}

func NewStationInfo(sid ID, name string, host string, port uint16, chosen bool) *StationInfo {
	return &StationInfo{
		ID:     sid,
		Name:   name,
		Host:   host,
		Port:   port,
		Chosen: chosen,
	}
}

/**
 *  Service Provider
 *  ~~~~~~~~~~~~~~~~
 */
type ProviderTable interface {

	/**
	 *  Get all providers
	 *
	 * @return provider list
	 */
	GetProviders() []*ProviderInfo

	/**
	 *  Add provider info
	 *
	 * @param identifier - sp ID
	 * @param name       - sp name
	 * @param url        - entrance URL
	 * @param chosen     - whether current sp
	 * @return false on failed
	 */
	AddProvider(identifier ID, name string, url string, chosen bool) bool

	/**
	 *  Update provider info
	 *
	 * @param identifier - sp ID
	 * @param name       - sp name
	 * @param url        - entrance URL
	 * @param chosen     - whether current sp
	 * @return false on failed
	 */
	UpdateProvider(identifier ID, name string, url string, chosen bool) bool

	/**
	 *  Remove provider info
	 *
	 * @param identifier - sp ID
	 * @return false on failed
	 */
	RemoveProvider(identifier ID) bool
}

/**
 *  Station
 *  ~~~~~~~
 */
type StationTable interface {

	/**
	 *  Get all stations of this sp
	 *
	 * @param sp - sp ID
	 * @return station list
	 */
	GetStations(sp ID) []*StationInfo

	/**
	 *  Add station info with sp ID
	 *
	 * @param sp      - sp ID
	 * @param station - station ID
	 * @param host    - station IP
	 * @param port    - station port
	 * @param name    - station name
	 * @param chosen  - whether current station
	 * @return false on failed
	 */
	AddStation(sp ID, station ID, host string, port uint16, name string, chosen bool) bool

	/**
	 *  Update station info
	 *
	 * @param sp      - sp ID
	 * @param station - station ID
	 * @param host    - station IP
	 * @param port    - station port
	 * @param name    - station name
	 * @param chosen  - whether current station
	 * @return false on failed
	 */
	UpdateStation(sp ID, station ID, host string, port uint16, name string, chosen bool) bool

	/**
	 *  Set this station as current station
	 *
	 * @param sp      - sp ID
	 * @param station - station ID
	 * @return false on failed
	 */
	ChooseStation(sp ID, station ID) bool

	/**
	 *  Remove this station
	 *
	 * @param sp      - sp ID
	 * @param station - station ID
	 * @return false on failed
	 */
	RemoveStation(sp ID, station ID) bool

	/**
	 *  Remove all station of the sp
	 *
	 * @param sp - sp ID
	 * @return false on failed
	 */
	RemoveStations(sp ID) bool
}
