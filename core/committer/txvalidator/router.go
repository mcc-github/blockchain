/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package txvalidator

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/channelconfig"
)




type Validator interface {
	
	
	
	Validate(block *common.Block) error
}




type CapabilityProvider interface {
	
	Capabilities() channelconfig.ApplicationCapabilities
}



type ValidationRouter struct {
	CapabilityProvider
	V20Validator Validator
	V14Validator Validator
}




func (v *ValidationRouter) Validate(block *common.Block) error {
	switch {
	case v.Capabilities().V2_0Validation():
		return v.V20Validator.Validate(block)
	default:
		return v.V14Validator.Validate(block)
	}
}
