/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common.channelconfig")



const RootGroupKey = "Channel"






type Bundle struct {
	policyManager   policies.Manager
	channelConfig   *ChannelConfig
	configtxManager configtx.Validator
}


func (b *Bundle) PolicyManager() policies.Manager {
	return b.policyManager
}


func (b *Bundle) MSPManager() msp.MSPManager {
	return b.channelConfig.MSPManager()
}


func (b *Bundle) ChannelConfig() Channel {
	return b.channelConfig
}



func (b *Bundle) OrdererConfig() (Orderer, bool) {
	result := b.channelConfig.OrdererConfig()
	return result, result != nil
}



func (b *Bundle) ConsortiumsConfig() (Consortiums, bool) {
	result := b.channelConfig.ConsortiumsConfig()
	return result, result != nil
}



func (b *Bundle) ApplicationConfig() (Application, bool) {
	result := b.channelConfig.ApplicationConfig()
	return result, result != nil
}


func (b *Bundle) ConfigtxValidator() configtx.Validator {
	return b.configtxManager
}



func (b *Bundle) ValidateNew(nb Resources) error {
	if oc, ok := b.OrdererConfig(); ok {
		noc, ok := nb.OrdererConfig()
		if !ok {
			return errors.New("current config has orderer section, but new config does not")
		}

		
		if !b.channelConfig.Capabilities().ConsensusTypeMigration() {
			if oc.ConsensusType() != noc.ConsensusType() {
				return errors.Errorf("attempted to change consensus type from %s to %s",
					oc.ConsensusType(), noc.ConsensusType())
			}
		}

		for orgName, org := range oc.Organizations() {
			norg, ok := noc.Organizations()[orgName]
			if !ok {
				continue
			}
			mspID := org.MSPID()
			if mspID != norg.MSPID() {
				return errors.Errorf("orderer org %s attempted to change MSP ID from %s to %s", orgName, mspID, norg.MSPID())
			}
		}
	}

	if ac, ok := b.ApplicationConfig(); ok {
		nac, ok := nb.ApplicationConfig()
		if !ok {
			return errors.New("current config has application section, but new config does not")
		}

		for orgName, org := range ac.Organizations() {
			norg, ok := nac.Organizations()[orgName]
			if !ok {
				continue
			}
			mspID := org.MSPID()
			if mspID != norg.MSPID() {
				return errors.Errorf("application org %s attempted to change MSP ID from %s to %s", orgName, mspID, norg.MSPID())
			}
		}
	}

	if cc, ok := b.ConsortiumsConfig(); ok {
		ncc, ok := nb.ConsortiumsConfig()
		if !ok {
			return errors.Errorf("current config has consortiums section, but new config does not")
		}

		for consortiumName, consortium := range cc.Consortiums() {
			nconsortium, ok := ncc.Consortiums()[consortiumName]
			if !ok {
				continue
			}

			for orgName, org := range consortium.Organizations() {
				norg, ok := nconsortium.Organizations()[orgName]
				if !ok {
					continue
				}
				mspID := org.MSPID()
				if mspID != norg.MSPID() {
					return errors.Errorf("consortium %s org %s attempted to change MSP ID from %s to %s", consortiumName, orgName, mspID, norg.MSPID())
				}
			}
		}
	} else if _, okNew := nb.ConsortiumsConfig(); okNew {
		return errors.Errorf("current config has no consortiums section, but new config does")
	}

	return nil
}



func NewBundleFromEnvelope(env *cb.Envelope, bccsp bccsp.BCCSP) (*Bundle, error) {
	payload, err := protoutil.UnmarshalPayload(env.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload from envelope")
	}

	configEnvelope, err := configtx.UnmarshalConfigEnvelope(payload.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config envelope from payload")
	}

	if payload.Header == nil {
		return nil, errors.Errorf("envelope header cannot be nil")
	}

	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal channel header")
	}

	return NewBundle(chdr.ChannelId, configEnvelope.Config, bccsp)
}


func NewBundle(channelID string, config *cb.Config, bccsp bccsp.BCCSP) (*Bundle, error) {
	if err := preValidate(config); err != nil {
		return nil, err
	}

	channelConfig, err := NewChannelConfig(config.ChannelGroup, bccsp)
	if err != nil {
		return nil, errors.Wrap(err, "initializing channelconfig failed")
	}

	policyProviderMap := make(map[int32]policies.Provider)
	for pType := range cb.Policy_PolicyType_name {
		rtype := cb.Policy_PolicyType(pType)
		switch rtype {
		case cb.Policy_UNKNOWN:
			
		case cb.Policy_SIGNATURE:
			policyProviderMap[pType] = cauthdsl.NewPolicyProvider(channelConfig.MSPManager())
		case cb.Policy_MSP:
			
		}
	}

	policyManager, err := policies.NewManagerImpl(RootGroupKey, policyProviderMap, config.ChannelGroup)
	if err != nil {
		return nil, errors.Wrap(err, "initializing policymanager failed")
	}

	configtxManager, err := configtx.NewValidatorImpl(channelID, config, RootGroupKey, policyManager)
	if err != nil {
		return nil, errors.Wrap(err, "initializing configtx manager failed")
	}

	return &Bundle{
		policyManager:   policyManager,
		channelConfig:   channelConfig,
		configtxManager: configtxManager,
	}, nil
}

func preValidate(config *cb.Config) error {
	if config == nil {
		return errors.New("channelconfig Config cannot be nil")
	}

	if config.ChannelGroup == nil {
		return errors.New("config must contain a channel group")
	}

	if og, ok := config.ChannelGroup.Groups[OrdererGroupKey]; ok {
		if _, ok := og.Values[CapabilitiesKey]; !ok {
			if _, ok := config.ChannelGroup.Values[CapabilitiesKey]; ok {
				return errors.New("cannot enable channel capabilities without orderer support first")
			}

			if ag, ok := config.ChannelGroup.Groups[ApplicationGroupKey]; ok {
				if _, ok := ag.Values[CapabilitiesKey]; ok {
					return errors.New("cannot enable application capabilities without orderer support first")
				}
			}
		}
	}

	return nil
}
