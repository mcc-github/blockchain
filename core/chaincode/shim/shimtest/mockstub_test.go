/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package shimtest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/chaincode/shim/shimtest/mock"
	"github.com/stretchr/testify/assert"
)


type chaincode interface {
	shim.Chaincode
}

func TestMockStateRangeQueryIterator(t *testing.T) {
	stub := NewMockStub("rangeTest", nil)
	stub.MockTransactionStart("init")
	stub.PutState("1", []byte{61})
	stub.PutState("0", []byte{62})
	stub.PutState("5", []byte{65})
	stub.PutState("3", []byte{63})
	stub.PutState("4", []byte{64})
	stub.PutState("6", []byte{66})
	stub.MockTransactionEnd("init")

	expectKeys := []string{"3", "4"}
	expectValues := [][]byte{{63}, {64}}

	rqi := NewMockStateRangeQueryIterator(stub, "2", "5")

	fmt.Println("Running loop")
	for i := 0; i < 2; i++ {
		response, err := rqi.Next()
		fmt.Println("Loop", i, "got", response.Key, response.Value, err)
		if expectKeys[i] != response.Key {
			fmt.Println("Expected key", expectKeys[i], "got", response.Key)
			t.FailNow()
		}
		if expectValues[i][0] != response.Value[0] {
			fmt.Println("Expected value", expectValues[i], "got", response.Value)
		}
	}
}



func TestMockStateRangeQueryIterator_openEnded(t *testing.T) {
	stub := NewMockStub("rangeTest", nil)
	stub.MockTransactionStart("init")
	stub.PutState("1", []byte{61})
	stub.PutState("0", []byte{62})
	stub.PutState("5", []byte{65})
	stub.PutState("3", []byte{63})
	stub.PutState("4", []byte{64})
	stub.PutState("6", []byte{66})
	stub.MockTransactionEnd("init")

	rqi := NewMockStateRangeQueryIterator(stub, "", "")

	count := 0
	for rqi.HasNext() {
		rqi.Next()
		count++
	}

	if count != rqi.Stub.Keys.Len() {
		t.FailNow()
	}
}

type Marble struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Color      string `json:"color"`
	Size       int    `json:"size"`
	Owner      string `json:"owner"`
}


func jsonBytesEqual(expected []byte, actual []byte) bool {
	var infExpected, infActual interface{}
	if err := json.Unmarshal(expected, &infExpected); err != nil {
		return false
	}
	if err := json.Unmarshal(actual, &infActual); err != nil {
		return false
	}
	return reflect.DeepEqual(infActual, infExpected)
}

func TestGetStateByPartialCompositeKey(t *testing.T) {
	stub := NewMockStub("GetStateByPartialCompositeKeyTest", nil)
	stub.MockTransactionStart("init")

	marble1 := &Marble{"marble", "set-1", "red", 5, "tom"}
	
	compositeKey1, _ := stub.CreateCompositeKey(marble1.ObjectType, []string{marble1.Name, marble1.Color})
	marbleJSONBytes1, _ := json.Marshal(marble1)
	
	stub.PutState(compositeKey1, marbleJSONBytes1)

	marble2 := &Marble{"marble", "set-1", "blue", 5, "jerry"}
	compositeKey2, _ := stub.CreateCompositeKey(marble2.ObjectType, []string{marble2.Name, marble2.Color})
	marbleJSONBytes2, _ := json.Marshal(marble2)
	stub.PutState(compositeKey2, marbleJSONBytes2)

	marble3 := &Marble{"marble", "set-2", "red", 5, "tom-jerry"}
	compositeKey3, _ := stub.CreateCompositeKey(marble3.ObjectType, []string{marble3.Name, marble3.Color})
	marbleJSONBytes3, _ := json.Marshal(marble3)
	stub.PutState(compositeKey3, marbleJSONBytes3)

	stub.MockTransactionEnd("init")
	
	expectKeys := []string{compositeKey2, compositeKey1}
	expectKeysAttributes := [][]string{{"set-1", "blue"}, {"set-1", "red"}}
	expectValues := [][]byte{marbleJSONBytes2, marbleJSONBytes1}

	rqi, _ := stub.GetStateByPartialCompositeKey("marble", []string{"set-1"})
	fmt.Println("Running loop")
	for i := 0; i < 2; i++ {
		response, err := rqi.Next()
		fmt.Println("Loop", i, "got", response.Key, response.Value, err)
		if expectKeys[i] != response.Key {
			fmt.Println("Expected key", expectKeys[i], "got", response.Key)
			t.FailNow()
		}
		objectType, attributes, _ := stub.SplitCompositeKey(response.Key)
		if objectType != "marble" {
			fmt.Println("Expected objectType", "marble", "got", objectType)
			t.FailNow()
		}
		fmt.Println(attributes)
		for index, attr := range attributes {
			if expectKeysAttributes[i][index] != attr {
				fmt.Println("Expected keys attribute", expectKeysAttributes[index][i], "got", attr)
				t.FailNow()
			}
		}
		if jsonBytesEqual(expectValues[i], response.Value) != true {
			fmt.Println("Expected value", expectValues[i], "got", response.Value)
			t.FailNow()
		}
	}
}

func TestGetStateByPartialCompositeKeyCollision(t *testing.T) {
	stub := NewMockStub("GetStateByPartialCompositeKeyCollisionTest", nil)
	stub.MockTransactionStart("init")

	vehicle1Bytes := []byte("vehicle1")
	compositeKeyVehicle1, _ := stub.CreateCompositeKey("Vehicle", []string{"VIN_1234"})
	stub.PutState(compositeKeyVehicle1, vehicle1Bytes)

	vehicleListing1Bytes := []byte("vehicleListing1")
	compositeKeyVehicleListing1, _ := stub.CreateCompositeKey("VehicleListing", []string{"LIST_1234"})
	stub.PutState(compositeKeyVehicleListing1, vehicleListing1Bytes)

	stub.MockTransactionEnd("init")

	
	rqi, _ := stub.GetStateByPartialCompositeKey("Vehicle", []string{})
	i := 0
	fmt.Println("Running loop")
	for rqi.HasNext() {
		i++
		response, err := rqi.Next()
		fmt.Println("Loop", i, "got", response.Key, response.Value, err)
	}
	
	if i != 1 {
		fmt.Println("Expected 1, got", i)
		t.FailNow()
	}
}

func TestGetTxTimestamp(t *testing.T) {
	stub := NewMockStub("GetTxTimestamp", nil)
	stub.MockTransactionStart("init")

	timestamp, err := stub.GetTxTimestamp()
	if timestamp == nil || err != nil {
		t.FailNow()
	}

	stub.MockTransactionEnd("init")
}



func TestPutEmptyState(t *testing.T) {
	stub := NewMockStub("FAB-12545", nil)

	
	stub.MockTransactionStart("1")
	err := stub.PutState("empty", []byte{})
	assert.NoError(t, err)
	err = stub.PutState("nil", nil)
	assert.NoError(t, err)
	stub.MockTransactionEnd("1")

	
	stub.MockTransactionStart("2")
	val, err := stub.GetState("empty")
	assert.NoError(t, err)
	assert.Nil(t, val)
	val, err = stub.GetState("nil")
	assert.NoError(t, err)
	assert.Nil(t, val)
	
	err = stub.PutState("empty", []byte{0})
	assert.NoError(t, err)
	err = stub.PutState("nil", []byte{0})
	assert.NoError(t, err)
	stub.MockTransactionEnd("2")

	
	stub.MockTransactionStart("3")
	val, err = stub.GetState("empty")
	assert.NoError(t, err)
	assert.Equal(t, val, []byte{0})
	val, err = stub.GetState("nil")
	assert.NoError(t, err)
	assert.Equal(t, val, []byte{0})
	stub.MockTransactionEnd("3")

	
	stub.MockTransactionStart("4")
	err = stub.PutState("empty", []byte{})
	assert.NoError(t, err)
	err = stub.PutState("nil", nil)
	assert.NoError(t, err)
	stub.MockTransactionEnd("4")

	
	stub.MockTransactionStart("5")
	val, err = stub.GetState("empty")
	assert.NoError(t, err)
	assert.Nil(t, val)
	val, err = stub.GetState("nil")
	assert.NoError(t, err)
	assert.Nil(t, val)
	stub.MockTransactionEnd("5")

}







func TestMockMock(t *testing.T) {
	stub := NewMockStub("MOCKMOCK", &mock.Chaincode{})
	stub.args = [][]byte{[]byte("a"), []byte("b")}
	stub.MockInit("id", nil)
	stub.GetArgs()
	stub.GetStringArgs()
	stub.GetFunctionAndParameters()
	stub.GetTxID()
	stub.GetChannelID()
	stub.MockInvoke("id", nil)
	stub.MockInvokeWithSignedProposal("id", nil, nil)
	stub.DelState("dummy")
	stub.GetStateByRange("start", "end")
	stub.GetQueryResult("q")

	stub2 := NewMockStub("othercc", &mock.Chaincode{})
	stub2.MockPeerChaincode("othercc", stub2, "mychan")
	stub2.InvokeChaincode("othercc", nil, "mychan")
	stub2.GetCreator()
	stub2.GetTransient()
	stub2.GetBinding()
	stub2.GetSignedProposal()
	stub2.GetArgsSlice()
	stub2.SetEvent("e", nil)
	stub2.GetHistoryForKey("k")
	iter := &MockStateRangeQueryIterator{}
	iter.HasNext()
	iter.Close()
	getBytes("f", []string{"a", "b"})
	getFuncArgs([][]byte{[]byte("a")})
}
