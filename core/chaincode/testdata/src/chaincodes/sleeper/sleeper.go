/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main







import (
	"fmt"
	"strconv"
	"time"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)


type SleeperChaincode struct {
}

func (t *SleeperChaincode) sleep(sleepTime string) {
	st, _ := strconv.Atoi(sleepTime)
	if st >= 0 {
		time.Sleep(time.Duration(st) * time.Millisecond)
	}
}


func (t *SleeperChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	sleepTime := args[0]

	t.sleep(sleepTime)

	return shim.Success(nil)
}


func (t *SleeperChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "put" {
		if len(args) != 3 {
			return shim.Error("Incorrect number of arguments. Expecting 3")
		}

		
		return t.invoke(stub, args)
	} else if function == "get" {
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2")
		}

		
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"put\" or \"get\"")
}


func (t *SleeperChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	key := args[0]
	val := args[1]

	err := stub.PutState(key, []byte(val))
	if err != nil {
		return shim.Error(err.Error())
	}

	sleepTime := args[2]

	
	t.sleep(sleepTime)

	return shim.Success([]byte("OK"))
}


func (t *SleeperChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := args[0]

	
	val, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	sleepTime := args[1]

	
	t.sleep(sleepTime)

	return shim.Success(val)
}

func main() {
	err := shim.Start(new(SleeperChaincode))
	if err != nil {
		fmt.Printf("Error starting Sleeper chaincode: %s", err)
	}
}
