package dimp

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

type ICheckEmitter interface {
	IEntityRequest
	IEntityRespond
}

type CheckEmitter struct {
	//ICheckEmitter
}

// Override
func (checker CheckEmitter) QueryMeta(did ID) bool {
	return true
}

// Override
func (checker CheckEmitter) QueryDocuments(did ID, docs []Document) bool {
	return true
}

// Override
func (checker CheckEmitter) QueryMembers(gid ID, members []ID) bool {
	return true
}

// Override
func (checker CheckEmitter) SendVisa(visa Visa, receiver ID, updated bool) bool {
	return true
}
