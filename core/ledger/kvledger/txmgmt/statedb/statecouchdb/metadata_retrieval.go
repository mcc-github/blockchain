/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package statecouchdb

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/mcc-github/blockchain/core/ledger/util/couchdb"
)



type nsMetadataRetriever struct {
	ns              string
	db              *couchdb.CouchDatabase
	keys            []string
	executionResult []*couchdb.DocMetadata
}




type subNsMetadataRetriever nsMetadataRetriever


func (vdb *VersionedDB) retrieveMetadata(nsKeysMap map[string][]string) (map[string][]*couchdb.DocMetadata, error) {
	
	nsMetadataRetrievers := []batch{}
	for ns, keys := range nsKeysMap {
		db, err := vdb.getNamespaceDBHandle(ns)
		if err != nil {
			return nil, err
		}
		nsMetadataRetrievers = append(nsMetadataRetrievers, &nsMetadataRetriever{ns: ns, db: db, keys: keys})
	}
	if err := executeBatches(nsMetadataRetrievers); err != nil {
		return nil, err
	}
	
	executionResults := make(map[string][]*couchdb.DocMetadata)
	for _, r := range nsMetadataRetrievers {
		nsMetadataRetriever := r.(*nsMetadataRetriever)
		executionResults[nsMetadataRetriever.ns] = nsMetadataRetriever.executionResult
	}
	return executionResults, nil
}


func retrieveNsMetadata(db *couchdb.CouchDatabase, keys []string) ([]*couchdb.DocMetadata, error) {
	
	maxBacthSize := ledgerconfig.GetMaxBatchUpdateSize()
	batches := []batch{}
	remainingKeys := keys
	for {
		numKeys := minimum(maxBacthSize, len(remainingKeys))
		if numKeys == 0 {
			break
		}
		batch := &subNsMetadataRetriever{db: db, keys: remainingKeys[:numKeys]}
		batches = append(batches, batch)
		remainingKeys = remainingKeys[numKeys:]
	}
	if err := executeBatches(batches); err != nil {
		return nil, err
	}
	
	var executionResults []*couchdb.DocMetadata
	for _, b := range batches {
		executionResults = append(executionResults, b.(*subNsMetadataRetriever).executionResult...)
	}
	return executionResults, nil
}

func (r *nsMetadataRetriever) execute() error {
	var err error
	if r.executionResult, err = retrieveNsMetadata(r.db, r.keys); err != nil {
		return err
	}
	return nil
}

func (r *nsMetadataRetriever) String() string {
	return fmt.Sprintf("nsMetadataRetriever:ns=%s, num keys=%d", r.ns, len(r.keys))
}

func (b *subNsMetadataRetriever) execute() error {
	var err error
	if b.executionResult, err = b.db.BatchRetrieveDocumentMetadata(b.keys); err != nil {
		return err
	}
	return nil
}

func (b *subNsMetadataRetriever) String() string {
	return fmt.Sprintf("subNsMetadataRetriever:ns=%s, num keys=%d", b.ns, len(b.keys))
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
