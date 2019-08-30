/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lscc_test

import (
	"testing"

	"github.com/mcc-github/blockchain-chaincode-go/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type systemChaincodeProvider interface {
	sysccprovider.SystemChaincodeProvider
}


type queryExecutor interface {
	ledger.QueryExecutor
}


type fileSystemSupport interface {
	lscc.FilesystemSupport
}


type ccPackage interface {
	ccprovider.CCPackage
}


type chaincodeStub interface {
	shim.ChaincodeStubInterface
}


type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestLscc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lscc Suite")
}
