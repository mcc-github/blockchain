/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package config

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)


type Config interface {
	
	ConfigProto() *cb.Config

	
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)
}


type Manager interface {
	
	GetChannelConfig(channel string) Config
}
