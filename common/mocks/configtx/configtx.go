/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	cb "github.com/mcc-github/blockchain/protos/common"
)


type Validator struct {
	
	ChannelIDVal string

	
	SequenceVal uint64

	
	ApplyVal error

	
	AppliedConfigUpdateEnvelope *cb.ConfigEnvelope

	
	ValidateVal error

	
	ProposeConfigUpdateError error

	
	ProposeConfigUpdateVal *cb.ConfigEnvelope

	
	ConfigProtoVal *cb.Config
}


func (cm *Validator) ConfigProto() *cb.Config {
	return cm.ConfigProtoVal
}


func (cm *Validator) ChannelID() string {
	return cm.ChannelIDVal
}


func (cm *Validator) Sequence() uint64 {
	return cm.SequenceVal
}


func (cm *Validator) ProposeConfigUpdate(update *cb.Envelope) (*cb.ConfigEnvelope, error) {
	return cm.ProposeConfigUpdateVal, cm.ProposeConfigUpdateError
}


func (cm *Validator) Apply(configEnv *cb.ConfigEnvelope) error {
	cm.AppliedConfigUpdateEnvelope = configEnv
	return cm.ApplyVal
}


func (cm *Validator) Validate(configEnv *cb.ConfigEnvelope) error {
	return cm.ValidateVal
}
