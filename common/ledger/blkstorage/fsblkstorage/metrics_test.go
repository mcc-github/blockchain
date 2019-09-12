/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/common/metrics/metricsfakes"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/stretchr/testify/assert"
)

func TestStatsBlockchainHeight(t *testing.T) {
	testMetricProvider := testutilConstructMetricProvider()
	env := newTestEnvWithMetricsProvider(t, NewConf(testPath(), 0), testMetricProvider.fakeProvider)
	defer env.Cleanup()

	provider := env.provider
	ledgerid := "ledger-stats"
	store, err := provider.OpenBlockStore(ledgerid)
	assert.NoError(t, err)
	defer store.Shutdown()

	
	blockGenerator, genesisBlock := testutil.NewBlockGenerator(t, util.GetTestChannelID(), false)
	err = store.AddBlock(genesisBlock)
	assert.NoError(t, err)

	
	b1 := blockGenerator.NextBlock([][]byte{})
	err = store.AddBlock(b1)
	assert.NoError(t, err)

	
	fakeBlockchainHeightGauge := testMetricProvider.fakeBlockchainHeightGauge
	expectedCallCount := 3
	assert.Equal(t, expectedCallCount, fakeBlockchainHeightGauge.SetCallCount())

	
	assert.Equal(t, float64(0), fakeBlockchainHeightGauge.SetArgsForCall(0))
	assert.Equal(t, []string{"channel", ledgerid}, fakeBlockchainHeightGauge.WithArgsForCall(0))

	
	assert.Equal(t, float64(1), fakeBlockchainHeightGauge.SetArgsForCall(1))
	assert.Equal(t, []string{"channel", ledgerid}, fakeBlockchainHeightGauge.WithArgsForCall(1))

	
	assert.Equal(t, float64(2), fakeBlockchainHeightGauge.SetArgsForCall(2))
	assert.Equal(t, []string{"channel", ledgerid}, fakeBlockchainHeightGauge.WithArgsForCall(2))

	
	store.Shutdown()
	store, err = provider.OpenBlockStore(ledgerid)
	assert.NoError(t, err)

	
	assert.Equal(t, float64(2), fakeBlockchainHeightGauge.SetArgsForCall(3))
	assert.Equal(t, []string{"channel", ledgerid}, fakeBlockchainHeightGauge.WithArgsForCall(3))

	
	store.(*fsBlockStore).updateBlockStats(10, 1*time.Second)
	assert.Equal(t, float64(11), fakeBlockchainHeightGauge.SetArgsForCall(4))
	assert.Equal(t, []string{"channel", ledgerid}, fakeBlockchainHeightGauge.WithArgsForCall(4))
}

func TestStatsBlockCommit(t *testing.T) {
	testMetricProvider := testutilConstructMetricProvider()
	env := newTestEnvWithMetricsProvider(t, NewConf(testPath(), 0), testMetricProvider.fakeProvider)
	defer env.Cleanup()

	provider := env.provider
	ledgerid := "ledger-stats"
	store, err := provider.OpenBlockStore(ledgerid)
	assert.NoError(t, err)
	defer store.Shutdown()

	
	blockGenerator, genesisBlock := testutil.NewBlockGenerator(t, util.GetTestChannelID(), false)
	err = store.AddBlock(genesisBlock)
	assert.NoError(t, err)

	
	for i := 1; i <= 3; i++ {
		b := blockGenerator.NextBlock([][]byte{})
		err = store.AddBlock(b)
		assert.NoError(t, err)
	}

	fakeBlockstorageCommitTimeHist := testMetricProvider.fakeBlockstorageCommitTimeHist

	
	expectedCallCount := 1 + 3
	assert.Equal(t, expectedCallCount, fakeBlockstorageCommitTimeHist.ObserveCallCount())

	
	for i := 0; i < expectedCallCount; i++ {
		assert.Equal(t, []string{"channel", ledgerid}, fakeBlockstorageCommitTimeHist.WithArgsForCall(i))
	}

	
	store.(*fsBlockStore).updateBlockStats(4, 10*time.Second)
	assert.Equal(t,
		[]string{"channel", ledgerid},
		testMetricProvider.fakeBlockstorageCommitTimeHist.WithArgsForCall(4),
	)
	assert.Equal(t,
		float64(10),
		testMetricProvider.fakeBlockstorageCommitTimeHist.ObserveArgsForCall(4),
	)
}

type testMetricProvider struct {
	fakeProvider                   *metricsfakes.Provider
	fakeBlockchainHeightGauge      *metricsfakes.Gauge
	fakeBlockstorageCommitTimeHist *metricsfakes.Histogram
}

func testutilConstructMetricProvider() *testMetricProvider {
	fakeProvider := &metricsfakes.Provider{}
	fakeBlockchainHeightGauge := testutilConstructGauge()
	fakeBlockstorageCommitTimeHist := testutilConstructHist()
	fakeProvider.NewGaugeStub = func(opts metrics.GaugeOpts) metrics.Gauge {
		switch opts.Name {
		case blockchainHeightOpts.Name:
			return fakeBlockchainHeightGauge
		default:
			return nil
		}
	}
	fakeProvider.NewHistogramStub = func(opts metrics.HistogramOpts) metrics.Histogram {
		switch opts.Name {
		case blockstorageCommitTimeOpts.Name:
			return fakeBlockstorageCommitTimeHist
		default:
			return nil
		}
	}

	return &testMetricProvider{
		fakeProvider,
		fakeBlockchainHeightGauge,
		fakeBlockstorageCommitTimeHist,
	}
}

func testutilConstructGauge() *metricsfakes.Gauge {
	fakeGauge := &metricsfakes.Gauge{}
	fakeGauge.WithStub = func(lableValues ...string) metrics.Gauge {
		return fakeGauge
	}
	return fakeGauge
}

func testutilConstructHist() *metricsfakes.Histogram {
	fakeHist := &metricsfakes.Histogram{}
	fakeHist.WithStub = func(lableValues ...string) metrics.Histogram {
		return fakeHist
	}
	return fakeHist
}
