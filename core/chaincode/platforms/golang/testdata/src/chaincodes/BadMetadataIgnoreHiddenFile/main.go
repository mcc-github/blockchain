

package main

import (
	"fmt"

	"github.com/mcc-github/blockchain-chaincode-go/shim"
)


type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
