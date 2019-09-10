/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package marbles_private

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/mcc-github/blockchain-chaincode-go/shim"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
)


type MarblesPrivateChaincode struct {
}

type marble struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Color      string `json:"color"`
	Size       int    `json:"size"`
	Owner      string `json:"owner"`
}

type marblePrivateDetails struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Price      int    `json:"price"`
}



func (t *MarblesPrivateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}



func (t *MarblesPrivateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	
	switch function {
	case "initMarble":
		
		return t.initMarble(stub, args)
	case "readMarble":
		
		return t.readMarble(stub, args)
	case "readMarblePrivateDetails":
		
		return t.readMarblePrivateDetails(stub, args)
	case "transferMarble":
		
		return t.transferMarble(stub, args)
	case "delete":
		
		return t.delete(stub, args)
	case "getMarblesByRange":
		
		return t.getMarblesByRange(stub, args)
	case "getMarbleHash":
		
		return t.getMarbleHash(stub, args)
	case "getMarblePrivateDetailsHash":
		
		return t.getMarblePrivateDetailsHash(stub, args)
	default:
		
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}




func (t *MarblesPrivateChaincode) initMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type marbleTransientInput struct {
		Name  string `json:"name"` 
		Color string `json:"color"`
		Size  int    `json:"size"`
		Owner string `json:"owner"`
		Price int    `json:"price"`
	}

	
	fmt.Println("- start init marble")

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	marbleJsonBytes, ok := transMap["marble"]
	if !ok {
		return shim.Error("marble must be a key in the transient map")
	}

	if len(marbleJsonBytes) == 0 {
		return shim.Error("marble value in the transient map must be a non-empty JSON string")
	}

	var marbleInput marbleTransientInput
	err = json.Unmarshal(marbleJsonBytes, &marbleInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(marbleJsonBytes))
	}

	if len(marbleInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(marbleInput.Color) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if marbleInput.Size <= 0 {
		return shim.Error("size field must be a positive integer")
	}
	if len(marbleInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if marbleInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}

	
	marbleAsBytes, err := stub.GetPrivateData("collectionMarbles", marbleInput.Name)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if marbleAsBytes != nil {
		fmt.Println("This marble already exists: " + marbleInput.Name)
		return shim.Error("This marble already exists: " + marbleInput.Name)
	}

	
	marble := &marble{
		ObjectType: "marble",
		Name:       marbleInput.Name,
		Color:      marbleInput.Color,
		Size:       marbleInput.Size,
		Owner:      marbleInput.Owner,
	}
	marbleJSONasBytes, err := json.Marshal(marble)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	err = stub.PutPrivateData("collectionMarbles", marbleInput.Name, marbleJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	marblePrivateDetails := &marblePrivateDetails{
		ObjectType: "marblePrivateDetails",
		Name:       marbleInput.Name,
		Price:      marbleInput.Price,
	}
	marblePrivateDetailsBytes, err := json.Marshal(marblePrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionMarblePrivateDetails", marbleInput.Name, marblePrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	
	
	
	
	
	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{marble.Color, marble.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	
	
	value := []byte{0x00}
	stub.PutPrivateData("collectionMarbles", colorNameIndexKey, value)

	
	fmt.Println("- end init marble")
	return shim.Success(nil)
}




func (t *MarblesPrivateChaincode) readMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionMarbles", name) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}




func (t *MarblesPrivateChaincode) readMarblePrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionMarblePrivateDetails", name) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}




func (t *MarblesPrivateChaincode) getMarbleHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateDataHash("collectionMarbles", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get marble private data hash for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble private marble data hash does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}




func (t *MarblesPrivateChaincode) getMarblePrivateDetailsHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateDataHash("collectionMarblePrivateDetails", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get marble private details hash for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble private details hash does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}




func (t *MarblesPrivateChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start delete marble")

	type marbleDeleteTransientInput struct {
		Name string `json:"name"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private marble name must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	marbleDeleteJsonBytes, ok := transMap["marble_delete"]
	if !ok {
		return shim.Error("marble_delete must be a key in the transient map")
	}

	if len(marbleDeleteJsonBytes) == 0 {
		return shim.Error("marble_delete value in the transient map must be a non-empty JSON string")
	}

	var marbleDeleteInput marbleDeleteTransientInput
	err = json.Unmarshal(marbleDeleteJsonBytes, &marbleDeleteInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(marbleDeleteJsonBytes))
	}

	if len(marbleDeleteInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}

	
	valAsbytes, err := stub.GetPrivateData("collectionMarbles", marbleDeleteInput.Name) 
	if err != nil {
		return shim.Error("Failed to get state for " + marbleDeleteInput.Name)
	} else if valAsbytes == nil {
		return shim.Error("Marble does not exist: " + marbleDeleteInput.Name)
	}

	var marbleToDelete marble
	err = json.Unmarshal([]byte(valAsbytes), &marbleToDelete)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(valAsbytes))
	}

	
	err = stub.DelPrivateData("collectionMarbles", marbleDeleteInput.Name)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	
	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{marbleToDelete.Color, marbleToDelete.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.DelPrivateData("collectionMarbles", colorNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	
	err = stub.DelPrivateData("collectionMarblePrivateDetails", marbleDeleteInput.Name)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}




func (t *MarblesPrivateChaincode) transferMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- start transfer marble")

	type marbleTransferTransientInput struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	marbleOwnerJsonBytes, ok := transMap["marble_owner"]
	if !ok {
		return shim.Error("marble_owner must be a key in the transient map")
	}

	if len(marbleOwnerJsonBytes) == 0 {
		return shim.Error("marble_owner value in the transient map must be a non-empty JSON string")
	}

	var marbleTransferInput marbleTransferTransientInput
	err = json.Unmarshal(marbleOwnerJsonBytes, &marbleTransferInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(marbleOwnerJsonBytes))
	}

	if len(marbleTransferInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(marbleTransferInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}

	marbleAsBytes, err := stub.GetPrivateData("collectionMarbles", marbleTransferInput.Name)
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if marbleAsBytes == nil {
		return shim.Error("Marble does not exist: " + marbleTransferInput.Name)
	}

	marbleToTransfer := marble{}
	err = json.Unmarshal(marbleAsBytes, &marbleToTransfer) 
	if err != nil {
		return shim.Error(err.Error())
	}
	marbleToTransfer.Owner = marbleTransferInput.Owner 

	marbleJSONasBytes, _ := json.Marshal(marbleToTransfer)
	err = stub.PutPrivateData("collectionMarbles", marbleToTransfer.Name, marbleJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferMarble (success)")
	return shim.Success(nil)
}












func (t *MarblesPrivateChaincode) getMarblesByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetPrivateDataByRange("collectionMarbles", startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		
		if bArrayMemberAlreadyWritten {
			buffer.WriteString(",")
		}

		buffer.WriteString(
			fmt.Sprintf(
				`{"Key":"%s", "Record":%s}`,
				queryResponse.Key, queryResponse.Value,
			),
		)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
