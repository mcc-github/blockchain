/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package rwsetutil

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/ledger/rwset/kvrwset"
	"github.com/mcc-github/blockchain/bccsp"
	bccspfactory "github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/pkg/errors"
)


type MerkleTreeLevel uint32


type Hash []byte

const (
	leafLevel = MerkleTreeLevel(1)
)

var (
	hashOpts = &bccsp.SHA256Opts{}
)





















type RangeQueryResultsHelper struct {
	pendingResults []*kvrwset.KVRead
	mt             *merkleTree
	maxDegree      uint32
	hashingEnabled bool
}


func NewRangeQueryResultsHelper(enableHashing bool, maxDegree uint32) (*RangeQueryResultsHelper, error) {
	helper := &RangeQueryResultsHelper{pendingResults: nil,
		hashingEnabled: enableHashing,
		maxDegree:      maxDegree,
		mt:             nil}
	if enableHashing {
		var err error
		if helper.mt, err = newMerkleTree(maxDegree); err != nil {
			return nil, err
		}
	}
	return helper, nil
}




func (helper *RangeQueryResultsHelper) AddResult(kvRead *kvrwset.KVRead) error {
	logger.Debug("Adding a result")
	helper.pendingResults = append(helper.pendingResults, kvRead)
	if helper.hashingEnabled && uint32(len(helper.pendingResults)) > helper.maxDegree {
		logger.Debug("Processing the accumulated results")
		if err := helper.processPendingResults(); err != nil {
			return err
		}
	}
	return nil
}






func (helper *RangeQueryResultsHelper) Done() ([]*kvrwset.KVRead, *kvrwset.QueryReadsMerkleSummary, error) {
	
	
	if !helper.hashingEnabled || helper.mt.isEmpty() {
		return helper.pendingResults, nil, nil
	}
	if len(helper.pendingResults) != 0 {
		logger.Debug("Processing the pending results")
		if err := helper.processPendingResults(); err != nil {
			return helper.pendingResults, nil, err
		}
	}
	helper.mt.done()
	return helper.pendingResults, helper.mt.getSummery(), nil
}





func (helper *RangeQueryResultsHelper) GetMerkleSummary() *kvrwset.QueryReadsMerkleSummary {
	if !helper.hashingEnabled {
		return nil
	}
	return helper.mt.getSummery()
}

func (helper *RangeQueryResultsHelper) processPendingResults() error {
	var b []byte
	var err error
	if b, err = serializeKVReads(helper.pendingResults); err != nil {
		return err
	}
	helper.pendingResults = nil
	hash, err := bccspfactory.GetDefault().Hash(b, hashOpts)
	if err != nil {
		return err
	}
	helper.mt.update(hash)
	return nil
}

func serializeKVReads(kvReads []*kvrwset.KVRead) ([]byte, error) {
	return proto.Marshal(&kvrwset.QueryReads{KvReads: kvReads})
}



type merkleTree struct {
	tree      map[MerkleTreeLevel][]Hash
	maxLevel  MerkleTreeLevel
	maxDegree uint32
}

func newMerkleTree(maxDegree uint32) (*merkleTree, error) {
	if maxDegree < 2 {
		return nil, errors.Errorf("maxDegree [%d] should not be less than 2 in the merkle tree", maxDegree)
	}
	return &merkleTree{make(map[MerkleTreeLevel][]Hash), 1, maxDegree}, nil
}




func (m *merkleTree) update(nextLeafLevelHash Hash) error {
	logger.Debugf("Before update() = %s", m)
	defer logger.Debugf("After update() = %s", m)
	m.tree[leafLevel] = append(m.tree[leafLevel], nextLeafLevelHash)
	currentLevel := leafLevel
	for {
		currentLevelHashes := m.tree[currentLevel]
		if uint32(len(currentLevelHashes)) <= m.maxDegree {
			return nil
		}
		nextLevelHash, err := computeCombinedHash(currentLevelHashes)
		if err != nil {
			return err
		}
		delete(m.tree, currentLevel)
		nextLevel := currentLevel + 1
		m.tree[nextLevel] = append(m.tree[nextLevel], nextLevelHash)
		if nextLevel > m.maxLevel {
			m.maxLevel = nextLevel
		}
		currentLevel = nextLevel
	}
}




func (m *merkleTree) done() error {
	logger.Debugf("Before done() = %s", m)
	defer logger.Debugf("After done() = %s", m)
	currentLevel := leafLevel
	var h Hash
	var err error
	for currentLevel < m.maxLevel {
		currentLevelHashes := m.tree[currentLevel]
		switch len(currentLevelHashes) {
		case 0:
			currentLevel++
			continue
		case 1:
			h = currentLevelHashes[0]
		default:
			if h, err = computeCombinedHash(currentLevelHashes); err != nil {
				return err
			}
		}
		delete(m.tree, currentLevel)
		currentLevel++
		m.tree[currentLevel] = append(m.tree[currentLevel], h)
	}

	finalHashes := m.tree[m.maxLevel]
	if uint32(len(finalHashes)) > m.maxDegree {
		delete(m.tree, m.maxLevel)
		m.maxLevel++
		combinedHash, err := computeCombinedHash(finalHashes)
		if err != nil {
			return err
		}
		m.tree[m.maxLevel] = []Hash{combinedHash}
	}
	return nil
}

func (m *merkleTree) getSummery() *kvrwset.QueryReadsMerkleSummary {
	return &kvrwset.QueryReadsMerkleSummary{MaxDegree: m.maxDegree,
		MaxLevel:       uint32(m.getMaxLevel()),
		MaxLevelHashes: hashesToBytes(m.getMaxLevelHashes())}
}

func (m *merkleTree) getMaxLevel() MerkleTreeLevel {
	return m.maxLevel
}

func (m *merkleTree) getMaxLevelHashes() []Hash {
	return m.tree[m.maxLevel]
}

func (m *merkleTree) isEmpty() bool {
	return m.maxLevel == 1 && len(m.tree[m.maxLevel]) == 0
}

func (m *merkleTree) String() string {
	return fmt.Sprintf("tree := %#v", m.tree)
}

func computeCombinedHash(hashes []Hash) (Hash, error) {
	combinedHash := []byte{}
	for _, h := range hashes {
		combinedHash = append(combinedHash, h...)
	}
	return bccspfactory.GetDefault().Hash(combinedHash, hashOpts)
}

func hashesToBytes(hashes []Hash) [][]byte {
	b := [][]byte{}
	for _, hash := range hashes {
		b = append(b, hash)
	}
	return b
}
