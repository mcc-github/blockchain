

package main

import (
	"chaincodes/AutoVendor/directdep"
	"fmt"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("NOT IMPL")
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("NOT IMPL")
}

func main() {
	directdep.PointlessFunction()

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
