package ext

import (
	. "github.com/dimchat/plugins-go/ext"
	. "github.com/dimpart/demo-go/sdk/common/ext"
)

type ILibraryLoader interface {
	Run()
}

type LibraryLoader struct {
	extensionLoader IExtensionLoader
	pluginLoader    IPluginLoader

	// flag
	loaded bool
}

func NewLibraryLoader() ILibraryLoader {
	return &LibraryLoader{
		extensionLoader: &ClientExtensionLoader{},
		pluginLoader:    &CommonPluginLoader{},
		loaded:          false,
	}
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
