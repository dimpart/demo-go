package db

import (
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

var GSP = NewID("gsp", EVERYWHERE, "")

type ProviderInfo struct {
	ID     ID
	Chosen int
}

func ConvertProviderInfo(array []StringKeyMap) []*ProviderInfo {
	providers := make([]*ProviderInfo, 0, len(array))
	var info *ProviderInfo
	var did ID
	var chosen int
	for _, item := range array {
		did = ParseID(item["did"])
		if did == nil {
			did = ParseID(item["ID"])
			if did == nil {
				// SP ID error
				continue
			}
		}
		chosen = ConvertInt(item["chosen"], 0)
		info = &ProviderInfo{
			ID:     did,
			Chosen: chosen,
		}
		providers = append(providers, info)
	}
	return providers
}

func RevertProviderInfo(providers []*ProviderInfo) []StringKeyMap {
	array := make([]StringKeyMap, len(providers))
	for index, item := range providers {
		array[index] = StringKeyMap{
			"ID":     item.ID.String(),
			"did":    item.ID.String(),
			"chosen": item.Chosen,
		}
	}
	return array
}

type ProviderDBI interface {

	// get list of (SP_ID, chosen)
	AllProviders() []*ProviderInfo

	AddProvider(pid ID, chosen bool) bool

	UpdateProvider(pid ID, chosen bool) bool

	RemoveProvider(pid ID) bool
}
