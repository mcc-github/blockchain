/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"testing"

	"github.com/mcc-github/blockchain/common/channelconfig"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/core/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/ledger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChaincode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chaincode Suite")
}


type txSimulator interface {
	ledger.TxSimulator
}


type historyQueryExecutor interface {
	ledger.HistoryQueryExecutor
}


type queryResultsIterator interface {
	commonledger.QueryResultsIterator
}


type chaincodeRuntime interface {
	chaincode.Runtime
}


type certGenerator interface {
	chaincode.CertGenerator
}


type processor interface {
	chaincode.Processor
}


type invoker interface {
	chaincode.Invoker
}


type packageProvider interface {
	chaincode.PackageProvider
}




type chaincodeDefinition interface {
	ccprovider.ChaincodeDefinition
}


type chaincodeStream interface {
	ccintf.ChaincodeStream
}


type transactionRegistry interface {
	chaincode.TransactionRegistry
}


type systemCCProvider interface {
	chaincode.SystemCCProvider
}


type aclProvider interface {
	chaincode.ACLProvider
}


type chaincodeDefinitionGetter interface {
	chaincode.ChaincodeDefinitionGetter
}


type instantiationPolicyChecker interface {
	chaincode.InstantiationPolicyChecker
}


type ledgerGetter interface {
	chaincode.LedgerGetter
}


type peerLedger interface {
	ledger.PeerLedger
}




type launchRegistry interface {
	chaincode.LaunchRegistry
}


type messageHandler interface {
	chaincode.MessageHandler
}


type contextRegistry interface {
	chaincode.ContextRegistry
}


type queryResponseBuilder interface {
	chaincode.QueryResponseBuilder
}


type registry interface {
	chaincode.Registry
}


type applicationConfigRetriever interface {
	chaincode.ApplicationConfigRetriever
}


type collectionStore interface {
	privdata.CollectionStore
}


type applicationCapabilities interface {
	channelconfig.ApplicationCapabilities
}


type applicationConfig interface {
	channelconfig.Application
}
