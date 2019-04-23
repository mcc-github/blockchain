/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/mcc-github/blockchain/internal/peer/chaincode"
	"github.com/mcc-github/blockchain/internal/peer/channel"
	"github.com/mcc-github/blockchain/internal/peer/clilogging"
	"github.com/mcc-github/blockchain/internal/peer/common"
	"github.com/mcc-github/blockchain/internal/peer/lifecycle"
	"github.com/mcc-github/blockchain/internal/peer/node"
	"github.com/mcc-github/blockchain/internal/peer/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



var mainCmd = &cobra.Command{Use: "peer"}

func main() {
	
	viper.SetEnvPrefix(common.CmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	
	
	mainFlags := mainCmd.PersistentFlags()

	mainFlags.String("logging-level", "", "Legacy logging level flag")
	viper.BindPFlag("logging_level", mainFlags.Lookup("logging-level"))
	mainFlags.MarkHidden("logging-level")

	mainCmd.AddCommand(version.Cmd())
	mainCmd.AddCommand(node.Cmd())
	mainCmd.AddCommand(chaincode.Cmd(nil))
	mainCmd.AddCommand(clilogging.Cmd(nil))
	mainCmd.AddCommand(channel.Cmd(nil))
	mainCmd.AddCommand(lifecycle.Cmd())

	
	
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}
