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
	"fmt"

	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/crypto"
	. "github.com/dimchat/sdk-go/mkm"
)

//goland:noinspection GoSnakeCaseUsage
var (
	ANY_STATION   = NewID("station", ANYWHERE, "")
	EVERY_STATION = NewID("stations", EVERYWHERE, "")
)

/**
 *  DIM Server
 *  ~~~~~~~~~~
 */
type Station interface {
	User

	/**
	 *  Station Document
	 */
	Profile() Document

	/**
	 *  Get provider ID
	 *
	 * @return ISP ID, station group
	 */
	Provider() ID

	/**
	 *  Station IP
	 */
	Host() string

	/**
	 *  Station Port
	 */
	Port() uint16

	SetID(sid ID)
}

func NewStation(did ID) Station {
	station := &BaseStation{}
	return station.Init(did, "", 0)
}

type BaseStation struct {
	//User

	_user User

	_host string
	_port uint16
}

func (server *BaseStation) Init(did ID, host string, port uint16) Station {
	if did == nil {
		did = ANY_STATION
	}
	user := &BaseUser{}
	server._user = user.Init(did)
	server._host = host
	server._port = port
	return server
}

// Override
func (server *BaseStation) Equal(other interface{}) bool {
	if station, ok := other.(Station); ok {
		return sameStation(station, server)
	}
	// others?
	return server._user.Equal(other)
}

// Override
func (server *BaseStation) String() string {
	clazz := "BaseStation"
	sid := server.ID()
	network := sid.Address().Network()
	return fmt.Sprintf("<%s id=\"%s\" network=%d host=\"%s\" port=%d />",
		clazz, sid.String(), network, server.Host(), server.Port())
}

// Override
func (server *BaseStation) Profile() Document {
	docs := server.Documents()
	return GetLastVisa(docs)
}

// Override
func (server *BaseStation) Provider() ID {
	doc := server.Profile()
	if doc == nil {
		return nil
	}
	pid := doc.GetProperty("provider")
	return ParseID(pid)
}

// Override
func (server *BaseStation) Host() string {
	if server._host == "" {
		doc := server.Profile()
		if doc != nil {
			host := doc.GetProperty("host")
			server._host = ConvertString(host, "")
		}
		if server._host == "" {
			server._host = "0.0.0.0"
		}
	}
	return server._host
}

// Override
func (server *BaseStation) Port() uint16 {
	if server._port == 0 {
		doc := server.Profile()
		if doc != nil {
			port := doc.GetProperty("port")
			server._port = ConvertUInt16(port, 0)
		}
		if server._port == 0 {
			server._port = 9394
		}
	}
	return server._port
}

func (server *BaseStation) SetID(sid ID) {
	delegate := server.DataSource()
	user := &BaseUser{}
	user.Init(sid)
	user.SetDataSource(delegate)
	server._user = user
}

//-------- Entity

// Override
func (server *BaseStation) ID() ID {
	return server._user.ID()
}

// Override
func (server *BaseStation) Type() EntityType {
	return server._user.Type()
}

// Override
func (server *BaseStation) DataSource() EntityDataSource {
	return server._user.DataSource()
}

// Override
func (server *BaseStation) SetDataSource(facebook EntityDataSource) {
	server._user.SetDataSource(facebook)
}

// Override
func (server *BaseStation) Meta() Meta {
	return server._user.Meta()
}

// Override
func (server *BaseStation) Documents() []Document {
	return server._user.Documents()
}

//-------- User

// Override
func (server *BaseStation) Contacts() []ID {
	return server._user.Contacts()
}

// Override
func (server *BaseStation) Terminals() []string {
	return server._user.Terminals()
}

// Override
func (server *BaseStation) Verify(data []byte, signature []byte) bool {
	return server._user.Verify(data, signature)
}

// Override
func (server *BaseStation) EncryptBundle(plaintext []byte) EncryptedBundle {
	return server._user.EncryptBundle(plaintext)
}

// Override
func (server *BaseStation) Sign(data []byte) []byte {
	return server._user.Sign(data)
}

// Override
func (server *BaseStation) DecryptBundle(bundle EncryptedBundle) []byte {
	return server._user.DecryptBundle(bundle)
}

// Override
func (server *BaseStation) SignVisa(visa Visa) Visa {
	return server._user.SignVisa(visa)
}

// Override
func (server *BaseStation) VerifyVisa(visa Visa) bool {
	return server._user.VerifyVisa(visa)
}
