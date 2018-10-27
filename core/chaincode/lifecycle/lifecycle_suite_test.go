/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle_test

import (
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


type chaincodeStub interface {
	shim.ChaincodeStubInterface
}


type chaincodeStore interface {
	lifecycle.ChaincodeStore
}


type packageParser interface {
	lifecycle.PackageParser
}


type sccFunctions interface {
	lifecycle.SCCFunctions
}

func TestLifecycle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifecycle Suite")
}
