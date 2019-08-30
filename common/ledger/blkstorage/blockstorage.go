/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blkstorage

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/ledger"
	l "github.com/mcc-github/blockchain/core/ledger"
	"github.com/pkg/errors"
)


type IndexableAttr string


const (
	IndexableAttrBlockNum         = IndexableAttr("BlockNum")
	IndexableAttrBlockHash        = IndexableAttr("BlockHash")
	IndexableAttrTxID             = IndexableAttr("TxID")
	IndexableAttrBlockNumTranNum  = IndexableAttr("BlockNumTranNum")
	IndexableAttrBlockTxID        = IndexableAttr("BlockTxID")
	IndexableAttrTxValidationCode = IndexableAttr("TxValidationCode")
)


type IndexConfig struct {
	AttrsToIndex []IndexableAttr
}

var (
	
	ErrNotFoundInIndex = l.NotFoundInIndexErr("")

	
	ErrAttrNotIndexed = errors.New("attribute not indexed")
)


type BlockStoreProvider interface {
	CreateBlockStore(ledgerid string) (BlockStore, error)
	OpenBlockStore(ledgerid string) (BlockStore, error)
	Exists(ledgerid string) (bool, error)
	List() ([]string, error)
	Close()
}




type BlockStore interface {
	AddBlock(block *common.Block) error
	GetBlockchainInfo() (*common.BlockchainInfo, error)
	RetrieveBlocks(startNum uint64) (ledger.ResultsIterator, error)
	RetrieveBlockByHash(blockHash []byte) (*common.Block, error)
	RetrieveBlockByNumber(blockNum uint64) (*common.Block, error) 
	RetrieveTxByID(txID string) (*common.Envelope, error)
	RetrieveTxByBlockNumTranNum(blockNum uint64, tranNum uint64) (*common.Envelope, error)
	RetrieveBlockByTxID(txID string) (*common.Block, error)
	RetrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error)
	Shutdown()
}
