package ext

import (
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/client/cpu"
	. "github.com/dimpart/demo-go/sdk/common/ext"
)

type ClientExtensionLoader struct {
	CommonExtensionLoader
}

// Override
func (loader ClientExtensionLoader) Load() {
	loader.CommonExtensionLoader.Load()

	registerCustomizedHandlers()

}

func registerCustomizedHandlers() {
	filter := NewAppCustomizedFilter()

	// 'chat.dim.group:history'
	handler := &GroupHistoryHandler{}
	filter.SetContentHandler("chat.dim.group", "history", handler)

	SetCustomizedContentFilter(filter)
}
