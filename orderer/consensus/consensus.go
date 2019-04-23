/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package consensus

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
)


type Consenter interface {
	
	
	
	
	
	
	HandleChain(support ConsenterSupport, metadata *cb.Metadata) (Chain, error)
}







type Chain interface {
	
	
	
	
	Order(env *cb.Envelope, configSeq uint64) error

	
	
	
	
	
	
	Configure(config *cb.Envelope, configSeq uint64) error

	
	
	
	
	
	WaitReady() error

	
	
	
	Errored() <-chan struct{}

	
	
	
	Start()

	
	Halt()
}




type ConsenterSupport interface {
	identity.SignerSerializer
	msgprocessor.Processor

	
	
	VerifyBlockSignature([]*protoutil.SignedData, *cb.ConfigEnvelope) error

	
	BlockCutter() blockcutter.Receiver

	
	SharedConfig() channelconfig.Orderer

	
	ChannelConfig() channelconfig.Channel

	
	
	CreateNextBlock(messages []*cb.Envelope) *cb.Block

	
	
	Block(number uint64) *cb.Block

	
	WriteBlock(block *cb.Block, encodedMetadataValue []byte)

	
	WriteConfigBlock(block *cb.Block, encodedMetadataValue []byte)

	
	Sequence() uint64

	
	ChainID() string

	
	Height() uint64

	
	
	Append(block *cb.Block) error
}
