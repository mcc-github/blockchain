/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/dispatcher"
	pb "github.com/mcc-github/blockchain/protos/peer"
	lb "github.com/mcc-github/blockchain/protos/peer/lifecycle"
)

const (
	
	InstallChaincodeFuncName = "InstallChaincode"

	
	QueryInstalledChaincodeFuncName = "QueryInstalledChaincode"

	
	QueryInstalledChaincodesFuncName = "QueryInstalledChaincodes"
)



type SCCFunctions interface {
	
	InstallChaincode(name, version string, chaincodePackage []byte) (hash []byte, err error)

	
	QueryInstalledChaincode(name, version string) (hash []byte, err error)

	
	QueryInstalledChaincodes() (chaincodes []chaincode.InstalledChaincode, err error)
}


type InstalledChaincode struct {
	Name    string
	Version string
	Hash    []byte
}



type SCC struct {
	
	Functions SCCFunctions

	
	
	Dispatcher *dispatcher.Dispatcher
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

	

	outputBytes, err := scc.Dispatcher.Dispatch(args[1], string(args[0]), scc)
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to invoke backing implementation of '%s': %s", string(args[0]), err.Error()))
	}

	return shim.Success(outputBytes)
}



func (scc *SCC) InstallChaincode(input *lb.InstallChaincodeArgs) (proto.Message, error) {
	hash, err := scc.Functions.InstallChaincode(input.Name, input.Version, input.ChaincodeInstallPackage)
	if err != nil {
		return nil, err
	}

	return &lb.InstallChaincodeResult{
		Hash: hash,
	}, nil
}



func (scc *SCC) QueryInstalledChaincode(input *lb.QueryInstalledChaincodeArgs) (proto.Message, error) {
	hash, err := scc.Functions.QueryInstalledChaincode(input.Name, input.Version)
	if err != nil {
		return nil, err
	}

	return &lb.QueryInstalledChaincodeResult{
		Hash: hash,
	}, nil
}



func (scc *SCC) QueryInstalledChaincodes(input *lb.QueryInstalledChaincodesArgs) (proto.Message, error) {
	chaincodes, err := scc.Functions.QueryInstalledChaincodes()
	if err != nil {
		return nil, err
	}

	result := &lb.QueryInstalledChaincodesResult{}
	for _, chaincode := range chaincodes {
		result.InstalledChaincodes = append(
			result.InstalledChaincodes,
			&lb.QueryInstalledChaincodesResult_InstalledChaincode{
				Name:    chaincode.Name,
				Version: chaincode.Version,
				Hash:    chaincode.Id,
			})
	}
	return result, nil
}
