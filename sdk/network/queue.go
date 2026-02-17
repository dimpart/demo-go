package network

import . "github.com/dimchat/dkd-go/protocol"

type Departure interface {
}

type IMessageQueue interface {

	/**
	 *  Append message with departure ship
	 *
	 * @param rMsg - outgoing message
	 * @param ship - departure ship
	 * @return false on duplicated
	 */
	Append(rMsg ReliableMessage, ship Departure) bool
}

type MessageQueue struct {
	//IMessageQueue
}
