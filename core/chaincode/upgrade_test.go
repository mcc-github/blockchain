/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chaincode

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
	putils "github.com/mcc-github/blockchain/protos/utils"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)


func getUpgradeLSCCSpec(chainID string, cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeInvocationSpec, error) {
	b, err := proto.Marshal(cds)
	if err != nil {
		return nil, err
	}

	
	lsccSpec := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: &pb.ChaincodeID{Name: "lscc"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("upgrade"), []byte(chainID), b}}}}

	return lsccSpec, nil
}


func upgrade(ctx context.Context, cccid *ccprovider.CCContext, spec *pb.ChaincodeSpec, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (*ccprovider.CCContext, error) {
	
	chaincodeDeploymentSpec, err := getDeploymentSpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	return upgrade2(ctx, cccid, chaincodeDeploymentSpec, blockNumber, chaincodeSupport)
}

func upgrade2(ctx context.Context, cccid *ccprovider.CCContext,
	chaincodeDeploymentSpec *pb.ChaincodeDeploymentSpec, blockNumber uint64, chaincodeSupport *ChaincodeSupport) (newcccid *ccprovider.CCContext, err error) {
	cis, err := getUpgradeLSCCSpec(cccid.ChainID, chaincodeDeploymentSpec)
	if err != nil {
		return nil, fmt.Errorf("Error creating lscc spec : %s\n", err)
	}

	uuid := util.GenerateUUID()
	cccid.TxID = uuid
	ctx, txsim, err := startTxSimulation(ctx, cccid.ChainID, cccid.TxID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get handle to simulator: %s ", err)
	}

	defer func() {
		
		if err == nil {
			
			err = endTxSimulationCDS(cccid.ChainID, uuid, txsim, []byte("upgraded"), true, chaincodeDeploymentSpec, blockNumber)
		} else {
			
			endTxSimulationCDS(cccid.ChainID, uuid, txsim, []byte("upgraded"), false, chaincodeDeploymentSpec, blockNumber)
		}
	}()

	
	ccprovider.PutChaincodeIntoFS(chaincodeDeploymentSpec)

	sysCCVers := util.GetSysCCVersion()
	sprop, prop := putils.MockSignedEndorserProposal2OrPanic(cccid.ChainID, cis.ChaincodeSpec, signer)
	lsccid := ccprovider.NewCCContext(cccid.ChainID, cis.ChaincodeSpec.ChaincodeId.Name, sysCCVers, uuid, true, sprop, prop)

	
	var resp *pb.Response
	if resp, _, err = chaincodeSupport.Execute(ctx, lsccid, cis); err != nil {
		return nil, fmt.Errorf("Error executing LSCC for upgrade: %s", err)
	}

	cdbytes := resp.Payload
	if cdbytes == nil {
		return nil, fmt.Errorf("Expected ChaincodeData back from LSCC but got nil")
	}

	cd := &ccprovider.ChaincodeData{}
	if err = proto.Unmarshal(cdbytes, cd); err != nil {
		return nil, fmt.Errorf("getting  ChaincodeData failed")
	}

	newVersion := string(cd.Version)
	if newVersion == cccid.Version {
		return nil, fmt.Errorf("Expected new version from LSCC but got same %s(%s)", newVersion, cccid.Version)
	}

	newcccid = ccprovider.NewCCContext(cccid.ChainID, chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeId.Name, newVersion, uuid, false, nil, nil)

	if _, _, err = chaincodeSupport.Execute(ctx, newcccid, chaincodeDeploymentSpec); err != nil {
		return nil, fmt.Errorf("Error deploying chaincode for upgrade: %s", err)
	}
	return
}










func TestUpgradeCC(t *testing.T) {
	testForSkip(t)
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	ccName := "mycc"
	url := "github.com/mcc-github/blockchain/examples/chaincode/go/chaincode_example01"
	chaincodeID := &pb.ChaincodeID{Name: ccName, Path: url, Version: "0"}

	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := ccprovider.NewCCContext(chainID, ccName, "0", "", false, nil, nil)
	var nextBlockNumber uint64 = 1
	_, err = deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)

	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		t.Fail()
		t.Logf("Error deploying chaincode %s(%s)", chaincodeID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	
	qArgs := util.ToChaincodeArgs("query", "a")

	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: qArgs}}

	
	_, _, _, err = invoke(ctxt, chainID, spec, nextBlockNumber, nil, chaincodeSupport)
	if err == nil {
		t.Fail()
		t.Logf("querying chaincode exampl01 should fail transaction: %s", err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	} else if !strings.Contains(err.Error(), "Invalid invoke function name. Expecting \"invoke\"") {
		t.Fail()
		t.Logf("expected <Invalid invoke function name. Expecting \"invoke\"> found <%s>", err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: chaincodeID}})
		return
	}

	
	
	url = "github.com/mcc-github/blockchain/examples/chaincode/go/example02/cmd"

	
	chaincodeID = &pb.ChaincodeID{Name: ccName, Path: url, Version: "1"}
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: args}}

	
	nextBlockNumber++
	cccid2, err := upgrade(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	if err != nil {
		t.Fail()
		t.Logf("Error upgrading chaincode %s(%s)", chaincodeID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		if cccid2 != nil {
			chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		}
		return
	}

	
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: qArgs}}
	nextBlockNumber++
	_, _, _, err = invokeWithVersion(ctxt, chainID, cccid2.Version, spec, nextBlockNumber, nil, chaincodeSupport)

	if err != nil {
		t.Fail()
		t.Logf("querying chaincode exampl02 did not succeed: %s", err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		return
	}
	chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
	chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}




func TestInvalUpgradeCC(t *testing.T) {
	testForSkip(t)
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	ccName := "mycc"
	url := "github.com/mcc-github/blockchain/examples/chaincode/go/example02/cmd"

	f := "init"
	args := util.ToChaincodeArgs(f, "a", "100", "b", "200")

	cccid := ccprovider.NewCCContext(chainID, ccName, "0", "", false, nil, nil)

	
	chaincodeID := &pb.ChaincodeID{Name: ccName, Path: url, Version: "1"}
	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: args}}

	
	var nextBlockNumber uint64
	cccid2, err := upgrade(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	if err == nil {
		t.Fail()
		t.Logf("Error expected upgrading to fail but it succeeded%s(%s)", chaincodeID, err)
		chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		if cccid2 != nil {
			chaincodeSupport.Stop(ctxt, cccid2, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
		}
		return
	}

	chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})
}
