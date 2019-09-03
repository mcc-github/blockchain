/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ramledger

import (
	"testing"

	cb "github.com/mcc-github/blockchain-protos-go/common"
	ab "github.com/mcc-github/blockchain-protos-go/orderer"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/blockledger"
	genesisconfig "github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
	"github.com/mcc-github/blockchain/protoutil"
)

var genesisBlock = protoutil.NewBlock(0, nil)

func init() {
	flogging.ActivateSpec("common.ledger.blockledger.ram=DEBUG")
}

func newTestChain(maxSize int) *ramLedger {
	rlf := New(maxSize)
	chain, err := rlf.GetOrCreate(genesisconfig.TestChannelID)
	if err != nil {
		panic(err)
	}
	chain.Append(genesisBlock)
	return chain.(*ramLedger)
}



func TestAppend(t *testing.T) {
	maxSize := 3
	rl := newTestChain(maxSize)
	var blocks []*cb.Block
	for i := 0; i < 3; i++ {
		blocks = append(blocks, &cb.Block{Header: &cb.BlockHeader{Number: uint64(i + 1)}})
		rl.appendBlock(blocks[i])
	}
	item := rl.oldest
	for i := 0; i < 3; i++ {
		if item.block == nil {
			t.Fatalf("Block for item %d should not be nil", i)
		}
		if item.block.Header.Number != blocks[i].Header.Number {
			t.Errorf("Expected block %d to be %d but got %d", i, blocks[i].Header.Number, item.block.Header.Number)
		}
		if i != 2 && item.next == nil {
			t.Fatalf("Next item should not be nil")
		} else {
			item = item.next
		}
	}
}


func TestSignal(t *testing.T) {
	maxSize := 3
	rl := newTestChain(maxSize)
	item := rl.newest
	select {
	case <-item.signal:
		t.Fatalf("There is no successor, there should be no signal to continue")
	default:
	}
	rl.appendBlock(&cb.Block{Header: &cb.BlockHeader{Number: 1}})
	select {
	case <-item.signal:
	default:
		t.Fatalf("There is a successor, there should be a signal to continue")
	}
}




func TestTruncationSafety(t *testing.T) {
	maxSize := 3
	newBlocks := 10
	rl := newTestChain(maxSize)
	item := rl.newest
	for i := 0; i < newBlocks; i++ {
		rl.appendBlock(&cb.Block{Header: &cb.BlockHeader{Number: uint64(i + 1)}})
	}
	count := 0
	for item.next != nil {
		item = item.next
		count++
	}

	if count != newBlocks {
		t.Fatalf("The iterator should have found %d new blocks but found %d", newBlocks, count)
	}
}

func TestRetrieval(t *testing.T) {
	rl := newTestChain(3)
	rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}}))
	it, num := rl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Oldest{}})
	defer it.Close()
	if num != 0 {
		t.Fatalf("Expected genesis block iterator, but got %d", num)
	}

	block, status := it.Next()
	if status != cb.Status_SUCCESS {
		t.Fatalf("Expected to successfully read the genesis block")
	}
	if block.Header.Number != 0 {
		t.Fatalf("Expected to successfully retrieve the genesis block")
	}

	block, status = it.Next()
	if status != cb.Status_SUCCESS {
		t.Fatalf("Expected to successfully read the second block")
	}
	if block.Header.Number != 1 {
		t.Fatalf("Expected to successfully retrieve the second block but got block number %d", block.Header.Number)
	}
}

func TestBlockedRetrieval(t *testing.T) {
	rl := newTestChain(3)
	it, num := rl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 1}}})
	defer it.Close()
	if num != 1 {
		t.Fatalf("Expected block iterator at 1, but got %d", num)
	}

	rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}}))

	block, status := it.Next()
	if status != cb.Status_SUCCESS {
		t.Fatalf("Expected to successfully read the second block")
	}
	if block.Header.Number != 1 {
		t.Fatalf("Expected to successfully retrieve the second block")
	}
}

func TestIteratorPastEnd(t *testing.T) {
	rl := newTestChain(3)
	it, _ := rl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 2}}})
	defer it.Close()
	if _, status := it.Next(); status != cb.Status_NOT_FOUND {
		t.Fatalf("Expected block with status NOT_FOUND, but got %d", status)
	}
}

func TestIteratorOldest(t *testing.T) {
	rl := newTestChain(3)
	
	rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}}))
	rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}}))
	rl.Append(blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}}))
	it, num := rl.Iterator(&ab.SeekPosition{Type: &ab.SeekPosition_Specified{Specified: &ab.SeekSpecified{Number: 1}}})
	defer it.Close()
	if num != 1 {
		t.Fatalf("Expected block iterator at 1, but got %d", num)
	}
}

func TestAppendBadBLock(t *testing.T) {
	rl := newTestChain(3)
	t.Run("BadBlockNumber", func(t *testing.T) {
		nextBlock := blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}})
		nextBlock.Header.Number = nextBlock.Header.Number + 1
		if err := rl.Append(nextBlock); err == nil {
			t.Fatalf("Expected Append to fail.")
		}
	})
	t.Run("BadPreviousHash", func(t *testing.T) {
		nextBlock := blockledger.CreateNextBlock(rl, []*cb.Envelope{{Payload: []byte("My Data")}})
		nextBlock.Header.PreviousHash = []byte("bad hash")
		if err := rl.Append(nextBlock); err == nil {
			t.Fatalf("Expected Append to fail.")
		}
	})
}
