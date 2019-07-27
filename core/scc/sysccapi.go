/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
)

var sysccLogger = flogging.MustGetLogger("sccapi")




type SystemChaincode struct {
	
	Name string

	
	Path string

	
	InitArgs [][]byte

	
	Chaincode shim.Chaincode

	
	
	
	InvokableExternal bool

	
	
	
	
	InvokableCC2CC bool

	
	
	Enabled bool
}

type SysCCWrapper struct {
	SCC *SystemChaincode
}

func (sccw *SysCCWrapper) Name() string              { return sccw.SCC.Name }
func (sccw *SysCCWrapper) Path() string              { return sccw.SCC.Path }
func (sccw *SysCCWrapper) InitArgs() [][]byte        { return sccw.SCC.InitArgs }
func (sccw *SysCCWrapper) Chaincode() shim.Chaincode { return sccw.SCC.Chaincode }
func (sccw *SysCCWrapper) InvokableExternal() bool   { return sccw.SCC.InvokableExternal }
func (sccw *SysCCWrapper) InvokableCC2CC() bool      { return sccw.SCC.InvokableCC2CC }
func (sccw *SysCCWrapper) Enabled() bool             { return sccw.SCC.Enabled }

type SelfDescribingSysCC interface {
	
	Name() string

	
	Path() string

	
	InitArgs() [][]byte

	
	Chaincode() shim.Chaincode

	
	
	
	InvokableExternal() bool

	
	
	
	
	InvokableCC2CC() bool

	
	
	Enabled() bool
}
