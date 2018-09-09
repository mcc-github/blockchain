/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package keylevelep

import (
	"encoding/json"
	"fmt"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/chaincode/shim/ext/statebased"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type EndorsementCC struct {
}


func (cc *EndorsementCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := stub.PutState("endorsed_state", []byte("foo"))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func (cc *EndorsementCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, _ := stub.GetFunctionAndParameters()
	if function, ok := functions[funcName]; ok {
		return function(stub)
	}
	return shim.Error(fmt.Sprintf("Unknown function %s", funcName))
}


var functions = map[string]func(stub shim.ChaincodeStubInterface) pb.Response{
	"addorgs":  addOrgs,
	"delorgs":  delOrgs,
	"listorgs": listOrgs,
	"delep":    delEP,
	"setval":   setVal,
	"getval":   getVal,
	"cc2cc":    invokeCC,
}



func addOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) == 0 {
		return shim.Error("No orgs to add specified")
	}

	
	epBytes, err := stub.GetStateValidationParameter("endorsed_state")
	if err != nil {
		return shim.Error(err.Error())
	}
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	err = ep.AddOrgs(statebased.RoleTypePeer, parameters...)
	if err != nil {
		return shim.Error(err.Error())
	}
	epBytes, err = ep.Policy()
	if err != nil {
		return shim.Error(err.Error())
	}

	
	err = stub.SetStateValidationParameter("endorsed_state", epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}



func delOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) == 0 {
		return shim.Error("No orgs to delete specified")
	}

	
	epBytes, err := stub.GetStateValidationParameter("endorsed_state")
	if err != nil {
		return shim.Error(err.Error())
	}
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	ep.DelOrgs(parameters...)
	epBytes, err = ep.Policy()
	if err != nil {
		return shim.Error(err.Error())
	}

	
	err = stub.SetStateValidationParameter("endorsed_state", epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}



func listOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	
	epBytes, err := stub.GetStateValidationParameter("endorsed_state")
	if err != nil {
		return shim.Error(err.Error())
	}
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	orgs := ep.ListOrgs()
	orgsList, err := json.Marshal(orgs)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(orgsList)
}


func delEP(stub shim.ChaincodeStubInterface) pb.Response {
	
	err := stub.SetStateValidationParameter("endorsed_state", nil)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}


func setVal(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) != 2 {
		return shim.Error("setval expects one argument")
	}
	err := stub.PutState("endorsed_state", args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte{})
}


func getVal(stub shim.ChaincodeStubInterface) pb.Response {
	val, err := stub.GetState("endorsed_state")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(val)
}


func invokeCC(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) < 3 {
		return shim.Error("cc2cc expects at least two arguments (channel and chaincode)")
	}
	channel := string(args[1])
	cc := string(args[2])
	nargs := args[3:]
	resp := stub.InvokeChaincode(cc, nargs, channel)
	return resp
}
