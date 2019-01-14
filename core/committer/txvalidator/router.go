/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"github.com/mcc-github/blockchain/core/committer/txvalidator/plugin"
	validatorv14 "github.com/mcc-github/blockchain/core/committer/txvalidator/v14"
	validatorv20 "github.com/mcc-github/blockchain/core/committer/txvalidator/v20"
	"github.com/mcc-github/blockchain/core/common/sysccprovider"
	"github.com/mcc-github/blockchain/protos/common"
)


type Validator interface {
	
	
	
	Validate(block *common.Block) error
}

type routingValidator struct {
	validatorv14.ChannelResources
	validator_v20 Validator
	validator_v14 Validator
}

func (v *routingValidator) Validate(block *common.Block) error {
	switch {
	case v.Capabilities().V2_0Validation():
		return v.validator_v20.Validate(block)
	default:
		return v.validator_v14.Validate(block)
	}
}

func NewTxValidator(chainID string, sem validatorv14.Semaphore, cr validatorv14.ChannelResources, sccp sysccprovider.SystemChaincodeProvider, pm plugin.Mapper) *routingValidator {
	return &routingValidator{
		ChannelResources: cr,
		validator_v14:    validatorv14.NewTxValidator(chainID, sem, cr, sccp, pm),
		validator_v20:    validatorv20.NewTxValidator(chainID, sem, cr, sccp, pm),
	}
}
