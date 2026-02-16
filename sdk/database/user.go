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

//-------- UserTable

// Override
func (db *Storage) LoadLocalUsers() []ID {
	return db._users
}

// Override
func (db *Storage) SaveLocalUsers(users []ID) bool {
	db._users = users
	return true
}

func (db *Storage) AddUser(user ID) bool {
	for _, id := range db._users {
		if user.Equal(id) {
			return false
		}
	}
	users := make([]ID, 0, len(db._users)+1)
	users = append(users, user)
	users = append(users, db._users...)
	db._users = users
	return true
}

func (db *Storage) RemoveUser(user ID) bool {
	var pos = -1
	for index, id := range db._users {
		if user.Equal(id) {
			pos = index
			break
		}
	}
	if pos == -1 {
		// user ID not found
		return false
	}
	db._users = append(db._users[:pos], db._users[pos+1:]...)
	return true
}

func (db *Storage) SetCurrentUser(user ID) {
	db.RemoveUser(user)
	db.AddUser(user)
}

func (db *Storage) GetCurrentUser() ID {
	if len(db._users) == 0 {
		return nil
	}
	return db._users[0]
}
