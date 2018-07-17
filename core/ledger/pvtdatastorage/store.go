/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
)



type Provider interface {
	OpenStore(id string) (Store, error)
	Close()
}











type Store interface {
	
	Init(btlPolicy pvtdatapolicy.BTLPolicy)
	
	
	
	
	
	
	
	InitLastCommittedBlock(blockNum uint64) error
	
	
	
	GetPvtDataByBlockNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error)
	
	
	
	
	
	Prepare(blockNum uint64, pvtData []*ledger.TxPvtData) error
	
	Commit() error
	
	Rollback() error
	
	IsEmpty() (bool, error)
	
	LastCommittedBlockHeight() (uint64, error)
	
	HasPendingBatch() (bool, error)
	
	Shutdown()
}


type ErrIllegalCall struct {
	msg string
}

func (err *ErrIllegalCall) Error() string {
	return err.msg
}


type ErrIllegalArgs struct {
	msg string
}

func (err *ErrIllegalArgs) Error() string {
	return err.msg
}


type ErrOutOfRange struct {
	msg string
}

func (err *ErrOutOfRange) Error() string {
	return err.msg
}
