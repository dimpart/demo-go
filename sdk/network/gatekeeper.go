package network

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimpart/demo-go/sdk/common"
)

type IGateKeeper interface {

	// IP+Port
	GetRemoteAddress() SocketAddress

	// protected
	PackMessage(payload []byte, priority int) Departure
	QueueAppend(msg ReliableMessage, ship Departure) bool
}

type GateKeeper struct {
	//IGateKeeper

	RemoteAddress SocketAddress
}
