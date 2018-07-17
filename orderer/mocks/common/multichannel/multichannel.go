/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	mockconfig "github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/orderer/common/blockcutter"
	"github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	mockblockcutter "github.com/mcc-github/blockchain/orderer/mocks/common/blockcutter"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
)



type ConsenterSupport struct {
	
	SharedConfigVal *mockconfig.Orderer

	
	BlockCutterVal *mockblockcutter.Receiver

	
	Blocks chan *cb.Block

	
	ChainIDVal string

	
	HeightVal uint64

	
	NextBlockVal *cb.Block

	
	ClassifyMsgVal msgprocessor.Classification

	
	ConfigSeqVal uint64

	
	ProcessNormalMsgErr error

	
	ProcessConfigUpdateMsgVal *cb.Envelope

	
	ProcessConfigUpdateMsgErr error

	
	ProcessConfigMsgVal *cb.Envelope

	
	ProcessConfigMsgErr error

	
	SequenceVal uint64
}


func (mcs *ConsenterSupport) BlockCutter() blockcutter.Receiver {
	return mcs.BlockCutterVal
}


func (mcs *ConsenterSupport) SharedConfig() channelconfig.Orderer {
	return mcs.SharedConfigVal
}


func (mcs *ConsenterSupport) CreateNextBlock(data []*cb.Envelope) *cb.Block {
	block := cb.NewBlock(0, nil)
	mtxs := make([][]byte, len(data))
	for i := range data {
		mtxs[i] = utils.MarshalOrPanic(data[i])
	}
	block.Data = &cb.BlockData{Data: mtxs}
	mcs.NextBlockVal = block
	return block
}


func (mcs *ConsenterSupport) WriteBlock(block *cb.Block, encodedMetadataValue []byte) {
	if encodedMetadataValue != nil {
		block.Metadata.Metadata[cb.BlockMetadataIndex_ORDERER] = utils.MarshalOrPanic(&cb.Metadata{Value: encodedMetadataValue})
	}
	mcs.HeightVal++
	mcs.Blocks <- block
}


func (mcs *ConsenterSupport) WriteConfigBlock(block *cb.Block, encodedMetadataValue []byte) {
	mcs.WriteBlock(block, encodedMetadataValue)
}


func (mcs *ConsenterSupport) ChainID() string {
	return mcs.ChainIDVal
}


func (mcs *ConsenterSupport) Height() uint64 {
	return mcs.HeightVal
}


func (mcs *ConsenterSupport) Sign(message []byte) ([]byte, error) {
	return message, nil
}


func (mcs *ConsenterSupport) NewSignatureHeader() (*cb.SignatureHeader, error) {
	return &cb.SignatureHeader{}, nil
}


func (mcs *ConsenterSupport) ClassifyMsg(chdr *cb.ChannelHeader) msgprocessor.Classification {
	return mcs.ClassifyMsgVal
}


func (mcs *ConsenterSupport) ProcessNormalMsg(env *cb.Envelope) (configSeq uint64, err error) {
	return mcs.ConfigSeqVal, mcs.ProcessNormalMsgErr
}


func (mcs *ConsenterSupport) ProcessConfigUpdateMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	return mcs.ProcessConfigUpdateMsgVal, mcs.ConfigSeqVal, mcs.ProcessConfigUpdateMsgErr
}


func (mcs *ConsenterSupport) ProcessConfigMsg(env *cb.Envelope) (*cb.Envelope, uint64, error) {
	return mcs.ProcessConfigMsgVal, mcs.ConfigSeqVal, mcs.ProcessConfigMsgErr
}


func (mcs *ConsenterSupport) Sequence() uint64 {
	return mcs.SequenceVal
}
