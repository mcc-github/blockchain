/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
)



type Validator interface {
	
	Validate(configEnv *cb.ConfigEnvelope) error

	
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)

	
	ChannelID() string

	
	ConfigProto() *cb.Config

	
	Sequence() uint64
}
