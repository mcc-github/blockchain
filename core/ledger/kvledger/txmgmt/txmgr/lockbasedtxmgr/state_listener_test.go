/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package lockbasedtxmgr

import (
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/privacyenabledstate"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/stretchr/testify/assert"
)

func TestStateListener(t *testing.T) {
	testLedgerid := "testLedger"
	ml1 := new(mock.StateListener)
	ml1.InterestedInNamespacesStub = func() []string { return []string{"ns1", "ns2"} }

	ml2 := new(mock.StateListener)
	ml2.InterestedInNamespacesStub = func() []string { return []string{"ns2", "ns3"} }

	ml3 := new(mock.StateListener)
	ml3.InterestedInNamespacesStub = func() []string { return []string{"ns4"} }

	testEnv := testEnvsMap[levelDBtestEnvName]
	testEnv.init(t, testLedgerid, nil)
	defer testEnv.cleanup()
	txmgr := testEnv.getTxMgr().(*LockBasedTxMgr)
	txmgr.stateListeners = []ledger.StateListener{ml1, ml2, ml3}

	
	
	sampleBatch := privacyenabledstate.NewUpdateBatch()
	sampleBatch.PubUpdates.Put("ns1", "key1_1", []byte("value1_1"), version.NewHeight(1, 1))
	sampleBatch.PubUpdates.Put("ns1", "key1_2", []byte("value1_2"), version.NewHeight(1, 2))
	sampleBatch.PubUpdates.Put("ns2", "key2_1", []byte("value2_1"), version.NewHeight(1, 3))
	sampleBatch.PubUpdates.Put("ns3", "key3_1", []byte("value3_1"), version.NewHeight(1, 4))
	dummyBlock := common.NewBlock(1, []byte("dummyHash"))
	txmgr.current = &current{block: dummyBlock, batch: sampleBatch}
	txmgr.invokeNamespaceListeners()
	assert.Equal(t, 1, ml1.HandleStateUpdatesCallCount())
	assert.Equal(t, 1, ml2.HandleStateUpdatesCallCount())
	assert.Equal(t, 0, ml3.HandleStateUpdatesCallCount())
	expectedLedgerid, expectedStateUpdate, expectedHt :=
		testLedgerid,
		ledger.StateUpdates{
			"ns1": []*kvrwset.KVWrite{
				{Key: "key1_1", Value: []byte("value1_1")}, {Key: "key1_2", Value: []byte("value1_2")}},
			"ns2": []*kvrwset.KVWrite{{Key: "key2_1", Value: []byte("value2_1")}},
		},
		uint64(1)
	checkHandleStateUpdatesCallback(t, ml1, 0, expectedLedgerid, expectedStateUpdate, expectedHt)
	expectedLedgerid, expectedStateUpdate, expectedHt =
		testLedgerid,
		ledger.StateUpdates{
			"ns2": []*kvrwset.KVWrite{{Key: "key2_1", Value: []byte("value2_1")}},
			"ns3": []*kvrwset.KVWrite{{Key: "key3_1", Value: []byte("value3_1")}},
		},
		uint64(1)
	checkHandleStateUpdatesCallback(t, ml2, 0, expectedLedgerid, expectedStateUpdate, expectedHt)
	txmgr.Commit()
	assert.Equal(t, 1, ml1.StateCommitDoneCallCount())
	assert.Equal(t, 1, ml2.StateCommitDoneCallCount())
	assert.Equal(t, 0, ml3.StateCommitDoneCallCount())

	
	
	sampleBatch = privacyenabledstate.NewUpdateBatch()
	sampleBatch.PubUpdates.Put("ns4", "key4_1", []byte("value4_1"), version.NewHeight(2, 1))
	txmgr.current = &current{block: common.NewBlock(2, []byte("anotherDummyHash")), batch: sampleBatch}
	txmgr.invokeNamespaceListeners()
	assert.Equal(t, 1, ml1.HandleStateUpdatesCallCount())
	assert.Equal(t, 1, ml2.HandleStateUpdatesCallCount())
	assert.Equal(t, 1, ml3.HandleStateUpdatesCallCount())

	expectedLedgerid, expectedStateUpdate, expectedHt =
		testLedgerid,
		ledger.StateUpdates{
			"ns4": []*kvrwset.KVWrite{{Key: "key4_1", Value: []byte("value4_1")}},
		},
		uint64(2)

	checkHandleStateUpdatesCallback(t, ml3, 0, expectedLedgerid, expectedStateUpdate, expectedHt)

	txmgr.Commit()
	assert.Equal(t, 1, ml1.StateCommitDoneCallCount())
	assert.Equal(t, 1, ml2.StateCommitDoneCallCount())
	assert.Equal(t, 1, ml3.StateCommitDoneCallCount())
}

func checkHandleStateUpdatesCallback(t *testing.T, ml *mock.StateListener, callNumber int,
	expectedLedgerid string,
	expectedUpdates ledger.StateUpdates,
	expectedCommitHt uint64) {
	actualNs, actualStateUpdate, actualHt := ml.HandleStateUpdatesArgsForCall(callNumber)
	assert.Equal(t, expectedLedgerid, actualNs)
	checkEqualUpdates(t, expectedUpdates, actualStateUpdate)
	assert.Equal(t, expectedCommitHt, actualHt)
}

func checkEqualUpdates(t *testing.T, expected, actual ledger.StateUpdates) {
	assert.Equal(t, len(expected), len(actual))
	for ns, expectedUpdates := range expected {
		assert.ElementsMatch(t, expectedUpdates, actual[ns])
	}
}
