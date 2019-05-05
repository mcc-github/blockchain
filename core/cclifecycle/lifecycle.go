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


type Lifecycle struct {
	sync.RWMutex
	listeners              []LifecycleChangeListener
	installedCCs           []chaincode.InstalledChaincode
	deployedCCsByChannel   map[string]*chaincode.MetadataMapping
	queryCreatorsByChannel map[string]QueryCreator
}





type LifecycleChangeListener interface {
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


func NewLifecycle(installedChaincodes Enumerator) (*Lifecycle, error) {
	installedCCs, err := installedChaincodes.Enumerate()
	if err != nil {
		return nil, errors.Wrap(err, "failed listing installed chaincodes")
	}

	lc := &Lifecycle{
		installedCCs:           installedCCs,
		deployedCCsByChannel:   map[string]*chaincode.MetadataMapping{},
		queryCreatorsByChannel: map[string]QueryCreator{},
	}

	return lc, nil
}



func (lc *Lifecycle) Metadata(channel string, cc string, collections bool) *chaincode.Metadata {
	queryCreator := lc.queryCreatorsByChannel[channel]
	if queryCreator == nil {
		Logger.Warning("Requested Metadata for non-existent channel", channel)
		return nil
	}
	
	
	if md, found := lc.deployedCCsByChannel[channel].Lookup(cc); found && !collections {
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

func (lc *Lifecycle) initMetadataForChannel(channel string, queryCreator QueryCreator) error {
	if lc.isChannelMetadataInitialized(channel) {
		return nil
	}
	
	query, err := queryCreator.NewQuery()
	if err != nil {
		return errors.WithStack(err)
	}
	ccs, err := queryChaincodeDefinitions(query, lc.installedCCs, DeployedChaincodes)
	if err != nil {
		return errors.WithStack(err)
	}
	lc.createMetadataForChannel(channel, queryCreator)
	lc.updateState(channel, ccs)
	return nil
}

func (lc *Lifecycle) createMetadataForChannel(channel string, newQuery QueryCreator) {
	lc.Lock()
	defer lc.Unlock()
	lc.deployedCCsByChannel[channel] = chaincode.NewMetadataMapping()
	lc.queryCreatorsByChannel[channel] = newQuery
}

func (lc *Lifecycle) isChannelMetadataInitialized(channel string) bool {
	lc.RLock()
	defer lc.RUnlock()
	_, exists := lc.deployedCCsByChannel[channel]
	return exists
}

func (lc *Lifecycle) updateState(channel string, ccUpdate chaincode.MetadataSet) {
	lc.RLock()
	defer lc.RUnlock()
	for _, cc := range ccUpdate {
		lc.deployedCCsByChannel[channel].Update(cc)
	}
}

func (lc *Lifecycle) fireChangeListeners(channel string) {
	lc.RLock()
	md := lc.deployedCCsByChannel[channel]
	lc.RUnlock()
	for _, listener := range lc.listeners {
		aggregatedMD := md.Aggregate()
		listener.HandleMetadataUpdate(channel, aggregatedMD)
	}
	Logger.Debug("Listeners for channel", channel, "invoked")
}


func (lc *Lifecycle) NewChannelSubscription(channel string, queryCreator QueryCreator) (*Subscription, error) {
	sub := &Subscription{
		lc:           lc,
		channel:      channel,
		queryCreator: queryCreator,
	}
	
	
	if err := lc.initMetadataForChannel(channel, queryCreator); err != nil {
		return nil, errors.WithStack(err)
	}
	lc.fireChangeListeners(channel)
	return sub, nil
}


func (lc *Lifecycle) AddListener(listener LifecycleChangeListener) {
	lc.Lock()
	defer lc.Unlock()
	lc.listeners = append(lc.listeners, listener)
}
