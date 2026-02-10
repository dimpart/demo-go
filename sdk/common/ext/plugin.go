package ext

import (
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/plugins-go/ext"
	. "github.com/dimpart/demo-go/sdk/common/digest"
)

type CommonPluginLoader struct {
	PluginLoader
}

// Override
func (loader CommonPluginLoader) Load() {
	loader.PluginLoader.Load()

	registerDigesters()
}

func registerDigesters() {

	// RipeMD-160
	SetRIPEMD160Digester(NewRIPEMD160Digester())

	// Keccak-256
	SetKECCAK256Digester(NewKECCAK256Digester())

	// SHA-1
	SetSHA1Digester(NewSHA1Digester())

	// MD-5
	SetMD5Digester(NewMD5Digester())

}
