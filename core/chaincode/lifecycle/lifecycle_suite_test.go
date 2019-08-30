/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle_test

import (
	"strings"
	"testing"

	"github.com/mcc-github/blockchain-chaincode-go/shim"
	"github.com/mcc-github/blockchain/common/channelconfig"
	commonledger "github.com/mcc-github/blockchain/common/ledger"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/chaincode/lifecycle"
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api/state"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/msp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)




type convertiblePolicy interface {
	policies.Policy
	policies.Converter
}


type inconvertiblePolicy interface {
	policies.Policy
}


type aclProvider interface {
	aclmgmt.ACLProvider
}


type chaincodeStub interface {
	shim.ChaincodeStubInterface
}


type stateIterator interface {
	shim.StateQueryIteratorInterface
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


type readWritableState interface {
	lifecycle.ReadWritableState
	lifecycle.OpaqueState
	lifecycle.RangeableState
}


type simpleQueryExecutor interface {
	ledger.SimpleQueryExecutor
}


type resultsIterator interface {
	commonledger.ResultsIterator
}


type channelConfig interface {
	channelconfig.Resources
}


type applicationConfig interface {
	channelconfig.Application
}


type applicationOrgConfig interface {
	channelconfig.ApplicationOrg
}


type policyManager interface {
	policies.Manager
}


type applicationCapabilities interface {
	channelconfig.ApplicationCapabilities
}


type validationState interface {
	validation.State
}


type metadataUpdateListener interface {
	lifecycle.MetadataUpdateListener
}


type chaincodeInfoProvider interface {
	lifecycle.ChaincodeInfoProvider
}


type legacyMetadataProvider interface {
	lifecycle.LegacyMetadataProvider
}


type metadataHandler interface {
	lifecycle.MetadataHandler
}


type mspManager interface {
	msp.MSPManager
}


type msp1 interface {
	msp.MSP
}

func TestLifecycle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifecycle Suite")
}



type MapLedgerShim map[string][]byte

func (m MapLedgerShim) PutState(key string, value []byte) error {
	m[key] = value
	return nil
}

func (m MapLedgerShim) DelState(key string) error {
	delete(m, key)
	return nil
}

func (m MapLedgerShim) GetState(key string) (value []byte, err error) {
	return m[key], nil
}

func (m MapLedgerShim) GetStateHash(key string) (value []byte, err error) {
	if val, ok := m[key]; ok {
		return util.ComputeSHA256(val), nil
	}
	return nil, nil
}

func (m MapLedgerShim) GetStateRange(prefix string) (map[string][]byte, error) {
	result := map[string][]byte{}
	for key, value := range m {
		if strings.HasPrefix(key, prefix) {
			result[key] = value
		}

	}
	return result, nil
}
