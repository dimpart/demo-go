package cpu

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/sdk"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/*
Command Transform:

	+===============================+===============================+
	|      Customized Content       |      Group Query Command      |
	+-------------------------------+-------------------------------+
	|   "type" : i2s(0xCC)          |   "type" : i2s(0x88)          |
	|   "sn"   : 123                |   "sn"   : 123                |
	|   "time" : 123.456            |   "time" : 123.456            |
	|   "app"  : "chat.dim.group"   |                               |
	|   "mod"  : "history"          |                               |
	|   "act"  : "query"            |                               |
	|                               |   "command"   : "query"       |
	|   "group"     : "{GROUP_ID}"  |   "group"     : "{GROUP_ID}"  |
	|   "last_time" : 0             |   "last_time" : 0             |
	+===============================+===============================+
*/
type GroupHistoryHandler struct {
	BaseCustomizedHandler
}

// Override
func (handler GroupHistoryHandler) HandleContent(content CustomizedContent, rMsg ReliableMessage, messenger Messenger) []Content {
	if content.Group() == nil {
		return handler.RespondReceipt("group command error.", rMsg.Envelope(), content, nil)
	}
	act := content.Action()
	if act == "query" {
		return handler.transformQueryCommand(content, rMsg, messenger)
	}
	//panic("unknown action: " + act)
	return handler.BaseCustomizedHandler.HandleContent(content, rMsg, messenger)
}

func (handler GroupHistoryHandler) transformQueryCommand(content CustomizedContent, rMsg ReliableMessage, messenger Messenger) []Content {
	//info := content.CopyMap(false)
	//info["type"] = ContentType.COMMAND
	//info["command"] = "query"
	//query := ParseContent(info)
	//if command, ok := query.(QueryCommand); ok {
	//	return messenger.ProcessContent(command, rMsg)
	//}
	return handler.RespondReceipt("Query Command error.", rMsg.Envelope(), content, nil)
}

type AppCustomizedFilter struct {
	//CustomizedContentFilter

	defaultHandler CustomizedContentHandler
	handlers       map[string]CustomizedContentHandler
}

func NewAppCustomizedFilter() *AppCustomizedFilter {
	return &AppCustomizedFilter{
		defaultHandler: &BaseCustomizedHandler{},
		handlers:       make(map[string]CustomizedContentHandler, 8),
	}
}

func (filter *AppCustomizedFilter) SetContentHandler(app, mod string, handler CustomizedContentHandler) {
	key := app + ":" + mod
	filter.handlers[key] = handler
}

// protected
func (filter *AppCustomizedFilter) GetContentHandler(app, mod string) CustomizedContentHandler {
	key := app + ":" + mod
	return filter.handlers[key]
}

// Override
func (filter *AppCustomizedFilter) FilterContent(content CustomizedContent, _ ReliableMessage) CustomizedContentHandler {
	//app := content.Application()
	app := content.GetString("app", "")
	mod := content.Module()
	handler := filter.GetContentHandler(app, mod)
	if handler == nil {
		return filter.defaultHandler
	}
	return handler
}
