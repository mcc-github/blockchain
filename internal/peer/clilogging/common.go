/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package clilogging

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/internal/peer/common"
	common2 "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type envelopeWrapper func(msg proto.Message) *common2.Envelope


type LoggingCmdFactory struct {
	AdminClient      pb.AdminClient
	wrapWithEnvelope envelopeWrapper
}


func InitCmdFactory() (*LoggingCmdFactory, error) {
	var err error
	var adminClient pb.AdminClient

	adminClient, err = common.GetAdminClient()
	if err != nil {
		return nil, err
	}

	signer, err := common.GetDefaultSignerFnc()
	if err != nil {
		return nil, errors.Errorf("failed obtaining default signer: %v", err)
	}

	wrapEnv := func(msg proto.Message) *common2.Envelope {
		env, err := protoutil.CreateSignedEnvelope(common2.HeaderType_PEER_ADMIN_OPERATION, "", signer, msg, 0, 0)
		if err != nil {
			logger.Panicf("Failed signing: %v", err)
		}
		return env
	}

	return &LoggingCmdFactory{
		AdminClient:      adminClient,
		wrapWithEnvelope: wrapEnv,
	}, nil
}

func checkLoggingCmdParams(cmd *cobra.Command, args []string) error {
	var err error
	if cmd.Name() == "revertlevels" || cmd.Name() == "getlogspec" {
		if len(args) > 0 {
			err = errors.Errorf("more parameters than necessary were provided. Expected 0, received %d", len(args))
			return err
		}
	} else {
		
		if len(args) == 0 {
			err = errors.New("no parameters provided")
			return err
		}
	}

	if cmd.Name() == "setlevel" {
		
		if len(args) == 1 {
			err = errors.New("no log level provided")
		} else {
			
			err = common.CheckLogLevel(args[1])
		}
	}

	return err
}
