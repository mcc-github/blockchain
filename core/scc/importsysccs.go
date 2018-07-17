/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	
	"github.com/mcc-github/blockchain/core/aclmgmt"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/scc/cscc"
	"github.com/mcc-github/blockchain/core/scc/lscc"
	"github.com/mcc-github/blockchain/core/scc/qscc"
)

func builtInSystemChaincodes(ccp ccprovider.ChaincodeProvider, p *Provider, aclProvider aclmgmt.ACLProvider) []*SystemChaincode {
	return []*SystemChaincode{
		{
			Enabled:           true,
			Name:              "cscc",
			Path:              "github.com/mcc-github/blockchain/core/scc/cscc",
			InitArgs:          nil,
			Chaincode:         cscc.New(ccp, p, aclProvider),
			InvokableExternal: true, 
		},
		{
			Enabled:           true,
			Name:              "lscc",
			Path:              "github.com/mcc-github/blockchain/core/scc/lscc",
			InitArgs:          nil,
			Chaincode:         lscc.New(p, aclProvider),
			InvokableExternal: true, 
			InvokableCC2CC:    true, 
		},
		{
			Enabled:           true,
			Name:              "qscc",
			Path:              "github.com/mcc-github/blockchain/core/scc/qscc",
			InitArgs:          nil,
			Chaincode:         qscc.New(aclProvider),
			InvokableExternal: true, 
			InvokableCC2CC:    true, 
		},
	}
}



func (p *Provider) DeploySysCCs(chainID string, ccp ccprovider.ChaincodeProvider) {
	for _, sysCC := range p.SysCCs {
		sysCC.deploySysCC(chainID, ccp)
	}
}




func (p *Provider) DeDeploySysCCs(chainID string, ccp ccprovider.ChaincodeProvider) {
	for _, sysCC := range p.SysCCs {
		sysCC.deDeploySysCC(chainID, ccp)
	}
}
