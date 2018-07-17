/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package config

import (
	"fmt"
	"net"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	mspconstants "github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/discovery"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("discovery/config")


type CurrentConfigBlockGetter interface {
	
	GetCurrConfigBlock(channel string) *common.Block
}


type CurrentConfigBlockGetterFunc func(channel string) *common.Block


func (f CurrentConfigBlockGetterFunc) GetCurrConfigBlock(channel string) *common.Block {
	return f(channel)
}



type DiscoverySupport struct {
	CurrentConfigBlockGetter
}


func NewDiscoverySupport(getLastConfigBlock CurrentConfigBlockGetter) *DiscoverySupport {
	return &DiscoverySupport{
		CurrentConfigBlockGetter: getLastConfigBlock,
	}
}


func (s *DiscoverySupport) Config(channel string) (*discovery.ConfigResult, error) {
	block := s.GetCurrConfigBlock(channel)
	if block == nil {
		return nil, errors.Errorf("could not get last config block for channel %s", channel)
	}
	if block.Data == nil || len(block.Data.Data) == 0 {
		return nil, errors.Errorf("no transactions in block")
	}
	env := &common.Envelope{}
	if err := proto.Unmarshal(block.Data.Data[0], env); err != nil {
		return nil, errors.Wrap(err, "failed unmarshaling envelope")
	}
	pl := &common.Payload{}
	if err := proto.Unmarshal(env.Payload, pl); err != nil {
		return nil, errors.Wrap(err, "failed unmarshaling payload")
	}
	ce := &common.ConfigEnvelope{}
	if err := proto.Unmarshal(pl.Data, ce); err != nil {
		return nil, errors.Wrap(err, "failed unmarshaling config envelope")
	}

	if err := ValidateConfigEnvelope(ce); err != nil {
		return nil, errors.Wrap(err, "config envelope is invalid")
	}

	res := &discovery.ConfigResult{
		Msps:     make(map[string]*msp.FabricMSPConfig),
		Orderers: make(map[string]*discovery.Endpoints),
	}
	ordererGrp := ce.Config.ChannelGroup.Groups[channelconfig.OrdererGroupKey].Groups
	appGrp := ce.Config.ChannelGroup.Groups[channelconfig.ApplicationGroupKey].Groups

	ordererAddresses := &common.OrdererAddresses{}
	if err := proto.Unmarshal(ce.Config.ChannelGroup.Values[channelconfig.OrdererAddressesKey].Value, ordererAddresses); err != nil {
		return nil, errors.Wrap(err, "failed unmarshaling orderer addresses")
	}

	ordererEndpoints, err := computeOrdererEndpoints(ordererGrp, ordererAddresses)
	if err != nil {
		return nil, errors.Wrap(err, "failed computing orderer addresses")
	}
	res.Orderers = ordererEndpoints

	if err := appendMSPConfigs(ordererGrp, appGrp, res.Msps); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil

}

func computeOrdererEndpoints(ordererGrp map[string]*common.ConfigGroup, ordererAddresses *common.OrdererAddresses) (map[string]*discovery.Endpoints, error) {
	res := make(map[string]*discovery.Endpoints)
	for name, group := range ordererGrp {
		mspConfig := &msp.MSPConfig{}
		if err := proto.Unmarshal(group.Values[channelconfig.MSPKey].Value, mspConfig); err != nil {
			return nil, errors.Wrap(err, "failed parsing MSPConfig")
		}
		
		
		
		if mspConfig.Type != int32(mspconstants.FABRIC) {
			logger.Error("Orderer group", name, "is not a FABRIC MSP, but is of type", mspConfig.Type)
			continue
		}
		blockchainConfig := &msp.FabricMSPConfig{}
		if err := proto.Unmarshal(mspConfig.Config, blockchainConfig); err != nil {
			return nil, errors.Wrap(err, "failed marshaling FabricMSPConfig")
		}
		res[blockchainConfig.Name] = &discovery.Endpoints{}
		for _, endpoint := range ordererAddresses.Addresses {
			host, portStr, err := net.SplitHostPort(endpoint)
			if err != nil {
				return nil, errors.Errorf("failed parsing orderer endpoint %s", endpoint)
			}
			port, err := strconv.ParseInt(portStr, 10, 32)
			if err != nil {
				return nil, errors.Errorf("%s is not a valid port number", portStr)
			}
			res[blockchainConfig.Name].Endpoint = append(res[blockchainConfig.Name].Endpoint, &discovery.Endpoint{
				Host: host,
				Port: uint32(port),
			})
		}
	}
	return res, nil
}

func appendMSPConfigs(ordererGrp, appGrp map[string]*common.ConfigGroup, output map[string]*msp.FabricMSPConfig) error {
	for _, group := range []map[string]*common.ConfigGroup{ordererGrp, appGrp} {
		for _, grp := range group {
			mspConfig := &msp.MSPConfig{}
			if err := proto.Unmarshal(grp.Values[channelconfig.MSPKey].Value, mspConfig); err != nil {
				return errors.Wrap(err, "failed parsing MSPConfig")
			}
			
			if mspConfig.Type != int32(mspconstants.FABRIC) {
				continue
			}
			blockchainConfig := &msp.FabricMSPConfig{}
			if err := proto.Unmarshal(mspConfig.Config, blockchainConfig); err != nil {
				return errors.Wrap(err, "failed marshaling FabricMSPConfig")
			}
			if _, exists := output[blockchainConfig.Name]; exists {
				continue
			}
			output[blockchainConfig.Name] = blockchainConfig
		}
	}

	return nil
}

func ValidateConfigEnvelope(ce *common.ConfigEnvelope) error {
	if ce.Config == nil {
		return fmt.Errorf("field Config is nil")
	}
	if ce.Config.ChannelGroup == nil {
		return fmt.Errorf("field Config.ChannelGroup is nil")
	}
	grps := ce.Config.ChannelGroup.Groups
	if grps == nil {
		return fmt.Errorf("field Config.ChannelGroup.Groups is nil")
	}
	for _, field := range []string{channelconfig.OrdererGroupKey, channelconfig.ApplicationGroupKey} {
		grp, exists := grps[field]
		if !exists {
			return fmt.Errorf("key Config.ChannelGroup.Groups[%s] is missing", field)
		}
		if grp.Groups == nil {
			return fmt.Errorf("key Config.ChannelGroup.Groups[%s].Groups is nil", field)
		}
	}
	if ce.Config.ChannelGroup.Values == nil {
		return fmt.Errorf("field Config.ChannelGroup.Values is nil")
	}
	if _, exists := ce.Config.ChannelGroup.Values[channelconfig.OrdererAddressesKey]; !exists {
		return fmt.Errorf("field Config.ChannelGroup.Values is empty")
	}
	return nil
}
