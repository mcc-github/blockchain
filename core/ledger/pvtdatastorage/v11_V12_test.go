/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pvtdatastorage

import (
	"io/ioutil"
	"path/filepath"

	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mcc-github/blockchain/common/ledger/testutil"
	"github.com/mcc-github/blockchain/core/ledger"
	btltestutil "github.com/mcc-github/blockchain/core/ledger/pvtdatapolicy/testutil"
	"github.com/stretchr/testify/assert"
)






func TestV11v12(t *testing.T) {
	testWorkingDir, err := ioutil.TempDir("", "pdstore")
	if err != nil {
		t.Fatalf("Failed to create private data storage directory: %s", err)
	}
	defer os.RemoveAll(testWorkingDir)
	testutil.CopyDir("testdata/v11_v12/ledgersData/pvtdataStore", testWorkingDir, false)

	ledgerid := "ch1"
	btlPolicy := btltestutil.SampleBTLPolicy(
		map[[2]string]uint64{
			{"marbles_private", "collectionMarbles"}:              0,
			{"marbles_private", "collectionMarblePrivateDetails"}: 0,
		},
	)
	conf := &PrivateDataConfig{
		PrivateDataConfig: &ledger.PrivateDataConfig{
			BatchesInterval: 1000,
			MaxBatchSize:    5000,
			PurgeInterval:   100,
		},
		StorePath: filepath.Join(testWorkingDir, "pvtdataStore"),
	}
	p, err := NewProvider(conf)
	assert.NoError(t, err)
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
