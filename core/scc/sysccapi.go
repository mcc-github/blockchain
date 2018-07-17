/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/peer"

	"github.com/spf13/viper"

	pb "github.com/mcc-github/blockchain/protos/peer"
)

var sysccLogger = flogging.MustGetLogger("sccapi")


type Registrar interface {
	
	Register(ccid *ccintf.CCID, cc shim.Chaincode) error
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


func (p *Provider) registerSysCC(syscc *SystemChaincode) (bool, error) {
	if !syscc.Enabled || !syscc.isWhitelisted() {
		sysccLogger.Info(fmt.Sprintf("system chaincode (%s,%s,%t) disabled", syscc.Name, syscc.Path, syscc.Enabled))
		return false, nil
	}

	
	version := util.GetSysCCVersion()

	ccid := &ccintf.CCID{
		Name:    syscc.Name,
		Version: version,
	}
	err := p.Registrar.Register(ccid, syscc.Chaincode)
	if err != nil {
		
		if _, ok := err.(inproccontroller.SysCCRegisteredErr); !ok {
			errStr := fmt.Sprintf("could not register (%s,%v): %s", syscc.Path, syscc, err)
			sysccLogger.Error(errStr)
			return false, fmt.Errorf(errStr)
		}
	}

	sysccLogger.Infof("system chaincode %s(%s) registered", syscc.Name, syscc.Path)
	return true, err
}


func (syscc *SystemChaincode) deploySysCC(chainID string, ccprov ccprovider.ChaincodeProvider) error {
	if !syscc.Enabled || !syscc.isWhitelisted() {
		sysccLogger.Info(fmt.Sprintf("system chaincode (%s,%s) disabled", syscc.Name, syscc.Path))
		return nil
	}

	txid := util.GenerateUUID()

	ctxt := context.Background()
	if chainID != "" {
		lgr := peer.GetLedger(chainID)
		if lgr == nil {
			panic(fmt.Sprintf("syschain %s start up failure - unexpected nil ledger for channel %s", syscc.Name, chainID))
		}

		
		
		ctxt2, txsim, err := ccprov.GetContext(lgr, txid)
		if err != nil {
			return err
		}

		ctxt = ctxt2

		defer txsim.Done()
	}

	chaincodeID := &pb.ChaincodeID{Path: syscc.Path, Name: syscc.Name}
	spec := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: syscc.InitArgs}}

	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ExecEnv: pb.ChaincodeDeploymentSpec_SYSTEM, ChaincodeSpec: spec}

	
	version := util.GetSysCCVersion()

	cccid := ccprovider.NewCCContext(chainID, chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeId.Name, version, txid, true, nil, nil)

	resp, _, err := ccprov.Execute(ctxt, cccid, chaincodeDeploymentSpec)
	if err == nil && resp.Status != shim.OK {
		err = errors.New(resp.Message)
	}

	sysccLogger.Infof("system chaincode %s/%s(%s) deployed", syscc.Name, chainID, syscc.Path)

	return err
}


func (syscc *SystemChaincode) deDeploySysCC(chainID string, ccprov ccprovider.ChaincodeProvider) error {
	chaincodeID := &pb.ChaincodeID{Path: syscc.Path, Name: syscc.Name}
	spec := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]), ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: syscc.InitArgs}}

	ctx := context.Background()
	
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ExecEnv: pb.ChaincodeDeploymentSpec_SYSTEM, ChaincodeSpec: spec}

	
	version := util.GetSysCCVersion()

	cccid := ccprovider.NewCCContext(chainID, syscc.Name, version, "", true, nil, nil)

	err := ccprov.Stop(ctx, cccid, chaincodeDeploymentSpec)

	return err
}

func (syscc *SystemChaincode) isWhitelisted() bool {
	chaincodes := viper.GetStringMapString("chaincode.system")
	val, ok := chaincodes[syscc.Name]
	enabled := val == "enable" || val == "true" || val == "yes"
	return ok && enabled
}
