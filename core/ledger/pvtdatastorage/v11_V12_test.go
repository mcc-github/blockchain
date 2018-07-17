/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)






func TestV11v12(t *testing.T) {
	testWorkingDir := "test-working-dir"
	testutil.CopyDir("testdata/v11_v12/ledgersData", testWorkingDir)
	defer os.RemoveAll(testWorkingDir)

	viper.Set("peer.fileSystemPath", testWorkingDir)
	defer viper.Reset()

	ledgerid := "ch1"
	cs := btltestutil.NewMockCollectionStore()
	cs.SetBTL("marbles_private", "collectionMarbles", 0)
	cs.SetBTL("marbles_private", "collectionMarblePrivateDetails", 0)
	btlPolicy := pvtdatapolicy.ConstructBTLPolicy(cs)

	p := NewProvider()
	defer p.Close()
	s, err := p.OpenStore(ledgerid)
	assert.NoError(t, err)
	s.Init(btlPolicy)

	for blk := 0; blk < 10; blk++ {
		checkDataNotExists(t, s, blk)
	}
	checkDataExists(t, s, 10)
	for blk := 11; blk < 14; blk++ {
		checkDataNotExists(t, s, blk)
	}
	checkDataExists(t, s, 14)

	_, err = s.GetPvtDataByBlockNum(uint64(15), nil)
	_, ok := err.(*ErrOutOfRange)
	assert.True(t, ok)
}

func checkDataNotExists(t *testing.T, s Store, blkNum int) {
	data, err := s.GetPvtDataByBlockNum(uint64(blkNum), nil)
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func checkDataExists(t *testing.T, s Store, blkNum int) {
	data, err := s.GetPvtDataByBlockNum(uint64(blkNum), nil)
	assert.NoError(t, err)
	assert.NotNil(t, data)
	t.Logf("pvtdata = %s\n", spew.Sdump(data))
}
