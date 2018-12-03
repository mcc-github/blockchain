/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	cb "github.com/mcc-github/blockchain/protos/common"
)




const createdBlocksBuffersize = 20







type blockCreator struct {
	CreatedBlocks      chan *cb.Block
	LastCreatedBlock   *cb.Block
	LastCommittedBlock *cb.Block
	logger             *flogging.FabricLogger
}

func newBlockCreator(lastBlock *cb.Block, logger *flogging.FabricLogger) *blockCreator {
	if lastBlock == nil {
		logger.Panic("block creator initialized with nil last block")
	}
	bc := &blockCreator{
		CreatedBlocks:      make(chan *cb.Block, createdBlocksBuffersize),
		LastCreatedBlock:   lastBlock,
		LastCommittedBlock: lastBlock,
		logger:             logger,
	}

	logger.Debugf("Initialized block creator with (lastblockNumber=%d)", lastBlock.Header.Number)
	return bc
}






func (bc *blockCreator) createNextBlock(messages []*cb.Envelope) *cb.Block {
	previousBlockHash := bc.LastCreatedBlock.Header.Hash()

	data := &cb.BlockData{
		Data: make([][]byte, len(messages)),
	}

	var err error
	for i, msg := range messages {
		data.Data[i], err = proto.Marshal(msg)
		if err != nil {
			bc.logger.Panicf("Could not marshal envelope: %s", err)
		}
	}

	block := cb.NewBlock(bc.LastCreatedBlock.Header.Number+1, previousBlockHash)
	block.Header.DataHash = data.Hash()
	block.Data = data

	select {
	case bc.CreatedBlocks <- block:
		bc.LastCreatedBlock = block
		bc.logger.Debugf("Created block %d", bc.LastCreatedBlock.Header.Number)
		return block
	default:
		return nil
	}
}




func (bc *blockCreator) resetCreatedBlocks() {
	
	
	for len(bc.CreatedBlocks) > 0 {
		
		<-bc.CreatedBlocks
	}
	bc.LastCreatedBlock = bc.LastCommittedBlock
	bc.logger.Debug("Reset created blocks")
}




func (bc *blockCreator) commitBlock(block *cb.Block) {
	bc.LastCommittedBlock = block

	
	select {
	case b := <-bc.CreatedBlocks:
		if !bytes.Equal(b.Header.Bytes(), block.Header.Bytes()) {
			
			
			bc.resetCreatedBlocks()
		}
		
	default:
		
		
		
		
		bc.LastCreatedBlock = block
	}
}
