/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"sync"
	"sync/atomic"

	"github.com/mcc-github/blockchain/gossip/common"
)




type ChannelDeMultiplexer struct {
	channels []*channel
	lock     *sync.RWMutex
	closed   int32
}


func NewChannelDemultiplexer() *ChannelDeMultiplexer {
	return &ChannelDeMultiplexer{
		channels: make([]*channel, 0),
		lock:     &sync.RWMutex{},
		closed:   int32(0),
	}
}

type channel struct {
	pred common.MessageAcceptor
	ch   chan interface{}
}

func (m *ChannelDeMultiplexer) isClosed() bool {
	return atomic.LoadInt32(&m.closed) == int32(1)
}



func (m *ChannelDeMultiplexer) Close() {
	defer func() {
		
		recover()
	}()
	atomic.StoreInt32(&m.closed, int32(1))
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, ch := range m.channels {
		close(ch.ch)
	}
}


func (m *ChannelDeMultiplexer) AddChannel(predicate common.MessageAcceptor) chan interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	ch := &channel{ch: make(chan interface{}, 10), pred: predicate}
	m.channels = append(m.channels, ch)
	return ch.ch
}



func (m *ChannelDeMultiplexer) DeMultiplex(msg interface{}) {
	defer func() {
		recover()
	}() 

	if m.isClosed() {
		return
	}

	m.lock.RLock()
	channels := m.channels
	m.lock.RUnlock()

	for _, ch := range channels {
		if ch.pred(msg) {
			ch.ch <- msg
		}
	}
}