/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockledger

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
	ab "github.com/mcc-github/blockchain-protos-go/orderer"
)


type Factory interface {
	
	
	GetOrCreate(chainID string) (ReadWriter, error)

	
	ChainIDs() []string

	
	Close()
}


type Iterator interface {
	
	
	Next() (*cb.Block, cb.Status)
	
	Close()
}


type Reader interface {
	
	
	Iterator(startType *ab.SeekPosition) (Iterator, uint64)
	
	Height() uint64
}


type Writer interface {
	
	Append(block *cb.Block) error
}




type ReadWriter interface {
	Reader
	Writer
}
