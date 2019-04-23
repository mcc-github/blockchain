/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"os"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
)

func TestMain(m *testing.M) {
	flogging.ActivateSpec("lockbasedtxmgr,statevalidator,statebasedval,statecouchdb,valimpl,pvtstatepurgemgmt,confighistory,kvledger=debug")
	os.Exit(m.Run())
}

func TestLedgerAPIs(t *testing.T) {
	env := newEnv(t)
	defer env.cleanup()
	env.initLedgerMgmt()

	
	h1 := newTestHelperCreateLgr("ledger1", t)
	h2 := newTestHelperCreateLgr("ledger2", t)

	
	dataHelper := newSampleDataHelper(t)
	dataHelper.populateLedger(h1)
	dataHelper.populateLedger(h2)

	
	dataHelper.verifyLedgerContent(h1)
	dataHelper.verifyLedgerContent(h2)
}
