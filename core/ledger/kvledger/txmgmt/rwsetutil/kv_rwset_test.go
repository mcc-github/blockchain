/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rwsetutil

import (
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	"github.com/stretchr/testify/assert"
)

const kvrwsetV1ProtoBytesFile = "testdata/kvrwsetV1ProtoBytes"




func TestKVRWSetV1BackwardCompatible(t *testing.T) {
	protoBytes, err := ioutil.ReadFile(kvrwsetV1ProtoBytesFile)
	assert.NoError(t, err)
	kvrwset1 := &kvrwset.KVRWSet{}
	assert.NoError(t, proto.Unmarshal(protoBytes, kvrwset1))
	kvrwset2 := constructSampleKVRWSet()
	t.Logf("kvrwset1=%s, kvrwset2=%s", spew.Sdump(kvrwset1), spew.Sdump(kvrwset2))
	assert.Equal(t, kvrwset2, kvrwset1)
}





func PrepareBinaryFileSampleKVRWSetV1(t *testing.T) {
	b, err := proto.Marshal(constructSampleKVRWSet())
	assert.NoError(t, err)
	assert.NoError(t, ioutil.WriteFile(kvrwsetV1ProtoBytesFile, b, 0644))
}

func constructSampleKVRWSet() *kvrwset.KVRWSet {
	rqi1 := &kvrwset.RangeQueryInfo{StartKey: "k0", EndKey: "k9", ItrExhausted: true}
	SetRawReads(rqi1, []*kvrwset.KVRead{
		{Key: "k1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}},
		{Key: "k2", Version: &kvrwset.Version{BlockNum: 1, TxNum: 2}},
	})

	rqi2 := &kvrwset.RangeQueryInfo{StartKey: "k00", EndKey: "k90", ItrExhausted: true}
	SetMerkelSummary(rqi2, &kvrwset.QueryReadsMerkleSummary{MaxDegree: 5, MaxLevel: 4, MaxLevelHashes: [][]byte{[]byte("Hash-1"), []byte("Hash-2")}})
	return &kvrwset.KVRWSet{
		Reads:            []*kvrwset.KVRead{{Key: "key1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}}},
		RangeQueriesInfo: []*kvrwset.RangeQueryInfo{rqi1, rqi2},
		Writes:           []*kvrwset.KVWrite{{Key: "key2", IsDelete: false, Value: []byte("value2")}},
	}
}
