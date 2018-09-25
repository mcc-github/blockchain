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

package testutils

import (
	"fmt"
	"os"

	"github.com/mcc-github/blockchain/common/crypto"
	mmsp "github.com/mcc-github/blockchain/common/mocks/msp"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	putils "github.com/mcc-github/blockchain/protos/utils"
)

var (
	signer msp.SigningIdentity
)

func init() {
	var err error
	
	err = msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		fmt.Printf("Could not load msp config, err %s", err)
		os.Exit(-1)
		return
	}
	signer, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not initialize msp/signer")
		return
	}
}


func ConstructBytesProposalResponsePayload(chainID string, ccid *pb.ChaincodeID, pResponse *pb.Response, simulationResults []byte) ([]byte, error) {
	ss, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	prop, _, err := putils.CreateChaincodeProposal(common.HeaderType_ENDORSER_TRANSACTION, chainID, &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccid}}, ss)
	if err != nil {
		return nil, err
	}

	presp, err := putils.CreateProposalResponse(prop.Header, prop.Payload, pResponse, simulationResults, nil, ccid, nil, signer)
	if err != nil {
		return nil, err
	}

	return presp.Payload, nil
}



func ConstructSignedTxEnvWithDefaultSigner(chainID string, ccid *pb.ChaincodeID, response *pb.Response, simulationResults []byte, txid string, events []byte, visibility []byte) (*common.Envelope, string, error) {
	return ConstructSignedTxEnv(chainID, ccid, response, simulationResults, txid, events, visibility, signer)
}


func ConstructSignedTxEnv(chainID string, ccid *pb.ChaincodeID, pResponse *pb.Response, simulationResults []byte, txid string, events []byte, visibility []byte, signer msp.SigningIdentity) (*common.Envelope, string, error) {
	ss, err := signer.Serialize()
	if err != nil {
		return nil, "", err
	}

	var prop *pb.Proposal
	if txid == "" {
		
		prop, txid, err = putils.CreateChaincodeProposal(common.HeaderType_ENDORSER_TRANSACTION, chainID, &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccid}}, ss)

	} else {
		
		nonce, err := crypto.GetRandomNonce()
		if err != nil {
			return nil, "", err
		}
		prop, txid, err = putils.CreateChaincodeProposalWithTxIDNonceAndTransient(txid, common.HeaderType_ENDORSER_TRANSACTION, chainID, &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{ChaincodeId: ccid}}, nonce, ss, nil)
	}
	if err != nil {
		return nil, "", err
	}

	presp, err := putils.CreateProposalResponse(prop.Header, prop.Payload, pResponse, simulationResults, nil, ccid, nil, signer)
	if err != nil {
		return nil, "", err
	}

	env, err := putils.CreateSignedTx(prop, signer, presp)
	if err != nil {
		return nil, "", err
	}
	return env, txid, nil
}

var mspLcl msp.MSP
var sigId msp.SigningIdentity


func ConstructUnsignedTxEnv(chainID string, ccid *pb.ChaincodeID, response *pb.Response, simulationResults []byte, txid string, events []byte, visibility []byte) (*common.Envelope, string, error) {
	if mspLcl == nil {
		mspLcl = mmsp.NewNoopMsp()
		sigId, _ = mspLcl.GetDefaultSigningIdentity()
	}

	return ConstructSignedTxEnv(chainID, ccid, response, simulationResults, txid, events, visibility, sigId)
}
