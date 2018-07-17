/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os"

	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/cmd/common"
	"github.com/mcc-github/blockchain/discovery/cmd"
)

func main() {
	factory.InitFactories(nil)
	cli := common.NewCLI("discover", "Command line client for blockchain discovery service")
	discovery.AddCommands(cli)
	cli.Run(os.Args[1:])
}
