/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/metrics/disabled"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/mock"
	"github.com/mcc-github/blockchain/protos/ledger/queryresult"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
	"github.com/stretchr/testify/assert"
)

func TestStateListener(t *testing.T) {
	env := newTestEnv(t)
	defer env.cleanup()
	provider, _ := NewProvider()

	
	channelid := "testLedger"
	namespace := "testchaincode"
	mockListener := &mockStateListener{namespace: namespace}
	provider.Initialize(&ledger.Initializer{
		DeployedChaincodeInfoProvider: &mock.DeployedChaincodeInfoProvider{},
		StateListeners:                []ledger.StateListener{mockListener},
		MetricsProvider:               &disabled.Provider{},
		Config:                        &ledger.Config{},
	})

	bg, gb := testutil.NewBlockGenerator(t, channelid, false)
	lgr, err := provider.Create(gb)
	
	sim1, err := lgr.NewTxSimulator("test_tx_1")
	assert.NoError(t, err)
	sim1.GetState(namespace, "key1")
	sim1.SetState(namespace, "key1", []byte("value1"))
	sim1.SetState(namespace, "key2", []byte("value2"))
	sim1.Done()

	
	sim2, err := lgr.NewTxSimulator("test_tx_2")
	assert.NoError(t, err)
	sim2.GetState(namespace, "key1")
	sim2.SetState(namespace, "key3", []byte("value3"))
	sim2.Done()

	
	sim3, err := lgr.NewTxSimulator("test_tx_3")
	assert.NoError(t, err)
	sim3.SetState(namespace, "key4", []byte("value4"))
	sim3.Done()

	
	mockListener.reset()
	sim1Res, _ := sim1.GetTxSimulationResults()
	sim1ResBytes, _ := sim1Res.GetPubSimulationBytes()
	assert.NoError(t, err)
	blk1 := bg.NextBlock([][]byte{sim1ResBytes})
	assert.NoError(t, lgr.CommitWithPvtData(&ledger.BlockAndPvtData{Block: blk1}))
	assert.Equal(t, channelid, mockListener.channelName)
	assert.Contains(t, mockListener.kvWrites, &kvrwset.KVWrite{Key: "key1", Value: []byte("value1")})
	assert.Contains(t, mockListener.kvWrites, &kvrwset.KVWrite{Key: "key2", Value: []byte("value2")})
	
	
	mockListener.reset()
	sim2Res, _ := sim2.GetTxSimulationResults()
	sim2ResBytes, _ := sim2Res.GetPubSimulationBytes()
	assert.NoError(t, err)
	blk2 := bg.NextBlock([][]byte{sim2ResBytes})
	assert.NoError(t, lgr.CommitWithPvtData(&ledger.BlockAndPvtData{Block: blk2}))
	assert.Equal(t, "", mockListener.channelName)
	assert.Nil(t, mockListener.kvWrites)

	
	mockListener.reset()
	sim3Res, _ := sim3.GetTxSimulationResults()
	sim3ResBytes, _ := sim3Res.GetPubSimulationBytes()
	assert.NoError(t, err)
	blk3 := bg.NextBlock([][]byte{sim3ResBytes})
	assert.NoError(t, lgr.CommitWithPvtData(&ledger.BlockAndPvtData{Block: blk3}))
	assert.Equal(t, channelid, mockListener.channelName)
	assert.Equal(t, []*kvrwset.KVWrite{
		{Key: "key4", Value: []byte("value4")},
	}, mockListener.kvWrites)

	provider.Close()
	provider, _ = NewProvider()
	provider.Initialize(&ledger.Initializer{
		DeployedChaincodeInfoProvider: &mock.DeployedChaincodeInfoProvider{},
		StateListeners:                []ledger.StateListener{mockListener},
		MetricsProvider:               &disabled.Provider{},
		Config:                        &ledger.Config{},
	})
	defer provider.Close()
	lgr, err = provider.Open(channelid)
	assert.NoError(t, err)
	defer lgr.Close()
	assert.NoError(t, err)
	assert.Equal(t,
		[]*queryresult.KV{
			{
				Namespace: namespace,
				Key:       "key1",
				Value:     []byte("value1"),
			},
			{
				Namespace: namespace,
				Key:       "key2",
				Value:     []byte("value2"),
			},
			{
				Namespace: namespace,
				Key:       "key4",
				Value:     []byte("value4"),
			},
		},
		mockListener.queryResultsInInitializeFunc,
	)
}

type mockStateListener struct {
	channelName                  string
	namespace                    string
	kvWrites                     []*kvrwset.KVWrite
	queryResultsInInitializeFunc []*queryresult.KV
}

func (l *mockStateListener) Initialize(ledgerID string, qe ledger.SimpleQueryExecutor) error {
	_, err := qe.GetPrivateDataHash(l.namespace, "random-coll", "random-key")
	if err != nil {
		return err
	}
	l.channelName = ledgerID
	itr, err := qe.GetStateRangeScanIterator(l.namespace, "", "")
	if err != nil {
		return err
	}
	for {
		res, err := itr.Next()
		if err != nil {
			return err
		}
		if res == nil {
			break
		}
		kv := res.(*queryresult.KV)
		l.queryResultsInInitializeFunc = append(l.queryResultsInInitializeFunc,
			kv,
		)
	}
	return nil
}

func (l *mockStateListener) InterestedInNamespaces() []string {
	return []string{l.namespace}
}

func (l *mockStateListener) HandleStateUpdates(trigger *ledger.StateUpdateTrigger) error {
	channelName, stateUpdates := trigger.LedgerID, trigger.StateUpdates
	l.channelName = channelName
	l.kvWrites = stateUpdates[l.namespace].PublicUpdates
	return nil
}

func (l *mockStateListener) StateCommitDone(channelID string) {
	
}

func (l *mockStateListener) reset() {
	l.channelName = ""
	l.kvWrites = nil
	l.queryResultsInInitializeFunc = nil
}
