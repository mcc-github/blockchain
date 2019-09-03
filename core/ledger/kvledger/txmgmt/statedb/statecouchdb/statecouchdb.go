/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("statecouchdb")


const lsccCacheSize = 50


type VersionedDBProvider struct {
	couchInstance      *couchdb.CouchInstance
	databases          map[string]*VersionedDB
	mux                sync.Mutex
	openCounts         uint64
	redoLoggerProvider *redoLoggerProvider
}


func NewVersionedDBProvider(config *couchdb.Config, metricsProvider metrics.Provider) (*VersionedDBProvider, error) {
	logger.Debugf("constructing CouchDB VersionedDBProvider")
	couchInstance, err := couchdb.CreateCouchInstance(config, metricsProvider)
	if err != nil {
		return nil, err
	}
	p, err := newRedoLoggerProvider(config.RedoLogPath)
	if err != nil {
		return nil, err
	}
	return &VersionedDBProvider{
			couchInstance:      couchInstance,
			databases:          make(map[string]*VersionedDB),
			mux:                sync.Mutex{},
			openCounts:         0,
			redoLoggerProvider: p,
		},
		nil
}


func (provider *VersionedDBProvider) GetDBHandle(dbName string) (statedb.VersionedDB, error) {
	provider.mux.Lock()
	defer provider.mux.Unlock()
	vdb := provider.databases[dbName]
	if vdb == nil {
		var err error
		vdb, err = newVersionedDB(
			provider.couchInstance,
			provider.redoLoggerProvider.newRedoLogger(dbName),
			dbName,
		)
		if err != nil {
			return nil, err
		}
		provider.databases[dbName] = vdb
	}
	return vdb, nil
}


func (provider *VersionedDBProvider) Close() {
	
	provider.redoLoggerProvider.close()
}


func (provider *VersionedDBProvider) HealthCheck(ctx context.Context) error {
	return provider.couchInstance.HealthCheck(ctx)
}


type VersionedDB struct {
	couchInstance      *couchdb.CouchInstance
	metadataDB         *couchdb.CouchDatabase            
	chainName          string                            
	namespaceDBs       map[string]*couchdb.CouchDatabase 
	committedDataCache *versionsCache                    
	verCacheLock       sync.RWMutex
	mux                sync.RWMutex
	lsccStateCache     *lsccStateCache
	redoLogger         *redoLogger
}

type lsccStateCache struct {
	cache   map[string]*statedb.VersionedValue
	rwMutex sync.RWMutex
}

func (l *lsccStateCache) getState(key string) *statedb.VersionedValue {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()

	if versionedValue, ok := l.cache[key]; ok {
		logger.Debugf("key:[%s] found in the lsccStateCache", key)
		return versionedValue
	}
	return nil
}

func (l *lsccStateCache) updateState(key string, value *statedb.VersionedValue) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	if _, ok := l.cache[key]; ok {
		logger.Debugf("key:[%s] is updated in lsccStateCache", key)
		l.cache[key] = value
	}
}

func (l *lsccStateCache) setState(key string, value *statedb.VersionedValue) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	if l.isCacheFull() {
		l.evictARandomEntry()
	}

	logger.Debugf("key:[%s] is stoed in lsccStateCache", key)
	l.cache[key] = value
}

func (l *lsccStateCache) isCacheFull() bool {
	return len(l.cache) == lsccCacheSize
}

func (l *lsccStateCache) evictARandomEntry() {
	for key := range l.cache {
		delete(l.cache, key)
		return
	}
}


func newVersionedDB(couchInstance *couchdb.CouchInstance, redoLogger *redoLogger, dbName string) (*VersionedDB, error) {
	
	chainName := dbName
	dbName = couchdb.ConstructMetadataDBName(dbName)

	metadataDB, err := couchdb.CreateCouchDatabase(couchInstance, dbName)
	if err != nil {
		return nil, err
	}
	namespaceDBMap := make(map[string]*couchdb.CouchDatabase)
	vdb := &VersionedDB{
		couchInstance:      couchInstance,
		metadataDB:         metadataDB,
		chainName:          chainName,
		namespaceDBs:       namespaceDBMap,
		committedDataCache: newVersionCache(),
		lsccStateCache: &lsccStateCache{
			cache: make(map[string]*statedb.VersionedValue),
		},
		redoLogger: redoLogger,
	}
	logger.Debugf("chain [%s]: checking for redolog record", chainName)
	redologRecord, err := redoLogger.load()
	if err != nil {
		return nil, err
	}
	savepoint, err := vdb.GetLatestSavePoint()
	if err != nil {
		return nil, err
	}

	
	
	
	if redologRecord == nil || savepoint == nil {
		logger.Debugf("chain [%s]: No redo-record or save point present", chainName)
		return vdb, nil
	}

	logger.Debugf("chain [%s]: save point = %#v, version of redolog record = %#v",
		chainName, savepoint, redologRecord.Version)

	if redologRecord.Version.BlockNum-savepoint.BlockNum == 1 {
		logger.Debugf("chain [%s]: Re-applying last batch", chainName)
		if err := vdb.applyUpdates(redologRecord.UpdateBatch, redologRecord.Version); err != nil {
			return nil, err
		}
	}
	return vdb, nil
}


func (vdb *VersionedDB) getNamespaceDBHandle(namespace string) (*couchdb.CouchDatabase, error) {
	vdb.mux.RLock()
	db := vdb.namespaceDBs[namespace]
	vdb.mux.RUnlock()
	if db != nil {
		return db, nil
	}
	namespaceDBName := couchdb.ConstructNamespaceDBName(vdb.chainName, namespace)
	vdb.mux.Lock()
	defer vdb.mux.Unlock()
	db = vdb.namespaceDBs[namespace]
	if db == nil {
		var err error
		db, err = couchdb.CreateCouchDatabase(vdb.couchInstance, namespaceDBName)
		if err != nil {
			return nil, err
		}
		vdb.namespaceDBs[namespace] = db
	}
	return db, nil
}


func (vdb *VersionedDB) ProcessIndexesForChaincodeDeploy(namespace string, fileEntries []*ccprovider.TarFileEntry) error {
	db, err := vdb.getNamespaceDBHandle(namespace)
	if err != nil {
		return err
	}
	for _, fileEntry := range fileEntries {
		indexData := fileEntry.FileContent
		filename := fileEntry.FileHeader.Name
		_, err = db.CreateIndex(string(indexData))
		if err != nil {
			return errors.WithMessagef(err, "error creating index from file [%s] for channel [%s]", filename, namespace)
		}
	}
	return nil
}


func (vdb *VersionedDB) GetDBType() string {
	return "couchdb"
}





func (vdb *VersionedDB) LoadCommittedVersions(keys []*statedb.CompositeKey) error {
	nsKeysMap := map[string][]string{}
	committedDataCache := newVersionCache()
	for _, compositeKey := range keys {
		ns, key := compositeKey.Namespace, compositeKey.Key
		committedDataCache.setVerAndRev(ns, key, nil, "")
		logger.Debugf("Load into version cache: %s~%s", ns, key)
		nsKeysMap[ns] = append(nsKeysMap[ns], key)
	}
	nsMetadataMap, err := vdb.retrieveMetadata(nsKeysMap)
	logger.Debugf("nsKeysMap=%s", nsKeysMap)
	logger.Debugf("nsMetadataMap=%s", nsMetadataMap)
	if err != nil {
		return err
	}
	for ns, nsMetadata := range nsMetadataMap {
		for _, keyMetadata := range nsMetadata {
			
			if len(keyMetadata.Version) != 0 {
				version, _, err := decodeVersionAndMetadata(keyMetadata.Version)
				if err != nil {
					return err
				}
				committedDataCache.setVerAndRev(ns, keyMetadata.ID, version, keyMetadata.Rev)
			}
		}
	}
	vdb.verCacheLock.Lock()
	defer vdb.verCacheLock.Unlock()
	vdb.committedDataCache = committedDataCache
	return nil
}


func (vdb *VersionedDB) GetVersion(namespace string, key string) (*version.Height, error) {
	returnVersion, keyFound := vdb.GetCachedVersion(namespace, key)
	if !keyFound {
		
		
		vv, err := vdb.GetState(namespace, key)
		if err != nil || vv == nil {
			return nil, err
		}
		returnVersion = vv.Version
	}
	return returnVersion, nil
}


func (vdb *VersionedDB) GetCachedVersion(namespace string, key string) (*version.Height, bool) {
	logger.Debugf("Retrieving cached version: %s~%s", key, namespace)
	vdb.verCacheLock.RLock()
	defer vdb.verCacheLock.RUnlock()
	return vdb.committedDataCache.getVersion(namespace, key)
}


func (vdb *VersionedDB) ValidateKeyValue(key string, value []byte) error {
	err := validateKey(key)
	if err != nil {
		return err
	}
	return validateValue(value)
}


func (vdb *VersionedDB) BytesKeySupported() bool {
	return false
}


func (vdb *VersionedDB) GetState(namespace string, key string) (*statedb.VersionedValue, error) {
	logger.Debugf("GetState(). ns=%s, key=%s", namespace, key)
	if namespace == "lscc" {
		if value := vdb.lsccStateCache.getState(key); value != nil {
			return value, nil
		}
	}

	db, err := vdb.getNamespaceDBHandle(namespace)
	if err != nil {
		return nil, err
	}
	couchDoc, _, err := db.ReadDoc(key)
	if err != nil {
		return nil, err
	}
	if couchDoc == nil {
		return nil, nil
	}
	kv, err := couchDocToKeyValue(couchDoc)
	if err != nil {
		return nil, err
	}

	if namespace == "lscc" {
		vdb.lsccStateCache.setState(key, kv.VersionedValue)
	}

	return kv.VersionedValue, nil
}


func (vdb *VersionedDB) GetStateMultipleKeys(namespace string, keys []string) ([]*statedb.VersionedValue, error) {
	vals := make([]*statedb.VersionedValue, len(keys))
	for i, key := range keys {
		val, err := vdb.GetState(namespace, key)
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return vals, nil
}




func (vdb *VersionedDB) GetStateRangeScanIterator(namespace string, startKey string, endKey string) (statedb.ResultsIterator, error) {
	return vdb.GetStateRangeScanIteratorWithMetadata(namespace, startKey, endKey, nil)
}

const optionBookmark = "bookmark"
const optionLimit = "limit"





func (vdb *VersionedDB) GetStateRangeScanIteratorWithMetadata(namespace string, startKey string, endKey string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	logger.Debugf("Entering GetStateRangeScanIteratorWithMetadata  namespace: %s  startKey: %s  endKey: %s  metadata: %v", namespace, startKey, endKey, metadata)
	
	internalQueryLimit := vdb.couchInstance.InternalQueryLimit()
	requestedLimit := int32(0)
	
	if metadata != nil {
		
		err := statedb.ValidateRangeMetadata(metadata)
		if err != nil {
			return nil, err
		}
		if limitOption, ok := metadata[optionLimit]; ok {
			requestedLimit = limitOption.(int32)
		}
	}
	db, err := vdb.getNamespaceDBHandle(namespace)
	if err != nil {
		return nil, err
	}
	return newQueryScanner(namespace, db, "", internalQueryLimit, requestedLimit, "", startKey, endKey)
}

func (scanner *queryScanner) getNextStateRangeScanResults() error {
	queryLimit := scanner.queryDefinition.internalQueryLimit
	if scanner.paginationInfo.requestedLimit > 0 {
		moreResultsNeeded := scanner.paginationInfo.requestedLimit - scanner.resultsInfo.totalRecordsReturned
		if moreResultsNeeded < scanner.queryDefinition.internalQueryLimit {
			queryLimit = moreResultsNeeded
		}
	}
	queryResult, nextStartKey, err := rangeScanFilterCouchInternalDocs(scanner.db,
		scanner.queryDefinition.startKey, scanner.queryDefinition.endKey, queryLimit)
	if err != nil {
		return err
	}
	scanner.resultsInfo.results = queryResult
	scanner.queryDefinition.startKey = nextStartKey
	scanner.paginationInfo.cursor = 0
	return nil
}

func rangeScanFilterCouchInternalDocs(db *couchdb.CouchDatabase,
	startKey, endKey string, queryLimit int32,
) ([]*couchdb.QueryResult, string, error) {
	var finalResults []*couchdb.QueryResult
	var finalNextStartKey string
	for {
		results, nextStartKey, err := db.ReadDocRange(startKey, endKey, queryLimit)
		if err != nil {
			logger.Debugf("Error calling ReadDocRange(): %s\n", err.Error())
			return nil, "", err
		}
		var filteredResults []*couchdb.QueryResult
		for _, doc := range results {
			if !isCouchInternalKey(doc.ID) {
				filteredResults = append(filteredResults, doc)
			}
		}

		finalResults = append(finalResults, filteredResults...)
		finalNextStartKey = nextStartKey
		queryLimit = int32(len(results) - len(filteredResults))
		if queryLimit == 0 || finalNextStartKey == "" {
			break
		}
		startKey = finalNextStartKey
	}
	var err error
	for i := 0; isCouchInternalKey(finalNextStartKey); i++ {
		_, finalNextStartKey, err = db.ReadDocRange(finalNextStartKey, endKey, 1)
		logger.Debugf("i=%d, finalNextStartKey=%s", i, finalNextStartKey)
		if err != nil {
			return nil, "", err
		}
	}
	return finalResults, finalNextStartKey, nil
}

func isCouchInternalKey(key string) bool {
	return len(key) != 0 && key[0] == '_'
}


func (vdb *VersionedDB) ExecuteQuery(namespace, query string) (statedb.ResultsIterator, error) {
	queryResult, err := vdb.ExecuteQueryWithMetadata(namespace, query, nil)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}


func (vdb *VersionedDB) ExecuteQueryWithMetadata(namespace, query string, metadata map[string]interface{}) (statedb.QueryResultsIterator, error) {
	logger.Debugf("Entering ExecuteQueryWithMetadata  namespace: %s,  query: %s,  metadata: %v", namespace, query, metadata)
	
	internalQueryLimit := vdb.couchInstance.InternalQueryLimit()
	bookmark := ""
	requestedLimit := int32(0)
	
	if metadata != nil {
		err := validateQueryMetadata(metadata)
		if err != nil {
			return nil, err
		}
		if limitOption, ok := metadata[optionLimit]; ok {
			requestedLimit = limitOption.(int32)
		}
		if bookmarkOption, ok := metadata[optionBookmark]; ok {
			bookmark = bookmarkOption.(string)
		}
	}
	queryString, err := applyAdditionalQueryOptions(query, internalQueryLimit, bookmark)
	if err != nil {
		logger.Errorf("Error calling applyAdditionalQueryOptions(): %s", err.Error())
		return nil, err
	}
	db, err := vdb.getNamespaceDBHandle(namespace)
	if err != nil {
		return nil, err
	}
	return newQueryScanner(namespace, db, queryString, internalQueryLimit, requestedLimit, bookmark, "", "")
}



func (scanner *queryScanner) executeQueryWithBookmark() error {
	queryLimit := scanner.queryDefinition.internalQueryLimit
	if scanner.paginationInfo.requestedLimit > 0 {
		if scanner.paginationInfo.requestedLimit-scanner.resultsInfo.totalRecordsReturned < scanner.queryDefinition.internalQueryLimit {
			queryLimit = scanner.paginationInfo.requestedLimit - scanner.resultsInfo.totalRecordsReturned
		}
	}
	queryString, err := applyAdditionalQueryOptions(scanner.queryDefinition.query,
		queryLimit, scanner.paginationInfo.bookmark)
	if err != nil {
		logger.Debugf("Error calling applyAdditionalQueryOptions(): %s\n", err.Error())
		return err
	}
	queryResult, bookmark, err := scanner.db.QueryDocuments(queryString)
	if err != nil {
		logger.Debugf("Error calling QueryDocuments(): %s\n", err.Error())
		return err
	}
	scanner.resultsInfo.results = queryResult
	scanner.paginationInfo.bookmark = bookmark
	scanner.paginationInfo.cursor = 0
	return nil
}

func validateQueryMetadata(metadata map[string]interface{}) error {
	for key, keyVal := range metadata {
		switch key {
		case optionBookmark:
			
			if _, ok := keyVal.(string); ok {
				continue
			}
			return fmt.Errorf("Invalid entry, \"bookmark\" must be a string")

		case optionLimit:
			
			if _, ok := keyVal.(int32); ok {
				continue
			}
			return fmt.Errorf("Invalid entry, \"limit\" must be an int32")

		default:
			return fmt.Errorf("Invalid entry, option %s not recognized", key)
		}
	}
	return nil
}


func (vdb *VersionedDB) ApplyUpdates(updates *statedb.UpdateBatch, height *version.Height) error {
	if height != nil && updates.ContainsPostOrderWrites {
		
		r := &redoRecord{
			UpdateBatch: updates,
			Version:     height,
		}
		if err := vdb.redoLogger.persist(r); err != nil {
			return err
		}
	}
	return vdb.applyUpdates(updates, height)
}

func (vdb *VersionedDB) applyUpdates(updates *statedb.UpdateBatch, height *version.Height) error {
	
	
	

	
	
	var updateBatches []batch
	var err error
	if updateBatches, err = vdb.buildCommitters(updates); err != nil {
		return err
	}
	
	if err = executeBatches(updateBatches); err != nil {
		return err
	}

	
	namespaces := updates.GetUpdatedNamespaces()
	
	if err = vdb.ensureFullCommitAndRecordSavepoint(height, namespaces); err != nil {
		logger.Errorf("Error during recordSavepoint: %s", err.Error())
		return err
	}

	lsccUpdates := updates.GetUpdates("lscc")
	for key, value := range lsccUpdates {
		vdb.lsccStateCache.updateState(key, value)
	}
	return nil
}


func (vdb *VersionedDB) ClearCachedVersions() {
	logger.Debugf("Clear Cache")
	vdb.verCacheLock.Lock()
	defer vdb.verCacheLock.Unlock()
	vdb.committedDataCache = newVersionCache()
}


func (vdb *VersionedDB) Open() error {
	
	return nil
}


func (vdb *VersionedDB) Close() {
	
}


const savepointDocID = "statedb_savepoint"







func (vdb *VersionedDB) ensureFullCommitAndRecordSavepoint(height *version.Height, namespaces []string) error {
	
	
	var dbs []*couchdb.CouchDatabase
	for _, ns := range namespaces {
		db, err := vdb.getNamespaceDBHandle(ns)
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}
	if err := vdb.ensureFullCommit(dbs); err != nil {
		return err
	}

	
	
	
	if height == nil {
		return nil
	}

	
	savepointCouchDoc, err := encodeSavepoint(height)
	if err != nil {
		return err
	}
	_, err = vdb.metadataDB.SaveDoc(savepointDocID, "", savepointCouchDoc)
	if err != nil {
		logger.Errorf("Failed to save the savepoint to DB %s", err.Error())
		return err
	}
	
	
	
	
	return nil
}


func (vdb *VersionedDB) GetLatestSavePoint() (*version.Height, error) {
	var err error
	couchDoc, _, err := vdb.metadataDB.ReadDoc(savepointDocID)
	if err != nil {
		logger.Errorf("Failed to read savepoint data %s", err.Error())
		return nil, err
	}
	
	if couchDoc == nil || couchDoc.JSONValue == nil {
		return nil, nil
	}
	return decodeSavepoint(couchDoc)
}


func applyAdditionalQueryOptions(queryString string, queryLimit int32, queryBookmark string) (string, error) {
	const jsonQueryFields = "fields"
	const jsonQueryLimit = "limit"
	const jsonQueryBookmark = "bookmark"
	
	jsonQueryMap := make(map[string]interface{})
	
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(queryString)))
	decoder.UseNumber()
	err := decoder.Decode(&jsonQueryMap)
	if err != nil {
		return "", err
	}
	if fieldsJSONArray, ok := jsonQueryMap[jsonQueryFields]; ok {
		switch fieldsJSONArray.(type) {
		case []interface{}:
			
			jsonQueryMap[jsonQueryFields] = append(fieldsJSONArray.([]interface{}),
				idField, versionField)
		default:
			return "", errors.New("fields definition must be an array")
		}
	}
	
	
	
	jsonQueryMap[jsonQueryLimit] = queryLimit
	
	if queryBookmark != "" {
		jsonQueryMap[jsonQueryBookmark] = queryBookmark
	}
	
	editedQuery, err := json.Marshal(jsonQueryMap)
	if err != nil {
		return "", err
	}
	logger.Debugf("Rewritten query: %s", editedQuery)
	return string(editedQuery), nil
}

type queryScanner struct {
	namespace       string
	db              *couchdb.CouchDatabase
	queryDefinition *queryDefinition
	paginationInfo  *paginationInfo
	resultsInfo     *resultsInfo
}

type queryDefinition struct {
	startKey           string
	endKey             string
	query              string
	internalQueryLimit int32
}

type paginationInfo struct {
	cursor         int32
	requestedLimit int32
	bookmark       string
}

type resultsInfo struct {
	totalRecordsReturned int32
	results              []*couchdb.QueryResult
}

func newQueryScanner(namespace string, db *couchdb.CouchDatabase, query string, internalQueryLimit,
	limit int32, bookmark, startKey, endKey string) (*queryScanner, error) {
	scanner := &queryScanner{namespace, db, &queryDefinition{startKey, endKey, query, internalQueryLimit}, &paginationInfo{-1, limit, bookmark}, &resultsInfo{0, nil}}
	var err error
	
	if scanner.queryDefinition.query != "" {
		err = scanner.executeQueryWithBookmark()
	} else {
		err = scanner.getNextStateRangeScanResults()
	}
	if err != nil {
		return nil, err
	}
	scanner.paginationInfo.cursor = -1
	return scanner, nil
}

func (scanner *queryScanner) Next() (statedb.QueryResult, error) {
	
	if len(scanner.resultsInfo.results) == 0 {
		return nil, nil
	}
	
	scanner.paginationInfo.cursor++
	
	
	if scanner.paginationInfo.cursor >= scanner.queryDefinition.internalQueryLimit {
		var err error
		
		if scanner.queryDefinition.query != "" {
			err = scanner.executeQueryWithBookmark()
		} else {
			err = scanner.getNextStateRangeScanResults()
		}
		if err != nil {
			return nil, err
		}
		
		if len(scanner.resultsInfo.results) == 0 {
			return nil, nil
		}
	}
	
	if scanner.paginationInfo.cursor >= int32(len(scanner.resultsInfo.results)) {
		return nil, nil
	}
	selectedResultRecord := scanner.resultsInfo.results[scanner.paginationInfo.cursor]
	key := selectedResultRecord.ID
	
	kv, err := couchDocToKeyValue(&couchdb.CouchDoc{JSONValue: selectedResultRecord.Value, Attachments: selectedResultRecord.Attachments})
	if err != nil {
		return nil, err
	}
	scanner.resultsInfo.totalRecordsReturned++
	return &statedb.VersionedKV{
		CompositeKey:   statedb.CompositeKey{Namespace: scanner.namespace, Key: key},
		VersionedValue: *kv.VersionedValue}, nil
}

func (scanner *queryScanner) Close() {
	scanner = nil
}

func (scanner *queryScanner) GetBookmarkAndClose() string {
	retval := ""
	if scanner.queryDefinition.query != "" {
		retval = scanner.paginationInfo.bookmark
	} else {
		retval = scanner.queryDefinition.startKey
	}
	scanner.Close()
	return retval
}
