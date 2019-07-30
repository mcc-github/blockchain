/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package state

import (
	"sync"
	"testing"

	"github.com/mcc-github/blockchain/gossip/discovery"
	"github.com/mcc-github/blockchain/gossip/metrics"
	gmetricsmocks "github.com/mcc-github/blockchain/gossip/metrics/mocks"
	"github.com/mcc-github/blockchain/gossip/protoext"
	"github.com/mcc-github/blockchain/gossip/state/mocks"
	proto "github.com/mcc-github/blockchain/protos/gossip"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetrics(t *testing.T) {
	t.Parallel()
	mc := &mockCommitter{Mock: &mock.Mock{}}
	mc.On("CommitLegacy", mock.Anything).Return(nil)
	mc.On("LedgerHeight", mock.Anything).Return(uint64(100), nil).Twice()
	mc.On("DoesPvtDataInfoExistInLedger", mock.Anything).Return(false, nil)
	g := &mocks.GossipMock{}
	g.On("PeersOfChannel", mock.Anything).Return([]discovery.NetworkMember{})
	g.On("Accept", mock.Anything, false).Return(make(<-chan *proto.GossipMessage), nil)
	g.On("Accept", mock.Anything, true).Return(nil, make(chan protoext.ReceivedMessage))

	heightWG := sync.WaitGroup{}
	heightWG.Add(1)
	committedDurationWG := sync.WaitGroup{}
	committedDurationWG.Add(1)

	testMetricProvider := gmetricsmocks.TestUtilConstructMetricProvider()

	testMetricProvider.FakeHeightGauge.SetStub = func(delta float64) {
		heightWG.Done()
	}
	testMetricProvider.FakeCommitDurationHist.ObserveStub = func(value float64) {
		committedDurationWG.Done()
	}

	gossipMetrics := metrics.NewGossipMetrics(testMetricProvider.FakeProvider)

	
	p := newPeerNodeWithGossipWithMetrics(0, mc, noopPeerIdentityAcceptor, g, gossipMetrics)
	defer p.shutdown()

	
	err := p.s.AddPayload(&proto.Payload{
		SeqNum: 100,
		Data:   protoutil.MarshalOrPanic(protoutil.NewBlock(100, []byte{})),
	})
	assert.NoError(t, err)

	
	mc.On("LedgerHeight", mock.Anything).Return(uint64(101), nil)

	
	heightWG.Wait()
	committedDurationWG.Wait()

	
	assert.Equal(t,
		[]string{"channel", "testchainid"},
		testMetricProvider.FakeHeightGauge.WithArgsForCall(0),
	)
	assert.EqualValues(t,
		101,
		testMetricProvider.FakeHeightGauge.SetArgsForCall(0),
	)

	
	assert.Equal(t,
		[]string{"channel", "testchainid"},
		testMetricProvider.FakePayloadBufferSizeGauge.WithArgsForCall(0),
	)
	assert.Equal(t,
		[]string{"channel", "testchainid"},
		testMetricProvider.FakePayloadBufferSizeGauge.WithArgsForCall(1),
	)
	
	size := testMetricProvider.FakePayloadBufferSizeGauge.SetArgsForCall(0)
	assert.True(t, size == 1 || size == 0)
	size = testMetricProvider.FakePayloadBufferSizeGauge.SetArgsForCall(1)
	assert.True(t, size == 1 || size == 0)

}
