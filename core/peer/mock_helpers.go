/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package peer

import (
	configtxtest "github.com/mcc-github/blockchain/common/configtx/test"
	mockchannelconfig "github.com/mcc-github/blockchain/common/mocks/config"
	mockconfigtx "github.com/mcc-github/blockchain/common/mocks/configtx"
	mockpolicies "github.com/mcc-github/blockchain/common/mocks/policies"
	"github.com/mcc-github/blockchain/core/ledger"
)

func CreateMockChannel(p *Peer, cid string) error {
	var ledger ledger.PeerLedger
	var err error

	if ledger = p.GetLedger(cid); ledger == nil {
		gb, _ := configtxtest.MakeGenesisBlock(cid)
		if ledger, err = p.LedgerMgr.CreateLedger(gb); err != nil {
			return err
		}
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.channels == nil {
		p.channels = map[string]*Channel{}
	}

	p.channels[cid] = &Channel{
		ledger: ledger,
		resources: &mockchannelconfig.Resources{
			PolicyManagerVal: &mockpolicies.Manager{
				Policy: &mockpolicies.Policy{},
			},
			ConfigtxValidatorVal: &mockconfigtx.Validator{},
			ApplicationConfigVal: &mockchannelconfig.MockApplication{CapabilitiesRv: &mockchannelconfig.MockApplicationCapabilities{}},
		},
	}

	return nil
}
