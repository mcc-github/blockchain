/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"errors"
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
)



type nsCommittersBuilder struct {
	updates         map[string]*statedb.VersionedValue
	db              *couchdb.CouchDatabase
	revisions       map[string]string
	subNsCommitters []batch
}


type subNsCommitter struct {
	db             *couchdb.CouchDatabase
	batchUpdateMap map[string]*batchableDocument
}


func (vdb *VersionedDB) buildCommitters(updates *statedb.UpdateBatch) ([]batch, error) {
	namespaces := updates.GetUpdatedNamespaces()
	var nsCommitterBuilder []batch
	for _, ns := range namespaces {
		nsUpdates := updates.GetUpdates(ns)
		db, err := vdb.getNamespaceDBHandle(ns)
		if err != nil {
			return nil, err
		}
		nsRevs := vdb.committedDataCache.revs[ns]
		if nsRevs == nil {
			nsRevs = make(nsRevisions)
		}
		
		
		nsCommitterBuilder = append(nsCommitterBuilder, &nsCommittersBuilder{updates: nsUpdates, db: db, revisions: nsRevs})
	}
	if err := executeBatches(nsCommitterBuilder); err != nil {
		return nil, err
	}
	
	var combinedSubNsCommitters []batch
	for _, b := range nsCommitterBuilder {
		combinedSubNsCommitters = append(combinedSubNsCommitters, b.(*nsCommittersBuilder).subNsCommitters...)
	}
	return combinedSubNsCommitters, nil
}



func (builder *nsCommittersBuilder) execute() error {
	if err := addRevisionsForMissingKeys(builder.revisions, builder.db, builder.updates); err != nil {
		return err
	}
	maxBacthSize := ledgerconfig.GetMaxBatchUpdateSize()
	batchUpdateMap := make(map[string]*batchableDocument)
	for key, vv := range builder.updates {
		couchDoc, err := keyValToCouchDoc(&keyValue{key: key, VersionedValue: vv}, builder.revisions[key])
		if err != nil {
			return err
		}
		batchUpdateMap[key] = &batchableDocument{CouchDoc: *couchDoc, Deleted: vv.Value == nil}
		if len(batchUpdateMap) == maxBacthSize {
			builder.subNsCommitters = append(builder.subNsCommitters, &subNsCommitter{builder.db, batchUpdateMap})
			batchUpdateMap = make(map[string]*batchableDocument)
		}
	}
	if len(batchUpdateMap) > 0 {
		builder.subNsCommitters = append(builder.subNsCommitters, &subNsCommitter{builder.db, batchUpdateMap})
	}
	return nil
}


func (committer *subNsCommitter) execute() error {
	return commitUpdates(committer.db, committer.batchUpdateMap)
}


func commitUpdates(db *couchdb.CouchDatabase, batchUpdateMap map[string]*batchableDocument) error {
	
	batchUpdateDocs := []*couchdb.CouchDoc{}
	for _, updateDocument := range batchUpdateMap {
		batchUpdateDocument := updateDocument
		batchUpdateDocs = append(batchUpdateDocs, &batchUpdateDocument.CouchDoc)
	}

	
	batchUpdateResp, err := db.BatchUpdateDocuments(batchUpdateDocs)
	if err != nil {
		return err
	}
	
	
	for _, respDoc := range batchUpdateResp {
		
		if respDoc.Ok != true {
			batchUpdateDocument := batchUpdateMap[respDoc.ID]
			var err error
			
			
			
			if batchUpdateDocument.CouchDoc.JSONValue != nil {
				err = removeJSONRevision(&batchUpdateDocument.CouchDoc.JSONValue)
				if err != nil {
					return err
				}
			}
			
			if batchUpdateDocument.Deleted {
				logger.Warningf("CouchDB batch document delete encountered an problem. Retrying delete for document ID:%s", respDoc.ID)
				
				
				
				err = db.DeleteDoc(respDoc.ID, "")
			} else {
				logger.Warningf("CouchDB batch document update encountered an problem. Retrying update for document ID:%s", respDoc.ID)
				
				
				_, err = db.SaveDoc(respDoc.ID, "", &batchUpdateDocument.CouchDoc)
			}

			
			if err != nil {
				errorString := fmt.Sprintf("Error occurred while saving document ID = %v  Error: %s  Reason: %s\n",
					respDoc.ID, respDoc.Error, respDoc.Reason)

				logger.Errorf(errorString)
				return fmt.Errorf(errorString)
			}
		}
	}
	return nil
}


type nsFlusher struct {
	db *couchdb.CouchDatabase
}

func (vdb *VersionedDB) ensureFullCommit(dbs []*couchdb.CouchDatabase) error {
	var flushers []batch
	for _, db := range dbs {
		flushers = append(flushers, &nsFlusher{db})
	}
	return executeBatches(flushers)
}

func (f *nsFlusher) execute() error {
	dbResponse, err := f.db.EnsureFullCommit()
	if err != nil || dbResponse.Ok != true {
		logger.Errorf("Failed to perform full commit\n")
		return errors.New("Failed to perform full commit")
	}
	return nil
}

func addRevisionsForMissingKeys(revisions map[string]string, db *couchdb.CouchDatabase, nsUpdates map[string]*statedb.VersionedValue) error {
	var missingKeys []string
	for key := range nsUpdates {
		_, ok := revisions[key]
		if !ok {
			missingKeys = append(missingKeys, key)
		}
	}
	logger.Debugf("Pulling revisions for the [%d] keys for namsespace [%s] that were not part of the readset", len(missingKeys), db.DBName)
	retrievedMetadata, err := retrieveNsMetadata(db, missingKeys)
	if err != nil {
		return err
	}
	for _, metadata := range retrievedMetadata {
		revisions[metadata.ID] = metadata.Rev
	}
	return nil
}


type batchableDocument struct {
	CouchDoc couchdb.CouchDoc
	Deleted  bool
}
