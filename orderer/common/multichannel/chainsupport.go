/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	"github.com/mcc-github/blockchain/orderer/consensus"
	cb "github.com/mcc-github/blockchain/protos/common"
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
