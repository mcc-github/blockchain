/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package storageutil

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
)


func SerializeMetadata(metadataEntries []*kvrwset.KVMetadataEntry) ([]byte, error) {
	metadata := &kvrwset.KVMetadataWrite{Entries: metadataEntries}
	return proto.Marshal(metadata)
}


func SerializeMetadataByMap(metadataMap map[string][]byte) ([]byte, error) {
	if metadataMap == nil {
		return nil, nil
	}
	names := util.GetSortedKeys(metadataMap)
	metadataEntries := []*kvrwset.KVMetadataEntry{}
	for _, k := range names {
		metadataEntries = append(metadataEntries, &kvrwset.KVMetadataEntry{Name: k, Value: metadataMap[k]})
	}
	return SerializeMetadata(metadataEntries)
}


func DeserializeMetadata(metadataBytes []byte) (map[string][]byte, error) {
	if metadataBytes == nil {
		return nil, nil
	}
	metadata := &kvrwset.KVMetadataWrite{}
	if err := proto.Unmarshal(metadataBytes, metadata); err != nil {
		return nil, err
	}
	m := make(map[string][]byte, len(metadata.Entries))
	for _, metadataEntry := range metadata.Entries {
		m[metadataEntry.Name] = metadataEntry.Value
	}
	return m, nil
}
