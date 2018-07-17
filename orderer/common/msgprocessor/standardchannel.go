/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
)


type StandardChannelSupport interface {
	
	Sequence() uint64

	
	ChainID() string

	
	Signer() crypto.LocalSigner

	
	
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)
}


type StandardChannel struct {
	support StandardChannelSupport
	filters *RuleSet
}


func NewStandardChannel(support StandardChannelSupport, filters *RuleSet) *StandardChannel {
	return &StandardChannel{
		filters: filters,
		support: support,
	}
}


func CreateStandardChannelFilters(filterSupport channelconfig.Resources) *RuleSet {
	ordererConfig, ok := filterSupport.OrdererConfig()
	if !ok {
		logger.Panicf("Missing orderer config")
	}
	return NewRuleSet([]Rule{
		EmptyRejectRule,
		NewExpirationRejectRule(filterSupport),
		NewSizeFilter(ordererConfig),
		NewSigFilter(policies.ChannelWriters, filterSupport),
	})
}


func (s *StandardChannel) ClassifyMsg(chdr *cb.ChannelHeader) Classification {
	switch chdr.Type {
	case int32(cb.HeaderType_CONFIG_UPDATE):
		return ConfigUpdateMsg
	case int32(cb.HeaderType_ORDERER_TRANSACTION):
		
		return ConfigMsg
	case int32(cb.HeaderType_CONFIG):
		
		return ConfigMsg
	default:
		return NormalMsg
	}
}



func (s *StandardChannel) ProcessNormalMsg(env *cb.Envelope) (configSeq uint64, err error) {
	configSeq = s.support.Sequence()
	err = s.filters.Apply(env)
	return
}




func (s *StandardChannel) ProcessConfigUpdateMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	logger.Debugf("Processing config update message for channel %s", s.support.ChainID())

	
	
	seq := s.support.Sequence()
	err = s.filters.Apply(env)
	if err != nil {
		return nil, 0, err
	}

	configEnvelope, err := s.support.ProposeConfigUpdate(env)
	if err != nil {
		return nil, 0, err
	}

	config, err = utils.CreateSignedEnvelope(cb.HeaderType_CONFIG, s.support.ChainID(), s.support.Signer(), configEnvelope, msgVersion, epoch)
	if err != nil {
		return nil, 0, err
	}

	
	
	
	
	
	err = s.filters.Apply(config)
	if err != nil {
		return nil, 0, err
	}

	return config, seq, nil
}



func (s *StandardChannel) ProcessConfigMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	logger.Debugf("Processing config message for channel %s", s.support.ChainID())

	configEnvelope := &cb.ConfigEnvelope{}
	_, err = utils.UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, configEnvelope)
	if err != nil {
		return
	}

	return s.ProcessConfigUpdateMsg(configEnvelope.LastUpdate)
}
