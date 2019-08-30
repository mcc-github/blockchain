/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/



package msgprocessor

import (
	"errors"

	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/flogging"
)

const (
	
	msgVersion = int32(0)
	epoch      = 0
)

var logger = flogging.MustGetLogger("orderer.common.msgprocessor")



var ErrChannelDoesNotExist = errors.New("channel does not exist")



var ErrPermissionDenied = errors.New("permission denied")



var ErrMaintenanceMode = errors.New("maintenance mode")


type Classification int

const (
	
	
	NormalMsg Classification = iota

	
	
	ConfigUpdateMsg

	
	
	ConfigMsg
)



type Processor interface {
	
	ClassifyMsg(chdr *cb.ChannelHeader) Classification

	
	
	ProcessNormalMsg(env *cb.Envelope) (configSeq uint64, err error)

	
	
	
	ProcessConfigUpdateMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error)

	
	
	
	ProcessConfigMsg(env *cb.Envelope) (*cb.Envelope, uint64, error)
}
