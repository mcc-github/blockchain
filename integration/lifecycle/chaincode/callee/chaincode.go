/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package callee

import (
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type CC struct {
}

func (t *CC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := stub.PutState("foo", []byte("foo"))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *CC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, _ := stub.GetFunctionAndParameters()
	switch fn {
	case "INVOKE":
		err := stub.PutState("foo", []byte("bar"))
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	case "QUERY":
		val, err := stub.GetState("foo")
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(val)
	default:
		return shim.Error("unknown function")
	}
}
