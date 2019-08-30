/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	cb "github.com/mcc-github/blockchain-protos-go/common"
	mspproto "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain-protos-go/peer"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/genesis"
	"github.com/mcc-github/blockchain/core/ledger/util"
	"github.com/mcc-github/blockchain/internal/configtxgen/configtxgentest"
	"github.com/mcc-github/blockchain/internal/configtxgen/encoder"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/protoutil"
)

var logger = flogging.MustGetLogger("common.configtx.test")


func MakeGenesisBlock(chainID string) (*cb.Block, error) {
	profile := configtxgentest.Load(genesisconfig.SampleDevModeSoloProfile)
	channelGroup, err := encoder.NewChannelGroup(profile)
	if err != nil {
		logger.Panicf("Error creating channel config: %s", err)
	}

	gb := genesis.NewFactoryImpl(channelGroup).Block(chainID)
	if gb == nil {
		return gb, nil
	}

	txsFilter := util.NewTxValidationFlagsSetValue(len(gb.Data.Data), peer.TxValidationCode_VALID)
	gb.Metadata.Metadata[cb.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter

	return gb, nil
}


func MakeGenesisBlockFromMSPs(chainID string, appMSPConf, ordererMSPConf *mspproto.MSPConfig, appOrgID, ordererOrgID string) (*cb.Block, error) {
	profile := configtxgentest.Load(genesisconfig.SampleDevModeSoloProfile)
	profile.Orderer.Organizations = nil
	channelGroup, err := encoder.NewChannelGroup(profile)
	if err != nil {
		logger.Panicf("Error creating channel config: %s", err)
	}

	ordererOrg := protoutil.NewConfigGroup()
	ordererOrg.ModPolicy = channelconfig.AdminsPolicyKey
	ordererOrg.Values[channelconfig.MSPKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(channelconfig.MSPValue(ordererMSPConf).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}

	applicationOrg := protoutil.NewConfigGroup()
	applicationOrg.ModPolicy = channelconfig.AdminsPolicyKey
	applicationOrg.Values[channelconfig.MSPKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(channelconfig.MSPValue(appMSPConf).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}
	applicationOrg.Values[channelconfig.AnchorPeersKey] = &cb.ConfigValue{
		Value:     protoutil.MarshalOrPanic(channelconfig.AnchorPeersValue([]*pb.AnchorPeer{}).Value()),
		ModPolicy: channelconfig.AdminsPolicyKey,
	}

	channelGroup.Groups[channelconfig.OrdererGroupKey].Groups[ordererOrgID] = ordererOrg
	channelGroup.Groups[channelconfig.ApplicationGroupKey].Groups[appOrgID] = applicationOrg

	return genesis.NewFactoryImpl(channelGroup).Block(chainID), nil
}
