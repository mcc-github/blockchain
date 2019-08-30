/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protoutil"
)

const (
	channelConfigKey = "resourcesconfigtx.CHANNEL_CONFIG_KEY"
	peerNamespace    = ""
)


type ConfigTxProcessor struct{}









func (tp *ConfigTxProcessor) GenerateSimulationResults(txEnv *common.Envelope, simulator ledger.TxSimulator, initializingLedger bool) error {
	payload := protoutil.UnmarshalPayloadOrPanic(txEnv.Payload)
	channelHdr := protoutil.UnmarshalChannelHeaderOrPanic(payload.Header.ChannelHeader)
	txType := common.HeaderType(channelHdr.GetType())

	switch txType {
	case common.HeaderType_CONFIG:
		peerLogger.Debugf("Processing CONFIG")
		return processChannelConfigTx(txEnv, simulator)

	default:
		return fmt.Errorf("tx type [%s] is not expected", txType)
	}
}

func processChannelConfigTx(txEnv *common.Envelope, simulator ledger.TxSimulator) error {
	configEnvelope := &common.ConfigEnvelope{}
	if _, err := protoutil.UnmarshalEnvelopeOfType(txEnv, common.HeaderType_CONFIG, configEnvelope); err != nil {
		return err
	}
	channelConfig := configEnvelope.Config
	if channelConfig == nil {
		return fmt.Errorf("channel config found nil")
	}

	if err := persistConf(simulator, channelConfigKey, channelConfig); err != nil {
		return err
	}

	peerLogger.Debugf("channelConfig=%s", channelConfig)

	return nil
}

func persistConf(simulator ledger.TxSimulator, key string, config *common.Config) error {
	serializedConfig, err := proto.Marshal(config)
	if err != nil {
		return err
	}
	return simulator.SetState(peerNamespace, key, serializedConfig)
}

func retrievePersistedConf(queryExecuter ledger.QueryExecutor, key string) (*common.Config, error) {
	serializedConfig, err := queryExecuter.GetState(peerNamespace, key)
	if err != nil {
		return nil, err
	}
	if serializedConfig == nil {
		return nil, nil
	}
	return deserialize(serializedConfig)
}

func deserialize(serializedConf []byte) (*common.Config, error) {
	conf := &common.Config{}
	if err := proto.Unmarshal(serializedConf, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
