/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	subscriptionBuffSize = 50
)





type PubSub struct {
	sync.RWMutex

	
	subscriptions map[string]*Set
}



type Subscription interface {
	
	
	
	Listen() (interface{}, error)
}

type subscription struct {
	top string
	ttl time.Duration
	c   chan interface{}
}




func (s *subscription) Listen() (interface{}, error) {
	select {
	case <-time.After(s.ttl):
		return nil, errors.New("timed out")
	case item := <-s.c:
		return item, nil
	}
}



func NewPubSub() *PubSub {
	return &PubSub{
		subscriptions: make(map[string]*Set),
	}
}


func (ps *PubSub) Publish(topic string, item interface{}) error {
	ps.RLock()
	defer ps.RUnlock()
	s, subscribed := ps.subscriptions[topic]
	if !subscribed {
		return errors.New("no subscribers")
	}
	for _, sub := range s.ToArray() {
		c := sub.(*subscription).c
		
		if len(c) == subscriptionBuffSize {
			continue
		}
		c <- item
	}
	return nil
}


func (ps *PubSub) Subscribe(topic string, ttl time.Duration) Subscription {
	sub := &subscription{
		top: topic,
		ttl: ttl,
		c:   make(chan interface{}, subscriptionBuffSize),
	}

	ps.Lock()
	
	s, exists := ps.subscriptions[topic]
	
	if !exists {
		s = NewSet()
		ps.subscriptions[topic] = s
	}
	ps.Unlock()

	
	s.Add(sub)

	
	time.AfterFunc(ttl, func() {
		ps.unSubscribe(sub)
	})
	return sub
}

func (ps *PubSub) unSubscribe(sub *subscription) {
	ps.Lock()
	defer ps.Unlock()
	ps.subscriptions[sub.top].Remove(sub)
	if ps.subscriptions[sub.top].Size() != 0 {
		return
	}
	
	
	delete(ps.subscriptions, sub.top)

}
