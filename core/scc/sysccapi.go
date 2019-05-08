/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"errors"
	"fmt"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/peer"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var sysccLogger = flogging.MustGetLogger("sccapi")


type Registrar interface {
	
	Register(ccid ccintf.CCID, cc shim.Chaincode) error
}




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


func (p *Provider) registerSysCC(syscc SelfDescribingSysCC) (bool, error) {
	if !syscc.Enabled() || !p.isWhitelisted(syscc) {
		sysccLogger.Info(fmt.Sprintf("system chaincode (%s,%s,%t) disabled", syscc.Name(), syscc.Path(), syscc.Enabled()))
		return false, nil
	}

	
	version := util.GetSysCCVersion()

	ccid := ccintf.CCID(syscc.Name() + ":" + version)
	err := p.Registrar.Register(ccid, syscc.Chaincode())
	if err != nil {
		
		if _, ok := err.(inproccontroller.SysCCRegisteredErr); !ok {
			errStr := fmt.Sprintf("could not register (%s,%v): %s", syscc.Path(), syscc, err)
			sysccLogger.Error(errStr)
			return false, fmt.Errorf(errStr)
		}
	}

	sysccLogger.Infof("system chaincode %s(%s) registered", syscc.Name(), syscc.Path())
	return true, err
}


func (p *Provider) deploySysCC(chainID string, ccprov ccprovider.ChaincodeProvider, syscc SelfDescribingSysCC) error {
	if !syscc.Enabled() || !p.isWhitelisted(syscc) {
		sysccLogger.Info(fmt.Sprintf("system chaincode (%s,%s) disabled", syscc.Name(), syscc.Path()))
		return nil
	}

	txid := util.GenerateUUID()

	
	
	
	txParams := &ccprovider.TransactionParams{
		TxID:      txid,
		ChannelID: chainID,
	}

	if chainID != "" {
		lgr := peer.GetLedger(chainID)
		if lgr == nil {
			panic(fmt.Sprintf("syschain %s start up failure - unexpected nil ledger for channel %s", syscc.Name(), chainID))
		}

		txsim, err := lgr.NewTxSimulator(txid)
		if err != nil {
			return err
		}

		txParams.TXSimulator = txsim
		defer txsim.Done()
	}

	chaincodeID := &pb.ChaincodeID{Path: syscc.Path(), Name: syscc.Name()}
	spec := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: syscc.InitArgs()}}

	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ExecEnv: pb.ChaincodeDeploymentSpec_SYSTEM, ChaincodeSpec: spec}

	
	version := util.GetSysCCVersion()

	cccid := &ccprovider.CCContext{
		Name:    chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeId.Name,
		Version: version,
	}

	resp, _, err := ccprov.ExecuteLegacyInit(txParams, cccid, chaincodeDeploymentSpec)
	if err == nil && resp.Status != shim.OK {
		err = errors.New(resp.Message)
	}

	sysccLogger.Infof("system chaincode %s/%s(%s) deployed", syscc.Name(), chainID, syscc.Path())

	return err
}


func deDeploySysCC(chainID string, ccprov ccprovider.ChaincodeProvider, syscc SelfDescribingSysCC) error {
	
	version := util.GetSysCCVersion()

	ccci := &ccprovider.ChaincodeContainerInfo{
		Type:          "GOLANG",
		Name:          syscc.Name(),
		Path:          syscc.Path(),
		Version:       version,
		ContainerType: inproccontroller.ContainerType,
	}

	err := ccprov.Stop(ccci)

	return err
}

func (p *Provider) isWhitelisted(syscc SelfDescribingSysCC) bool {
	enabled, ok := p.Whitelist[syscc.Name()]
	return ok && enabled
}
