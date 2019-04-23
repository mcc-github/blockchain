/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/customtx"
	"github.com/mcc-github/blockchain/core/ledger/customtx/mock"
	"github.com/mcc-github/blockchain/protos/common"
	protopeer "github.com/mcc-github/blockchain/protos/peer"
	"github.com/stretchr/testify/assert"
)

func TestReadWriteCustomTxProcessor(t *testing.T) {
	fakeTxProcessor := &mock.Processor{}
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	customtx.InitializeTestEnv(
		customtx.Processors{
			100: fakeTxProcessor,
		},
	)
	h := newTestHelperCreateLgr("ledger1", t)
	h.simulateDataTx("tx0", func(s *simulator) {
		s.setState("ns", "key1", "value1")
		s.setState("ns", "key2", "value2")
		s.setState("ns", "key3", "value3")
	})
	h.cutBlockAndCommitWithPvtdata() 

	valueCounter := 0
	fakeTxProcessor.GenerateSimulationResultsStub =
		
		func(txEnvelop *common.Envelope, s ledger.TxSimulator, initializingLedger bool) error {
			valKey1, err := s.GetState("ns", "key1")
			assert.NoError(t, err)
			assert.Equal(t, []byte("value1"), valKey1)
			valueCounter++
			return s.SetState("ns", "key1", []byte(fmt.Sprintf("value1_%d", valueCounter)))
		}

	
	h.addPostOrderTx("tx1", 100)
	h.addPostOrderTx("tx2", 100)
	h.cutBlockAndCommitWithPvtdata()

	
	
	h.verifyTxValidationCode("tx1", protopeer.TxValidationCode_VALID)
	h.verifyTxValidationCode("tx2", protopeer.TxValidationCode_MVCC_READ_CONFLICT)
	h.verifyPubState("ns", "key1", "value1_1")
}

func TestRangeReadAndWriteCustomTxProcessor(t *testing.T) {
	fakeTxProcessor1 := &mock.Processor{}
	fakeTxProcessor2 := &mock.Processor{}
	fakeTxProcessor3 := &mock.Processor{}
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()
	customtx.InitializeTestEnv(
		customtx.Processors{
			101: fakeTxProcessor1,
			102: fakeTxProcessor2,
			103: fakeTxProcessor3,
		},
	)
	h := newTestHelperCreateLgr("ledger1", t)
	h.simulateDataTx("tx0", func(s *simulator) {
		s.setState("ns", "key1", "value1")
		s.setState("ns", "key2", "value2")
		s.setState("ns", "key3", "value3")
	})
	h.cutBlockAndCommitWithPvtdata() 

	fakeTxProcessor1.GenerateSimulationResultsStub =
		
		func(txEnvelop *common.Envelope, s ledger.TxSimulator, initializingLedger bool) error {
			return s.SetState("ns", "key1", []byte("value1_new"))
		}

	fakeTxProcessor2.GenerateSimulationResultsStub =
		
		func(txEnvelop *common.Envelope, s ledger.TxSimulator, initializingLedger bool) error {
			itr, err := s.GetStateRangeScanIterator("ns", "key1", "key2")
			assert.NoError(t, err)
			for {
				res, err := itr.Next()
				assert.NoError(t, err)
				if res == nil {
					break
				}
			}
			return s.SetState("ns", "key2", []byte("value2_new"))
		}

	fakeTxProcessor3.GenerateSimulationResultsStub =
		
		func(txEnvelop *common.Envelope, s ledger.TxSimulator, initializingLedger bool) error {
			itr, err := s.GetStateRangeScanIterator("ns", "key2", "key3")
			assert.NoError(t, err)
			for {
				res, err := itr.Next()
				assert.NoError(t, err)
				if res == nil {
					break
				}
			}
			return s.SetState("ns", "key3", []byte("value3_new"))
		}

	
	h.addPostOrderTx("tx1", 101)
	h.addPostOrderTx("tx2", 102)
	h.addPostOrderTx("tx3", 103)
	h.cutBlockAndCommitWithPvtdata()

	
	
	
	h.verifyTxValidationCode("tx1", protopeer.TxValidationCode_VALID)
	h.verifyTxValidationCode("tx2", protopeer.TxValidationCode_PHANTOM_READ_CONFLICT)
	h.verifyTxValidationCode("tx3", protopeer.TxValidationCode_VALID)
	h.verifyPubState("ns", "key1", "value1_new")
	h.verifyPubState("ns", "key2", "value2")
	h.verifyPubState("ns", "key3", "value3_new")
}
