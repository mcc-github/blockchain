/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privacyenabledstate

import (
	"encoding/base64"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/common/privdata"
	"github.com/mcc-github/blockchain/core/ledger/cceventmgmt"
	"github.com/mcc-github/blockchain/protos/common"

	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/statecouchdb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb/stateleveldb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("privacyenabledstate")

const (
	nsJoiner       = "$$"
	pvtDataPrefix  = "p"
	hashDataPrefix = "h"
)


type CommonStorageDBProvider struct {
	statedb.VersionedDBProvider
}


func NewCommonStorageDBProvider() (DBProvider, error) {
	var vdbProvider statedb.VersionedDBProvider
	var err error
	if ledgerconfig.IsCouchDBEnabled() {
		if vdbProvider, err = statecouchdb.NewVersionedDBProvider(); err != nil {
			return nil, err
		}
	} else {
		vdbProvider = stateleveldb.NewVersionedDBProvider()
	}
	return &CommonStorageDBProvider{vdbProvider}, nil
}


func (p *CommonStorageDBProvider) GetDBHandle(id string) (DB, error) {
	vdb, err := p.VersionedDBProvider.GetDBHandle(id)
	if err != nil {
		return nil, err
	}
	return NewCommonStorageDB(vdb, id)
}


func (p *CommonStorageDBProvider) Close() {
	p.VersionedDBProvider.Close()
}



type CommonStorageDB struct {
	statedb.VersionedDB
}



func NewCommonStorageDB(vdb statedb.VersionedDB, ledgerid string) (DB, error) {
	return &CommonStorageDB{VersionedDB: vdb}, nil
}


func (s *CommonStorageDB) IsBulkOptimizable() bool {
	_, ok := s.VersionedDB.(statedb.BulkOptimizable)
	return ok
}


func (s *CommonStorageDB) LoadCommittedVersionsOfPubAndHashedKeys(pubKeys []*statedb.CompositeKey,
	hashedKeys []*HashedCompositeKey) error {

	bulkOptimizable, ok := s.VersionedDB.(statedb.BulkOptimizable)
	if !ok {
		return nil
	}
	
	for _, key := range hashedKeys {
		ns := deriveHashedDataNs(key.Namespace, key.CollectionName)
		
		var keyHashStr string
		if !s.BytesKeySuppoted() {
			keyHashStr = base64.StdEncoding.EncodeToString([]byte(key.KeyHash))
		} else {
			keyHashStr = key.KeyHash
		}
		pubKeys = append(pubKeys, &statedb.CompositeKey{
			Namespace: ns,
			Key:       keyHashStr,
		})
	}

	err := bulkOptimizable.LoadCommittedVersions(pubKeys)
	if err != nil {
		return err
	}

	return nil
}


func (s *CommonStorageDB) ClearCachedVersions() {
	bulkOptimizable, ok := s.VersionedDB.(statedb.BulkOptimizable)
	if ok {
		bulkOptimizable.ClearCachedVersions()
	}
}


func (s *CommonStorageDB) GetChaincodeEventListener() cceventmgmt.ChaincodeLifecycleEventListener {
	_, ok := s.VersionedDB.(statedb.IndexCapable)
	if ok {
		return s
	}
	return nil
}


func (s *CommonStorageDB) GetPrivateData(namespace, collection, key string) (*statedb.VersionedValue, error) {
	return s.GetState(derivePvtDataNs(namespace, collection), key)
}


func (s *CommonStorageDB) GetValueHash(namespace, collection string, keyHash []byte) (*statedb.VersionedValue, error) {
	keyHashStr := string(keyHash)
	if !s.BytesKeySuppoted() {
		keyHashStr = base64.StdEncoding.EncodeToString(keyHash)
	}
	return s.GetState(deriveHashedDataNs(namespace, collection), keyHashStr)
}


func (s *CommonStorageDB) GetKeyHashVersion(namespace, collection string, keyHash []byte) (*version.Height, error) {
	keyHashStr := string(keyHash)
	if !s.BytesKeySuppoted() {
		keyHashStr = base64.StdEncoding.EncodeToString(keyHash)
	}
	return s.GetVersion(deriveHashedDataNs(namespace, collection), keyHashStr)
}


func (s *CommonStorageDB) GetCachedKeyHashVersion(namespace, collection string, keyHash []byte) (*version.Height, bool) {
	bulkOptimizable, ok := s.VersionedDB.(statedb.BulkOptimizable)
	if !ok {
		return nil, false
	}

	keyHashStr := string(keyHash)
	if !s.BytesKeySuppoted() {
		keyHashStr = base64.StdEncoding.EncodeToString(keyHash)
	}
	return bulkOptimizable.GetCachedVersion(deriveHashedDataNs(namespace, collection), keyHashStr)
}


func (s *CommonStorageDB) GetPrivateDataMultipleKeys(namespace, collection string, keys []string) ([]*statedb.VersionedValue, error) {
	return s.GetStateMultipleKeys(derivePvtDataNs(namespace, collection), keys)
}


func (s *CommonStorageDB) GetPrivateDataRangeScanIterator(namespace, collection, startKey, endKey string) (statedb.ResultsIterator, error) {
	return s.GetStateRangeScanIterator(derivePvtDataNs(namespace, collection), startKey, endKey)
}


func (s CommonStorageDB) ExecuteQueryOnPrivateData(namespace, collection, query string) (statedb.ResultsIterator, error) {
	return s.ExecuteQuery(derivePvtDataNs(namespace, collection), query)
}



func (s *CommonStorageDB) ApplyUpdates(batch *statedb.UpdateBatch, height *version.Height) error {
	return errors.New("this function should not be invoked on this type. Please invoke function ApplyPrivacyAwareUpdates")
}


func (s *CommonStorageDB) ApplyPrivacyAwareUpdates(updates *UpdateBatch, height *version.Height) error {
	
	combinedUpdates := updates.PubUpdates
	addPvtUpdates(combinedUpdates, updates.PvtUpdates)
	addHashedUpdates(combinedUpdates, updates.HashUpdates, !s.BytesKeySuppoted())
	return s.VersionedDB.ApplyUpdates(combinedUpdates.UpdateBatch, height)
}

func (s *CommonStorageDB) getCollectionConfigMap(chaincodeDefinition *cceventmgmt.ChaincodeDefinition) (map[string]bool, error) {
	var collectionConfigsBytes []byte
	collectionConfigsMap := make(map[string]bool)

	
	
	
	if chaincodeDefinition.CollectionConfigs != nil {
		
		
		collectionConfigsBytes = chaincodeDefinition.CollectionConfigs
	} else {
		
		
		
		lsccNamespace := "lscc"
		collectionConfigKey := privdata.BuildCollectionKVSKey(chaincodeDefinition.Name)

		versionedValue, err := s.VersionedDB.GetState(lsccNamespace, collectionConfigKey)
		if err != nil {
			return nil, err
		}
		
		
		if versionedValue != nil {
			collectionConfigsBytes = versionedValue.Value
		}
	}

	if collectionConfigsBytes != nil {
		collectionConfigs := &common.CollectionConfigPackage{}
		if err := proto.Unmarshal(collectionConfigsBytes, collectionConfigs); err != nil {
			return nil, err
		}

		for _, config := range collectionConfigs.Config {
			sConfig := config.GetStaticCollectionConfig()
			if sConfig == nil {
				continue
			}
			collectionConfigsMap[sConfig.Name] = true
		}
	}

	return collectionConfigsMap, nil
}







func (s *CommonStorageDB) HandleChaincodeDeploy(chaincodeDefinition *cceventmgmt.ChaincodeDefinition, dbArtifactsTar []byte) error {

	
	indexCapable, ok := s.VersionedDB.(statedb.IndexCapable)
	if !ok {
		return nil
	}

	if chaincodeDefinition == nil {
		return errors.New("chaincode definition not found while creating couchdb index")
	}

	dbArtifacts, err := ccprovider.ExtractFileEntries(dbArtifactsTar, indexCapable.GetDBType())
	if err != nil {
		logger.Errorf("Index creation: error extracting db artifacts from tar for chaincode [%s]: %s", chaincodeDefinition.Name, err)
		return nil
	}

	collectionConfigMap, err := s.getCollectionConfigMap(chaincodeDefinition)
	if err != nil {
		logger.Errorf("Error while retrieving collection config for chaincode=[%s]: %s",
			chaincodeDefinition.Name, err)
		return nil
	}

	for directoryPath, archiveDirectoryEntries := range dbArtifacts {
		
		directoryPathArray := strings.Split(directoryPath, "/")
		
		if directoryPathArray[3] == "indexes" {
			err := indexCapable.ProcessIndexesForChaincodeDeploy(chaincodeDefinition.Name, archiveDirectoryEntries)
			if err != nil {
				logger.Errorf("Error processing index for chaincode [%s]: %s", chaincodeDefinition.Name, err)
			}
			continue
		}
		
		if directoryPathArray[3] == "collections" && directoryPathArray[5] == "indexes" {
			collectionName := directoryPathArray[4]
			_, ok := collectionConfigMap[collectionName]
			if !ok {
				logger.Errorf("Error processing index for chaincode [%s]: cannot create an index for an undefined collection=[%s]", chaincodeDefinition.Name, collectionName)
			} else {

				err := indexCapable.ProcessIndexesForChaincodeDeploy(derivePvtDataNs(chaincodeDefinition.Name, collectionName),
					archiveDirectoryEntries)
				if err != nil {
					logger.Errorf("Error processing collection index for chaincode [%s]: %s", chaincodeDefinition.Name, err)
				}
			}
		}
	}
	return nil
}


func (s *CommonStorageDB) ChaincodeDeployDone(succeeded bool) {
	
}

func derivePvtDataNs(namespace, collection string) string {
	return namespace + nsJoiner + pvtDataPrefix + collection
}

func deriveHashedDataNs(namespace, collection string) string {
	return namespace + nsJoiner + hashDataPrefix + collection
}

func addPvtUpdates(pubUpdateBatch *PubUpdateBatch, pvtUpdateBatch *PvtUpdateBatch) {
	for ns, nsBatch := range pvtUpdateBatch.UpdateMap {
		for _, coll := range nsBatch.GetCollectionNames() {
			for key, vv := range nsBatch.GetUpdates(coll) {
				pubUpdateBatch.Update(derivePvtDataNs(ns, coll), key, vv)
			}
		}
	}
}

func addHashedUpdates(pubUpdateBatch *PubUpdateBatch, hashedUpdateBatch *HashedUpdateBatch, base64Key bool) {
	for ns, nsBatch := range hashedUpdateBatch.UpdateMap {
		for _, coll := range nsBatch.GetCollectionNames() {
			for key, vv := range nsBatch.GetUpdates(coll) {
				if base64Key {
					key = base64.StdEncoding.EncodeToString([]byte(key))
				}
				pubUpdateBatch.Update(deriveHashedDataNs(ns, coll), key, vv)
			}
		}
	}
}
