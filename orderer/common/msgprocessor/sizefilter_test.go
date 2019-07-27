/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"testing"

	"github.com/golang/protobuf/proto"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/stretchr/testify/assert"
)

func TestMaxBytesRule(t *testing.T) {
	dataSize := uint32(100)
	maxBytes := calcMessageBytesForPayloadDataSize(dataSize)
	mcr := &mockconfig.Resources{
		OrdererConfigVal: &mockconfig.Orderer{
			BatchSizeVal: &ab.BatchSize{AbsoluteMaxBytes: maxBytes},
		},
	}
	msf := NewSizeFilter(mcr)

	t.Run("Less Than", func(t *testing.T) {
		assert.Nil(t, msf.Apply(makeMessage(make([]byte, dataSize-1))))
	})

	t.Run("Exact", func(t *testing.T) {
		assert.Nil(t, msf.Apply(makeMessage(make([]byte, dataSize))))
	})

	t.Run("Too Big", func(t *testing.T) {
		assert.NotNil(t, msf.Apply(makeMessage(make([]byte, dataSize+1))))
	})

	t.Run("Dynamic Resources", func(t *testing.T) {
		assert.NotNil(t, msf.Apply(makeMessage(make([]byte, dataSize+1))))
		mcr.OrdererConfigVal = &mockconfig.Orderer{
			BatchSizeVal: &ab.BatchSize{AbsoluteMaxBytes: maxBytes + 2},
		}
		assert.Nil(t, msf.Apply(makeMessage(make([]byte, dataSize+1))))
	})
}

func calcMessageBytesForPayloadDataSize(dataSize uint32) uint32 {
	return messageByteSize(makeMessage(make([]byte, dataSize)))
}

func makeMessage(data []byte) *cb.Envelope {
	data, err := proto.Marshal(&cb.Payload{Data: data})
	if err != nil {
		panic(err)
	}
	return &cb.Envelope{Payload: data}
}
