/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)



type Validator interface {
	
	Validate(configEnv *cb.ConfigEnvelope) error

	
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)

	
	ChainID() string

	
	ConfigProto() *cb.Config

	
	Sequence() uint64
}
