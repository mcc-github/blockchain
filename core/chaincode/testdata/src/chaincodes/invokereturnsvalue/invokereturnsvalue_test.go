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
package main

import (
	"fmt"
	"testing"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/chaincode/shim/shimtest"
)

func checkInit(t *testing.T, stub *shimtest.MockStub, args [][]byte, retval []byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
	if retval != nil {
		if res.Payload == nil {
			fmt.Printf("Init returned nil, expected %s", string(retval))
			t.FailNow()
		}
		if string(res.Payload) != string(retval) {
			fmt.Printf("Init returned %s, expected %s", string(res.Payload), string(retval))
			t.FailNow()
		}
	}
}

func checkState(t *testing.T, stub *shimtest.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shimtest.MockStub, args [][]byte, retval []byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}

	if retval != nil {
		if res.Payload == nil {
			fmt.Printf("Invoke returned nil, expected %s", string(retval))
			t.FailNow()
		}
		if string(res.Payload) != string(retval) {
			fmt.Printf("Invoke returned %s, expected %s", string(res.Payload), string(retval))
			t.FailNow()
		}
	}
}

func Test_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shimtest.NewMockStub("ex02", scc)

	
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("123"), []byte("B"), []byte("234")}, []byte("OK"))

	checkState(t, stub, "A", "123")
	checkState(t, stub, "B", "234")
}

func Test_Invoke(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shimtest.NewMockStub("ex02", scc)

	
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("567"), []byte("B"), []byte("678")}, []byte("OK"))

	
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("123")}, []byte("{444,801}"))

	
	checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("B"), []byte("A"), []byte("234")}, []byte("{567,678}"))
}
