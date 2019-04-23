/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


func New() shim.Chaincode {
	return &scc{}
}

type scc struct{}


func (s *scc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (s *scc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func main() {}
