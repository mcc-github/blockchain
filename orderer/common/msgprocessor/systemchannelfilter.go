/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)


type ChainCreator interface {
	
	NewChannelConfig(envConfigUpdate *cb.Envelope) (channelconfig.Resources, error)

	
	CreateBundle(channelID string, config *cb.Config) (channelconfig.Resources, error)

	
	ChannelsCount() int
}


type LimitedSupport interface {
	OrdererConfig() (channelconfig.Orderer, bool)
}


type SystemChainFilter struct {
	cc        ChainCreator
	support   LimitedSupport
	validator MetadataValidator
}


func NewSystemChannelFilter(ls LimitedSupport, cc ChainCreator, validator MetadataValidator) *SystemChainFilter {
	return &SystemChainFilter{
		support:   ls,
		cc:        cc,
		validator: validator,
	}
}


func (scf *SystemChainFilter) Apply(env *cb.Envelope) error {
	msgData := &cb.Payload{}

	err := proto.Unmarshal(env.Payload, msgData)
	if err != nil {
		return errors.Errorf("bad payload: %s", err)
	}

	if msgData.Header == nil {
		return errors.Errorf("missing payload header")
	}

	chdr, err := protoutil.UnmarshalChannelHeader(msgData.Header.ChannelHeader)
	if err != nil {
		return errors.Errorf("bad channel header: %s", err)
	}

	if chdr.Type != int32(cb.HeaderType_ORDERER_TRANSACTION) {
		return nil
	}

	ordererConfig, ok := scf.support.OrdererConfig()
	if !ok {
		logger.Panicf("System channel does not have orderer config")
	}

	maxChannels := ordererConfig.MaxChannelsCount()
	if maxChannels > 0 {
		
		if uint64(scf.cc.ChannelsCount()) > maxChannels {
			return errors.Errorf("channel creation would exceed maximimum number of channels: %d", maxChannels)
		}
	}

	if ordererConfig.ConsensusState() != orderer.ConsensusType_STATE_NORMAL {
		return errors.WithMessage(ErrMaintenanceMode, "channel creation is not permitted")
	}

	configTx := &cb.Envelope{}
	err = proto.Unmarshal(msgData.Data, configTx)
	if err != nil {
		return errors.Errorf("payload data error unmarshaling to envelope: %s", err)
	}

	return scf.authorizeAndInspect(configTx)
}

func (scf *SystemChainFilter) authorizeAndInspect(configTx *cb.Envelope) error {
	payload := &cb.Payload{}
	err := proto.Unmarshal(configTx.Payload, payload)
	if err != nil {
		return errors.Errorf("error unmarshaling wrapped configtx envelope payload: %s", err)
	}

	if payload.Header == nil {
		return errors.Errorf("wrapped configtx envelope missing header")
	}

	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return errors.Errorf("error unmarshaling wrapped configtx envelope channel header: %s", err)
	}

	if chdr.Type != int32(cb.HeaderType_CONFIG) {
		return errors.Errorf("wrapped configtx envelope not a config transaction")
	}

	configEnvelope := &cb.ConfigEnvelope{}
	err = proto.Unmarshal(payload.Data, configEnvelope)
	if err != nil {
		return errors.Errorf("error unmarshalling wrapped configtx config envelope from payload: %s", err)
	}

	if configEnvelope.LastUpdate == nil {
		return errors.Errorf("updated config does not include a config update")
	}

	res, err := scf.cc.NewChannelConfig(configEnvelope.LastUpdate)
	if err != nil {
		return errors.Errorf("error constructing new channel config from update: %s", err)
	}

	
	newChannelConfigEnv, err := res.ConfigtxValidator().ProposeConfigUpdate(configEnvelope.LastUpdate)
	if err != nil {
		return errors.Errorf("error proposing channel update to new channel config: %s", err)
	}

	
	if !proto.Equal(newChannelConfigEnv, configEnvelope) {
		return errors.Errorf("config proposed by the channel creation request did not match the config received with the channel creation request")
	}

	bundle, err := scf.cc.CreateBundle(res.ConfigtxValidator().ChannelID(), newChannelConfigEnv.Config)
	if err != nil {
		return errors.Wrap(err, "config does not validly parse")
	}

	if err = res.ValidateNew(bundle); err != nil {
		return errors.Wrap(err, "new bundle invalid")
	}

	oc, ok := bundle.OrdererConfig()
	if !ok {
		return errors.New("config is missing orderer group")
	}
	newMetadata := oc.ConsensusMetadata()

	oldOrdererConfig, ok := scf.support.OrdererConfig()
	if !ok {
		logger.Panic("old config is missing orderer group")
	}
	oldMetadata := oldOrdererConfig.ConsensusMetadata()

	if err = scf.validator.ValidateConsensusMetadata(oldMetadata, newMetadata, true); err != nil {
		return errors.Wrap(err, "consensus metadata update for channel creation is invalid")
	}

	if err = oc.Capabilities().Supported(); err != nil {
		return errors.Wrap(err, "config update is not compatible")
	}

	if err = bundle.ChannelConfig().Capabilities().Supported(); err != nil {
		return errors.Wrap(err, "config update is not compatible")
	}

	return nil
}
