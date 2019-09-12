/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"sync"

	pb "github.com/mcc-github/blockchain-protos-go/peer"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/pkg/errors"
)


type TransactionContexts struct {
	mutex    sync.Mutex
	contexts map[string]*TransactionContext
}


func NewTransactionContexts() *TransactionContexts {
	return &TransactionContexts{
		contexts: map[string]*TransactionContext{},
	}
}


func contextID(channelID, txID string) string {
	return channelID + txID
}




func (c *TransactionContexts) Create(txParams *ccprovider.TransactionParams) (*TransactionContext, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctxID := contextID(txParams.ChannelID, txParams.TxID)
	if c.contexts[ctxID] != nil {
		return nil, errors.Errorf("txid: %s(%s) exists", txParams.TxID, txParams.ChannelID)
	}

	txctx := &TransactionContext{
		NamespaceID:          txParams.NamespaceID,
		ChannelID:            txParams.ChannelID,
		SignedProp:           txParams.SignedProp,
		Proposal:             txParams.Proposal,
		ResponseNotifier:     make(chan *pb.ChaincodeMessage, 1),
		TXSimulator:          txParams.TXSimulator,
		HistoryQueryExecutor: txParams.HistoryQueryExecutor,
		CollectionStore:      txParams.CollectionStore,
		IsInitTransaction:    txParams.IsInitTransaction,

		queryIteratorMap:    map[string]commonledger.ResultsIterator{},
		pendingQueryResults: map[string]*PendingQueryResult{},
	}
	txctx.InitializeCollectionACLCache()

	c.contexts[ctxID] = txctx

	return txctx, nil
}



func (c *TransactionContexts) Get(channelID, txID string) *TransactionContext {
	ctxID := contextID(channelID, txID)
	c.mutex.Lock()
	tc := c.contexts[ctxID]
	c.mutex.Unlock()
	return tc
}



func (c *TransactionContexts) Delete(channelID, txID string) {
	ctxID := contextID(channelID, txID)
	c.mutex.Lock()
	delete(c.contexts, ctxID)
	c.mutex.Unlock()
}


func (c *TransactionContexts) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, txctx := range c.contexts {
		txctx.CloseQueryIterators()
	}
}
