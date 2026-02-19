package ext

import (
	. "github.com/dimchat/plugins-go/ext"
	. "github.com/dimpart/demo-go/sdk/common/ext"
)

type ILibraryLoader interface {
	Run()
}

type LibraryLoader struct {

	// flag
	loaded bool

	extensionLoader IExtensionLoader
	pluginLoader    IPluginLoader
}

func (loader *LibraryLoader) Init(extensionLoader IExtensionLoader, pluginLoader IPluginLoader) ILibraryLoader {
	if extensionLoader == nil {
		extensionLoader = &ClientExtensionLoader{}
	}
	if pluginLoader == nil {
		pluginLoader = &CommonPluginLoader{}
	}
	loader.extensionLoader = extensionLoader
	loader.pluginLoader = pluginLoader
	loader.loaded = false
	return loader
}

func (loader *LibraryLoader) Run() {
	if loader.loaded {
		// no need to load it again
		return
	}
	// mark it to loaded
	loader.loaded = true
	// try to load all extensions
	loader.Load()
}

// protected
func (loader *LibraryLoader) Load() {
	loader.extensionLoader.Load()
	loader.pluginLoader.Load()
}
