/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm

import (
	"sync"

	"github.com/mcc-github/blockchain/gossip/common"
)




type ChannelDeMultiplexer struct {
	channels []*channel
	lock     *sync.RWMutex
	closed   bool
}


func NewChannelDemultiplexer() *ChannelDeMultiplexer {
	return &ChannelDeMultiplexer{
		channels: make([]*channel, 0),
		lock:     &sync.RWMutex{},
	}
}

type channel struct {
	pred common.MessageAcceptor
	ch   chan interface{}
}

func (m *ChannelDeMultiplexer) isClosed() bool {
	return m.closed
}



func (m *ChannelDeMultiplexer) Close() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.closed = true
	for _, ch := range m.channels {
		close(ch.ch)
	}
	m.channels = nil
}


func (m *ChannelDeMultiplexer) AddChannel(predicate common.MessageAcceptor) chan interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	ch := &channel{ch: make(chan interface{}, 10), pred: predicate}
	m.channels = append(m.channels, ch)
	return ch.ch
}



func (m *ChannelDeMultiplexer) DeMultiplex(msg interface{}) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.isClosed() {
		return
	}
	for _, ch := range m.channels {
		if ch.pred(msg) {
			ch.ch <- msg
		}
	}
}
