/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func getSeedBlock() *cb.Block {
	seedBlock := cb.NewBlock(0, []byte("firsthash"))
	seedBlock.Data.Data = [][]byte{[]byte("somebytes")}
	return seedBlock
}

func TestValidCreatedBlocksQueue(t *testing.T) {
	seedBlock := getSeedBlock()
	logger := flogging.NewFabricLogger(zap.NewNop())
	bc := newBlockCreator(seedBlock, logger)

	t.Run("correct creation of a single block", func(t *testing.T) {
		block := bc.createNextBlock([]*cb.Envelope{{Payload: []byte("some other bytes")}})

		assert.Equal(t, seedBlock.Header.Number+1, block.Header.Number)
		assert.Equal(t, block.Data.Hash(), block.Header.DataHash)
		assert.Equal(t, seedBlock.Header.Hash(), block.Header.PreviousHash)
		
		assert.Len(t, bc.CreatedBlocks, 1)

		bc.commitBlock(block)

		assert.Empty(t, bc.CreatedBlocks)
		assert.Equal(t, bc.LastCommittedBlock.Header.Hash(), block.Header.Hash())
	})

	t.Run("ResetCreatedBlocks resets the queue of created blocks", func(t *testing.T) {
		numBlocks := 10
		blocks := []*cb.Block{}
		for i := 0; i < numBlocks; i++ {
			blocks = append(blocks, bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}}))
		}

		bc.resetCreatedBlocks()
		assert.True(t,
			bytes.Equal(bc.LastCommittedBlock.Header.Bytes(), bc.LastCreatedBlock.Header.Bytes()),
			"resetting the created blocks queue should leave the lastCommittedBlock and the lastCreatedBlock equal",
		)
		assert.Empty(t, bc.CreatedBlocks)
	})

	t.Run("commit of block without any created blocks sets the lastCreatedBlock correctly", func(t *testing.T) {
		block := bc.createNextBlock([]*cb.Envelope{{Payload: []byte("some other bytes")}})
		bc.resetCreatedBlocks()

		bc.commitBlock(block)

		assert.True(t,
			bytes.Equal(block.Header.Bytes(), bc.LastCommittedBlock.Header.Bytes()),
			"resetting the created blocks queue should leave the lastCommittedBlock and the lastCreatedBlock equal",
		)
		assert.True(t,
			bytes.Equal(block.Header.Bytes(), bc.LastCreatedBlock.Header.Bytes()),
			"resetting the created blocks queue should leave the lastCommittedBlock and the lastCreatedBlock equal",
		)
		assert.Empty(t, bc.CreatedBlocks)
	})
	t.Run("propose multiple blocks without having to commit them; tests the optimistic block creation", func(t *testing.T) {
		
		blocks := []*cb.Block{}
		
		for i := 0; i < 5; i++ {
			blocks = append(blocks, bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}}))
		}
		assert.Len(t, bc.CreatedBlocks, 5)

		
		for i := 0; i < 2; i++ {
			bc.commitBlock(blocks[i])
		}
		assert.Len(t, bc.CreatedBlocks, 3)

		
		for i := 0; i < 5; i++ {
			blocks = append(blocks, bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}}))
		}
		assert.Len(t, bc.CreatedBlocks, 8)

		
		for i := 2; i < 10; i++ {
			bc.commitBlock(blocks[i])
		}
		assert.Empty(t, bc.CreatedBlocks)

		
		for i := 0; i < 9; i++ {
			assertNextBlock(t, blocks[i], blocks[i+1])
		}
	})

	t.Run("createNextBlock returns nil after createdBlocksBuffersize blocks have been created", func(t *testing.T) {
		numBlocks := createdBlocksBuffersize
		blocks := []*cb.Block{}

		for i := 0; i < numBlocks; i++ {
			blocks = append(blocks, bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}}))
		}

		block := bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}})

		assert.Nil(t, block)
	})

	t.Run("created blocks should always be over committed blocks", func(t *testing.T) {
		

		baseLastCreatedBlock := bc.LastCreatedBlock

		
		blocks := []*cb.Block{}
		for i := 0; i < 5; i++ {
			blocks = append(blocks, bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}}))
		}

		
		alternateBlock := createBlockOverSpecifiedBlock(baseLastCreatedBlock, []*cb.Envelope{{Payload: []byte("alternate test envelope")}})
		bc.commitBlock(alternateBlock)

		
		createdBlockOverAlternateBlock := bc.createNextBlock([]*cb.Envelope{{Payload: []byte("test envelope")}})
		synthesizedBlockOverAlternateBlock := createBlockOverSpecifiedBlock(alternateBlock, []*cb.Envelope{{Payload: []byte("test envelope")}})
		assert.True(t,
			bytes.Equal(createdBlockOverAlternateBlock.Header.Bytes(), synthesizedBlockOverAlternateBlock.Header.Bytes()),
			"created and synthesized blocks should be equal",
		)
		bc.commitBlock(createdBlockOverAlternateBlock)
	})

}

func createBlockOverSpecifiedBlock(baseBlock *cb.Block, messages []*cb.Envelope) *cb.Block {
	previousBlockHash := baseBlock.Header.Hash()

	data := &cb.BlockData{
		Data: make([][]byte, len(messages)),
	}

	var err error
	for i, msg := range messages {
		data.Data[i], err = proto.Marshal(msg)
		if err != nil {
			panic(fmt.Sprintf("Could not marshal envelope: %s", err))
		}
	}

	block := cb.NewBlock(baseBlock.Header.Number+1, previousBlockHash)
	block.Header.DataHash = data.Hash()
	block.Data = data

	return block
}



func assertNextBlock(t *testing.T, block, nextBlock *cb.Block) {
	assert.Equal(t, block.Header.Number+1, nextBlock.Header.Number)
	assert.Equal(t, block.Header.Hash(), nextBlock.Header.PreviousHash)
}
