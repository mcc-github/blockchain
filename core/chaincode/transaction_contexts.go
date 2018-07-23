/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"sync"

	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type key string

const (
	
	
	TXSimulatorKey key = "txsimulatorkey"

	
	
	HistoryQueryExecutorKey key = "historyqueryexecutorkey"
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


func contextID(chainID, txID string) string {
	return chainID + txID
}




func (c *TransactionContexts) Create(txParams *ccprovider.TransactionParams) (*TransactionContext, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctxID := contextID(txParams.ChannelID, txParams.TxID)
	if c.contexts[ctxID] != nil {
		return nil, errors.Errorf("txid: %s(%s) exists", txParams.TxID, txParams.ChannelID)
	}

	txctx := &TransactionContext{
		ChainID:              txParams.ChannelID,
		SignedProp:           txParams.SignedProp,
		Proposal:             txParams.Proposal,
		ResponseNotifier:     make(chan *pb.ChaincodeMessage, 1),
		TXSimulator:          txParams.TXSimulator,
		HistoryQueryExecutor: txParams.HistoryQueryExecutor,
		queryIteratorMap:     map[string]commonledger.ResultsIterator{},
		pendingQueryResults:  map[string]*PendingQueryResult{},
	}
	c.contexts[ctxID] = txctx

	return txctx, nil
}

func getTxSimulator(ctx context.Context) ledger.TxSimulator {
	if txsim, ok := ctx.Value(TXSimulatorKey).(ledger.TxSimulator); ok {
		return txsim
	}
	return nil
}

func getHistoryQueryExecutor(ctx context.Context) ledger.HistoryQueryExecutor {
	if historyQueryExecutor, ok := ctx.Value(HistoryQueryExecutorKey).(ledger.HistoryQueryExecutor); ok {
		return historyQueryExecutor
	}
	return nil
}



func (c *TransactionContexts) Get(chainID, txID string) *TransactionContext {
	ctxID := contextID(chainID, txID)
	c.mutex.Lock()
	tc := c.contexts[ctxID]
	c.mutex.Unlock()
	return tc
}



func (c *TransactionContexts) Delete(chainID, txID string) {
	ctxID := contextID(chainID, txID)
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