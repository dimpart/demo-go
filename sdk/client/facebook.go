package sdk

import (
	. "github.com/dimpart/demo-go/sdk/common"
	. "github.com/dimpart/demo-go/sdk/database"
)

type IClientFacebook interface {
	ICommonFacebook
}

/**
 *  Client Facebook with Address Name Service
 */

type ClientFacebook struct {
	//IClientFacebook
	*CommonFacebook
}

func NewClientFacebook(db Database) *ClientFacebook {
	super := NewCommonFacebook(db)
	super.Database = db
	return &ClientFacebook{
		CommonFacebook: super,
	}
}

//
//  GroupDataSource
//
