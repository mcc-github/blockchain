/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"sync"

	"github.com/mcc-github/blockchain/gossip/common"
	proto "github.com/mcc-github/blockchain/protos/gossip"
)



type MembershipStore struct {
	m map[string]*proto.SignedGossipMessage
	sync.RWMutex
}


func NewMembershipStore() *MembershipStore {
	return &MembershipStore{m: make(map[string]*proto.SignedGossipMessage)}
}



func (m *MembershipStore) MsgByID(pkiID common.PKIidType) *proto.SignedGossipMessage {
	m.RLock()
	defer m.RUnlock()
	if msg, exists := m.m[string(pkiID)]; exists {
		return msg
	}
	return nil
}


func (m *MembershipStore) Size() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}


func (m *MembershipStore) Put(pkiID common.PKIidType, msg *proto.SignedGossipMessage) {
	m.Lock()
	defer m.Unlock()
	m.m[string(pkiID)] = msg
}


func (m *MembershipStore) Remove(pkiID common.PKIidType) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, string(pkiID))
}



func (m *MembershipStore) ToSlice() []*proto.SignedGossipMessage {
	m.RLock()
	defer m.RUnlock()
	members := make([]*proto.SignedGossipMessage, len(m.m))
	i := 0
	for _, member := range m.m {
		members[i] = member
		i++
	}
	return members
}
