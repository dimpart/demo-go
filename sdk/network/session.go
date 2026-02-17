package network

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimpart/demo-go/sdk/common"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/utils"
)

type IBaseSession interface {
	Session

	QueueMessagePackage(msg ReliableMessage, data []byte, priority int) bool
}

type BaseSession struct {
	//Session

	// private
	ID         ID
	Database   SessionDBI
	Messenger  ICommonMessenger
	GateKeeper IGateKeeper
}

func (session *BaseSession) SetID(user ID) bool {
	did := session.ID
	if did == nil {
		if user == nil {
			return false
		}
	} else if did.Equal(user) {
		return false
	}
	session.ID = did
	return true
}

func (session *BaseSession) GetDatabase() SessionDBI {
	return session.Database
}

func (session *BaseSession) GetRemoteAddress() SocketAddress {
	keeper := session.GateKeeper
	return keeper.GetRemoteAddress()
}

// protected
func (session *BaseSession) QueueMessagePackage(msg ReliableMessage, data []byte, priority int) bool {
	keeper := session.GateKeeper
	ship := keeper.PackMessage(data, priority)
	return keeper.QueueAppend(msg, ship)
}

//
//  Transmitter
//

// Override
func (session *BaseSession) SendContent(content Content, sender, receiver ID, priority int) Pair[InstantMessage, ReliableMessage] {
	messenger := session.Messenger
	return messenger.SendContent(content, sender, receiver, priority)
}

// Override
func (session *BaseSession) SendInstantMessage(iMsg InstantMessage, priority int) ReliableMessage {
	messenger := session.Messenger
	return messenger.SendInstantMessage(iMsg, priority)
}

// Override
func (session *BaseSession) SendReliableMessage(rMsg ReliableMessage, priority int) bool {
	messenger := session.Messenger
	return messenger.SendReliableMessage(rMsg, priority)
}
