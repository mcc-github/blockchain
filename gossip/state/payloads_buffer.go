/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package state

import (
	"sync"
	"sync/atomic"

	proto "github.com/mcc-github/blockchain-protos-go/gossip"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/gossip/util"
)





type PayloadsBuffer interface {
	
	Push(payload *proto.Payload)

	
	Next() uint64

	
	Pop() *proto.Payload

	
	Size() int

	
	
	Ready() chan struct{}

	Close()
}



type PayloadsBufferImpl struct {
	next uint64

	buf map[uint64]*proto.Payload

	readyChan chan struct{}

	mutex sync.RWMutex

	logger util.Logger
}


func NewPayloadsBuffer(next uint64) PayloadsBuffer {
	return &PayloadsBufferImpl{
		buf:       make(map[uint64]*proto.Payload),
		readyChan: make(chan struct{}, 1),
		next:      next,
		logger:    util.GetLogger(util.StateLogger, ""),
	}
}




func (b *PayloadsBufferImpl) Ready() chan struct{} {
	return b.readyChan
}





func (b *PayloadsBufferImpl) Push(payload *proto.Payload) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	seqNum := payload.SeqNum

	if seqNum < b.next || b.buf[seqNum] != nil {
		logger.Debugf("Payload with sequence number = %d has been already processed", payload.SeqNum)
		return
	}

	b.buf[seqNum] = payload

	
	if seqNum == b.next && len(b.readyChan) == 0 {
		b.readyChan <- struct{}{}
	}
}


func (b *PayloadsBufferImpl) Next() uint64 {
	
	return atomic.LoadUint64(&b.next)
}



func (b *PayloadsBufferImpl) Pop() *proto.Payload {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	result := b.buf[b.Next()]

	if result != nil {
		
		delete(b.buf, b.Next())
		
		atomic.AddUint64(&b.next, 1)

		b.drainReadChannel()

	}

	return result
}




func (b *PayloadsBufferImpl) drainReadChannel() {
	if len(b.buf) == 0 {
		for {
			if len(b.readyChan) > 0 {
				<-b.readyChan
			} else {
				break
			}
		}
	}
}


func (b *PayloadsBufferImpl) Size() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return len(b.buf)
}


func (b *PayloadsBufferImpl) Close() {
	close(b.readyChan)
}

type metricsBuffer struct {
	PayloadsBuffer
	sizeMetrics metrics.Gauge
	chainID     string
}

func (mb *metricsBuffer) Push(payload *proto.Payload) {
	mb.PayloadsBuffer.Push(payload)
	mb.reportSize()
}

func (mb *metricsBuffer) Pop() *proto.Payload {
	pl := mb.PayloadsBuffer.Pop()
	mb.reportSize()
	return pl
}

func (mb *metricsBuffer) reportSize() {
	mb.sizeMetrics.With("channel", mb.chainID).Set(float64(mb.Size()))
}
