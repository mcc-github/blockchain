/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"github.com/mcc-github/blockchain/core/common/ccprovider"
)



func (p *Provider) DeploySysCCs(chainID string, ccp ccprovider.ChaincodeProvider) {
	for _, sysCC := range p.SysCCs {
		p.deploySysCC(chainID, ccp, sysCC)
	}
}




func (p *Provider) DeDeploySysCCs(chainID string, ccp ccprovider.ChaincodeProvider) {
	for _, sysCC := range p.SysCCs {
		deDeploySysCC(chainID, ccp, sysCC)
	}
}
