/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package privacyenabledstate

import (
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/statedb"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/pkg/errors"
)

type UpdatesBytesBuilder struct {
}











func (bb *UpdatesBytesBuilder) DeterministicBytesForPubAndHashUpdates(u *UpdateBatch) ([]byte, error) {
	pubUpdates := u.PubUpdates
	hashUpdates := u.HashUpdates.UpdateMap
	namespaces := dedupAndSort(
		pubUpdates.GetUpdatedNamespaces(),
		hashUpdates.getUpdatedNamespaces(),
	)

	kvWritesProto := []*KVWriteProto{}
	for _, ns := range namespaces {
		if ns == "" {
			
			
			
			continue
		}
		p := bb.buildForKeys(pubUpdates.GetUpdates(ns))
		collsForNs, ok := hashUpdates[ns]
		if ok {
			p = append(p, bb.buildForColls(collsForNs)...)
		}
		p[0].Namespace = ns
		kvWritesProto = append(kvWritesProto, p...)
	}
	batchProto := &KVWritesBatchProto{
		Kvwrites: kvWritesProto,
	}

	batchBytes, err := proto.Marshal(batchProto)
	return batchBytes, errors.Wrap(err, "error constructing deterministic bytes from update batch")
}

func (bb *UpdatesBytesBuilder) buildForColls(colls nsBatch) []*KVWriteProto {
	collNames := colls.GetCollectionNames()
	sort.Strings(collNames)
	collsProto := []*KVWriteProto{}
	for _, collName := range collNames {
		collUpdates := colls.getCollectionUpdates(collName)
		p := bb.buildForKeys(collUpdates)
		p[0].Collection = collName
		collsProto = append(collsProto, p...)
	}
	return collsProto
}

func (bb *UpdatesBytesBuilder) buildForKeys(kv map[string]*statedb.VersionedValue) []*KVWriteProto {
	keys := util.GetSortedKeys(kv)
	p := []*KVWriteProto{}
	for _, key := range keys {
		val := kv[key]
		p = append(
			p,
			&KVWriteProto{
				Key:          []byte(key),
				Value:        val.Value,
				IsDelete:     val.Value == nil,
				VersionBytes: val.Version.ToBytes(),
			},
		)
	}
	return p
}

func dedupAndSort(a, b []string) []string {
	m := map[string]struct{}{}
	for _, entry := range a {
		m[entry] = struct{}{}
	}
	for _, entry := range b {
		m[entry] = struct{}{}
	}
	return util.GetSortedKeys(m)
}
