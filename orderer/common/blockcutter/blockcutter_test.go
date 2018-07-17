/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockcutter

import (
	"testing"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter/mock"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/stretchr/testify/assert"
)


type ordererConfigFetcher interface {
	OrdererConfigFetcher
}


type ordererConfig interface {
	channelconfig.Orderer
}

func init() {
	flogging.SetModuleLevel(pkgLogID, "DEBUG")
}

var tx = &cb.Envelope{Payload: []byte("GOOD")}
var txLarge = &cb.Envelope{Payload: []byte("GOOD"), Signature: make([]byte, 1000)}

func TestNormalBatch(t *testing.T) {
	maxMessageCount := uint32(2)
	absoluteMaxBytes := uint32(1000)
	preferredMaxBytes := uint32(100)
	mockConfig := &mock.OrdererConfig{}
	mockConfig.BatchSizeReturns(&ab.BatchSize{
		MaxMessageCount:   maxMessageCount,
		AbsoluteMaxBytes:  absoluteMaxBytes,
		PreferredMaxBytes: preferredMaxBytes,
	})

	mockConfigFetcher := &mock.OrdererConfigFetcher{}
	mockConfigFetcher.OrdererConfigReturns(mockConfig, true)

	r := NewReceiverImpl(mockConfigFetcher)

	batches, pending := r.Ordered(tx)
	assert.Nil(t, batches, "Should not have created batch")
	assert.True(t, pending, "Should have message pending in the receiver")

	batches, pending = r.Ordered(tx)
	assert.NotNil(t, batches, "Should have created batch")
	assert.False(t, pending, "Should not have message pending in the receiver")
}

func TestBatchSizePreferredMaxBytesOverflow(t *testing.T) {
	txBytes := messageSizeBytes(tx)

	
	preferredMaxBytes := txBytes*10 - 1

	
	maxMessageCount := uint32(20)

	mockConfig := &mock.OrdererConfig{}
	mockConfig.BatchSizeReturns(&ab.BatchSize{
		MaxMessageCount:   maxMessageCount,
		AbsoluteMaxBytes:  preferredMaxBytes * 2,
		PreferredMaxBytes: preferredMaxBytes,
	})

	mockConfigFetcher := &mock.OrdererConfigFetcher{}
	mockConfigFetcher.OrdererConfigReturns(mockConfig, true)

	r := NewReceiverImpl(mockConfigFetcher)

	
	for i := 0; i < 9; i++ {
		batches, pending := r.Ordered(tx)
		assert.Nil(t, batches, "Should not have created batch")
		assert.True(t, pending, "Should have enqueued message into batch")
	}

	
	batches, pending := r.Ordered(tx)
	assert.NotNil(t, batches, "Should have created batch")
	assert.True(t, pending, "Should still have message pending")
	assert.Len(t, batches, 1, "Should have created one batch")
	assert.Len(t, batches[0], 9, "Should have had nine normal tx in the batch")

	
	messageBatch := r.Cut()
	assert.NotNil(t, batches, "Should have created batch")
	assert.Len(t, messageBatch, 1, "Should have had one tx in the batch")
}

func TestBatchSizePreferredMaxBytesOverflowNoPending(t *testing.T) {
	goodTxLargeBytes := messageSizeBytes(txLarge)

	
	preferredMaxBytes := goodTxLargeBytes - 1

	
	maxMessageCount := uint32(20)

	mockConfig := &mock.OrdererConfig{}
	mockConfig.BatchSizeReturns(&ab.BatchSize{
		MaxMessageCount:   maxMessageCount,
		AbsoluteMaxBytes:  preferredMaxBytes * 3,
		PreferredMaxBytes: preferredMaxBytes,
	})

	mockConfigFetcher := &mock.OrdererConfigFetcher{}
	mockConfigFetcher.OrdererConfigReturns(mockConfig, true)

	r := NewReceiverImpl(mockConfigFetcher)

	
	batches, pending := r.Ordered(tx)
	assert.Nil(t, batches, "Should not have created batch")
	assert.True(t, pending, "Should have enqueued message into batch")

	
	batches, pending = r.Ordered(txLarge)
	assert.NotNil(t, batches, "Should have created batch")
	assert.False(t, pending, "Should not have pending messages in receiver")
	assert.Len(t, batches, 2, "Should have created two batches")
	for i, batch := range batches {
		assert.Len(t, batch, 1, "Should have had one normal tx in batch %d", i)
	}
}

func TestPanicOnMissingConfig(t *testing.T) {
	mockConfigFetcher := &mock.OrdererConfigFetcher{}
	r := NewReceiverImpl(mockConfigFetcher)
	assert.Panics(t, func() { r.Ordered(tx) })
}
