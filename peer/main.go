/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/peer/chaincode"
	"github.com/mcc-github/blockchain/peer/channel"
	"github.com/mcc-github/blockchain/peer/clilogging"
	"github.com/mcc-github/blockchain/peer/common"
	"github.com/mcc-github/blockchain/peer/node"
	"github.com/mcc-github/blockchain/peer/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger = flogging.MustGetLogger("main")
var logOutput = os.Stderr


const cmdRoot = "core"



var mainCmd = &cobra.Command{
	Use: "peer"}

func main() {
	
	viper.SetEnvPrefix(cmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	
	
	mainFlags := mainCmd.PersistentFlags()

	mainFlags.String("logging-level", "", "Default logging level and overrides, see core.yaml for full syntax")
	viper.BindPFlag("logging_level", mainFlags.Lookup("logging-level"))

	mainCmd.AddCommand(version.Cmd())
	mainCmd.AddCommand(node.Cmd())
	mainCmd.AddCommand(chaincode.Cmd(nil))
	mainCmd.AddCommand(clilogging.Cmd(nil))
	mainCmd.AddCommand(channel.Cmd(nil))

	err := common.InitConfig(cmdRoot)
	if err != nil { 
		logger.Errorf("Fatal error when initializing %s config : %s", cmdRoot, err)
		os.Exit(1)
	}

	
	flogging.InitBackend(flogging.SetFormat(viper.GetString("logging.format")), logOutput)

	
	
	
	
	var loggingSpec string
	if viper.GetString("logging_level") != "" {
		loggingSpec = viper.GetString("logging_level")
	} else {
		loggingSpec = viper.GetString("logging.level")
	}
	flogging.InitFromSpec(loggingSpec)

	
	var mspMgrConfigDir = config.GetPath("peer.mspConfigPath")
	var mspID = viper.GetString("peer.localMspId")
	var mspType = viper.GetString("peer.localMspType")
	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}
	err = common.InitCrypto(mspMgrConfigDir, mspID, mspType)
	if err != nil { 
		logger.Errorf("Cannot run peer because %s", err.Error())
		os.Exit(1)
	}

	runtime.GOMAXPROCS(viper.GetInt("peer.gomaxprocs"))

	
	
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}
