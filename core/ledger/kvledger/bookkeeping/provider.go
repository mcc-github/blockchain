/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package bookkeeping

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/ledger/util/leveldbhelper"
)


type Category int

const (
	
	PvtdataExpiry Category = iota
	
	MetadataPresenceIndicator
)


type Provider interface {
	
	GetDBHandle(ledgerID string, cat Category) *leveldbhelper.DBHandle
	
	Close()
}

type provider struct {
	dbProvider *leveldbhelper.Provider
}


func NewProvider(dbPath string) Provider {
	dbProvider := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: dbPath})
	return &provider{dbProvider: dbProvider}
}


func (provider *provider) GetDBHandle(ledgerID string, cat Category) *leveldbhelper.DBHandle {
	return provider.dbProvider.GetDBHandle(fmt.Sprintf(ledgerID+"/%d", cat))
}


func (provider *provider) Close() {
	provider.dbProvider.Close()
}
