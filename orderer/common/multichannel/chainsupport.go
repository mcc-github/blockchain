/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/configtx"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	"github.com/mcc-github/blockchain/orderer/consensus"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)


type ChainSupport struct {
	*ledgerResources
	msgprocessor.Processor
	*BlockWriter
	consensus.Chain
	cutter blockcutter.Receiver
	identity.SignerSerializer
}

func newChainSupport(
	registrar *Registrar,
	ledgerResources *ledgerResources,
	consenters map[string]consensus.Consenter,
	signer identity.SignerSerializer,
	blockcutterMetrics *blockcutter.Metrics,
) *ChainSupport {
	
	lastBlock := blockledger.GetBlock(ledgerResources, ledgerResources.Height()-1)
	metadata, err := protoutil.GetMetadataFromBlock(lastBlock, cb.BlockMetadataIndex_ORDERER)
	
	
	if err != nil {
		logger.Fatalf("[channel: %s] Error extracting orderer metadata: %s", ledgerResources.ConfigtxValidator().ChainID(), err)
	}

	
	cs := &ChainSupport{
		ledgerResources:  ledgerResources,
		SignerSerializer: signer,
		cutter: blockcutter.NewReceiverImpl(
			ledgerResources.ConfigtxValidator().ChainID(),
			ledgerResources,
			blockcutterMetrics,
		),
	}

	
	cs.Processor = msgprocessor.NewStandardChannel(cs, msgprocessor.CreateStandardChannelFilters(cs))

	
	cs.BlockWriter = newBlockWriter(lastBlock, registrar, cs)

	
	consenterType := ledgerResources.SharedConfig().ConsensusType()
	consenter, ok := consenters[consenterType]
	if !ok {
		logger.Panicf("Error retrieving consenter of type: %s", consenterType)
	}

	cs.Chain, err = consenter.HandleChain(cs, metadata)
	if err != nil {
		logger.Panicf("[channel: %s] Error creating consenter: %s", cs.ChainID(), err)
	}

	logger.Debugf("[channel: %s] Done creating channel support resources", cs.ChainID())

	return cs
}







func (cs *ChainSupport) DetectConsensusMigration() bool {
	if !cs.ledgerResources.SharedConfig().Capabilities().ConsensusTypeMigration() {
		logger.Debugf("[channel: %s] Orderer capability ConsensusTypeMigration is disabled", cs.ChainID())
		return false
	}

	if cs.Height() <= 1 {
		return false
	}

	lastConfigIndex, err := protoutil.GetLastConfigIndexFromBlock(cs.lastBlock)
	if err != nil {
		logger.Panicf("[channel: %s] Chain did not have appropriately encoded last config in its latest block: %s",
			cs.ChainID(), err)
	}
	logger.Debugf("[channel: %s] lastBlockNumber=%d, lastConfigIndex=%d",
		cs.support.ChainID(), cs.lastBlock.Header.Number, lastConfigIndex)
	if lastConfigIndex != cs.lastBlock.Header.Number {
		return false
	}

	currentState := cs.support.SharedConfig().ConsensusState()
	logger.Debugf("[channel: %s] last block ConsensusState=%s", cs.ChainID(), currentState)
	if currentState != orderer.ConsensusType_STATE_MAINTENANCE {
		return false
	}

	metadata, err := protoutil.GetMetadataFromBlock(cs.lastBlock, cb.BlockMetadataIndex_ORDERER)
	if err != nil {
		logger.Panicf("[channel: %s] Error extracting orderer metadata: %s", cs.ChainID(), err)
	}

	metaLen := len(metadata.Value)
	logger.Debugf("[channel: %s] metadata.Value length=%d", cs.ChainID(), metaLen)
	if metaLen > 0 {
		return false
	}

	prevBlock := blockledger.GetBlock(cs.Reader(), cs.lastBlock.Header.Number-1)
	prevConfigIndex, err := protoutil.GetLastConfigIndexFromBlock(prevBlock)
	if err != nil {
		logger.Panicf("Chain did not have appropriately encoded last config in block %d: %s",
			prevBlock.Header.Number, err)
	}
	if prevConfigIndex != prevBlock.Header.Number {
		return false
	}

	prevPayload, prevChanHdr := cs.extractPayloadHeaderOrPanic(prevBlock)

	switch cb.HeaderType(prevChanHdr.Type) {
	case cb.HeaderType_ORDERER_TRANSACTION:
		return false

	case cb.HeaderType_CONFIG:
		return cs.isConsensusTypeChange(prevPayload, prevChanHdr.ChannelId, prevBlock.Header.Number)

	default:
		logger.Panicf("[channel: %s] config block with unknown header type in block %d: %v",
			cs.ChainID(), prevBlock.Header.Number, prevChanHdr.Type)
		return false
	}
}

func (cs *ChainSupport) extractPayloadHeaderOrPanic(prevBlock *cb.Block) (*cb.Payload, *cb.ChannelHeader) {
	configTx, err := protoutil.ExtractEnvelope(prevBlock, 0)
	if err != nil {
		logger.Panicf("[channel: %s] Error extracting configtx from block %d: %s",
			cs.ChainID(), prevBlock.Header.Number, err)
	}
	payload, err := protoutil.UnmarshalPayload(configTx.Payload)
	if err != nil {
		logger.Panicf("[channel: %s] configtx payload is invalid in block %d: %s",
			cs.ChainID(), prevBlock.Header.Number, err)
	}
	if payload.Header == nil {
		logger.Panicf("[channel: %s] configtx payload header is missing in block %d",
			cs.ChainID(), prevBlock.Header.Number)
	}
	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		logger.Panicf("[channel: %s] invalid channel header in block %d: %s",
			cs.ChainID(), prevBlock.Header.Number, err)
	}
	return payload, chdr
}

func (cs *ChainSupport) isConsensusTypeChange(payload *cb.Payload, channelId string, headerNumber uint64) bool {
	configEnvelope, err := configtx.UnmarshalConfigEnvelope(payload.Data)
	if err != nil {
		logger.Panicf("[channel: %s] Error extracting config envelope in block %d: %s",
			cs.ChainID(), headerNumber, err)
	}

	bundle, err := cs.CreateBundle(channelId, configEnvelope.Config)
	if err != nil {
		logger.Panicf("[channel: %s] Error converting config to a bundle in block %d: %s",
			cs.ChainID(), headerNumber, err)
	}

	oc, ok := bundle.OrdererConfig()
	if !ok {
		logger.Panicf("[channel: %s] OrdererConfig missing from bundle in block %d",
			cs.ChainID(), headerNumber)
	}

	prevState := oc.ConsensusState()
	logger.Debugf("[channel: %s] previous block ConsensusState=%s", cs.ChainID(), prevState)
	if prevState != orderer.ConsensusType_STATE_MAINTENANCE {
		return false
	}

	currentType := cs.SharedConfig().ConsensusType()
	prevType := oc.ConsensusType()
	logger.Debugf("[channel: %s] block ConsensusType: previous=%s, current=%s", cs.ChainID(), prevType, currentType)
	if currentType == prevType {
		return false
	}

	logger.Infof("[channel: %s] Consensus-type migration detected, ConsensusState=%s, ConsensusType changed from %s to %s",
		cs.ChainID(), cs.support.SharedConfig().ConsensusState(), prevType, currentType)

	return true
}



func (cs *ChainSupport) Block(number uint64) *cb.Block {
	if cs.Height() <= number {
		return nil
	}
	return blockledger.GetBlock(cs.Reader(), number)
}

func (cs *ChainSupport) Reader() blockledger.Reader {
	return cs
}


func (cs *ChainSupport) Signer() identity.SignerSerializer {
	return cs
}

func (cs *ChainSupport) start() {
	cs.Chain.Start()
}


func (cs *ChainSupport) BlockCutter() blockcutter.Receiver {
	return cs.cutter
}


func (cs *ChainSupport) Validate(configEnv *cb.ConfigEnvelope) error {
	return cs.ConfigtxValidator().Validate(configEnv)
}


func (cs *ChainSupport) ProposeConfigUpdate(configtx *cb.Envelope) (*cb.ConfigEnvelope, error) {
	env, err := cs.ConfigtxValidator().ProposeConfigUpdate(configtx)
	if err != nil {
		return nil, err
	}

	bundle, err := cs.CreateBundle(cs.ChainID(), env.Config)
	if err != nil {
		return nil, err
	}

	if err = checkResources(bundle); err != nil {
		return nil, errors.Wrap(err, "config update is not compatible")
	}

	return env, cs.ValidateNew(bundle)
}


func (cs *ChainSupport) ChainID() string {
	return cs.ConfigtxValidator().ChainID()
}


func (cs *ChainSupport) ConfigProto() *cb.Config {
	return cs.ConfigtxValidator().ConfigProto()
}


func (cs *ChainSupport) Sequence() uint64 {
	return cs.ConfigtxValidator().Sequence()
}



func (cs *ChainSupport) Append(block *cb.Block) error {
	return cs.ledgerResources.ReadWriter.Append(block)
}







func (cs *ChainSupport) VerifyBlockSignature(sd []*protoutil.SignedData, envelope *cb.ConfigEnvelope) error {
	policyMgr := cs.PolicyManager()
	
	if envelope != nil {
		bundle, err := channelconfig.NewBundle(cs.ChainID(), envelope.Config)
		if err != nil {
			return err
		}
		policyMgr = bundle.PolicyManager()
	}
	policy, exists := policyMgr.GetPolicy(policies.BlockValidation)
	if !exists {
		return errors.Errorf("policy %s wasn't found", policies.BlockValidation)
	}
	err := policy.Evaluate(sd)
	if err != nil {
		return errors.Wrap(err, "block verification failed")
	}
	return nil
}
