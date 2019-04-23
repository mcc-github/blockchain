/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type SimpleChaincode struct {
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	if err := stub.PutPrivateData("dummyColl", "dummyKey", []byte("dummyValue")); err != nil {
		return shim.Error(fmt.Sprintf("put operation failed. Error storing state: %s", err))
	}
	return shim.Success(nil)
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
