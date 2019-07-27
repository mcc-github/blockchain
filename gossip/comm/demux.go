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
	
	lock   sync.Mutex
	closed bool 
	stopCh chan struct{}
	
	
	
	deMuxInProgress sync.WaitGroup
	channels        []*channel
}


func NewChannelDemultiplexer() *ChannelDeMultiplexer {
	return &ChannelDeMultiplexer{stopCh: make(chan struct{})}
}

type channel struct {
	pred common.MessageAcceptor
	ch   chan<- interface{}
}



func (m *ChannelDeMultiplexer) Close() {
	m.lock.Lock()
	if m.closed {
		m.lock.Unlock()
		return
	}
	m.closed = true
	close(m.stopCh)
	m.deMuxInProgress.Wait()
	for _, ch := range m.channels {
		close(ch.ch)
	}
	m.channels = nil
	m.lock.Unlock()
}







func (m *ChannelDeMultiplexer) AddChannel(predicate common.MessageAcceptor) <-chan interface{} {
	m.lock.Lock()
	if m.closed { 
		m.lock.Unlock()
		ch := make(chan interface{})
		close(ch)
		return ch
	}
	bidirectionalCh := make(chan interface{}, 10)
	
	
	ch := &channel{ch: bidirectionalCh, pred: predicate}
	m.channels = append(m.channels, ch)
	m.lock.Unlock()
	return bidirectionalCh
}





func (m *ChannelDeMultiplexer) DeMultiplex(msg interface{}) {
	m.lock.Lock()
	if m.closed {
		m.lock.Unlock()
		return
	}
	channels := m.channels
	m.deMuxInProgress.Add(1)
	m.lock.Unlock()

	for _, ch := range channels {
		if ch.pred(msg) {
			select {
			case <-m.stopCh:
				m.deMuxInProgress.Done()
				return 
			case ch.ch <- msg:
			}
		}
	}
	m.deMuxInProgress.Done()
}
