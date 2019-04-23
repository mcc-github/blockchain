/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockledger

import (
	"github.com/golang/protobuf/proto"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
)

var closedChan chan struct{}

func init() {
	closedChan = make(chan struct{})
	close(closedChan)
}



type NotFoundErrorIterator struct{}


func (nfei *NotFoundErrorIterator) Next() (*cb.Block, cb.Status) {
	return nil, cb.Status_NOT_FOUND
}


func (nfei *NotFoundErrorIterator) ReadyChan() <-chan struct{} {
	return closedChan
}


func (nfei *NotFoundErrorIterator) Close() {}





func CreateNextBlock(rl Reader, messages []*cb.Envelope) *cb.Block {
	var nextBlockNumber uint64
	var previousBlockHash []byte

	if rl.Height() > 0 {
		it, _ := rl.Iterator(&ab.SeekPosition{
			Type: &ab.SeekPosition_Newest{
				Newest: &ab.SeekNewest{},
			},
		})
		block, status := it.Next()
		if status != cb.Status_SUCCESS {
			panic("Error seeking to newest block for chain with non-zero height")
		}
		nextBlockNumber = block.Header.Number + 1
		previousBlockHash = protoutil.BlockHeaderHash(block.Header)
	}

	data := &cb.BlockData{
		Data: make([][]byte, len(messages)),
	}

	var err error
	for i, msg := range messages {
		data.Data[i], err = proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
	}

	block := protoutil.NewBlock(nextBlockNumber, previousBlockHash)
	block.Header.DataHash = protoutil.BlockDataHash(data)
	block.Data = data

	return block
}


func GetBlock(rl Reader, index uint64) *cb.Block {
	iterator, _ := rl.Iterator(&ab.SeekPosition{
		Type: &ab.SeekPosition_Specified{
			Specified: &ab.SeekSpecified{Number: index},
		},
	})
	if iterator == nil {
		return nil
	}
	defer iterator.Close()
	block, status := iterator.Next()
	if status != cb.Status_SUCCESS {
		return nil
	}
	return block
}
