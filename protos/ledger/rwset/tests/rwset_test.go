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
	"github.com/mcc-github/blockchain/protos/ledger/rwset"
)

const (
	binaryTestFileName = "rwsetV1ProtoBytes"
)




func TestRWSetV1BackwardCompatible(t *testing.T) {
	protoBytes, err := ioutil.ReadFile(binaryTestFileName)
	testutil.AssertNoError(t, err, "")
	rwset1 := &rwset.TxReadWriteSet{}
	testutil.AssertNoError(t, proto.Unmarshal(protoBytes, rwset1), "")
	rwset2 := constructSampleRWSet()
	t.Logf("rwset1=%s, rwset2=%s", spew.Sdump(rwset1), spew.Sdump(rwset2))
	testutil.AssertEquals(t, rwset1, rwset2)
}





func testPrepareBinaryFileSampleRWSetV1(t *testing.T) {
	b, err := proto.Marshal(constructSampleRWSet())
	testutil.AssertNoError(t, err, "")
	testutil.AssertNoError(t, ioutil.WriteFile(binaryTestFileName, b, 0775), "")
}

func constructSampleRWSet() *rwset.TxReadWriteSet {
	rwset1 := &rwset.TxReadWriteSet{}
	rwset1.DataModel = rwset.TxReadWriteSet_KV
	rwset1.NsRwset = []*rwset.NsReadWriteSet{
		{Namespace: "ns-1", Rwset: []byte("ns-1-rwset")},
		{Namespace: "ns-2", Rwset: []byte("ns-2-rwset")},
	}
	return rwset1
}
