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

package rwsetutil

import (
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset"
	"github.com/stretchr/testify/assert"
)

const rwsetV1ProtoBytesFile = "testdata/rwsetV1ProtoBytes"




func TestRWSetV1BackwardCompatible(t *testing.T) {
	protoBytes, err := ioutil.ReadFile(rwsetV1ProtoBytesFile)
	assert.NoError(t, err)
	rwset1 := &rwset.TxReadWriteSet{}
	assert.NoError(t, proto.Unmarshal(protoBytes, rwset1))
	rwset2 := constructSampleRWSet()
	t.Logf("rwset1=%s, rwset2=%s", spew.Sdump(rwset1), spew.Sdump(rwset2))
	assert.Equal(t, rwset2, rwset1)
}





func PrepareBinaryFileSampleRWSetV1(t *testing.T) {
	b, err := proto.Marshal(constructSampleRWSet())
	assert.NoError(t, err)
	assert.NoError(t, ioutil.WriteFile(rwsetV1ProtoBytesFile, b, 0644))
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
