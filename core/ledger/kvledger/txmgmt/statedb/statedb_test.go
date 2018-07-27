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

package statedb

import (
	"sort"
	"testing"

	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger/kvledger/txmgmt/version"
)

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Nil value to Put() did not panic\n")
		}
	}()

	batch := NewUpdateBatch()
	
	batch.Put("ns1", "key1", nil, nil)
}


func TestPutGetDeleteExistsGetUpdates(t *testing.T) {
	batch := NewUpdateBatch()
	batch.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 1))

	
	actualVersionedValue := batch.Get("ns1", "key1")
	testutil.AssertEquals(t, actualVersionedValue, &VersionedValue{Value: []byte("value1"), Version: version.NewHeight(1, 1)})
	
	actualResult := batch.Exists("ns1", "key2")
	expectedResult := false
	testutil.AssertEquals(t, actualResult, expectedResult)

	
	actualResult = batch.Exists("ns3", "key2")
	expectedResult = false
	testutil.AssertEquals(t, actualResult, expectedResult)

	
	actualVersionedValue = batch.Get("ns1", "key2")
	testutil.AssertNil(t, actualVersionedValue)
	
	actualVersionedValue = batch.Get("ns3", "key2")
	testutil.AssertNil(t, actualVersionedValue)

	batch.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 2))
	
	actualResult = batch.Exists("ns1", "key2")
	expectedResult = true
	testutil.AssertEquals(t, actualResult, expectedResult)

	
	batch.Put("ns2", "key2", []byte("value2"), version.NewHeight(1, 2))
	batch.Put("ns3", "key2", []byte("value2"), version.NewHeight(1, 2))
	actualNamespaces := batch.GetUpdatedNamespaces()
	sort.Strings(actualNamespaces)
	expectedNamespaces := []string{"ns1", "ns2", "ns3"}
	testutil.AssertEquals(t, actualNamespaces, expectedNamespaces)

	
	expectedUpdates := make(map[string]*VersionedValue)
	expectedUpdates["key1"] = &VersionedValue{Value: []byte("value1"), Version: version.NewHeight(1, 1)}
	expectedUpdates["key2"] = &VersionedValue{Value: []byte("value2"), Version: version.NewHeight(1, 2)}
	actualUpdates := batch.GetUpdates("ns1")
	testutil.AssertEquals(t, actualUpdates, expectedUpdates)

	actualUpdates = batch.GetUpdates("ns4")
	testutil.AssertNil(t, actualUpdates)

	
	batch.Delete("ns1", "key2", version.NewHeight(1, 2))
	
	actualResult = batch.Exists("ns1", "key2")
	expectedResult = true
	testutil.AssertEquals(t, actualResult, expectedResult)

}

func TestUpdateBatchIterator(t *testing.T) {
	batch := NewUpdateBatch()
	batch.Put("ns1", "key1", []byte("value1"), version.NewHeight(1, 1))
	batch.Put("ns1", "key2", []byte("value2"), version.NewHeight(1, 2))
	batch.Put("ns1", "key3", []byte("value3"), version.NewHeight(1, 3))

	batch.Put("ns2", "key6", []byte("value6"), version.NewHeight(2, 3))
	batch.Put("ns2", "key5", []byte("value5"), version.NewHeight(2, 2))
	batch.Put("ns2", "key4", []byte("value4"), version.NewHeight(2, 1))

	checkItrResults(t, batch.GetRangeScanIterator("ns1", "key2", "key3"), []*VersionedKV{
		{CompositeKey{"ns1", "key2"}, VersionedValue{[]byte("value2"), nil, version.NewHeight(1, 2)}},
	})

	checkItrResults(t, batch.GetRangeScanIterator("ns2", "key0", "key8"), []*VersionedKV{
		{CompositeKey{"ns2", "key4"}, VersionedValue{[]byte("value4"), nil, version.NewHeight(2, 1)}},
		{CompositeKey{"ns2", "key5"}, VersionedValue{[]byte("value5"), nil, version.NewHeight(2, 2)}},
		{CompositeKey{"ns2", "key6"}, VersionedValue{[]byte("value6"), nil, version.NewHeight(2, 3)}},
	})

	checkItrResults(t, batch.GetRangeScanIterator("ns2", "", ""), []*VersionedKV{
		{CompositeKey{"ns2", "key4"}, VersionedValue{[]byte("value4"), nil, version.NewHeight(2, 1)}},
		{CompositeKey{"ns2", "key5"}, VersionedValue{[]byte("value5"), nil, version.NewHeight(2, 2)}},
		{CompositeKey{"ns2", "key6"}, VersionedValue{[]byte("value6"), nil, version.NewHeight(2, 3)}},
	})

	checkItrResults(t, batch.GetRangeScanIterator("non-existing-ns", "", ""), nil)
}

func checkItrResults(t *testing.T, itr ResultsIterator, expectedResults []*VersionedKV) {
	for i := 0; i < len(expectedResults); i++ {
		res, _ := itr.Next()
		testutil.AssertEquals(t, res, expectedResults[i])
	}
	lastRes, err := itr.Next()
	testutil.AssertNoError(t, err, "")
	testutil.AssertNil(t, lastRes)
	itr.Close()
}
