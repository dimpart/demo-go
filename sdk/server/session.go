package sdk

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/types"
)

// format "(IP, Port)"
type SessionAddress string

type SessionHandler interface {
	PushMessage(msg ReliableMessage) bool
}

type Session interface {

	// user ID
	ID() ID
	SetID(identifier ID)

	// session key
	Key() string

	// connection target "(IP, port)"
	ClientAddress() SessionAddress

	// when the client entered background, it should be set to False
	IsActive() bool
	SetActive(active bool)

	// Push message when session active
	PushMessage(msg ReliableMessage) bool
}

func generateSessionKey() string {
	return HexEncode(RandomBytes(32))
}

type BaseSession struct {
	//Session

	identifier ID
	key        string
	address    SessionAddress
	active     bool
	handler    SessionHandler
}

func NewBaseSession(address SessionAddress, handler SessionHandler) *BaseSession {
	return &BaseSession{
		identifier: nil,
		key:        generateSessionKey(),
		address:    address,
		active:     true,
		handler:    handler,
	}
}

//-------- Session

func (session *BaseSession) ID() ID {
	return session.identifier
}
func (session *BaseSession) SetID(identifier ID) {
	session.identifier = identifier
}

func (session *BaseSession) Key() string {
	return session.key
}

func (session *BaseSession) ClientAddress() SessionAddress {
	return session.address
}

func (session *BaseSession) IsActive() bool {
	return session.active
}
func (session *BaseSession) SetActive(active bool) {
	session.active = active
}

func (session *BaseSession) PushMessage(msg ReliableMessage) bool {
	if !session.active {
		return false
	}
	return session.handler.PushMessage(msg)
}

/**
 *  Session Server
 *  ~~~~~~~~~~~~~~
 */
type SessionServer struct {
	clientAddresses map[string][]SessionAddress
	sessions        map[SessionAddress]Session
}

func NewSessionServer() *SessionServer {
	return &SessionServer{
		clientAddresses: make(map[string][]SessionAddress, 1024),
		sessions:        make(map[SessionAddress]Session, 1024),
	}
}

// Session factory
func (server *SessionServer) GetSession(address SessionAddress, handler SessionHandler) Session {
	session := server.sessions[address]
	if session == nil && !ValueIsNil(handler) {
		// create a new session and cache it
		session := NewBaseSession(address, handler)
		server.sessions[address] = session
	}
	return session
}

func (server *SessionServer) insert(address SessionAddress, did ID) {
	identifier := did.String()
	array := server.clientAddresses[identifier]
	if array == nil {
		array = make([]SessionAddress, 0, 1)
		//} else {
		//	for _, item := range array {
		//		if item == address {
		//			// already exists
		//			return
		//		}
		//	}
	}
	server.clientAddresses[identifier] = append(array, address)
}

func (server *SessionServer) remove(address SessionAddress, did ID) {
	identifier := did.String()
	array := server.clientAddresses[identifier]
	if array == nil {
		// not exists
		return
	}
	pos := len(array)
	for pos > 0 {
		pos--
		if array[pos] == address {
			array = append(array[:pos], array[pos+1:]...)
		}
	}
	if len(array) == 0 {
		// all sessions removed
		delete(server.clientAddresses, identifier)
	}
}

// Insert a session with ID into memory cache
func (server *SessionServer) UpdateSession(session Session, identifier ID) {
	address := session.ClientAddress()
	old := session.ID()
	if old != nil {
		// 0. remove client_address from old ID
		server.remove(address, old)
	}
	// 1. insert client_address for new ID
	server.insert(address, identifier)
	// 2. update session ID
	session.SetID(identifier)
}

// Remove the session from memory cache
func (server *SessionServer) RemoveSession(session Session) {
	identifier := session.ID()
	address := session.ClientAddress()
	if identifier != nil {
		// 1. remove client_address with ID
		server.remove(address, identifier)
	}
	// 2. remove session with client_address
	session.SetActive(false)
	delete(server.sessions, address)
}

// Get all sessions of this user
func (server *SessionServer) AllSessions(did ID) []Session {
	identifier := did.String()
	results := make([]Session, 0, 1)
	// 1. get all client_address with ID
	array := server.clientAddresses[identifier]
	if array != nil {
		// 2. get session by each client_address
		var session Session
		for _, item := range array {
			session = server.sessions[item]
			if session != nil {
				results = append(results, session)
			}
		}
	}
	return results
}
func (server *SessionServer) ActiveSessions(identifier ID) []Session {
	results := make([]Session, 0, 1)
	// 1. get all sessions
	all := server.AllSessions(identifier)
	for _, item := range all {
		// 2. check session active
		if item.IsActive() {
			results = append(results, item)
		}
	}
	return results
}

//
//  Users
//

func (server *SessionServer) AllUsers() []ID {
	users := make([]ID, 0, 8)
	var did ID
	for key := range server.clientAddresses {
		did = ParseID(key)
		if did == nil {
			continue
		}
		users = append(users, did)
	}
	return users
}

func (server *SessionServer) IsActive(identifier ID) bool {
	sessions := server.AllSessions(identifier)
	for _, item := range sessions {
		if item.IsActive() {
			return true
		}
	}
	return false
}

func (server *SessionServer) ActiveUsers() []ID {
	users := make([]ID, 0, 8)
	all := server.AllUsers()
	for _, item := range all {
		if server.IsActive(item) {
			users = append(users, item)
		}
	}
	return users
}
