/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	"regexp"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common.configtx")


var (
	channelAllowedChars = "[a-z][a-z0-9.-]*"
	configAllowedChars  = "[a-zA-Z0-9.-]+"
	maxLength           = 249
	illegalNames        = map[string]struct{}{
		".":  {},
		"..": {},
	}
)


type ValidatorImpl struct {
	channelID   string
	sequence    uint64
	configMap   map[string]comparable
	configProto *cb.Config
	namespace   string
	pm          policies.Manager
}






func validateConfigID(configID string) error {
	re, _ := regexp.Compile(configAllowedChars)
	
	if len(configID) <= 0 {
		return errors.New("config ID illegal, cannot be empty")
	}
	if len(configID) > maxLength {
		return errors.Errorf("config ID illegal, cannot be longer than %d", maxLength)
	}
	
	if _, ok := illegalNames[configID]; ok {
		return errors.Errorf("name '%s' for config ID is not allowed", configID)
	}
	
	matched := re.FindString(configID)
	if len(matched) != len(configID) {
		return errors.Errorf("config ID '%s' contains illegal characters", configID)
	}

	return nil
}











func validateChannelID(channelID string) error {
	re, _ := regexp.Compile(channelAllowedChars)
	
	if len(channelID) <= 0 {
		return errors.Errorf("channel ID illegal, cannot be empty")
	}
	if len(channelID) > maxLength {
		return errors.Errorf("channel ID illegal, cannot be longer than %d", maxLength)
	}

	
	matched := re.FindString(channelID)
	if len(matched) != len(channelID) {
		return errors.Errorf("channel ID '%s' contains illegal characters", channelID)
	}

	return nil
}


func NewValidatorImpl(channelID string, config *cb.Config, namespace string, pm policies.Manager) (*ValidatorImpl, error) {
	if config == nil {
		return nil, errors.Errorf("nil config parameter")
	}

	if config.ChannelGroup == nil {
		return nil, errors.Errorf("nil channel group")
	}

	if err := validateChannelID(channelID); err != nil {
		return nil, errors.Errorf("bad channel ID: %s", err)
	}

	configMap, err := mapConfig(config.ChannelGroup, namespace)
	if err != nil {
		return nil, errors.Errorf("error converting config to map: %s", err)
	}

	return &ValidatorImpl{
		namespace:   namespace,
		pm:          pm,
		sequence:    config.Sequence,
		configMap:   configMap,
		channelID:   channelID,
		configProto: config,
	}, nil
}



func (vi *ValidatorImpl) ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error) {
	return vi.proposeConfigUpdate(configtx)
}

func (vi *ValidatorImpl) proposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error) {
	configUpdateEnv, err := utils.EnvelopeToConfigUpdate(configtx)
	if err != nil {
		return nil, errors.Errorf("error converting envelope to config update: %s", err)
	}

	configMap, err := vi.authorizeUpdate(configUpdateEnv)
	if err != nil {
		return nil, errors.Errorf("error authorizing update: %s", err)
	}

	channelGroup, err := configMapToConfig(configMap, vi.namespace)
	if err != nil {
		return nil, errors.Errorf("could not turn configMap back to channelGroup: %s", err)
	}

	return &cb.ConfigEnvelope{
		Config: &cb.Config{
			Sequence:     vi.sequence + 1,
			ChannelGroup: channelGroup,
		},
		LastUpdate: configtx,
	}, nil
}


func (vi *ValidatorImpl) Validate(configEnv *cb.ConfigEnvelope) error {
	if configEnv == nil {
		return errors.Errorf("config envelope is nil")
	}

	if configEnv.Config == nil {
		return errors.Errorf("config envelope has nil config")
	}

	if configEnv.Config.Sequence != vi.sequence+1 {
		return errors.Errorf("config currently at sequence %d, cannot validate config at sequence %d", vi.sequence, configEnv.Config.Sequence)
	}

	configUpdateEnv, err := utils.EnvelopeToConfigUpdate(configEnv.LastUpdate)
	if err != nil {
		return err
	}

	configMap, err := vi.authorizeUpdate(configUpdateEnv)
	if err != nil {
		return err
	}

	channelGroup, err := configMapToConfig(configMap, vi.namespace)
	if err != nil {
		return errors.Errorf("could not turn configMap back to channelGroup: %s", err)
	}

	
	if !proto.Equal(channelGroup, configEnv.Config.ChannelGroup) {
		return errors.Errorf("ConfigEnvelope LastUpdate did not produce the supplied config result")
	}

	return nil
}


func (vi *ValidatorImpl) ChainID() string {
	return vi.channelID
}


func (vi *ValidatorImpl) Sequence() uint64 {
	return vi.sequence
}


func (vi *ValidatorImpl) ConfigProto() *cb.Config {
	return vi.configProto
}
