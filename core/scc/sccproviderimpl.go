/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/pkg/errors"
)


type Provider struct {
	Peer      *peer.Peer
	SysCCs    []SelfDescribingSysCC
	Whitelist Whitelist
}


func (p *Provider) RegisterSysCC(scc SelfDescribingSysCC) error {
	for _, registeredSCC := range p.SysCCs {
		if scc.Name() == registeredSCC.Name() {
			return errors.Errorf("chaincode with name '%s' already registered", scc.Name())
		}
	}
	p.SysCCs = append(p.SysCCs, scc)
	return nil
}


func (p *Provider) IsSysCC(name string) bool {
	for _, sysCC := range p.SysCCs {
		if sysCC.Name() == name {
			return true
		}
	}
	if isDeprecatedSysCC(name) {
		return true
	}
	return false
}




func (p *Provider) IsSysCCAndNotInvokableCC2CC(name string) bool {
	for _, sysCC := range p.SysCCs {
		if sysCC.Name() == name {
			return !sysCC.InvokableCC2CC()
		}
	}

	if isDeprecatedSysCC(name) {
		return true
	}

	return false
}


func (p *Provider) GetQueryExecutorForLedger(cid string) (ledger.QueryExecutor, error) {
	l := p.Peer.GetLedger(cid)
	if l == nil {
		return nil, fmt.Errorf("Could not retrieve ledger for channel %s", cid)
	}

	return l.NewQueryExecutor()
}




func (p *Provider) IsSysCCAndNotInvokableExternal(name string) bool {
	for _, sysCC := range p.SysCCs {
		if sysCC.Name() == name {
			return !sysCC.InvokableExternal()
		}
	}

	if isDeprecatedSysCC(name) {
		return true
	}

	return false
}



func (p *Provider) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return p.Peer.GetApplicationConfig(cid)
}



func (p *Provider) PolicyManager(channelID string) (policies.Manager, bool) {
	m := p.Peer.GetPolicyManager(channelID)
	return m, (m != nil)
}

func isDeprecatedSysCC(name string) bool {
	return name == "vscc" || name == "escc"
}

func (p *Provider) isWhitelisted(syscc SelfDescribingSysCC) bool {
	enabled, ok := p.Whitelist[syscc.Name()]
	return ok && enabled
}
