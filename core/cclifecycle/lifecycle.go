/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cclifecycle

import (
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)

var (
	
	
	Logger = flogging.MustGetLogger("discovery.lifecycle")
)


type MetadataManager struct {
	sync.RWMutex
	listeners              []MetadataChangeListener
	installedCCs           []chaincode.InstalledChaincode
	deployedCCsByChannel   map[string]*chaincode.MetadataMapping
	queryCreatorsByChannel map[string]QueryCreator
}





type MetadataChangeListener interface {
	HandleMetadataUpdate(channel string, chaincodes chaincode.MetadataSet)
}


type HandleMetadataUpdateFunc func(channel string, chaincodes chaincode.MetadataSet)



func (handleMetadataUpdate HandleMetadataUpdateFunc) HandleMetadataUpdate(channel string, chaincodes chaincode.MetadataSet) {
	handleMetadataUpdate(channel, chaincodes)
}




type Enumerator interface {
	
	Enumerate() ([]chaincode.InstalledChaincode, error)
}


type EnumerateFunc func() ([]chaincode.InstalledChaincode, error)


func (enumerate EnumerateFunc) Enumerate() ([]chaincode.InstalledChaincode, error) {
	return enumerate()
}




type Query interface {
	
	GetState(namespace string, key string) ([]byte, error)

	
	Done()
}




type QueryCreator interface {
	
	NewQuery() (Query, error)
}


type QueryCreatorFunc func() (Query, error)


func (queryCreator QueryCreatorFunc) NewQuery() (Query, error) {
	return queryCreator()
}


func NewMetadataManager(installedChaincodes Enumerator) (*MetadataManager, error) {
	installedCCs, err := installedChaincodes.Enumerate()
	if err != nil {
		return nil, errors.Wrap(err, "failed listing installed chaincodes")
	}

	return &MetadataManager{
		installedCCs:           installedCCs,
		deployedCCsByChannel:   map[string]*chaincode.MetadataMapping{},
		queryCreatorsByChannel: map[string]QueryCreator{},
	}, nil
}



func (m *MetadataManager) Metadata(channel string, cc string, collections bool) *chaincode.Metadata {
	queryCreator := m.queryCreatorsByChannel[channel]
	if queryCreator == nil {
		Logger.Warning("Requested Metadata for non-existent channel", channel)
		return nil
	}
	
	
	if md, found := m.deployedCCsByChannel[channel].Lookup(cc); found && !collections {
		Logger.Debug("Returning metadata for channel", channel, ", chaincode", cc, ":", md)
		return &md
	}
	query, err := queryCreator.NewQuery()
	if err != nil {
		Logger.Error("Failed obtaining new query for channel", channel, ":", err)
		return nil
	}
	md, err := DeployedChaincodes(query, AcceptAll, collections, cc)
	if err != nil {
		Logger.Error("Failed querying LSCC for channel", channel, ":", err)
		return nil
	}
	if len(md) == 0 {
		Logger.Info("Chaincode", cc, "isn't defined in channel", channel)
		return nil
	}

	return &md[0]
}

func (m *MetadataManager) initMetadataForChannel(channel string, queryCreator QueryCreator) error {
	if m.isChannelMetadataInitialized(channel) {
		return nil
	}
	
	query, err := queryCreator.NewQuery()
	if err != nil {
		return errors.WithStack(err)
	}
	ccs, err := queryChaincodeDefinitions(query, m.installedCCs, DeployedChaincodes)
	if err != nil {
		return errors.WithStack(err)
	}
	m.createMetadataForChannel(channel, queryCreator)
	m.updateState(channel, ccs)
	return nil
}

func (m *MetadataManager) createMetadataForChannel(channel string, newQuery QueryCreator) {
	m.Lock()
	defer m.Unlock()
	m.deployedCCsByChannel[channel] = chaincode.NewMetadataMapping()
	m.queryCreatorsByChannel[channel] = newQuery
}

func (m *MetadataManager) isChannelMetadataInitialized(channel string) bool {
	m.RLock()
	defer m.RUnlock()
	_, exists := m.deployedCCsByChannel[channel]
	return exists
}

func (m *MetadataManager) updateState(channel string, ccUpdate chaincode.MetadataSet) {
	m.RLock()
	defer m.RUnlock()
	for _, cc := range ccUpdate {
		m.deployedCCsByChannel[channel].Update(cc)
	}
}

func (m *MetadataManager) fireChangeListeners(channel string) {
	m.RLock()
	md := m.deployedCCsByChannel[channel]
	m.RUnlock()
	for _, listener := range m.listeners {
		aggregatedMD := md.Aggregate()
		listener.HandleMetadataUpdate(channel, aggregatedMD)
	}
	Logger.Debug("Listeners for channel", channel, "invoked")
}


func (m *MetadataManager) NewChannelSubscription(channel string, queryCreator QueryCreator) (*Subscription, error) {
	sub := &Subscription{
		metadataManager: m,
		channel:         channel,
		queryCreator:    queryCreator,
	}
	
	
	if err := m.initMetadataForChannel(channel, queryCreator); err != nil {
		return nil, errors.WithStack(err)
	}
	m.fireChangeListeners(channel)
	return sub, nil
}


func (m *MetadataManager) AddListener(listener MetadataChangeListener) {
	m.Lock()
	defer m.Unlock()
	m.listeners = append(m.listeners, listener)
}
