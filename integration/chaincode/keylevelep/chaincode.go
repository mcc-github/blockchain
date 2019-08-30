/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package keylevelep

import (
	"encoding/json"
	"fmt"

	"github.com/mcc-github/blockchain-chaincode-go/pkg/statebased"
	"github.com/mcc-github/blockchain-chaincode-go/shim"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
)


type EndorsementCC struct {
}


func (cc *EndorsementCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := stub.PutState("pub", []byte("foo"))
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
	if len(parameters) < 2 {
		return shim.Error("No orgs to add specified")
	}

	
	var epBytes []byte
	var err error
	if parameters[0] == "pub" {
		epBytes, err = stub.GetStateValidationParameter("pub")
	} else if parameters[0] == "priv" {
		epBytes, err = stub.GetPrivateDataValidationParameter("col", "priv")
	} else {
		return shim.Error("Unknown key specified")
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	err = ep.AddOrgs(statebased.RoleTypePeer, parameters[1:]...)
	if err != nil {
		return shim.Error(err.Error())
	}
	epBytes, err = ep.Policy()
	if err != nil {
		return shim.Error(err.Error())
	}

	
	if parameters[0] == "pub" {
		err = stub.SetStateValidationParameter("pub", epBytes)
	} else if parameters[0] == "priv" {
		err = stub.SetPrivateDataValidationParameter("col", "priv", epBytes)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}



func delOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) < 2 {
		return shim.Error("No orgs to delete specified")
	}

	
	var epBytes []byte
	var err error
	if parameters[0] == "pub" {
		epBytes, err = stub.GetStateValidationParameter("pub")
	} else if parameters[0] == "priv" {
		epBytes, err = stub.GetPrivateDataValidationParameter("col", "priv")
	} else {
		return shim.Error("Unknown key specified")
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

	
	if parameters[0] == "pub" {
		err = stub.SetStateValidationParameter("pub", epBytes)
	} else if parameters[0] == "priv" {
		err = stub.SetPrivateDataValidationParameter("col", "priv", epBytes)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}



func listOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) < 1 {
		return shim.Error("No key specified")
	}

	
	var epBytes []byte
	var err error
	if parameters[0] == "pub" {
		epBytes, err = stub.GetStateValidationParameter("pub")
	} else if parameters[0] == "priv" {
		epBytes, err = stub.GetPrivateDataValidationParameter("col", "priv")
	} else {
		return shim.Error("Unknown key specified")
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
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) < 1 {
		return shim.Error("No key specified")
	}

	
	var err error
	if parameters[0] == "pub" {
		err = stub.SetStateValidationParameter("pub", nil)
	} else if parameters[0] == "priv" {
		err = stub.SetPrivateDataValidationParameter("col", "priv", nil)
	} else {
		return shim.Error("Unknown key specified")
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}


func setVal(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) != 3 {
		return shim.Error("setval expects two arguments")
	}
	var err error
	if string(args[1]) == "pub" {
		err = stub.PutState("pub", args[2])
	} else if string(args[1]) == "priv" {
		err = stub.PutPrivateData("col", "priv", args[2])
	} else {
		return shim.Error("Unknown key specified")
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte{})
}


func getVal(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetArgs()
	if len(args) != 2 {
		return shim.Error("No key specified")
	}
	var err error
	var val []byte
	if string(args[1]) == "pub" {
		val, err = stub.GetState("pub")
	} else if string(args[1]) == "priv" {
		val, err = stub.GetPrivateData("col", "priv")
	} else {
		return shim.Error("Unknown key specified")
	}
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
