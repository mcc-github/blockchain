/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"

	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/integration/e2e/lifecycle/chaincode/caller"
)

func main() {
	err := shim.Start(&caller.CC{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting caller chaincode: %s", err)
		os.Exit(2)
	}
}
