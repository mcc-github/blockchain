/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msptesttools

import (
	"testing"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp/mgmt"
)

func TestFakeSetup(t *testing.T) {
	err := LoadMSPSetupForTesting()
	if err != nil {
		t.Fatalf("LoadLocalMsp failed, err %s", err)
	}

	_, err = mgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		t.Fatalf("GetDefaultSigningIdentity failed, err %s", err)
	}

	msps, err := mgmt.GetManagerForChain(util.GetTestChannelID()).GetMSPs()
	if err != nil {
		t.Fatalf("EnlistedMSPs failed, err %s", err)
	}

	if msps == nil || len(msps) == 0 {
		t.Fatalf("There are no MSPS in the manager for chain %s", util.GetTestChannelID())
	}
}
