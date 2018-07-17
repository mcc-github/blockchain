/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kvrwset

import (
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/protos/ledger/rwset/kvrwset"
)

const (
	binaryTestFileName = "kvrwsetV1ProtoBytes"
)




func TestKVRWSetV1BackwardCompatible(t *testing.T) {
	protoBytes, err := ioutil.ReadFile(binaryTestFileName)
	testutil.AssertNoError(t, err, "")
	kvrwset1 := &kvrwset.KVRWSet{}
	testutil.AssertNoError(t, proto.Unmarshal(protoBytes, kvrwset1), "")
	kvrwset2 := constructSampleKVRWSet()
	t.Logf("kvrwset1=%s, kvrwset2=%s", spew.Sdump(kvrwset1), spew.Sdump(kvrwset2))
	testutil.AssertEquals(t, kvrwset1, kvrwset2)
}





func testPrepareBinaryFileSampleKVRWSetV1(t *testing.T) {
	b, err := proto.Marshal(constructSampleKVRWSet())
	testutil.AssertNoError(t, err, "")
	testutil.AssertNoError(t, ioutil.WriteFile(binaryTestFileName, b, 0775), "")
}

func constructSampleKVRWSet() *kvrwset.KVRWSet {
	rqi1 := &kvrwset.RangeQueryInfo{StartKey: "k0", EndKey: "k9", ItrExhausted: true}
	rqi1.SetRawReads([]*kvrwset.KVRead{
		{Key: "k1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}},
		{Key: "k2", Version: &kvrwset.Version{BlockNum: 1, TxNum: 2}},
	})

	rqi2 := &kvrwset.RangeQueryInfo{StartKey: "k00", EndKey: "k90", ItrExhausted: true}
	rqi2.SetMerkelSummary(&kvrwset.QueryReadsMerkleSummary{MaxDegree: 5, MaxLevel: 4, MaxLevelHashes: [][]byte{[]byte("Hash-1"), []byte("Hash-2")}})
	return &kvrwset.KVRWSet{
		Reads:            []*kvrwset.KVRead{{Key: "key1", Version: &kvrwset.Version{BlockNum: 1, TxNum: 1}}},
		RangeQueriesInfo: []*kvrwset.RangeQueryInfo{rqi1, rqi2},
		Writes:           []*kvrwset.KVWrite{{Key: "key2", IsDelete: false, Value: []byte("value2")}},
	}
}
