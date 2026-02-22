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

type MetaDBI interface {

	/**
	 *  Get meta from DB
	 */
	GetMeta(entity ID) Meta

	SaveMeta(meta Meta, entity ID) bool
}

type DocumentDBI interface {

	/**
	 *  Get documents from DB
	 */
	GetDocuments(entity ID) []Document

	SaveDocument(doc Document, entity ID) bool
}

type UserDBI interface {

	/**
	 *  Get local user ID list
	 */
	GetLocalUsers() []ID

	SaveLocalUsers(users []ID) bool
}

type ContactDBI interface {

	/**
	 *  Get contacts from DB
	 */
	GetContacts(user ID) []ID

	SaveContacts(contacts []ID, user ID) bool
}

type AccountDBI interface {
	PrivateKeyDBI

	MetaDBI
	DocumentDBI

	UserDBI
	ContactDBI

	GroupDBI
	GroupHistoryDBI
}
