package sdk

import (
	"fmt"

	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/plugins-go/mkm"
)

func getName(network EntityType) string {
	switch network {
	case BOT:
		return "Robot"
	case STATION:
		return "Station"
	case ICP:
		return "CP"
	case ISP:
		return "SP"
	}
	if EntityTypeIsUser(network) {
		return "User"
	}
	if EntityTypeIsGroup(network) {
		return "Group"
	}
	return "Unknown"
}

func AnonymousGetName(identifier ID) string {
	name := identifier.Name()
	if name == "" {
		name = getName(identifier.Type())
	}
	return name + " (" + AnonymousGetNumberString(identifier.Address()) + ")"
}

func AnonymousGetNumber(address Address) uint32 {
	btc, ok := address.(*BTCAddress)
	if ok {
		return btcNumber(btc.String())
	}
	eth, ok := address.(*ETHAddress)
	if ok {
		return ethNumber(eth.String())
	}
	panic(address)
}
func AnonymousGetNumberString(address Address) string {
	number := AnonymousGetNumber(address)
	str := fmt.Sprintf("%010d", number)
	return str[0:3] + "-" + str[3:6] + "-" + str[6:]
}

func btcNumber(address string) uint32 {
	data := Base58Decode(address)
	return userNumber(data[len(data)-4:])
}
func ethNumber(address string) uint32 {
	data := HexDecode(address[2:])
	return userNumber(data[len(data)-4:])
}

func userNumber(cc []byte) uint32 {
	length := len(cc)
	return (uint32(cc[length-4]) << 24) |
		(uint32(cc[length-3]) << 16) |
		(uint32(cc[length-2]) << 8) |
		uint32(cc[length-1])
}
