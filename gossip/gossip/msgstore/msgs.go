/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgstore

import (
	"sync"
	"time"

	"github.com/mcc-github/blockchain/gossip/common"
)

var noopLock = func() {}


func Noop(_ interface{}) {

}




type invalidationTrigger func(message interface{})



func NewMessageStore(pol common.MessageReplacingPolicy, trigger invalidationTrigger) MessageStore {
	return newMsgStore(pol, trigger)
}




func NewMessageStoreExpirable(pol common.MessageReplacingPolicy, trigger invalidationTrigger, msgTTL time.Duration, externalLock func(), externalUnlock func(), externalExpire func(interface{})) MessageStore {
	store := newMsgStore(pol, trigger)
	store.msgTTL = msgTTL

	if externalLock != nil {
		store.externalLock = externalLock
	}

	if externalUnlock != nil {
		store.externalUnlock = externalUnlock
	}

	if externalExpire != nil {
		store.expireMsgCallback = externalExpire
	}

	go store.expirationRoutine()
	return store
}

func newMsgStore(pol common.MessageReplacingPolicy, trigger invalidationTrigger) *messageStoreImpl {
	return &messageStoreImpl{
		pol:        pol,
		messages:   make([]*msg, 0),
		invTrigger: trigger,

		externalLock:      noopLock,
		externalUnlock:    noopLock,
		expireMsgCallback: func(m interface{}) {},
		expiredCount:      0,

		doneCh: make(chan struct{}),
	}
}







type MessageStore interface {
	
	
	Add(msg interface{}) bool

	
	
	CheckValid(msg interface{}) bool

	
	Size() int

	
	Get() []interface{}

	
	Stop()

	
	
	Purge(func(interface{}) bool)
}

type messageStoreImpl struct {
	pol               common.MessageReplacingPolicy
	lock              sync.RWMutex
	messages          []*msg
	invTrigger        invalidationTrigger
	msgTTL            time.Duration
	expiredCount      int
	externalLock      func()
	externalUnlock    func()
	expireMsgCallback func(msg interface{})
	doneCh            chan struct{}
	stopOnce          sync.Once
}

type msg struct {
	data    interface{}
	created time.Time
	expired bool
}


func (s *messageStoreImpl) Add(message interface{}) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	n := len(s.messages)
	for i := 0; i < n; i++ {
		m := s.messages[i]
		switch s.pol(message, m.data) {
		case common.MessageInvalidated:
			return false
		case common.MessageInvalidates:
			s.invTrigger(m.data)
			s.messages = append(s.messages[:i], s.messages[i+1:]...)
			n--
			i--
		}
	}

	s.messages = append(s.messages, &msg{data: message, created: time.Now()})
	return true
}

func (s *messageStoreImpl) Purge(shouldBePurged func(interface{}) bool) {
	shouldMsgBePurged := func(m *msg) bool {
		return shouldBePurged(m.data)
	}
	if !s.isPurgeNeeded(shouldMsgBePurged) {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	n := len(s.messages)
	for i := 0; i < n; i++ {
		if !shouldMsgBePurged(s.messages[i]) {
			continue
		}
		s.invTrigger(s.messages[i].data)
		s.messages = append(s.messages[:i], s.messages[i+1:]...)
		n--
		i--
	}
}


func (s *messageStoreImpl) CheckValid(message interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, m := range s.messages {
		if s.pol(message, m.data) == common.MessageInvalidated {
			return false
		}
	}
	return true
}


func (s *messageStoreImpl) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.messages) - s.expiredCount
}


func (s *messageStoreImpl) Get() []interface{} {
	res := make([]interface{}, 0)

	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, msg := range s.messages {
		if !msg.expired {
			res = append(res, msg.data)
		}
	}
	return res
}

func (s *messageStoreImpl) expireMessages() {
	s.externalLock()
	s.lock.Lock()
	defer s.lock.Unlock()
	defer s.externalUnlock()

	n := len(s.messages)
	for i := 0; i < n; i++ {
		m := s.messages[i]
		if !m.expired {
			if time.Since(m.created) > s.msgTTL {
				m.expired = true
				s.expireMsgCallback(m.data)
				s.expiredCount++
			}
		} else {
			if time.Since(m.created) > (s.msgTTL * 2) {
				s.messages = append(s.messages[:i], s.messages[i+1:]...)
				n--
				i--
				s.expiredCount--
			}

		}
	}
}

func (s *messageStoreImpl) isPurgeNeeded(shouldBePurged func(*msg) bool) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, m := range s.messages {
		if shouldBePurged(m) {
			return true
		}
	}
	return false
}

func (s *messageStoreImpl) expirationRoutine() {
	for {
		select {
		case <-s.doneCh:
			return
		case <-time.After(s.expirationCheckInterval()):
			hasMessageExpired := func(m *msg) bool {
				if !m.expired && time.Since(m.created) > s.msgTTL {
					return true
				} else if time.Since(m.created) > (s.msgTTL * 2) {
					return true
				}
				return false
			}
			if s.isPurgeNeeded(hasMessageExpired) {
				s.expireMessages()
			}
		}
	}
}

func (s *messageStoreImpl) Stop() {
	stopFunc := func() {
		close(s.doneCh)
	}
	s.stopOnce.Do(stopFunc)
}

func (s *messageStoreImpl) expirationCheckInterval() time.Duration {
	return s.msgTTL / 100
}
