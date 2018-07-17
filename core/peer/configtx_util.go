/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
)


func computeFullConfig(currentConfigBundle *channelconfig.Bundle, channelConfTx *common.Envelope) (*common.Config, error) {
	fullChannelConfigEnv, err := currentConfigBundle.ConfigtxValidator().ProposeConfigUpdate(channelConfTx)
	if err != nil {
		return nil, err
	}
	return fullChannelConfigEnv.Config, nil
}


func serialize(resConfig *common.Config) ([]byte, error) {
	return proto.Marshal(resConfig)
}

func deserialize(serializedConf []byte) (*common.Config, error) {
	conf := &common.Config{}
	if err := proto.Unmarshal(serializedConf, conf); err != nil {
		return nil, err
	}
	return conf, nil
}


func retrievePersistedChannelConfig(ledger ledger.PeerLedger) (*common.Config, error) {
	qe, err := ledger.NewQueryExecutor()
	if err != nil {
		return nil, err
	}
	defer qe.Done()
	return retrievePersistedConf(qe, channelConfigKey)
}
