/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"math/rand"
	"sort"
	"time"
)


type Filter interface {
	Filter(endorsers Endorsers) Endorsers
}



type ExclusionFilter interface {
	
	Exclude(Peer) bool
}

type selectionFunc func(Peer) bool

func (sf selectionFunc) Exclude(p Peer) bool {
	return sf(p)
}



type PrioritySelector interface {
	
	
	Compare(Peer, Peer) Priority
}






type Priority int

var (
	
	PrioritiesByHeight = &byHeight{}
	
	NoExclusion = selectionFunc(noExclusion)
	
	NoPriorities = &noPriorities{}
)

type noPriorities struct{}

func (nc noPriorities) Compare(_ Peer, _ Peer) Priority {
	return 0
}

type byHeight struct{}

func (*byHeight) Compare(left Peer, right Peer) Priority {
	leftHeight := left.StateInfoMessage.GetStateInfo().Properties.LedgerHeight
	rightHeight := right.StateInfoMessage.GetStateInfo().Properties.LedgerHeight

	if leftHeight > rightHeight {
		return 1
	}
	if rightHeight > leftHeight {
		return -1
	}
	return 0
}

func noExclusion(_ Peer) bool {
	return false
}


func ExcludeHosts(endpoints ...string) ExclusionFilter {
	m := make(map[string]struct{})
	for _, endpoint := range endpoints {
		m[endpoint] = struct{}{}
	}
	return ExcludeByHost(func(host string) bool {
		_, excluded := m[host]
		return excluded
	})
}


func ExcludeByHost(reject func(host string) bool) ExclusionFilter {
	return selectionFunc(func(p Peer) bool {
		endpoint := p.AliveMessage.GetAliveMsg().Membership.Endpoint
		var internalEndpoint string
		se := p.AliveMessage.GetSecretEnvelope()
		if se != nil {
			internalEndpoint = se.InternalEndpoint()
		}
		return reject(endpoint) || reject(internalEndpoint)
	})
}


func (endorsers Endorsers) Filter(f ExclusionFilter) Endorsers {
	var res Endorsers
	for _, e := range endorsers {
		if !f.Exclude(*e) {
			res = append(res, e)
		}
	}
	return res
}


func (endorsers Endorsers) Shuffle() Endorsers {
	res := make(Endorsers, len(endorsers))
	rand.Seed(time.Now().UnixNano())
	for i, index := range rand.Perm(len(endorsers)) {
		res[i] = endorsers[index]
	}
	return res
}

type endorserSort struct {
	Endorsers
	PrioritySelector
}


func (endorsers Endorsers) Sort(ps PrioritySelector) Endorsers {
	sort.Sort(&endorserSort{
		Endorsers:        endorsers,
		PrioritySelector: ps,
	})
	return endorsers
}

func (es *endorserSort) Len() int {
	return len(es.Endorsers)
}

func (es *endorserSort) Less(i, j int) bool {
	e1 := es.Endorsers[i]
	e2 := es.Endorsers[j]
	less := es.Compare(*e1, *e2)
	return less > Priority(0)
}

func (es *endorserSort) Swap(i, j int) {
	es.Endorsers[i], es.Endorsers[j] = es.Endorsers[j], es.Endorsers[i]
}
