/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"testing"

	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/ledger/ledgermgmt"
	"github.com/stretchr/testify/assert"
)







type testhelper struct {
	*client
	*committer
	*verifier
	lgr    ledger.PeerLedger
	lgrid  string
	assert *assert.Assertions
}


func newTestHelperCreateLgr(id string, t *testing.T) *testhelper {
	genesisBlk, err := constructTestGenesisBlock(id)
	assert.NoError(t, err)
	lgr, err := ledgermgmt.CreateLedger(genesisBlk)
	assert.NoError(t, err)
	client, committer, verifier := newClient(lgr, t), newCommitter(lgr, t), newVerifier(lgr, t)
	return &testhelper{client, committer, verifier, lgr, id, assert.New(t)}
}


func newTestHelperOpenLgr(id string, t *testing.T) *testhelper {
	lgr, err := ledgermgmt.OpenLedger(id)
	assert.NoError(t, err)
	client, committer, verifier := newClient(lgr, t), newCommitter(lgr, t), newVerifier(lgr, t)
	return &testhelper{client, committer, verifier, lgr, id, assert.New(t)}
}



func (h *testhelper) cutBlockAndCommitWithPvtdata() *ledger.BlockAndPvtData {
	defer func() { h.simulatedTrans = nil }()
	return h.committer.cutBlockAndCommitWithPvtdata(h.simulatedTrans...)
}

func (h *testhelper) cutBlockAndCommitExpectError() (*ledger.BlockAndPvtData, error) {
	defer func() { h.simulatedTrans = nil }()
	return h.committer.cutBlockAndCommitExpectError(h.simulatedTrans...)
}



func (h *testhelper) assertError(output ...interface{}) {
	lastParam := output[len(output)-1]
	assert.NotNil(h.t, lastParam)
	h.assert.Error(lastParam.(error))
}


func (h *testhelper) assertNoError(output ...interface{}) {
	h.assert.Nil(output[len(output)-1])
}
