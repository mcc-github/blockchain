/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"

	"github.com/mcc-github/blockchain/orderer/common/localconfig"
	"github.com/pkg/errors"
)



type signerSerializer interface {
	identity.SignerSerializer
}


type StandardChannelSupport interface {
	
	Sequence() uint64

	
	ChainID() string

	
	Signer() identity.SignerSerializer

	
	
	ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error)

	OrdererConfig() (channelconfig.Orderer, bool)
}


type StandardChannel struct {
	support           StandardChannelSupport
	filters           *RuleSet 
	maintenanceFilter Rule     
}


func NewStandardChannel(support StandardChannelSupport, filters *RuleSet) *StandardChannel {
	return &StandardChannel{
		filters:           filters,
		support:           support,
		maintenanceFilter: NewMaintenanceFilter(support),
	}
}





func CreateStandardChannelFilters(filterSupport channelconfig.Resources, config localconfig.TopLevel) *RuleSet {
	rules := []Rule{
		EmptyRejectRule,
		NewSizeFilter(filterSupport),
		NewSigFilter(policies.ChannelWriters, policies.ChannelOrdererWriters, filterSupport),
	}

	if !config.General.Authentication.NoExpirationChecks {
		expirationRule := NewExpirationRejectRule(filterSupport)
		
		rules = append(rules[:2], append([]Rule{expirationRule}, rules[2:]...)...)
	}

	return NewRuleSet(rules)
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
	oc, ok := s.support.OrdererConfig()
	if !ok {
		logger.Panicf("Missing orderer config")
	}
	if oc.Capabilities().ConsensusTypeMigration() {
		if oc.ConsensusState() != orderer.ConsensusType_STATE_NORMAL {
			return 0, errors.WithMessage(
				ErrMaintenanceMode, "normal transactions are rejected")
		}
	}

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
		return nil, 0, errors.WithMessagef(err, "error applying config update to existing channel '%s'", s.support.ChainID())
	}

	config, err = protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG, s.support.ChainID(), s.support.Signer(), configEnvelope, msgVersion, epoch)
	if err != nil {
		return nil, 0, err
	}

	
	
	
	
	
	err = s.filters.Apply(config)
	if err != nil {
		return nil, 0, err
	}

	err = s.maintenanceFilter.Apply(config)
	if err != nil {
		return nil, 0, err
	}

	return config, seq, nil
}



func (s *StandardChannel) ProcessConfigMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	logger.Debugf("Processing config message for channel %s", s.support.ChainID())

	configEnvelope := &cb.ConfigEnvelope{}
	_, err = protoutil.UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, configEnvelope)
	if err != nil {
		return
	}

	return s.ProcessConfigUpdateMsg(configEnvelope.LastUpdate)
}
