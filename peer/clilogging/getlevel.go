/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package clilogging

import (
	"context"

	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/spf13/cobra"
)

func getLevelCmd(cf *LoggingCmdFactory) *cobra.Command {
	var loggingGetLevelCmd = &cobra.Command{
		Use:   "getlevel <module>",
		Short: "Returns the logging level of the requested module logger.",
		Long:  `Returns the logging level of the requested module logger. Note: the module name should exactly match the name that is displayed in the logs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getLevel(cf, cmd, args)
		},
	}
	return loggingGetLevelCmd
}

func getLevel(cf *LoggingCmdFactory, cmd *cobra.Command, args []string) (err error) {
	err = checkLoggingCmdParams(cmd, args)
	if err == nil {
		
		cmd.SilenceUsage = true

		if cf == nil {
			cf, err = InitCmdFactory()
			if err != nil {
				return err
			}
		}
		op := &pb.AdminOperation{
			Content: &pb.AdminOperation_LogReq{
				LogReq: &pb.LogLevelRequest{
					LogModule: args[0],
				},
			},
		}
		env := cf.wrapWithEnvelope(op)
		logResponse, err := cf.AdminClient.GetModuleLogLevel(context.Background(), env)
		if err != nil {
			return err
		}
		logger.Infof("Current log level for module '%s': %s", logResponse.LogModule, logResponse.LogLevel)
	}
	return err
}
