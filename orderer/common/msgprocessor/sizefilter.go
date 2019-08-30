/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"fmt"

	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/channelconfig"
)


type SizeFilterResources interface {
	
	OrdererConfig() (channelconfig.Orderer, bool)
}


func NewSizeFilter(resources SizeFilterResources) *MaxBytesRule {
	return &MaxBytesRule{resources: resources}
}


type MaxBytesRule struct {
	resources SizeFilterResources
}


func (r *MaxBytesRule) Apply(message *common.Envelope) error {
	ordererConf, ok := r.resources.OrdererConfig()
	if !ok {
		logger.Panic("Programming error: orderer config not found")
	}

	maxBytes := ordererConf.BatchSize().AbsoluteMaxBytes
	if size := messageByteSize(message); size > maxBytes {
		return fmt.Errorf("message payload is %d bytes and exceeds maximum allowed %d bytes", size, maxBytes)
	}
	return nil
}

func messageByteSize(message *common.Envelope) uint32 {
	
	
	return uint32(len(message.Payload) + len(message.Signature))
}
