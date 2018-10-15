/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"

	"github.com/pkg/errors"
)

const (
	InstallChaincodeFuncName = "InstallChaincode"
)



type SCCFunctions interface {
	InstallChaincode(name, version string, chaincodePackage []byte) (hash []byte, err error)
}



type SCC struct {
	Protobuf  Protobuf
	Functions SCCFunctions
}


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

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("lifecycle scc operations require exactly two arguments but received %d", len(args)))
	}

	funcName := args[0]
	inputBytes := args[1]

	

	switch string(funcName) {
	
	case InstallChaincodeFuncName:
		input := &lb.InstallChaincodeArgs{}
		err := scc.Protobuf.Unmarshal(inputBytes, input)
		if err != nil {
			err = errors.WithMessage(err, "failed to decode input arg to InstallChaincode")
			return shim.Error(err.Error())
		}

		hash, err := scc.Functions.InstallChaincode(input.Name, input.Version, input.ChaincodeInstallPackage)
		if err != nil {
			err = errors.WithMessage(err, "failed to invoke backing InstallChaincode")
			return shim.Error(err.Error())
		}

		resultBytes, err := scc.Protobuf.Marshal(&lb.InstallChaincodeResult{
			Hash: hash,
		})
		if err != nil {
			err = errors.WithMessage(err, "failed to marshal result")
			return shim.Error(err.Error())
		}

		return shim.Success(resultBytes)
	default:
		return shim.Error(fmt.Sprintf("unknown lifecycle function: %s", funcName))
	}
}
