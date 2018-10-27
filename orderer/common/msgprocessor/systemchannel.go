/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
)


type ChannelConfigTemplator interface {
	
	NewChannelConfig(env *cb.Envelope) (channelconfig.Resources, error)
}


type SystemChannel struct {
	*StandardChannel
	templator ChannelConfigTemplator
}


func NewSystemChannel(support StandardChannelSupport, templator ChannelConfigTemplator, filters *RuleSet) *SystemChannel {
	logger.Debugf("Creating system channel msg processor for channel %s", support.ChainID())
	return &SystemChannel{
		StandardChannel: NewStandardChannel(support, filters),
		templator:       templator,
	}
}


func CreateSystemChannelFilters(chainCreator ChainCreator, ledgerResources channelconfig.Resources) *RuleSet {
	ordererConfig, ok := ledgerResources.OrdererConfig()
	if !ok {
		logger.Panicf("Cannot create system channel filters without orderer config")
	}
	return NewRuleSet([]Rule{
		EmptyRejectRule,
		NewExpirationRejectRule(ledgerResources),
		NewSizeFilter(ordererConfig),
		NewSigFilter(policies.ChannelWriters, ledgerResources),
		NewSystemChannelFilter(ledgerResources, chainCreator),
	})
}



func (s *SystemChannel) ProcessNormalMsg(msg *cb.Envelope) (configSeq uint64, err error) {
	channelID, err := utils.ChannelID(msg)
	if err != nil {
		return 0, err
	}

	
	
	
	
	if channelID != s.support.ChainID() {
		return 0, ErrChannelDoesNotExist
	}

	return s.StandardChannel.ProcessNormalMsg(msg)
}




func (s *SystemChannel) ProcessConfigUpdateMsg(envConfigUpdate *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	channelID, err := utils.ChannelID(envConfigUpdate)
	if err != nil {
		return nil, 0, err
	}

	logger.Debugf("Processing config update tx with system channel message processor for channel ID %s", channelID)

	if channelID == s.support.ChainID() {
		return s.StandardChannel.ProcessConfigUpdateMsg(envConfigUpdate)
	}

	

	logger.Debugf("Processing channel create tx for channel %s on system channel %s", channelID, s.support.ChainID())

	

	bundle, err := s.templator.NewChannelConfig(envConfigUpdate)
	if err != nil {
		return nil, 0, err
	}

	newChannelConfigEnv, err := bundle.ConfigtxValidator().ProposeConfigUpdate(envConfigUpdate)
	if err != nil {
		return nil, 0, err
	}

	newChannelEnvConfig, err := utils.CreateSignedEnvelope(cb.HeaderType_CONFIG, channelID, s.support.Signer(), newChannelConfigEnv, msgVersion, epoch)
	if err != nil {
		return nil, 0, err
	}

	wrappedOrdererTransaction, err := utils.CreateSignedEnvelope(cb.HeaderType_ORDERER_TRANSACTION, s.support.ChainID(), s.support.Signer(), newChannelEnvConfig, msgVersion, epoch)
	if err != nil {
		return nil, 0, err
	}

	
	
	
	
	
	err = s.StandardChannel.filters.Apply(wrappedOrdererTransaction)
	if err != nil {
		return nil, 0, err
	}

	return wrappedOrdererTransaction, s.support.Sequence(), nil
}






func (s *SystemChannel) ProcessConfigMsg(env *cb.Envelope) (*cb.Envelope, uint64, error) {
	payload, err := utils.UnmarshalPayload(env.Payload)
	if err != nil {
		return nil, 0, err
	}

	if payload.Header == nil {
		return nil, 0, fmt.Errorf("Abort processing config msg because no head was set")
	}

	if payload.Header.ChannelHeader == nil {
		return nil, 0, fmt.Errorf("Abort processing config msg because no channel header was set")
	}

	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, 0, fmt.Errorf("Abort processing config msg because channel header unmarshalling error: %s", err)
	}

	switch chdr.Type {
	case int32(cb.HeaderType_CONFIG):
		configEnvelope := &cb.ConfigEnvelope{}
		if err = proto.Unmarshal(payload.Data, configEnvelope); err != nil {
			return nil, 0, err
		}

		return s.StandardChannel.ProcessConfigUpdateMsg(configEnvelope.LastUpdate)

	case int32(cb.HeaderType_ORDERER_TRANSACTION):
		env, err := utils.UnmarshalEnvelope(payload.Data)
		if err != nil {
			return nil, 0, fmt.Errorf("Abort processing config msg because payload data unmarshalling error: %s", err)
		}

		configEnvelope := &cb.ConfigEnvelope{}
		_, err = utils.UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, configEnvelope)
		if err != nil {
			return nil, 0, fmt.Errorf("Abort processing config msg because payload data unmarshalling error: %s", err)
		}

		return s.ProcessConfigUpdateMsg(configEnvelope.LastUpdate)

	default:
		return nil, 0, fmt.Errorf("Panic processing config msg due to unexpected envelope type %s", cb.HeaderType_name[chdr.Type])
	}
}


type DefaultTemplatorSupport interface {
	
	ConsortiumsConfig() (channelconfig.Consortiums, bool)

	
	OrdererConfig() (channelconfig.Orderer, bool)

	
	ConfigtxValidator() configtx.Validator

	
	Signer() crypto.LocalSigner
}


type DefaultTemplator struct {
	support DefaultTemplatorSupport
}


func NewDefaultTemplator(support DefaultTemplatorSupport) *DefaultTemplator {
	return &DefaultTemplator{
		support: support,
	}
}


func (dt *DefaultTemplator) NewChannelConfig(envConfigUpdate *cb.Envelope) (channelconfig.Resources, error) {
	configUpdatePayload, err := utils.UnmarshalPayload(envConfigUpdate.Payload)
	if err != nil {
		return nil, fmt.Errorf("Failing initial channel config creation because of payload unmarshaling error: %s", err)
	}

	configUpdateEnv, err := configtx.UnmarshalConfigUpdateEnvelope(configUpdatePayload.Data)
	if err != nil {
		return nil, fmt.Errorf("Failing initial channel config creation because of config update envelope unmarshaling error: %s", err)
	}

	if configUpdatePayload.Header == nil {
		return nil, fmt.Errorf("Failed initial channel config creation because config update header was missing")
	}

	channelHeader, err := utils.UnmarshalChannelHeader(configUpdatePayload.Header.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("Failed initial channel config creation because channel header was malformed: %s", err)
	}

	configUpdate, err := configtx.UnmarshalConfigUpdate(configUpdateEnv.ConfigUpdate)
	if err != nil {
		return nil, fmt.Errorf("Failing initial channel config creation because of config update unmarshaling error: %s", err)
	}

	if configUpdate.ChannelId != channelHeader.ChannelId {
		return nil, fmt.Errorf("Failing initial channel config creation: mismatched channel IDs: '%s' != '%s'", configUpdate.ChannelId, channelHeader.ChannelId)
	}

	if configUpdate.WriteSet == nil {
		return nil, fmt.Errorf("Config update has an empty writeset")
	}

	if configUpdate.WriteSet.Groups == nil || configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey] == nil {
		return nil, fmt.Errorf("Config update has missing application group")
	}

	if uv := configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Version; uv != 1 {
		return nil, fmt.Errorf("Config update for channel creation does not set application group version to 1, was %d", uv)
	}

	consortiumConfigValue, ok := configUpdate.WriteSet.Values[channelconfig.ConsortiumKey]
	if !ok {
		return nil, fmt.Errorf("Consortium config value missing")
	}

	consortium := &cb.Consortium{}
	err = proto.Unmarshal(consortiumConfigValue.Value, consortium)
	if err != nil {
		return nil, fmt.Errorf("Error reading unmarshaling consortium name: %s", err)
	}

	applicationGroup := cb.NewConfigGroup()
	consortiumsConfig, ok := dt.support.ConsortiumsConfig()
	if !ok {
		return nil, fmt.Errorf("The ordering system channel does not appear to support creating channels")
	}

	consortiumConf, ok := consortiumsConfig.Consortiums()[consortium.Name]
	if !ok {
		return nil, fmt.Errorf("Unknown consortium name: %s", consortium.Name)
	}

	applicationGroup.Policies[channelconfig.ChannelCreationPolicyKey] = &cb.ConfigPolicy{
		Policy: consortiumConf.ChannelCreationPolicy(),
	}
	applicationGroup.ModPolicy = channelconfig.ChannelCreationPolicyKey

	
	systemChannelGroup := dt.support.ConfigtxValidator().ConfigProto().ChannelGroup

	
	
	if len(systemChannelGroup.Groups[channelconfig.ConsortiumsGroupKey].Groups[consortium.Name].Groups) > 0 &&
		len(configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups) == 0 {
		return nil, fmt.Errorf("Proposed configuration has no application group members, but consortium contains members")
	}

	
	
	if len(systemChannelGroup.Groups[channelconfig.ConsortiumsGroupKey].Groups[consortium.Name].Groups) > 0 {
		for orgName := range configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups {
			consortiumGroup, ok := systemChannelGroup.Groups[channelconfig.ConsortiumsGroupKey].Groups[consortium.Name].Groups[orgName]
			if !ok {
				return nil, fmt.Errorf("Attempted to include a member which is not in the consortium")
			}
			applicationGroup.Groups[orgName] = proto.Clone(consortiumGroup).(*cb.ConfigGroup)
		}
	}

	channelGroup := cb.NewConfigGroup()

	
	for key, value := range systemChannelGroup.Values {
		channelGroup.Values[key] = proto.Clone(value).(*cb.ConfigValue)
		if key == channelconfig.ConsortiumKey {
			
			continue
		}
	}

	for key, policy := range systemChannelGroup.Policies {
		channelGroup.Policies[key] = proto.Clone(policy).(*cb.ConfigPolicy)
	}

	
	channelGroup.Groups[channelconfig.OrdererGroupKey] = proto.Clone(systemChannelGroup.Groups[channelconfig.OrdererGroupKey]).(*cb.ConfigGroup)
	channelGroup.Groups[channelconfig.ApplicationGroupKey] = applicationGroup
	channelGroup.Values[channelconfig.ConsortiumKey] = &cb.ConfigValue{
		Value:     utils.MarshalOrPanic(channelconfig.ConsortiumValue(consortium.Name).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}

	
	
	if oc, ok := dt.support.OrdererConfig(); ok && oc.Capabilities().PredictableChannelTemplate() {
		channelGroup.ModPolicy = systemChannelGroup.ModPolicy
		zeroVersions(channelGroup)
	}

	bundle, err := channelconfig.NewBundle(channelHeader.ChannelId, &cb.Config{
		ChannelGroup: channelGroup,
	})

	if err != nil {
		return nil, err
	}

	return bundle, nil
}


func zeroVersions(cg *cb.ConfigGroup) {
	cg.Version = 0

	for _, value := range cg.Values {
		value.Version = 0
	}

	for _, policy := range cg.Policies {
		policy.Version = 0
	}

	for _, group := range cg.Groups {
		zeroVersions(group)
	}
}
