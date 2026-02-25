package db

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

type StationInfo struct {
	ID ID

	Host string
	Port uint16

	SP ID // provider

	Chosen int
}

func ConvertStationInfo(array []StringKeyMap) []*StationInfo {
	stations := make([]*StationInfo, 0, len(array))
	var info *StationInfo
	var did ID
	var host string
	var port uint16
	var provider ID
	var chosen int
	for _, item := range array {
		did = ParseID(item["did"])
		if did == nil {
			did = ParseID(item["ID"])
		}
		host = ConvertString(item["host"], "")
		port = ConvertUInt16(item["port"], 0)
		provider = ParseID(item["provider"])
		chosen = ConvertInt(item["chosen"], 0)
		if host == "" || port == 0 /* || provider == nil*/ {
			// station socket error
			continue
		}
		info = &StationInfo{
			ID:     did,
			Host:   host,
			Port:   port,
			SP:     provider,
			Chosen: chosen,
		}
		stations = append(stations, info)
	}
	return stations
}

func RevertStationInfo(stations []*StationInfo) []StringKeyMap {
	array := make([]StringKeyMap, len(stations))
	for index, item := range stations {
		array[index] = StringKeyMap{
			"ID":       item.ID.String(),
			"did":      item.ID.String(),
			"host":     item.Host,
			"port":     item.Port,
			"provider": item.SP.String(),
			"chosen":   item.Chosen,
		}
	}
	return array
}

type StationDBI interface {

	// get list of (SP_ID, chosen)
	AllStations(provider ID) []*StationInfo

	AddStation(sid ID, host string, port uint16, provider ID, chosen int) bool

	UpdateStation(sid ID, host string, port uint16, provider ID, chosen int) bool

	RemoveStation(host string, port uint16, provider ID) bool

	RemoveStations(provider ID) bool
}
