/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/ledger"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
	"github.com/mcc-github/blockchain/protoutil"

	ccpersistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/pkg/errors"
)

type LocalChaincodeInfo struct {
	Definition  *ChaincodeDefinition
	Approved    bool
	InstallInfo *ChaincodeInstallInfo
}

type ChaincodeInstallInfo struct {
	PackageID ccpersistence.PackageID
	Type      string
	Path      string
	Label     string
}

type CachedChaincodeDefinition struct {
	Definition  *ChaincodeDefinition
	Approved    bool
	InstallInfo *ChaincodeInstallInfo

	
	
	
	Hashes []string
}

type ChannelCache struct {
	Chaincodes map[string]*CachedChaincodeDefinition

	
	
	
	
	InterestingHashes map[string]string
}



type MetadataHandler interface {
	InitializeMetadata(channel string, chaincodes chaincode.MetadataSet)
	UpdateMetadata(channel string, chaincodes chaincode.MetadataSet)
}

type Cache struct {
	definedChaincodes map[string]*ChannelCache
	Resources         *Resources
	MyOrgMSPID        string

	
	
	
	
	
	
	mutex sync.RWMutex

	
	
	
	localChaincodes map[string]*LocalChaincode
	eventBroker     *EventBroker
	MetadataHandler MetadataHandler
}

type LocalChaincode struct {
	Info       *ChaincodeInstallInfo
	References map[string]map[string]*CachedChaincodeDefinition
}

func NewCache(resources *Resources, myOrgMSPID string, metadataManager MetadataHandler) *Cache {
	return &Cache{
		definedChaincodes: map[string]*ChannelCache{},
		localChaincodes:   map[string]*LocalChaincode{},
		Resources:         resources,
		MyOrgMSPID:        myOrgMSPID,
		eventBroker:       NewEventBroker(resources.ChaincodeStore, resources.PackageParser),
		MetadataHandler:   metadataManager,
	}
}





func (c *Cache) InitializeLocalChaincodes() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ccPackages, err := c.Resources.ChaincodeStore.ListInstalledChaincodes()
	if err != nil {
		return errors.WithMessage(err, "could not list installed chaincodes")
	}

	for _, ccPackage := range ccPackages {
		ccPackageBytes, err := c.Resources.ChaincodeStore.Load(ccPackage.PackageID)
		if err != nil {
			return errors.WithMessagef(err, "could not load chaincode with pakcage ID '%s'", ccPackage.PackageID.String())
		}
		parsedCCPackage, err := c.Resources.PackageParser.Parse(ccPackageBytes)
		if err != nil {
			return errors.WithMessagef(err, "could not parse chaincode with pakcage ID '%s'", ccPackage.PackageID.String())
		}
		c.handleChaincodeInstalledWhileLocked(true, parsedCCPackage.Metadata, ccPackage.PackageID)
	}

	logger.Infof("Initialized lifecycle cache with %d already installed chaincodes", len(c.localChaincodes))
	for channelID, chaincodeCache := range c.definedChaincodes {
		approved, installed, runnable := 0, 0, 0
		for _, cachedChaincode := range chaincodeCache.Chaincodes {
			if cachedChaincode.Approved {
				approved++
			}
			if cachedChaincode.InstallInfo != nil {
				installed++
			}
			if cachedChaincode.Approved && cachedChaincode.InstallInfo != nil {
				runnable++
			}
		}

		logger.Infof("Initialized lifecycle cache for channel '%s' with %d chaincodes runnable (%d approved, %d installed)", channelID, runnable, approved, installed)
	}

	return nil
}





func (c *Cache) Initialize(channelID string, qe ledger.SimpleQueryExecutor) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	publicState := &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	}

	metadatas, err := c.Resources.Serializer.DeserializeAllMetadata(NamespacesName, publicState)
	if err != nil {
		return errors.WithMessage(err, "could not query namespace metadata")
	}

	dirtyChaincodes := map[string]struct{}{}

	for namespace, metadata := range metadatas {
		switch metadata.Datatype {
		case ChaincodeDefinitionType:
			dirtyChaincodes[namespace] = struct{}{}
		default:
			
		}
	}

	return c.update(true, channelID, dirtyChaincodes, qe)
}


func (c *Cache) HandleChaincodeInstalled(md *persistence.ChaincodePackageMetadata, packageID ccpersistence.PackageID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handleChaincodeInstalledWhileLocked(false, md, packageID)
}

func (c *Cache) handleChaincodeInstalledWhileLocked(initializing bool, md *persistence.ChaincodePackageMetadata, packageID ccpersistence.PackageID) {
	
	
	encodedCCHash := protoutil.MarshalOrPanic(&lb.StateData{
		Type: &lb.StateData_String_{String_: packageID.String()},
	})
	hashOfCCHash := string(util.ComputeSHA256(encodedCCHash))
	localChaincode, ok := c.localChaincodes[hashOfCCHash]
	if !ok {
		localChaincode = &LocalChaincode{
			References: map[string]map[string]*CachedChaincodeDefinition{},
		}
		c.localChaincodes[hashOfCCHash] = localChaincode
	}
	localChaincode.Info = &ChaincodeInstallInfo{
		PackageID: packageID,
		Type:      md.Type,
		Path:      md.Path,
		Label:     md.Label,
	}
	for channelID, channelCache := range localChaincode.References {
		for chaincodeName, cachedChaincode := range channelCache {
			cachedChaincode.InstallInfo = localChaincode.Info
			logger.Infof("Installed chaincode with package ID '%s' now available on channel %s for chaincode definition %s:%s", packageID, channelID, chaincodeName, cachedChaincode.Definition.EndorsementInfo.Version)
		}
	}

	if !initializing {
		c.eventBroker.ProcessInstallEvent(localChaincode)
		c.handleMetadataUpdates(localChaincode)
	}
}



func (c *Cache) HandleStateUpdates(trigger *ledger.StateUpdateTrigger) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	channelID := trigger.LedgerID
	updates, ok := trigger.StateUpdates[LifecycleNamespace]
	if !ok {
		return errors.Errorf("no state updates for promised namespace _lifecycle")
	}

	dirtyChaincodes := map[string]struct{}{}

	for _, publicUpdate := range updates.PublicUpdates {
		matches := SequenceMatcher.FindStringSubmatch(publicUpdate.Key)
		if len(matches) != 2 {
			continue
		}

		dirtyChaincodes[matches[1]] = struct{}{}
	}

	channelCache, ok := c.definedChaincodes[channelID]

	
	if ok {
		for collection, privateUpdates := range updates.CollHashUpdates {
			matches := ImplicitCollectionMatcher.FindStringSubmatch(collection)
			if len(matches) != 2 {
				
				continue
			}

			if matches[1] != c.MyOrgMSPID {
				
				continue
			}

			for _, privateUpdate := range privateUpdates {
				chaincodeName, ok := channelCache.InterestingHashes[string(privateUpdate.KeyHash)]
				if ok {
					dirtyChaincodes[chaincodeName] = struct{}{}
				}
			}
		}
	}

	err := c.update(false, channelID, dirtyChaincodes, trigger.PostCommitQueryExecutor)
	if err != nil {
		return errors.WithMessage(err, "error updating cache")
	}

	return nil
}


func (c *Cache) InterestedInNamespaces() []string {
	return []string{LifecycleNamespace}
}


func (c *Cache) StateCommitDone(channelName string) {
	
	
	
	
	
	
	
	
	
	c.eventBroker.ApproveOrDefineCommitted(channelName)
}



func (c *Cache) ChaincodeInfo(channelID, name string) (*LocalChaincodeInfo, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	channelChaincodes, ok := c.definedChaincodes[channelID]
	if !ok {
		return nil, errors.Errorf("unknown channel '%s'", channelID)
	}

	cachedChaincode, ok := channelChaincodes.Chaincodes[name]
	if !ok {
		return nil, errors.Errorf("unknown chaincode '%s' for channel '%s'", name, channelID)
	}

	return &LocalChaincodeInfo{
		Definition:  cachedChaincode.Definition,
		InstallInfo: cachedChaincode.InstallInfo,
		Approved:    cachedChaincode.Approved,
	}, nil
}


func (c *Cache) update(initializing bool, channelID string, dirtyChaincodes map[string]struct{}, qe ledger.SimpleQueryExecutor) error {
	channelCache, ok := c.definedChaincodes[channelID]
	if !ok {
		channelCache = &ChannelCache{
			Chaincodes:        map[string]*CachedChaincodeDefinition{},
			InterestingHashes: map[string]string{},
		}
		c.definedChaincodes[channelID] = channelCache
	}

	publicState := &SimpleQueryExecutorShim{
		Namespace:           LifecycleNamespace,
		SimpleQueryExecutor: qe,
	}

	orgState := &PrivateQueryExecutorShim{
		Namespace:  LifecycleNamespace,
		Collection: ImplicitCollectionNameForOrg(c.MyOrgMSPID),
		State:      qe,
	}

	for name := range dirtyChaincodes {
		logger.Infof("Updating cached definition for chaincode '%s' on channel '%s'", name, channelID)
		cachedChaincode, ok := channelCache.Chaincodes[name]
		if !ok {
			cachedChaincode = &CachedChaincodeDefinition{}
			channelCache.Chaincodes[name] = cachedChaincode
		}

		for _, hash := range cachedChaincode.Hashes {
			delete(channelCache.InterestingHashes, hash)
		}

		exists, chaincodeDefinition, err := c.Resources.ChaincodeDefinitionIfDefined(name, publicState)
		if err != nil {
			return errors.WithMessagef(err, "could not get chaincode definition for '%s' on channel '%s'", name, channelID)
		}

		if !exists {
			
			
			delete(channelCache.Chaincodes, name)
			continue
		}

		cachedChaincode.Definition = chaincodeDefinition
		cachedChaincode.Approved = false

		privateName := fmt.Sprintf("%s#%d", name, chaincodeDefinition.Sequence)

		cachedChaincode.Hashes = []string{
			string(util.ComputeSHA256([]byte(MetadataKey(NamespacesName, privateName)))),
			string(util.ComputeSHA256([]byte(FieldKey(NamespacesName, privateName, "EndorsementInfo")))),
			string(util.ComputeSHA256([]byte(FieldKey(NamespacesName, privateName, "ValidationInfo")))),
			string(util.ComputeSHA256([]byte(FieldKey(NamespacesName, privateName, "Collections")))),
			string(util.ComputeSHA256([]byte(FieldKey(ChaincodeSourcesName, privateName, "PackageID")))),
		}

		for _, hash := range cachedChaincode.Hashes {
			channelCache.InterestingHashes[hash] = name
		}

		ok, err = c.Resources.Serializer.IsSerialized(NamespacesName, privateName, chaincodeDefinition.Parameters(), orgState)

		if err != nil {
			return errors.WithMessagef(err, "could not check opaque org state for '%s' on channel '%s'", name, channelID)
		}
		if !ok {
			logger.Debugf("Channel %s for chaincode definition %s:%s does not have our org's approval", channelID, name, chaincodeDefinition.EndorsementInfo.Version)
			continue
		}

		cachedChaincode.Approved = true

		isLocalPackage, err := c.Resources.Serializer.IsMetadataSerialized(ChaincodeSourcesName, privateName, &ChaincodeLocalPackage{}, orgState)
		if err != nil {
			return errors.WithMessagef(err, "could not check opaque org state for chaincode source for '%s' on channel '%s'", name, channelID)
		}

		if !isLocalPackage {
			logger.Debugf("Channel %s for chaincode definition %s:%s does not have a chaincode source defined", channelID, name, chaincodeDefinition.EndorsementInfo.Version)
			continue
		}

		hashKey := FieldKey(ChaincodeSourcesName, privateName, "PackageID")
		hashOfCCHash, err := orgState.GetStateHash(hashKey)
		if err != nil {
			return errors.WithMessagef(err, "could not check opaque org state for chaincode source hash for '%s' on channel '%s'", name, channelID)
		}

		localChaincode, ok := c.localChaincodes[string(hashOfCCHash)]
		if !ok {
			localChaincode = &LocalChaincode{
				References: map[string]map[string]*CachedChaincodeDefinition{},
			}
			c.localChaincodes[string(hashOfCCHash)] = localChaincode
		}

		cachedChaincode.InstallInfo = localChaincode.Info
		if localChaincode.Info != nil {
			logger.Infof("Chaincode with package ID '%s' now available on channel %s for chaincode definition %s:%s", localChaincode.Info.PackageID, channelID, name, cachedChaincode.Definition.EndorsementInfo.Version)
		} else {
			logger.Debugf("Chaincode definition for chaincode '%s' on channel '%s' is approved, but not installed", name, channelID)
		}

		channelReferences, ok := localChaincode.References[channelID]
		if !ok {
			channelReferences = map[string]*CachedChaincodeDefinition{}
			localChaincode.References[channelID] = channelReferences
		}

		channelReferences[name] = cachedChaincode
		if !initializing {
			c.eventBroker.ProcessApproveOrDefineEvent(channelID, name, cachedChaincode)
		}
	}

	if !initializing {
		c.handleMetadataUpdatesForChannel(channelID)
	}

	return nil
}


func (c *Cache) RegisterListener(channelID string, listener ledger.ChaincodeLifecycleEventListener) {
	c.eventBroker.RegisterListener(channelID, listener)
}

func (c *Cache) InitializeMetadata(channel string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ms, err := c.retrieveChaincodesMetadataSetWhileLocked(channel)
	if err != nil {
		logger.Warningf("no metadata found on channel '%s', err %s", channel, err)
		return
	}

	c.MetadataHandler.InitializeMetadata(channel, ms)
}

func (c *Cache) retrieveChaincodesMetadataSetWhileLocked(channelID string) (chaincode.MetadataSet, error) {
	channelChaincodes, ok := c.definedChaincodes[channelID]
	if !ok {
		return nil, errors.Errorf("unknown channel '%s'", channelID)
	}

	keys := make([]string, 0, len(channelChaincodes.Chaincodes))
	for name := range channelChaincodes.Chaincodes {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	metadataSet := chaincode.MetadataSet{}
	for _, name := range keys {
		def := channelChaincodes.Chaincodes[name]

		metadataSet = append(metadataSet,
			chaincode.Metadata{
				Name:              name,
				Version:           def.Definition.EndorsementInfo.Version,
				Policy:            def.Definition.ValidationInfo.ValidationParameter,
				CollectionsConfig: def.Definition.Collections,
				Approved:          def.Approved,
				Installed:         def.InstallInfo != nil,
			},
		)
	}

	return metadataSet, nil
}

func (c *Cache) handleMetadataUpdates(localChaincode *LocalChaincode) {
	for channelID := range localChaincode.References {
		c.handleMetadataUpdatesForChannel(channelID)
	}
}

func (c *Cache) handleMetadataUpdatesForChannel(channelID string) {
	ms, err := c.retrieveChaincodesMetadataSetWhileLocked(channelID)
	if err != nil {
		logger.Warningf("no metadata found on channel '%s': %s", channelID, err)
		return
	}

	c.MetadataHandler.UpdateMetadata(channelID, ms)
}
