package sdk

import (
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/**
 *  ANS
 */

type AddressNameServer struct {
	//AddressNameService

	reserved   map[string]bool
	caches     map[string]ID
	namesTable map[ID][]string
}

func NewAddressNameServer() *AddressNameServer {
	// reserved names
	reserved := make(map[string]bool, len(KEYWORDS))
	for _, item := range KEYWORDS {
		reserved[item] = true
	}

	// constant ANS records
	caches := make(map[string]ID, 1024)
	caches["all"] = EVERYONE
	caches["everyone"] = EVERYONE
	caches["anyone"] = ANYONE
	caches["owner"] = ANYONE
	caches["founder"] = FOUNDER

	// temp
	namesTable := make(map[ID][]string, 128)

	return &AddressNameServer{
		reserved:   reserved,
		caches:     caches,
		namesTable: namesTable,
	}
}

func (ans *AddressNameServer) setID(name string, identifier ID) {
	if ValueIsNil(identifier) {
		delete(ans.caches, name)
	} else {
		ans.caches[name] = identifier
	}
}

// Override
func (ans *AddressNameServer) IsReserved(name string) bool {
	return ans.reserved[name]
}

// Override
func (ans *AddressNameServer) GetID(name string) ID {
	return ans.caches[name]
}

// Override
func (ans *AddressNameServer) GetNames(identifier ID) []string {
	array := ans.namesTable[identifier]
	if array == nil {
		array = make([]string, 0, 1)
		// TODO: update all tables?
		for key, value := range ans.caches {
			if identifier.Equal(value) {
				array = append(array, key)
			}
		}
	}
	return array
}
