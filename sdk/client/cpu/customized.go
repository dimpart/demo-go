/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package cpu

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimchat/sdk-go/sdk"
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
	app := content.Application()
	mod := content.Module()
	handler := filter.GetContentHandler(app, mod)
	if handler == nil {
		return filter.defaultHandler
	}
	return handler
}
