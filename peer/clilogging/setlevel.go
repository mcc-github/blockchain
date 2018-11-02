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

func setLevelCmd(cf *LoggingCmdFactory) *cobra.Command {
	var loggingSetLevelCmd = &cobra.Command{
		Use:   "setlevel <module> <log level>",
		Short: "Adds the module and log level to the current logging spec.",
		Long:  `Adds the module and log level to the current logging specification.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return setLevel(cf, cmd, args)
		},
	}
	return loggingSetLevelCmd
}

func setLevel(cf *LoggingCmdFactory, cmd *cobra.Command, args []string) (err error) {
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
					LogLevel:  args[1],
				},
			},
		}
		env := cf.wrapWithEnvelope(op)
		logResponse, err := cf.AdminClient.SetModuleLogLevel(context.Background(), env)
		if err != nil {
			return err
		}
		logger.Infof("Log level set for module name/prefix '%s': %s", logResponse.LogModule, logResponse.LogLevel)
	}
	return err
}
