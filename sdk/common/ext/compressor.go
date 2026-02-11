package ext

import (
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/core"
)

/**
 *  Message Shortener
 *  ~~~~~~~~~~~~~~~~~
 */
type compatibleShortener struct {
	MessageShortener
}

// Override
func (shortener *compatibleShortener) CompressContent(content StringKeyMap) StringKeyMap {
	// DON'T COMPRESS NOW
	return content
}

// Override
func (shortener *compatibleShortener) CompressSymmetricKey(key StringKeyMap) StringKeyMap {
	// DON'T COMPRESS NOW
	return key
}

// Override
func (shortener *compatibleShortener) CompressReliableMessage(msg StringKeyMap) StringKeyMap {
	// DON'T COMPRESS NOW
	return msg
}

/**
 *  Message Compressor
 *  ~~~~~~~~~~~~~~~~~~
 */
type compatibleCompressor struct {
	MessageCompressor
}

// Override
func (compressor *compatibleCompressor) CompressContent(content StringKeyMap, key StringKeyMap) []byte {
	// TODO: fix outgoing content
	return compressor.MessageCompressor.CompressContent(content, key)
}

// Override
func (compressor *compatibleCompressor) ExtractContent(data []byte, key StringKeyMap) StringKeyMap {
	content := compressor.MessageCompressor.ExtractContent(data, key)
	// TODO: fix incoming content
	return content
}

/**
 *  Compress Factory
 *  ~~~~~~~~~~~~~~~~
 */
type compressFactory struct {
	//CompressFactory
}

// Override
func (factory compressFactory) CreateCompressor() Compressor {
	shortener := &compatibleShortener{}
	shortener.Init(ShortKeys.ContentKeys, ShortKeys.CryptoKeys, ShortKeys.MessageKeys)
	compressor := &compatibleCompressor{}
	compressor.Init(shortener)
	return compressor
}

func init() {
	factory := &compressFactory{}
	SetCompressFactory(factory)
}
