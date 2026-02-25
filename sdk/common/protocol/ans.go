/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
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
package protocol

import . "github.com/dimchat/mkm-go/protocol"

var KEYWORDS = [...]string{
	"all", "everyone", "anyone", "owner", "founder",
	// --------------------------------
	"dkd", "mkm", "dimp", "dim", "dimt",
	"rsa", "ecc", "aes", "des", "btc", "eth",
	// --------------------------------
	"crypto", "key", "symmetric", "asymmetric",
	"public", "private", "secret", "password",
	"id", "address", "meta",
	"tai", "document", "profile", "visa", "bulletin",
	"entity", "user", "group", "contact",
	// --------------------------------
	"member", "admin", "administrator", "assistant",
	"main", "polylogue", "chatroom",
	"social", "organization",
	"company", "school", "government", "department",
	"provider", "station", "thing", "bot", "robot",
	// --------------------------------
	"message", "instant", "secure", "reliable",
	"envelope", "sender", "receiver", "time",
	"content", "forward", "command", "history",
	"keys", "data", "signature",
	// --------------------------------
	"type", "serial", "sn",
	"text", "file", "image", "audio", "video", "page",
	"handshake", "receipt", "block", "mute",
	"register", "suicide", "found", "abdicate",
	"invite", "expel", "join", "quit", "reset", "query",
	"hire", "fire", "resign",
	// --------------------------------
	"server", "client", "terminal", "local", "remote",
	"barrack", "cache", "transceiver",
	"ans", "facebook", "store", "messenger",
	"root", "supervisor",
}

// AddressNameService defines the interface for managing alias-to-ID mappings (address name resolution)
//
// Core functionality: Check alias reservation status, resolve names to IDs, and list names for an ID
type AddressNameService interface {

	// IsReserved checks if a given alias name is reserved (unavailable for use)
	//
	// Parameters:
	//   - name - Alias/short name to check for reservation status
	// Returns: true if the name is reserved (cannot be used), false if available
	IsReserved(name string) bool

	// GetID resolves a short name/alias to its corresponding user ID
	//
	// Parameters:
	//   - name - Short name/alias to resolve (must be non-reserved and registered)
	// Returns: User ID associated with the name (zero value if name not found)
	GetID(name string) ID

	// GetNames retrieves all short names/aliases associated with a specific user ID
	//
	// Parameters:
	//   - did - User ID to look up associated names
	// Returns: Slice of short names/aliases linked to the ID (empty slice if no names exist)
	GetNames(did ID) []string
}
