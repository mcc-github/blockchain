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

package validation

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/mocks/config"
	mmsp "github.com/mcc-github/blockchain/common/mocks/msp"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/stretchr/testify/assert"
)

func getProposal(channel string) (*peer.Proposal, error) {
	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: getChaincodeID(),
			Type:        peer.ChaincodeSpec_GOLANG}}

	proposal, _, err := utils.CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, channel, cis, signerSerialized)
	return proposal, err
}

func getChaincodeID() *peer.ChaincodeID {
	return &peer.ChaincodeID{Name: "foo", Version: "v1"}
}

func createSignedTxTwoActions(proposal *peer.Proposal, signer msp.SigningIdentity, resps ...*peer.ProposalResponse) (*common.Envelope, error) {
	if len(resps) == 0 {
		return nil, fmt.Errorf("At least one proposal response is necessary")
	}

	
	hdr, err := utils.GetHeader(proposal.Header)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal header")
	}

	
	pPayl, err := utils.GetChaincodeProposalPayload(proposal.Payload)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal payload")
	}

	
	endorsements := make([]*peer.Endorsement, len(resps))
	for n, r := range resps {
		endorsements[n] = r.Endorsement
	}

	
	cea := &peer.ChaincodeEndorsedAction{ProposalResponsePayload: resps[0].Payload, Endorsements: endorsements}

	
	propPayloadBytes, err := utils.GetBytesProposalPayloadForTx(pPayl, nil)
	if err != nil {
		return nil, err
	}

	
	cap := &peer.ChaincodeActionPayload{ChaincodeProposalPayload: propPayloadBytes, Action: cea}
	capBytes, err := utils.GetBytesChaincodeActionPayload(cap)
	if err != nil {
		return nil, err
	}

	
	taa := &peer.TransactionAction{Header: hdr.SignatureHeader, Payload: capBytes}
	taas := make([]*peer.TransactionAction, 2)
	taas[0] = taa
	taas[1] = taa
	tx := &peer.Transaction{Actions: taas}

	
	txBytes, err := utils.GetBytesTransaction(tx)
	if err != nil {
		return nil, err
	}

	
	payl := &common.Payload{Header: hdr, Data: txBytes}
	paylBytes, err := utils.GetBytesPayload(payl)
	if err != nil {
		return nil, err
	}

	
	sig, err := signer.Sign(paylBytes)
	if err != nil {
		return nil, err
	}

	
	return &common.Envelope{Payload: paylBytes, Signature: sig}, nil
}

func TestGoodPath(t *testing.T) {
	
	prop, err := getProposal(util.GetTestChainID())
	if err != nil {
		t.Fatalf("getProposal failed, err %s", err)
		return
	}

	
	sProp, err := utils.GetSignedProposal(prop, signer)
	if err != nil {
		t.Fatalf("GetSignedProposal failed, err %s", err)
		return
	}

	
	_, _, _, err = ValidateProposalMessage(sProp)
	if err != nil {
		t.Fatalf("ValidateProposalMessage failed, err %s", err)
		return
	}

	response := &peer.Response{Status: 200}
	simRes := []byte("simulation_result")

	
	presp, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response, simRes, nil, getChaincodeID(), nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	
	tx, err := utils.CreateSignedTx(prop, signer, presp)
	if err != nil {
		t.Fatalf("CreateSignedTx failed, err %s", err)
		return
	}

	
	payl, txResult := ValidateTransaction(tx, &config.MockApplicationCapabilities{})
	if txResult != peer.TxValidationCode_VALID {
		t.Fatalf("ValidateTransaction failed, err %s", err)
		return
	}

	txx, err := utils.GetTransaction(payl.Data)
	if err != nil {
		t.Fatalf("GetTransaction failed, err %s", err)
		return
	}

	act := txx.Actions

	
	if len(act) != 1 {
		t.Fatalf("Ivalid number of TransactionAction, expected 1, got %d", len(act))
		return
	}

	
	_, simResBack, err := utils.GetPayloads(act[0])
	if err != nil {
		t.Fatalf("GetPayloads failed, err %s", err)
		return
	}

	
	if string(simRes) != string(simResBack.Results) {
		t.Fatal("Simulation results are different")
		return
	}
}

func TestTXWithTwoActionsRejected(t *testing.T) {
	
	prop, err := getProposal(util.GetTestChainID())
	if err != nil {
		t.Fatalf("getProposal failed, err %s", err)
		return
	}

	response := &peer.Response{Status: 200}
	simRes := []byte("simulation_result")

	
	presp, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response, simRes, nil, &peer.ChaincodeID{Name: "somename", Version: "someversion"}, nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	
	tx, err := createSignedTxTwoActions(prop, signer, presp)
	if err != nil {
		t.Fatalf("CreateSignedTx failed, err %s", err)
		return
	}

	
	_, txResult := ValidateTransaction(tx, &config.MockApplicationCapabilities{})
	if txResult == peer.TxValidationCode_VALID {
		t.Fatalf("ValidateTransaction should have failed")
		return
	}
}

func TestBadProp(t *testing.T) {
	
	prop, err := getProposal(util.GetTestChainID())
	if err != nil {
		t.Fatalf("getProposal failed, err %s", err)
		return
	}

	
	sProp, err := utils.GetSignedProposal(prop, signer)
	if err != nil {
		t.Fatalf("GetSignedProposal failed, err %s", err)
		return
	}

	
	sigOrig := sProp.Signature
	for i := 0; i < len(sigOrig); i++ {
		sigCopy := make([]byte, len(sigOrig))
		copy(sigCopy, sigOrig)
		sigCopy[i] = byte(int(sigCopy[i]+1) % 255)
		
		_, _, _, err = ValidateProposalMessage(&peer.SignedProposal{ProposalBytes: sProp.ProposalBytes, Signature: sigCopy})
		if err == nil {
			t.Fatal("ValidateProposalMessage should have failed")
			return
		}
	}

	
	sProp, err = utils.GetSignedProposal(prop, signer)
	if err != nil {
		t.Fatalf("GetSignedProposal failed, err %s", err)
		return
	}

	
	pbytesOrig := sProp.ProposalBytes
	for i := 0; i < len(pbytesOrig); i++ {
		pbytesCopy := make([]byte, len(pbytesOrig))
		copy(pbytesCopy, pbytesOrig)
		pbytesCopy[i] = byte(int(pbytesCopy[i]+1) % 255)
		
		_, _, _, err = ValidateProposalMessage(&peer.SignedProposal{ProposalBytes: pbytesCopy, Signature: sProp.Signature})
		if err == nil {
			t.Fatal("ValidateProposalMessage should have failed")
			return
		}
	}

	
	badSigner, err := mmsp.NewNoopMsp().GetDefaultSigningIdentity()
	if err != nil {
		t.Fatal("Couldn't get noop signer")
		return
	}

	
	sProp, err = utils.GetSignedProposal(prop, badSigner)
	if err != nil {
		t.Fatalf("GetSignedProposal failed, err %s", err)
		return
	}

	
	_, _, _, err = ValidateProposalMessage(sProp)
	if err == nil {
		t.Fatal("ValidateProposalMessage should have failed")
		return
	}
}

func corrupt(bytes []byte) {
	rand.Seed(time.Now().UnixNano())
	bytes[rand.Intn(len(bytes))]--
}

func TestBadTx(t *testing.T) {
	
	prop, err := getProposal(util.GetTestChainID())
	if err != nil {
		t.Fatalf("getProposal failed, err %s", err)
		return
	}

	response := &peer.Response{Status: 200}
	simRes := []byte("simulation_result")

	
	presp, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response, simRes, nil, getChaincodeID(), nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	
	tx, err := utils.CreateSignedTx(prop, signer, presp)
	if err != nil {
		t.Fatalf("CreateSignedTx failed, err %s", err)
		return
	}

	
	paylOrig := tx.Payload
	for i := 0; i < len(paylOrig); i++ {
		paylCopy := make([]byte, len(paylOrig))
		copy(paylCopy, paylOrig)
		paylCopy[i] = byte(int(paylCopy[i]+1) % 255)
		
		_, txResult := ValidateTransaction(&common.Envelope{Signature: tx.Signature, Payload: paylCopy}, &config.MockApplicationCapabilities{})
		if txResult == peer.TxValidationCode_VALID {
			t.Fatal("ValidateTransaction should have failed")
			return
		}
	}

	
	tx, err = utils.CreateSignedTx(prop, signer, presp)
	if err != nil {
		t.Fatalf("CreateSignedTx failed, err %s", err)
		return
	}

	
	corrupt(tx.Signature)

	
	_, txResult := ValidateTransaction(tx, &config.MockApplicationCapabilities{})
	if txResult == peer.TxValidationCode_VALID {
		t.Fatal("ValidateTransaction should have failed")
		return
	}
}

func Test2EndorsersAgree(t *testing.T) {
	
	prop, err := getProposal(util.GetTestChainID())
	if err != nil {
		t.Fatalf("getProposal failed, err %s", err)
		return
	}

	response1 := &peer.Response{Status: 200}
	simRes1 := []byte("simulation_result")

	
	presp1, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response1, simRes1, nil, getChaincodeID(), nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	response2 := &peer.Response{Status: 200}
	simRes2 := []byte("simulation_result")

	
	presp2, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response2, simRes2, nil, getChaincodeID(), nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	
	tx, err := utils.CreateSignedTx(prop, signer, presp1, presp2)
	if err != nil {
		t.Fatalf("CreateSignedTx failed, err %s", err)
		return
	}

	
	_, txResult := ValidateTransaction(tx, &config.MockApplicationCapabilities{})
	if txResult != peer.TxValidationCode_VALID {
		t.Fatalf("ValidateTransaction failed, err %s", err)
		return
	}
}

func Test2EndorsersDisagree(t *testing.T) {
	
	prop, err := getProposal(util.GetTestChainID())
	if err != nil {
		t.Fatalf("getProposal failed, err %s", err)
		return
	}

	response1 := &peer.Response{Status: 200}
	simRes1 := []byte("simulation_result1")

	
	presp1, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response1, simRes1, nil, getChaincodeID(), nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	response2 := &peer.Response{Status: 200}
	simRes2 := []byte("simulation_result2")

	
	presp2, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response2, simRes2, nil, getChaincodeID(), nil, signer)
	if err != nil {
		t.Fatalf("CreateProposalResponse failed, err %s", err)
		return
	}

	
	_, err = utils.CreateSignedTx(prop, signer, presp1, presp2)
	if err == nil {
		t.Fatal("CreateSignedTx should have failed")
		return
	}
}

func TestInvocationsBadArgs(t *testing.T) {
	_, code := ValidateTransaction(nil, &config.MockApplicationCapabilities{})
	assert.Equal(t, code, peer.TxValidationCode_NIL_ENVELOPE)
	err := validateEndorserTransaction(nil, nil)
	assert.Error(t, err)
	err = validateConfigTransaction(nil, nil)
	assert.Error(t, err)
	_, _, err = validateCommonHeader(nil)
	assert.Error(t, err)
	err = validateChannelHeader(nil)
	assert.Error(t, err)
	err = validateChannelHeader(&common.ChannelHeader{})
	assert.Error(t, err)
	err = validateSignatureHeader(nil)
	assert.Error(t, err)
	err = validateSignatureHeader(&common.SignatureHeader{})
	assert.Error(t, err)
	err = validateSignatureHeader(&common.SignatureHeader{Nonce: []byte("a")})
	assert.Error(t, err)
	err = checkSignatureFromCreator(nil, nil, nil, "")
	assert.Error(t, err)
	_, _, _, err = ValidateProposalMessage(nil)
	assert.Error(t, err)
	_, err = validateChaincodeProposalMessage(nil, nil)
	assert.Error(t, err)
	_, err = validateChaincodeProposalMessage(&peer.Proposal{}, &common.Header{ChannelHeader: []byte("a"), SignatureHeader: []byte("a")})
	assert.Error(t, err)
}

var signer msp.SigningIdentity
var signerSerialized []byte
var signerMSPId string

func TestMain(m *testing.M) {
	
	
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		fmt.Printf("Could not initialize msp, err %s", err)
		os.Exit(-1)
		return
	}

	signer, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		fmt.Println("Could not get signer")
		os.Exit(-1)
		return
	}
	signerMSPId = signer.GetMSPIdentifier()

	signerSerialized, err = signer.Serialize()
	if err != nil {
		fmt.Println("Could not serialize identity")
		os.Exit(-1)
		return
	}

	os.Exit(m.Run())
}
