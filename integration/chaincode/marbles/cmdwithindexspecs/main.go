/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/mcc-github/blockchain-chaincode-go/shim"
	"github.com/mcc-github/blockchain/integration/chaincode/marbles"
)

func main() {
	err := shim.Start(&marbles.SimpleChaincode{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting marbles.SimpleChaincode: %s", err)
		os.Exit(2)
	}
}
