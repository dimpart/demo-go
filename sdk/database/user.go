package db

import . "github.com/dimchat/mkm-go/protocol"

//-------- UserTable

// Override
func (db *Storage) GetLocalUsers() []ID {
	return db.users
}

// Override
func (db *Storage) SaveLocalUsers(users []ID) bool {
	db.users = users
	return true
}

func (db *Storage) AddUser(user ID) bool {
	for _, id := range db.users {
		if user.Equal(id) {
			return false
		}
	}
	users := make([]ID, 0, len(db.users)+1)
	users = append(users, user)
	users = append(users, db.users...)
	db.users = users
	return true
}

func (db *Storage) RemoveUser(user ID) bool {
	var pos = -1
	for index, id := range db.users {
		if user.Equal(id) {
			pos = index
			break
		}
	}
	if pos == -1 {
		// user ID not found
		return false
	}
	db.users = append(db.users[:pos], db.users[pos+1:]...)
	return true
}

func (db *Storage) SetCurrentUser(user ID) {
	db.RemoveUser(user)
	db.AddUser(user)
}

func (db *Storage) GetCurrentUser() ID {
	if len(db.users) == 0 {
		return nil
	}
	return db.users[0]
}
