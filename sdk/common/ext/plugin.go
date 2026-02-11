package ext

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/plugins-go/ext"
	. "github.com/dimpart/demo-go/sdk/common/digest"
	. "github.com/dimpart/demo-go/sdk/common/mkm"
)

/**
 *  Plugin Loader
 *  ~~~~~~~~~~~~~
 */
type CommonPluginLoader struct {
	PluginLoader
}

// Override
func (loader CommonPluginLoader) Load() {
	loader.PluginLoader.Load()

	registerDigesters()

	registerAddressFactory()

	registerMetaFactories()

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

/**
 *  Address factory
 */
func registerAddressFactory() {

	// Address
	SetAddressFactory(NewCompatibleAddressFactory())

}

/**
 *  Meta factories
 */
func registerMetaFactories() {
	mkm := NewCompatibleMetaFactory(MKM)
	btc := NewCompatibleMetaFactory(BTC)
	eth := NewCompatibleMetaFactory(ETH)

	SetMetaFactory("1", mkm)
	SetMetaFactory("2", btc)
	SetMetaFactory("4", eth)

	SetMetaFactory("mkm", mkm)
	SetMetaFactory("btc", btc)
	SetMetaFactory("eth", eth)

	SetMetaFactory("MKM", mkm)
	SetMetaFactory("BTC", btc)
	SetMetaFactory("ETH", eth)
}
