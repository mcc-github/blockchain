/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package algo

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/mcc-github/blockchain/gossip/util"
	"github.com/spf13/viper"
)



const (
	defDigestWaitTime   = time.Duration(1000) * time.Millisecond
	defRequestWaitTime  = time.Duration(1500) * time.Millisecond
	defResponseWaitTime = time.Duration(2000) * time.Millisecond
)


func SetDigestWaitTime(time time.Duration) {
	viper.Set("peer.gossip.digestWaitTime", time)
}


func SetRequestWaitTime(time time.Duration) {
	viper.Set("peer.gossip.requestWaitTime", time)
}


func SetResponseWaitTime(time time.Duration) {
	viper.Set("peer.gossip.responseWaitTime", time)
}



type DigestFilter func(context interface{}) func(digestItem string) bool






type PullAdapter interface {
	
	SelectPeers() []string

	
	
	
	Hello(dest string, nonce uint64)

	
	
	SendDigest(digest []string, nonce uint64, context interface{})

	
	
	SendReq(dest string, items []string, nonce uint64)

	
	SendRes(items []string, context interface{}, nonce uint64)
}



type PullEngine struct {
	PullAdapter
	stopFlag           int32
	state              *util.Set
	item2owners        map[string][]string
	peers2nonces       map[string]uint64
	nonces2peers       map[uint64]string
	acceptingDigests   int32
	acceptingResponses int32
	lock               sync.Mutex
	outgoingNONCES     *util.Set
	incomingNONCES     *util.Set
	digFilter          DigestFilter

	digestWaitTime   time.Duration
	requestWaitTime  time.Duration
	responseWaitTime time.Duration
}



func NewPullEngineWithFilter(participant PullAdapter, sleepTime time.Duration, df DigestFilter) *PullEngine {
	engine := &PullEngine{
		PullAdapter:        participant,
		stopFlag:           int32(0),
		state:              util.NewSet(),
		item2owners:        make(map[string][]string),
		peers2nonces:       make(map[string]uint64),
		nonces2peers:       make(map[uint64]string),
		acceptingDigests:   int32(0),
		acceptingResponses: int32(0),
		incomingNONCES:     util.NewSet(),
		outgoingNONCES:     util.NewSet(),
		digFilter:          df,
		digestWaitTime:     util.GetDurationOrDefault("peer.gossip.digestWaitTime", defDigestWaitTime),
		requestWaitTime:    util.GetDurationOrDefault("peer.gossip.requestWaitTime", defRequestWaitTime),
		responseWaitTime:   util.GetDurationOrDefault("peer.gossip.responseWaitTime", defResponseWaitTime),
	}

	go func() {
		for !engine.toDie() {
			time.Sleep(sleepTime)
			if engine.toDie() {
				return
			}
			engine.initiatePull()
		}
	}()

	return engine
}



func NewPullEngine(participant PullAdapter, sleepTime time.Duration) *PullEngine {
	acceptAllFilter := func(_ interface{}) func(string) bool {
		return func(_ string) bool {
			return true
		}
	}
	return NewPullEngineWithFilter(participant, sleepTime, acceptAllFilter)
}

func (engine *PullEngine) toDie() bool {
	return atomic.LoadInt32(&(engine.stopFlag)) == int32(1)
}

func (engine *PullEngine) acceptResponses() {
	atomic.StoreInt32(&(engine.acceptingResponses), int32(1))
}

func (engine *PullEngine) isAcceptingResponses() bool {
	return atomic.LoadInt32(&(engine.acceptingResponses)) == int32(1)
}

func (engine *PullEngine) acceptDigests() {
	atomic.StoreInt32(&(engine.acceptingDigests), int32(1))
}

func (engine *PullEngine) isAcceptingDigests() bool {
	return atomic.LoadInt32(&(engine.acceptingDigests)) == int32(1)
}

func (engine *PullEngine) ignoreDigests() {
	atomic.StoreInt32(&(engine.acceptingDigests), int32(0))
}


func (engine *PullEngine) Stop() {
	atomic.StoreInt32(&(engine.stopFlag), int32(1))
}

func (engine *PullEngine) initiatePull() {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	engine.acceptDigests()
	for _, peer := range engine.SelectPeers() {
		nonce := engine.newNONCE()
		engine.outgoingNONCES.Add(nonce)
		engine.nonces2peers[nonce] = peer
		engine.peers2nonces[peer] = nonce
		engine.Hello(peer, nonce)
	}

	time.AfterFunc(engine.digestWaitTime, func() {
		engine.processIncomingDigests()
	})
}

func (engine *PullEngine) processIncomingDigests() {
	engine.ignoreDigests()

	engine.lock.Lock()
	defer engine.lock.Unlock()

	requestMapping := make(map[string][]string)
	for n, sources := range engine.item2owners {
		
		source := sources[util.RandomInt(len(sources))]
		if _, exists := requestMapping[source]; !exists {
			requestMapping[source] = make([]string, 0)
		}
		
		requestMapping[source] = append(requestMapping[source], n)
	}

	engine.acceptResponses()

	for dest, seqsToReq := range requestMapping {
		engine.SendReq(dest, seqsToReq, engine.peers2nonces[dest])
	}

	time.AfterFunc(engine.responseWaitTime, engine.endPull)
}

func (engine *PullEngine) endPull() {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	atomic.StoreInt32(&(engine.acceptingResponses), int32(0))
	engine.outgoingNONCES.Clear()

	engine.item2owners = make(map[string][]string)
	engine.peers2nonces = make(map[string]uint64)
	engine.nonces2peers = make(map[uint64]string)
}


func (engine *PullEngine) OnDigest(digest []string, nonce uint64, context interface{}) {
	if !engine.isAcceptingDigests() || !engine.outgoingNONCES.Exists(nonce) {
		return
	}

	engine.lock.Lock()
	defer engine.lock.Unlock()

	for _, n := range digest {
		if engine.state.Exists(n) {
			continue
		}

		if _, exists := engine.item2owners[n]; !exists {
			engine.item2owners[n] = make([]string, 0)
		}

		engine.item2owners[n] = append(engine.item2owners[n], engine.nonces2peers[nonce])
	}
}


func (engine *PullEngine) Add(seqs ...string) {
	for _, seq := range seqs {
		engine.state.Add(seq)
	}
}


func (engine *PullEngine) Remove(seqs ...string) {
	for _, seq := range seqs {
		engine.state.Remove(seq)
	}
}


func (engine *PullEngine) OnHello(nonce uint64, context interface{}) {
	engine.incomingNONCES.Add(nonce)

	time.AfterFunc(engine.requestWaitTime, func() {
		engine.incomingNONCES.Remove(nonce)
	})

	a := engine.state.ToArray()
	var digest []string
	filter := engine.digFilter(context)
	for _, item := range a {
		dig := item.(string)
		if !filter(dig) {
			continue
		}
		digest = append(digest, dig)
	}
	if len(digest) == 0 {
		return
	}
	engine.SendDigest(digest, nonce, context)
}


func (engine *PullEngine) OnReq(items []string, nonce uint64, context interface{}) {
	if !engine.incomingNONCES.Exists(nonce) {
		return
	}
	engine.lock.Lock()
	defer engine.lock.Unlock()

	filter := engine.digFilter(context)
	var items2Send []string
	for _, item := range items {
		if engine.state.Exists(item) && filter(item) {
			items2Send = append(items2Send, item)
		}
	}

	if len(items2Send) == 0 {
		return
	}

	go engine.SendRes(items2Send, context, nonce)
}


func (engine *PullEngine) OnRes(items []string, nonce uint64) {
	if !engine.outgoingNONCES.Exists(nonce) || !engine.isAcceptingResponses() {
		return
	}

	engine.Add(items...)
}

func (engine *PullEngine) newNONCE() uint64 {
	n := uint64(0)
	for {
		n = util.RandomUInt64()
		if !engine.outgoingNONCES.Exists(n) {
			return n
		}
	}
}
