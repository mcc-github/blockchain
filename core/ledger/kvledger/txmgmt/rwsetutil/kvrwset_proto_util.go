/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rwsetutil

import "github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"


func SetRawReads(rqi *kvrwset.RangeQueryInfo, kvReads []*kvrwset.KVRead) {
	rqi.ReadsInfo = &kvrwset.RangeQueryInfo_RawReads{
		RawReads: &kvrwset.QueryReads{
			KvReads: kvReads,
		},
	}
}


func SetMerkelSummary(rqi *kvrwset.RangeQueryInfo, merkleSummary *kvrwset.QueryReadsMerkleSummary) {
	rqi.ReadsInfo = &kvrwset.RangeQueryInfo_ReadsMerkleHashes{ReadsMerkleHashes: merkleSummary}
}
