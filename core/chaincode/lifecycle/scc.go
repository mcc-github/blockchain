/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)



type SCC struct{}


func (scc *SCC) Name() string {
	return "+lifecycle"
}


func (scc *SCC) Path() string {
	return "github.com/mcc-github/blockchain/core/chaincode/lifecycle"
}


func (scc *SCC) InitArgs() [][]byte {
	return nil
}


func (scc *SCC) Chaincode() shim.Chaincode {
	return scc
}


func (scc *SCC) InvokableExternal() bool {
	return true
}


func (scc *SCC) InvokableCC2CC() bool {
	return true
}


func (scc *SCC) Enabled() bool {
	return true
}


func (scc *SCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}



func (scc *SCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) == 0 {
		return shim.Error("lifecycle scc must be invoked with arguments")
	}

	funcName := args[0]

	switch funcName {
	
	default:
		return shim.Error(fmt.Sprintf("unknown lifecycle function: %s", funcName))
	}
}
