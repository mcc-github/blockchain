/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ramledger

import (
	"bytes"
	"sync"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

const pkgLogID = "orderer/ledger/ramledger"

var logger *logging.Logger

func init() {
	logger = flogging.MustGetLogger(pkgLogID)
}

type cursor struct {
	list *simpleList
}

type simpleList struct {
	next   *simpleList
	signal chan struct{}
	block  *cb.Block
}

type ramLedger struct {
	lock    sync.RWMutex
	maxSize int
	size    int
	oldest  *simpleList
	newest  *simpleList
}



func (cu *cursor) Next() (*cb.Block, cb.Status) {
	
	for {
		if cu.list.next != nil {
			cu.list = cu.list.next
			return cu.list.block, cb.Status_SUCCESS
		}
		<-cu.list.signal
	}
}


func (cu *cursor) Close() {}



func (rl *ramLedger) Iterator(startPosition *ab.SeekPosition) (blockledger.Iterator, uint64) {
	rl.lock.RLock()
	defer rl.lock.RUnlock()

	var list *simpleList
	switch start := startPosition.Type.(type) {
	case *ab.SeekPosition_Oldest:
		oldest := rl.oldest
		list = &simpleList{
			block:  &cb.Block{Header: &cb.BlockHeader{Number: oldest.block.Header.Number - 1}},
			next:   oldest,
			signal: make(chan struct{}),
		}
		close(list.signal)
	case *ab.SeekPosition_Newest:
		newest := rl.newest
		list = &simpleList{
			block:  &cb.Block{Header: &cb.BlockHeader{Number: newest.block.Header.Number - 1}},
			next:   newest,
			signal: make(chan struct{}),
		}
		close(list.signal)
	case *ab.SeekPosition_Specified:
		oldest := rl.oldest
		specified := start.Specified.Number
		logger.Debugf("Attempting to return block %d", specified)

		
		if specified+1 < oldest.block.Header.Number+1 || specified > rl.newest.block.Header.Number+1 {
			logger.Debugf("Returning error iterator because specified seek was %d with oldest %d and newest %d",
				specified, rl.oldest.block.Header.Number, rl.newest.block.Header.Number)
			return &blockledger.NotFoundErrorIterator{}, 0
		}

		if specified == oldest.block.Header.Number {
			list = &simpleList{
				block:  &cb.Block{Header: &cb.BlockHeader{Number: oldest.block.Header.Number - 1}},
				next:   oldest,
				signal: make(chan struct{}),
			}
			close(list.signal)
			break
		}

		list = oldest
		for {
			if list.block.Header.Number == specified-1 {
				break
			}
			list = list.next 
		}
	}
	cursor := &cursor{list: list}
	blockNum := list.block.Header.Number + 1

	
	if blockNum == ^uint64(0) {
		cursor.Next()
		blockNum++
	}

	return cursor, blockNum
}


func (rl *ramLedger) Height() uint64 {
	rl.lock.RLock()
	defer rl.lock.RUnlock()
	return rl.newest.block.Header.Number + 1
}


func (rl *ramLedger) Append(block *cb.Block) error {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if block.Header.Number != rl.newest.block.Header.Number+1 {
		return errors.Errorf("block number should have been %d but was %d",
			rl.newest.block.Header.Number+1, block.Header.Number)
	}

	if rl.newest.block.Header.Number+1 != 0 { 
		if !bytes.Equal(block.Header.PreviousHash, rl.newest.block.Header.Hash()) {
			return errors.Errorf("block should have had previous hash of %x but was %x",
				rl.newest.block.Header.Hash(), block.Header.PreviousHash)
		}
	}

	rl.appendBlock(block)
	return nil
}

func (rl *ramLedger) appendBlock(block *cb.Block) {
	rl.newest.next = &simpleList{
		signal: make(chan struct{}),
		block:  block,
	}

	lastSignal := rl.newest.signal
	logger.Debugf("Sending signal that block %d has a successor", rl.newest.block.Header.Number)
	rl.newest = rl.newest.next
	close(lastSignal)

	rl.size++

	if rl.size > rl.maxSize {
		logger.Debugf("RAM ledger max size about to be exceeded, removing oldest item: %d",
			rl.oldest.block.Header.Number)
		rl.oldest = rl.oldest.next
		rl.size--
	}
}
