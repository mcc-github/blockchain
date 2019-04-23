/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)


type committer struct {
	lgr    ledger.PeerLedger
	blkgen *blkGenerator
	assert *assert.Assertions
}

func newCommitter(lgr ledger.PeerLedger, t *testing.T) *committer {
	return &committer{lgr, newBlockGenerator(lgr, t), assert.New(t)}
}






func (c *committer) cutBlockAndCommitWithPvtdata(trans ...*txAndPvtdata) *ledger.BlockAndPvtData {
	blk := c.blkgen.nextBlockAndPvtdata(trans...)
	blkCopy := c.copyOfBlockAndPvtdata(blk)
	c.assert.NoError(
		c.lgr.CommitWithPvtData(blk),
	)
	return blkCopy
}

func (c *committer) cutBlockAndCommitExpectError(trans ...*txAndPvtdata) (*ledger.BlockAndPvtData, error) {
	blk := c.blkgen.nextBlockAndPvtdata(trans...)
	blkCopy := c.copyOfBlockAndPvtdata(blk)
	err := c.lgr.CommitWithPvtData(blk)
	c.assert.Error(err)
	return blkCopy, err
}

func (c *committer) copyOfBlockAndPvtdata(blk *ledger.BlockAndPvtData) *ledger.BlockAndPvtData {
	blkBytes, err := proto.Marshal(blk.Block)
	c.assert.NoError(err)
	blkCopy := &common.Block{}
	c.assert.NoError(proto.Unmarshal(blkBytes, blkCopy))
	return &ledger.BlockAndPvtData{Block: blkCopy, PvtData: blk.PvtData,
		MissingPvtData: blk.MissingPvtData}
}



type blkGenerator struct {
	lastNum  uint64
	lastHash []byte
	assert   *assert.Assertions
}




func newBlockGenerator(lgr ledger.PeerLedger, t *testing.T) *blkGenerator {
	assert := assert.New(t)
	info, err := lgr.GetBlockchainInfo()
	assert.NoError(err)
	return &blkGenerator{info.Height - 1, info.CurrentBlockHash, assert}
}


func (g *blkGenerator) nextBlockAndPvtdata(trans ...*txAndPvtdata) *ledger.BlockAndPvtData {
	block := protoutil.NewBlock(g.lastNum+1, g.lastHash)
	blockPvtdata := make(map[uint64]*ledger.TxPvtData)
	for i, tran := range trans {
		seq := uint64(i)
		envelopeBytes, _ := proto.Marshal(tran.Envelope)
		block.Data.Data = append(block.Data.Data, envelopeBytes)
		if tran.Pvtws != nil {
			blockPvtdata[seq] = &ledger.TxPvtData{SeqInBlock: seq, WriteSet: tran.Pvtws}
		}
	}
	block.Header.DataHash = protoutil.BlockDataHash(block.Data)
	g.lastNum++
	g.lastHash = protoutil.BlockHeaderHash(block.Header)
	setBlockFlagsToValid(block)
	return &ledger.BlockAndPvtData{Block: block, PvtData: blockPvtdata,
		MissingPvtData: make(ledger.TxMissingPvtDataMap)}
}
