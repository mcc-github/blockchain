/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package main


import (
	_ "github.com/mcc-github/blockchain/core/chaincode/shim"
	_ "github.com/mcc-github/blockchain/core/chaincode/shim/ext/attrmgr"
	_ "github.com/mcc-github/blockchain/core/chaincode/shim/ext/cid"
	_ "github.com/mcc-github/blockchain/core/chaincode/shim/ext/entities"
	_ "github.com/mcc-github/blockchain/core/chaincode/shim/ext/statebased"
)

func main() {
	return
}
