/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
)




type MetadataUpdateListener interface {
	HandleMetadataUpdate(channel string, metadata chaincode.MetadataSet)
}



type HandleMetadataUpdateFunc func(channel string, metadata chaincode.MetadataSet)




func (handleMetadataUpdate HandleMetadataUpdateFunc) HandleMetadataUpdate(channel string, metadata chaincode.MetadataSet) {
	handleMetadataUpdate(channel, metadata)
}





type MetadataManager struct {
	mutex             sync.Mutex
	listeners         []MetadataUpdateListener
	LegacyMetadataSet map[string]chaincode.MetadataSet
	MetadataSet       map[string]chaincode.MetadataSet
}

func NewMetadataManager() *MetadataManager {
	return &MetadataManager{
		LegacyMetadataSet: map[string]chaincode.MetadataSet{},
		MetadataSet:       map[string]chaincode.MetadataSet{},
	}
}







func (m *MetadataManager) HandleMetadataUpdate(channel string, metadata chaincode.MetadataSet) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.LegacyMetadataSet[channel] = metadata
	m.fireListenersForChannel(channel)
}





func (m *MetadataManager) UpdateMetadata(channel string, metadata chaincode.MetadataSet) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.MetadataSet[channel] = metadata
	m.fireListenersForChannel(channel)
}





func (m *MetadataManager) InitializeMetadata(channel string, metadata chaincode.MetadataSet) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.MetadataSet[channel] = metadata
}



func (m *MetadataManager) AddListener(listener MetadataUpdateListener) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.listeners = append(m.listeners, listener)
}


func (m *MetadataManager) fireListenersForChannel(channel string) {
	aggregatedMD := chaincode.MetadataSet{}
	mdMapNewLifecycle := map[string]struct{}{}

	for _, meta := range m.MetadataSet[channel] {
		mdMapNewLifecycle[meta.Name] = struct{}{}

		
		
		
		
		
		
		
		
		
		if meta.Installed && meta.Approved {
			aggregatedMD = append(aggregatedMD, meta)
		}
	}

	for _, meta := range m.LegacyMetadataSet[channel] {
		if _, in := mdMapNewLifecycle[meta.Name]; in {
			continue
		}

		aggregatedMD = append(aggregatedMD, meta)
	}

	for _, listener := range m.listeners {
		listener.HandleMetadataUpdate(channel, aggregatedMD)
	}
}
