/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"
	"math"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/common/capabilities"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp"
	cb "github.com/mcc-github/blockchain/protos/common"

	"github.com/pkg/errors"
)


const (
	
	ConsortiumKey = "Consortium"

	
	HashingAlgorithmKey = "HashingAlgorithm"

	
	BlockDataHashingStructureKey = "BlockDataHashingStructure"

	
	OrdererAddressesKey = "OrdererAddresses"

	
	ChannelGroupKey = "Channel"

	
	
	CapabilitiesKey = "Capabilities"
)


type ChannelValues interface {
	
	
	HashingAlgorithm() func(input []byte) []byte

	
	
	BlockDataHashingStructureWidth() uint32

	
	OrdererAddresses() []string
}


type ChannelProtos struct {
	HashingAlgorithm          *cb.HashingAlgorithm
	BlockDataHashingStructure *cb.BlockDataHashingStructure
	OrdererAddresses          *cb.OrdererAddresses
	Consortium                *cb.Consortium
	Capabilities              *cb.Capabilities
}


type ChannelConfig struct {
	protos *ChannelProtos

	hashingAlgorithm func(input []byte) []byte

	mspManager msp.MSPManager

	appConfig         *ApplicationConfig
	ordererConfig     *OrdererConfig
	consortiumsConfig *ConsortiumsConfig
}


func NewChannelConfig(channelGroup *cb.ConfigGroup) (*ChannelConfig, error) {
	cc := &ChannelConfig{
		protos: &ChannelProtos{},
	}

	if err := DeserializeProtoValuesFromGroup(channelGroup, cc.protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	if err := cc.Validate(); err != nil {
		return nil, err
	}

	capabilities := cc.Capabilities()
	mspConfigHandler := NewMSPConfigHandler(capabilities.MSPVersion())

	var err error
	for groupName, group := range channelGroup.Groups {
		switch groupName {
		case ApplicationGroupKey:
			cc.appConfig, err = NewApplicationConfig(group, mspConfigHandler)
		case OrdererGroupKey:
			cc.ordererConfig, err = NewOrdererConfig(group, mspConfigHandler)
		case ConsortiumsGroupKey:
			cc.consortiumsConfig, err = NewConsortiumsConfig(group, mspConfigHandler)
		default:
			return nil, fmt.Errorf("Disallowed channel group: %s", group)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "could not create channel %s sub-group config", groupName)
		}
	}

	if cc.mspManager, err = mspConfigHandler.CreateMSPManager(); err != nil {
		return nil, err
	}

	return cc, nil
}


func (cc *ChannelConfig) MSPManager() msp.MSPManager {
	return cc.mspManager
}


func (cc *ChannelConfig) OrdererConfig() *OrdererConfig {
	return cc.ordererConfig
}


func (cc *ChannelConfig) ApplicationConfig() *ApplicationConfig {
	return cc.appConfig
}


func (cc *ChannelConfig) ConsortiumsConfig() *ConsortiumsConfig {
	return cc.consortiumsConfig
}


func (cc *ChannelConfig) HashingAlgorithm() func(input []byte) []byte {
	return cc.hashingAlgorithm
}


func (cc *ChannelConfig) BlockDataHashingStructureWidth() uint32 {
	return cc.protos.BlockDataHashingStructure.Width
}


func (cc *ChannelConfig) OrdererAddresses() []string {
	return cc.protos.OrdererAddresses.Addresses
}


func (cc *ChannelConfig) ConsortiumName() string {
	return cc.protos.Consortium.Name
}


func (cc *ChannelConfig) Capabilities() ChannelCapabilities {
	return capabilities.NewChannelProvider(cc.protos.Capabilities.Capabilities)
}


func (cc *ChannelConfig) Validate() error {
	for _, validator := range []func() error{
		cc.validateHashingAlgorithm,
		cc.validateBlockDataHashingStructure,
		cc.validateOrdererAddresses,
	} {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

func (cc *ChannelConfig) validateHashingAlgorithm() error {
	switch cc.protos.HashingAlgorithm.Name {
	case bccsp.SHA256:
		cc.hashingAlgorithm = util.ComputeSHA256
	case bccsp.SHA3_256:
		cc.hashingAlgorithm = util.ComputeSHA3256
	default:
		return fmt.Errorf("Unknown hashing algorithm type: %s", cc.protos.HashingAlgorithm.Name)
	}

	return nil
}

func (cc *ChannelConfig) validateBlockDataHashingStructure() error {
	if cc.protos.BlockDataHashingStructure.Width != math.MaxUint32 {
		return fmt.Errorf("BlockDataHashStructure width only supported at MaxUint32 in this version")
	}
	return nil
}

func (cc *ChannelConfig) validateOrdererAddresses() error {
	if len(cc.protos.OrdererAddresses.Addresses) == 0 {
		return fmt.Errorf("Must set some OrdererAddresses")
	}
	return nil
}
