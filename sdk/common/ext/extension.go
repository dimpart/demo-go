package ext

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/plugins-go/ext"
	. "github.com/dimpart/demo-go/sdk/common/dkd"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

/**
 *  Extensions Loader
 *  ~~~~~~~~~~~~~~~~~
 */

type CommonExtensionLoader struct {
	ExtensionLoader
}

// Override
func (loader CommonExtensionLoader) Load() {
	loader.ExtensionLoader.Load()

	registerContentFactories()
	registerCustomizedFactories()

	registerCommandFactories()

}

/**
 *  Core content factories
 */
func registerContentFactories() {

	// Text
	copyContentFactory(ContentType.TEXT, "text")

	// File
	copyContentFactory(ContentType.FILE, "file")
	// Image
	copyContentFactory(ContentType.IMAGE, "image")
	// Audio
	copyContentFactory(ContentType.AUDIO, "audio")
	// Video
	copyContentFactory(ContentType.VIDEO, "video")

	// Web Page
	copyContentFactory(ContentType.PAGE, "page")

	// Name Card
	copyContentFactory(ContentType.NAME_CARD, "card")

	// Quote
	copyContentFactory(ContentType.QUOTE, "quote")

	// Money
	copyContentFactory(ContentType.MONEY, "money")
	copyContentFactory(ContentType.TRANSFER, "transfer")
	// ...

	// Command
	copyContentFactory(ContentType.COMMAND, "command")

	// History Command
	copyContentFactory(ContentType.HISTORY, "history")

	// Content Array
	copyContentFactory(ContentType.ARRAY, "array")

	// Combine and Forward
	copyContentFactory(ContentType.COMBINE_FORWARD, "combine")

	// Top-Secret
	copyContentFactory(ContentType.FORWARD, "forward")

	// unknown content type
	copyContentFactory(ContentType.ANY, "*")

}

/**
 *  Customized content factories
 */
func registerCustomizedFactories() {

	// Application Customized Content
	copyContentFactory(ContentType.CUSTOMIZED, "customized")
	copyContentFactory(ContentType.CUSTOMIZED, "application")
	copyContentFactory(ContentType.CUSTOMIZED, ContentType.APPLICATION)

}

func copyContentFactory(msgType, alias string) {
	factory := GetContentFactory(msgType)
	if factory == nil {
		panic("content factory not found:" + msgType)
	}
	SetContentFactory(alias, factory)
}

/**
 *  Core command factories
 */
func registerCommandFactories() {

	// Handshake
	registerCommandCreator(HANDSHAKE, func(dict StringKeyMap) Command {
		content := &BaseHandshakeCommand{}
		content.InitWithMap(dict)
		return content
	})

	// Login
	registerCommandCreator(LOGIN, func(dict StringKeyMap) Command {
		content := &BaseLoginCommand{}
		content.InitWithMap(dict)
		return content
	})

	// Mute
	registerCommandCreator(MUTE, func(dict StringKeyMap) Command {
		content := &BaseMuteCommand{}
		content.InitWithMap(dict)
		return content
	})

	// Block
	registerCommandCreator(BLOCK, func(dict StringKeyMap) Command {
		content := &BaseBlockCommand{}
		content.InitWithMap(dict)
		return content
	})

	// Report
	registerCommandCreator(REPORT, func(dict StringKeyMap) Command {
		content := &BaseReportCommand{}
		content.InitWithMap(dict)
		return content
	})
	registerCommandCreator(ONLINE, func(dict StringKeyMap) Command {
		content := &BaseReportCommand{}
		content.InitWithMap(dict)
		return content
	})
	registerCommandCreator(OFFLINE, func(dict StringKeyMap) Command {
		content := &BaseReportCommand{}
		content.InitWithMap(dict)
		return content
	})

}

func registerCommandCreator(cmd string, fn FuncCreateCommand) {
	factory := NewCommandFactory(fn)
	SetCommandFactory(cmd, factory)
}
