/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package cluster

import (
	"encoding/base64"
	"encoding/pem"
	"time"

	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

const (
	
	RetryTimeout = time.Second * 10
)


type PullerConfig struct {
	TLSKey              []byte
	TLSCert             []byte
	Timeout             time.Duration
	Signer              crypto.LocalSigner
	Channel             string
	MaxTotalBufferBytes int
}


func BlockPullerFromConfigBlock(conf PullerConfig, block *common.Block) (*BlockPuller, error) {
	if block == nil {
		return nil, errors.New("nil block")
	}

	endpointconfig, err := EndpointconfigFromConfigBlock(block)
	if err != nil {
		return nil, err
	}

	dialer := &StandardDialer{
		Dialer: NewTLSPinningDialer(comm.ClientConfig{
			Timeout: conf.Timeout,
			SecOpts: &comm.SecureOptions{
				ServerRootCAs:     endpointconfig.TLSRootCAs,
				Certificate:       conf.TLSCert,
				Key:               conf.TLSKey,
				RequireClientCert: true,
				UseTLS:            true,
			},
		})}

	tlsCertAsDER, _ := pem.Decode(conf.TLSCert)
	if tlsCertAsDER == nil {
		return nil, errors.Errorf("unable to decode TLS certificate PEM: %s", base64.StdEncoding.EncodeToString(conf.TLSCert))
	}

	return &BlockPuller{
		Logger:  flogging.MustGetLogger("orderer.common.cluster.replication"),
		Dialer:  dialer,
		TLSCert: tlsCertAsDER.Bytes,
		VerifyBlockSequence: func(blocks []*common.Block) error {
			return VerifyBlocks(blocks, &NoopBlockVerifier{})
		},
		MaxTotalBufferBytes: conf.MaxTotalBufferBytes,
		Endpoints:           endpointconfig.Endpoints,
		RetryTimeout:        RetryTimeout,
		FetchTimeout:        conf.Timeout,
		Channel:             conf.Channel,
		Signer:              conf.Signer,
	}, nil
}


type NoopBlockVerifier struct{}


func (*NoopBlockVerifier) VerifyBlockSignature(sd []*common.SignedData, config *common.ConfigEnvelope) error {
	return nil
}




type ChainPuller interface {
	
	PullBlock(seq uint64) *common.Block

	
	HeightsByEndpoints() map[string]uint64

	
	Close()
}


type ChainInspector struct {
	Logger          *flogging.FabricLogger
	Puller          ChainPuller
	LastConfigBlock *common.Block
}


var NotInChannelError = errors.New("not in the channel")


type selfMembershipPredicate func(configBlock *common.Block) error






func Participant(puller ChainPuller, analyzeLastConfBlock selfMembershipPredicate) error {
	endpoint, latestHeight := latestHeightAndEndpoint(puller)
	if endpoint == "" {
		return errors.New("no available orderer")
	}

	lastBlock := puller.PullBlock(latestHeight - 1)
	lastConfNumber, err := lastConfigFromBlock(lastBlock)
	if err != nil {
		return err
	}

	lastConfigBlock := puller.PullBlock(lastConfNumber)
	return analyzeLastConfBlock(lastConfigBlock)
}

func latestHeightAndEndpoint(puller ChainPuller) (string, uint64) {
	var maxHeight uint64
	var mostUpToDateEndpoint string
	for endpoint, height := range puller.HeightsByEndpoints() {
		if height > maxHeight {
			maxHeight = height
			mostUpToDateEndpoint = endpoint
		}
	}
	return mostUpToDateEndpoint, maxHeight
}

func lastConfigFromBlock(block *common.Block) (uint64, error) {
	if block.Metadata == nil || len(block.Metadata.Metadata) <= int(common.BlockMetadataIndex_LAST_CONFIG) {
		return 0, errors.New("no metadata in block")
	}
	return utils.GetLastConfigIndexFromBlock(block)
}



func (ci *ChainInspector) Channels() []string {
	channels := make(map[string]struct{})
	lastConfigBlockNum := ci.LastConfigBlock.Header.Number
	var block *common.Block
	for seq := uint64(1); seq < lastConfigBlockNum; seq++ {
		block = ci.Puller.PullBlock(seq)
		channel, err := IsNewChannelBlock(block)
		if err != nil {
			
			
			ci.Logger.Panic("Failed classifying block", seq, ":", err)
			continue
		}
		if channel == "" {
			ci.Logger.Info("Block", seq, "doesn't contain a new channel")
			continue
		}
		ci.Logger.Info("Block", seq, "contains channel", channel)
		channels[channel] = struct{}{}
	}
	
	
	
	
	
	last2Blocks := []*common.Block{block, ci.LastConfigBlock}
	if err := VerifyBlockHash(1, last2Blocks); err != nil {
		ci.Logger.Panic("System channel pulled doesn't match the boot last config block:", err)
	}

	return flattenChannelMap(channels)
}

func flattenChannelMap(m map[string]struct{}) []string {
	var res []string
	for channel := range m {
		res = append(res, channel)
	}
	return res
}



func IsNewChannelBlock(block *common.Block) (string, error) {
	if block == nil {
		return "", errors.New("nil block")
	}
	env, err := utils.ExtractEnvelope(block, 0)
	if err != nil {
		return "", err
	}
	payload, err := utils.ExtractPayload(env)
	if err != nil {
		return "", err
	}
	if payload.Header == nil {
		return "", errors.New("nil header in payload")
	}
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	
	if common.HeaderType(chdr.Type) != common.HeaderType_ORDERER_TRANSACTION {
		return "", nil
	}
	systemChannelName := chdr.ChannelId
	innerEnvelope, err := utils.UnmarshalEnvelope(payload.Data)
	if err != nil {
		return "", err
	}
	innerPayload, err := utils.UnmarshalPayload(innerEnvelope.Payload)
	if err != nil {
		return "", err
	}
	if innerPayload.Header == nil {
		return "", errors.New("inner payload's header is nil")
	}
	chdr, err = utils.UnmarshalChannelHeader(innerPayload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	
	if common.HeaderType(chdr.Type) != common.HeaderType_CONFIG {
		return "", nil
	}
	
	if chdr.ChannelId == systemChannelName {
		return "", nil
	}
	return chdr.ChannelId, nil
}
