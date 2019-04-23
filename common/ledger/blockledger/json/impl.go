/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package jsonledger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	cb "github.com/mcc-github/blockchain/protos/common"
	ab "github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common.ledger.blockledger.json")

var closedChan chan struct{}
var fileLock sync.Mutex

func init() {
	closedChan = make(chan struct{})
	close(closedChan)
}

const (
	blockFileFormatString      = "block_%020d.json"
	chainDirectoryFormatString = "chain_%s"
)

type cursor struct {
	jl          *jsonLedger
	blockNumber uint64
}

type jsonLedger struct {
	directory string
	height    uint64
	lastHash  []byte
	marshaler *jsonpb.Marshaler

	mutex  sync.Mutex
	signal chan struct{}
}


func (jl *jsonLedger) readBlock(number uint64) (*cb.Block, bool) {
	name := jl.blockFilename(number)

	
	
	fileLock.Lock()
	defer fileLock.Unlock()

	file, err := os.Open(name)
	if err == nil {
		defer file.Close()
		block := &cb.Block{}
		err = jsonpb.Unmarshal(file, block)
		if err != nil {
			return nil, true
		}
		logger.Debugf("Read block %d", block.Header.Number)
		return block, true
	}
	return nil, false
}



func (cu *cursor) Next() (*cb.Block, cb.Status) {
	
	
	for {
		block, found := cu.jl.readBlock(cu.blockNumber)
		if found {
			if block == nil {
				return nil, cb.Status_SERVICE_UNAVAILABLE
			}
			cu.blockNumber++
			return block, cb.Status_SUCCESS
		}

		
		
		cu.jl.mutex.Lock()
		signal := cu.jl.signal
		cu.jl.mutex.Unlock()
		<-signal
	}
}

func (cu *cursor) Close() {}



func (jl *jsonLedger) Iterator(startPosition *ab.SeekPosition) (blockledger.Iterator, uint64) {
	switch start := startPosition.Type.(type) {
	case *ab.SeekPosition_Oldest:
		return &cursor{jl: jl, blockNumber: 0}, 0
	case *ab.SeekPosition_Newest:
		high := jl.height - 1
		return &cursor{jl: jl, blockNumber: high}, high
	case *ab.SeekPosition_Specified:
		if start.Specified.Number > jl.height {
			return &blockledger.NotFoundErrorIterator{}, 0
		}
		return &cursor{jl: jl, blockNumber: start.Specified.Number}, start.Specified.Number
	default:
		return &blockledger.NotFoundErrorIterator{}, 0
	}
}


func (jl *jsonLedger) Height() uint64 {
	return jl.height
}


func (jl *jsonLedger) Append(block *cb.Block) error {
	if block.Header.Number != jl.height {
		return errors.Errorf("block number should have been %d but was %d", jl.height, block.Header.Number)
	}

	if !bytes.Equal(block.Header.PreviousHash, jl.lastHash) {
		return errors.Errorf("block should have had previous hash of %x but was %x", jl.lastHash, block.Header.PreviousHash)
	}

	jl.writeBlock(block)
	jl.lastHash = protoutil.BlockHeaderHash(block.Header)
	jl.height++

	
	jl.mutex.Lock()
	close(jl.signal)
	jl.signal = make(chan struct{})
	jl.mutex.Unlock()
	return nil
}


func (jl *jsonLedger) writeBlock(block *cb.Block) {
	name := jl.blockFilename(block.Header.Number)

	fileLock.Lock()
	defer fileLock.Unlock()

	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = jl.marshaler.Marshal(file, block)
	logger.Debugf("Wrote block %d", block.Header.Number)
	if err != nil {
		logger.Panicf("Error marshalling with block number [%d]: %s", block.Header.Number, err)
	}
}



func (jl *jsonLedger) blockFilename(number uint64) string {
	return filepath.Join(jl.directory, fmt.Sprintf(blockFileFormatString, number))
}
