/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
)

type Resources struct {
	
	ConfigtxValidatorVal configtx.Validator

	
	PolicyManagerVal policies.Manager

	
	ChannelConfigVal channelconfig.Channel

	
	OrdererConfigVal channelconfig.Orderer

	
	ApplicationConfigVal channelconfig.Application

	
	ConsortiumsConfigVal channelconfig.Consortiums

	
	MSPManagerVal msp.MSPManager

	
	ValidateNewErr error
}


func (r *Resources) ConfigtxValidator() configtx.Validator {
	return r.ConfigtxValidatorVal
}


func (r *Resources) PolicyManager() policies.Manager {
	return r.PolicyManagerVal
}


func (r *Resources) ChannelConfig() channelconfig.Channel {
	return r.ChannelConfigVal
}


func (r *Resources) OrdererConfig() (channelconfig.Orderer, bool) {
	return r.OrdererConfigVal, r.OrdererConfigVal != nil
}


func (r *Resources) ApplicationConfig() (channelconfig.Application, bool) {
	return r.ApplicationConfigVal, r.ApplicationConfigVal != nil
}

func (r *Resources) ConsortiumsConfig() (channelconfig.Consortiums, bool) {
	return r.ConsortiumsConfigVal, r.ConsortiumsConfigVal != nil
}


func (r *Resources) MSPManager() msp.MSPManager {
	return r.MSPManagerVal
}


func (r *Resources) ValidateNew(res channelconfig.Resources) error {
	return r.ValidateNewErr
}
