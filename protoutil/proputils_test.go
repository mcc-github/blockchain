/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoutil_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

func createCIS() *pb.ChaincodeInvocationSpec {
	return &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_GOLANG,
			ChaincodeId: &pb.ChaincodeID{Name: "chaincode_name"},
			Input:       &pb.ChaincodeInput{Args: [][]byte{[]byte("arg1"), []byte("arg2")}}}}
}

func TestNilProposal(t *testing.T) {
	
	_, err := protoutil.GetChaincodeInvocationSpec(nil)
	assert.Error(t, err, "Expected error with nil proposal")
	_, _, err = protoutil.GetChaincodeProposalContext(nil)
	assert.Error(t, err, "Expected error with nil proposal")
	_, err = protoutil.GetNonce(nil)
	assert.Error(t, err, "Expected error with nil proposal")
	_, err = protoutil.GetBytesProposal(nil)
	assert.Error(t, err, "Expected error with nil proposal")
	_, err = protoutil.ComputeProposalBinding(nil)
	assert.Error(t, err, "Expected error with nil proposal")
}

func TestBadProposalHeaders(t *testing.T) {
	
	
	

	
	prop := &pb.Proposal{
		Header: []byte{},
	}
	_, _, err := protoutil.GetChaincodeProposalContext(prop)
	assert.Error(t, err, "Expected error with empty proposal header")
	_, err = protoutil.ComputeProposalBinding(prop)
	assert.Error(t, err, "Expected error with empty proposal header")

	
	prop = &pb.Proposal{
		Header: []byte("header"),
	}
	_, _, err = protoutil.GetChaincodeProposalContext(prop)
	assert.Error(t, err, "Expected error with empty proposal payload")

	
	prop = &pb.Proposal{
		Header:  []byte("bad header"),
		Payload: []byte("payload"),
	}
	_, err = protoutil.GetHeader(prop.Header)
	assert.Error(t, err, "Expected error with malformed proposal header")
	_, err = protoutil.GetChaincodeInvocationSpec(prop)
	assert.Error(t, err, "Expected error with malformed proposal header")
	_, _, err = protoutil.GetChaincodeProposalContext(prop)
	assert.Error(t, err, "Expected error with malformed proposal header")
	_, err = protoutil.GetNonce(prop)
	assert.Error(t, err, "Expected error with malformed proposal header")
	_, err = protoutil.ComputeProposalBinding(prop)
	assert.Error(t, err, "Expected error with malformed proposal header")

	
	chdr, _ := proto.Marshal(&common.ChannelHeader{
		Type: int32(common.HeaderType_ENDORSER_TRANSACTION),
	})
	hdr := &common.Header{
		ChannelHeader:   chdr,
		SignatureHeader: []byte("bad signature header"),
	}
	_, err = protoutil.GetSignatureHeader(hdr.SignatureHeader)
	assert.Error(t, err, "Expected error with malformed signature header")
	hdrBytes, _ := proto.Marshal(hdr)
	prop.Header = hdrBytes
	_, err = protoutil.GetChaincodeInvocationSpec(prop)
	assert.Error(t, err, "Expected error with malformed signature header")
	_, _, err = protoutil.GetChaincodeProposalContext(prop)
	assert.Error(t, err, "Expected error with malformed signature header")
	_, err = protoutil.GetNonce(prop)
	assert.Error(t, err, "Expected error with malformed signature header")
	_, err = protoutil.ComputeProposalBinding(prop)
	assert.Error(t, err, "Expected error with malformed signature header")

	
	chdr, _ = proto.Marshal(&common.ChannelHeader{
		Type: int32(common.HeaderType_DELIVER_SEEK_INFO),
	})
	hdr.ChannelHeader = chdr
	hdrBytes, _ = proto.Marshal(hdr)
	prop.Header = hdrBytes
	_, _, err = protoutil.GetChaincodeProposalContext(prop)
	assert.Error(t, err, "Expected error with wrong header type")
	assert.Contains(t, err.Error(), "invalid proposal: invalid channel header type")
	_, err = protoutil.GetNonce(prop)
	assert.Error(t, err, "Expected error with wrong header type")

	
	hdr.ChannelHeader = []byte("bad channel header")
	hdrBytes, _ = proto.Marshal(hdr)
	prop.Header = hdrBytes
	_, _, err = protoutil.GetChaincodeProposalContext(prop)
	assert.Error(t, err, "Expected error with malformed channel header")
	_, err = protoutil.GetNonce(prop)
	assert.Error(t, err, "Expected error with malformed channel header")
	_, err = protoutil.GetChaincodeHeaderExtension(hdr)
	assert.Error(t, err, "Expected error with malformed channel header")
	_, err = protoutil.ComputeProposalBinding(prop)
	assert.Error(t, err, "Expected error with malformed channel header")

}

func TestGetNonce(t *testing.T) {
	chdr, _ := proto.Marshal(&common.ChannelHeader{
		Type: int32(common.HeaderType_ENDORSER_TRANSACTION),
	})
	hdr, _ := proto.Marshal(&common.Header{
		ChannelHeader:   chdr,
		SignatureHeader: []byte{},
	})
	prop := &pb.Proposal{
		Header: hdr,
	}
	_, err := protoutil.GetNonce(prop)
	assert.Error(t, err, "Expected error with nil signature header")

	shdr, _ := proto.Marshal(&common.SignatureHeader{
		Nonce: []byte("nonce"),
	})
	hdr, _ = proto.Marshal(&common.Header{
		ChannelHeader:   chdr,
		SignatureHeader: shdr,
	})
	prop = &pb.Proposal{
		Header: hdr,
	}
	nonce, err := protoutil.GetNonce(prop)
	assert.NoError(t, err, "Unexpected error getting nonce")
	assert.Equal(t, "nonce", string(nonce), "Failed to return the expected nonce")

}

func TestGetChaincodeDeploymentSpec(t *testing.T) {
	_, err := protoutil.GetChaincodeDeploymentSpec([]byte("bad spec"))
	assert.Error(t, err, "Expected error with malformed spec")

	cds, _ := proto.Marshal(&pb.ChaincodeDeploymentSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type: pb.ChaincodeSpec_GOLANG,
		},
	})
	_, err = protoutil.GetChaincodeDeploymentSpec(cds)
	assert.NoError(t, err, "Unexpected error getting deployment spec")
}

func TestCDSProposals(t *testing.T) {
	var prop *pb.Proposal
	var err error
	var txid string
	creator := []byte("creator")
	cds := &pb.ChaincodeDeploymentSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type: pb.ChaincodeSpec_GOLANG,
		},
	}
	policy := []byte("policy")
	escc := []byte("escc")
	vscc := []byte("vscc")
	chainID := "testchainid"

	
	prop, txid, err = protoutil.CreateInstallProposalFromCDS(cds, creator)
	assert.NotNil(t, prop, "Install proposal should not be nil")
	assert.NoError(t, err, "Unexpected error creating install proposal")
	assert.NotEqual(t, "", txid, "txid should not be empty")

	
	prop, txid, err = protoutil.CreateDeployProposalFromCDS(chainID, cds, creator, policy, escc, vscc, nil)
	assert.NotNil(t, prop, "Deploy proposal should not be nil")
	assert.NoError(t, err, "Unexpected error creating deploy proposal")
	assert.NotEqual(t, "", txid, "txid should not be empty")

	
	prop, txid, err = protoutil.CreateUpgradeProposalFromCDS(chainID, cds, creator, policy, escc, vscc, nil)
	assert.NotNil(t, prop, "Upgrade proposal should not be nil")
	assert.NoError(t, err, "Unexpected error creating upgrade proposal")
	assert.NotEqual(t, "", txid, "txid should not be empty")

}

func TestComputeProposalBinding(t *testing.T) {
	expectedDigestHex := "5093dd4f4277e964da8f4afbde0a9674d17f2a6a5961f0670fc21ae9b67f2983"
	expectedDigest, _ := hex.DecodeString(expectedDigestHex)
	chdr, _ := proto.Marshal(&common.ChannelHeader{
		Epoch: uint64(10),
	})
	shdr, _ := proto.Marshal(&common.SignatureHeader{
		Nonce:   []byte("nonce"),
		Creator: []byte("creator"),
	})
	hdr, _ := proto.Marshal(&common.Header{
		ChannelHeader:   chdr,
		SignatureHeader: shdr,
	})
	prop := &pb.Proposal{
		Header: hdr,
	}
	binding, _ := protoutil.ComputeProposalBinding(prop)
	assert.Equal(t, expectedDigest, binding, "Binding does not match expected digest")
}

func TestProposal(t *testing.T) {
	
	prop, _, err := protoutil.CreateChaincodeProposalWithTransient(
		common.HeaderType_ENDORSER_TRANSACTION,
		testChainID, createCIS(),
		[]byte("creator"),
		map[string][]byte{"certx": []byte("transient")})
	if err != nil {
		t.Fatalf("Could not create chaincode proposal, err %s\n", err)
		return
	}

	
	pBytes, err := protoutil.GetBytesProposal(prop)
	if err != nil {
		t.Fatalf("Could not serialize the chaincode proposal, err %s\n", err)
		return
	}

	
	propBack, err := protoutil.GetProposal(pBytes)
	if err != nil {
		t.Fatalf("Could not deserialize the chaincode proposal, err %s\n", err)
		return
	}
	if !proto.Equal(prop, propBack) {
		t.Fatalf("Proposal and deserialized proposals don't match\n")
		return
	}

	
	hdr, err := protoutil.GetHeader(prop.Header)
	if err != nil {
		t.Fatalf("Could not extract the header from the proposal, err %s\n", err)
	}

	hdrBytes, err := protoutil.GetBytesHeader(hdr)
	if err != nil {
		t.Fatalf("Could not marshal the header, err %s\n", err)
	}

	hdr, err = protoutil.GetHeader(hdrBytes)
	if err != nil {
		t.Fatalf("Could not unmarshal the header, err %s\n", err)
	}

	chdr, err := protoutil.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		t.Fatalf("Could not unmarshal channel header, err %s", err)
	}

	shdr, err := protoutil.GetSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		t.Fatalf("Could not unmarshal signature header, err %s", err)
	}

	_, err = protoutil.GetBytesSignatureHeader(shdr)
	if err != nil {
		t.Fatalf("Could not marshal signature header, err %s", err)
	}

	
	if chdr.Type != int32(common.HeaderType_ENDORSER_TRANSACTION) ||
		shdr.Nonce == nil ||
		string(shdr.Creator) != "creator" {
		t.Fatalf("Invalid header after unmarshalling\n")
		return
	}

	
	hdrExt, err := protoutil.GetChaincodeHeaderExtension(hdr)
	if err != nil {
		t.Fatalf("Could not extract the header extensions from the proposal, err %s\n", err)
		return
	}

	
	if string(hdrExt.ChaincodeId.Name) != "chaincode_name" {
		t.Fatalf("Invalid header extension after unmarshalling\n")
		return
	}

	
	cis, err := protoutil.GetChaincodeInvocationSpec(prop)
	if err != nil {
		t.Fatalf("Could not extract chaincode invocation spec from header, err %s\n", err)
		return
	}

	
	if cis.ChaincodeSpec.Type != pb.ChaincodeSpec_GOLANG ||
		cis.ChaincodeSpec.ChaincodeId.Name != "chaincode_name" ||
		len(cis.ChaincodeSpec.Input.Args) != 2 ||
		string(cis.ChaincodeSpec.Input.Args[0]) != "arg1" ||
		string(cis.ChaincodeSpec.Input.Args[1]) != "arg2" {
		t.Fatalf("Invalid chaincode invocation spec after unmarshalling\n")
		return
	}

	creator, transient, err := protoutil.GetChaincodeProposalContext(prop)
	if err != nil {
		t.Fatalf("Failed getting chaincode proposal context [%s]", err)
	}
	if string(creator) != "creator" {
		t.Fatalf("Failed checking Creator field. Invalid value, expectext 'creator', got [%s]", string(creator))
		return
	}
	value, ok := transient["certx"]
	if !ok || string(value) != "transient" {
		t.Fatalf("Failed checking Transient field. Invalid value, expectext 'transient', got [%s]", string(value))
		return
	}
}

func TestProposalWithTxID(t *testing.T) {
	
	prop, txid, err := protoutil.CreateChaincodeProposalWithTxIDAndTransient(
		common.HeaderType_ENDORSER_TRANSACTION,
		testChainID,
		createCIS(),
		[]byte("creator"),
		"testtx",
		map[string][]byte{"certx": []byte("transient")},
	)
	assert.Nil(t, err)
	assert.NotNil(t, prop)
	assert.Equal(t, txid, "testtx")

	prop, txid, err = protoutil.CreateChaincodeProposalWithTxIDAndTransient(
		common.HeaderType_ENDORSER_TRANSACTION,
		testChainID,
		createCIS(),
		[]byte("creator"),
		"",
		map[string][]byte{"certx": []byte("transient")},
	)
	assert.Nil(t, err)
	assert.NotNil(t, prop)
	assert.NotEmpty(t, txid)
}

func TestProposalResponse(t *testing.T) {
	events := &pb.ChaincodeEvent{
		ChaincodeId: "ccid",
		EventName:   "EventName",
		Payload:     []byte("EventPayload"),
		TxId:        "TxID"}
	ccid := &pb.ChaincodeID{
		Name:    "ccid",
		Version: "v1",
	}

	pHashBytes := []byte("proposal_hash")
	pResponse := &pb.Response{Status: 200}
	results := []byte("results")
	eventBytes, err := protoutil.GetBytesChaincodeEvent(events)
	if err != nil {
		t.Fatalf("Failure while marshalling the ProposalResponsePayload")
		return
	}

	
	pResponseBytes, err := protoutil.GetBytesResponse(pResponse)
	if err != nil {
		t.Fatalf("Failure while marshalling the Response")
		return
	}

	
	_, err = protoutil.GetResponse(pResponseBytes)
	if err != nil {
		t.Fatalf("Failure while unmarshalling the Response")
		return
	}

	
	prpBytes, err := protoutil.GetBytesProposalResponsePayload(pHashBytes, pResponse, results, eventBytes, ccid)
	if err != nil {
		t.Fatalf("Failure while marshalling the ProposalResponsePayload")
		return
	}

	
	prp, err := protoutil.GetProposalResponsePayload(prpBytes)
	if err != nil {
		t.Fatalf("Failure while unmarshalling the ProposalResponsePayload")
		return
	}

	
	act, err := protoutil.GetChaincodeAction(prp.Extension)
	if err != nil {
		t.Fatalf("Failure while unmarshalling the ChaincodeAction")
		return
	}

	
	if string(act.Results) != "results" {
		t.Fatalf("Invalid actions after unmarshalling")
		return
	}

	event, err := protoutil.GetChaincodeEvents(act.Events)
	if err != nil {
		t.Fatalf("Failure while unmarshalling the ChainCodeEvents")
		return
	}

	
	if string(event.ChaincodeId) != "ccid" {
		t.Fatalf("Invalid actions after unmarshalling")
		return
	}

	pr := &pb.ProposalResponse{
		Payload:     prpBytes,
		Endorsement: &pb.Endorsement{Endorser: []byte("endorser"), Signature: []byte("signature")},
		Version:     1, 
		Response:    &pb.Response{Status: 200, Message: "OK"}}

	
	prBytes, err := protoutil.GetBytesProposalResponse(pr)
	if err != nil {
		t.Fatalf("Failure while marshalling the ProposalResponse")
		return
	}

	
	prBack, err := protoutil.GetProposalResponse(prBytes)
	if err != nil {
		t.Fatalf("Failure while unmarshalling the ProposalResponse")
		return
	}

	
	if prBack.Response.Status != 200 ||
		string(prBack.Endorsement.Signature) != "signature" ||
		string(prBack.Endorsement.Endorser) != "endorser" ||
		bytes.Compare(prBack.Payload, prpBytes) != 0 {
		t.Fatalf("Invalid ProposalResponse after unmarshalling")
		return
	}
}

func TestEnvelope(t *testing.T) {
	
	prop, _, err := protoutil.CreateChaincodeProposal(common.HeaderType_ENDORSER_TRANSACTION, testChainID, createCIS(), signerSerialized)
	if err != nil {
		t.Fatalf("Could not create chaincode proposal, err %s\n", err)
		return
	}

	response := &pb.Response{Status: 200, Payload: []byte("payload")}
	result := []byte("res")
	ccid := &pb.ChaincodeID{Name: "foo", Version: "v1"}

	presp, err := protoutil.CreateProposalResponse(prop.Header, prop.Payload, response, result, nil, ccid, nil, signer)
	if err != nil {
		t.Fatalf("Could not create proposal response, err %s\n", err)
		return
	}

	tx, err := protoutil.CreateSignedTx(prop, signer, presp)
	if err != nil {
		t.Fatalf("Could not create signed tx, err %s\n", err)
		return
	}

	envBytes, err := protoutil.GetBytesEnvelope(tx)
	if err != nil {
		t.Fatalf("Could not marshal envelope, err %s\n", err)
		return
	}

	tx, err = protoutil.GetEnvelopeFromBlock(envBytes)
	if err != nil {
		t.Fatalf("Could not unmarshal envelope, err %s\n", err)
		return
	}

	act2, err := protoutil.GetActionFromEnvelope(envBytes)
	if err != nil {
		t.Fatalf("Could not extract actions from envelop, err %s\n", err)
		return
	}

	if act2.Response.Status != response.Status {
		t.Fatalf("response staus don't match")
		return
	}
	if bytes.Compare(act2.Response.Payload, response.Payload) != 0 {
		t.Fatalf("response payload don't match")
		return
	}

	if bytes.Compare(act2.Results, result) != 0 {
		t.Fatalf("results don't match")
		return
	}

	txpayl, err := protoutil.GetPayload(tx)
	if err != nil {
		t.Fatalf("Could not unmarshal payload, err %s\n", err)
		return
	}

	tx2, err := protoutil.GetTransaction(txpayl.Data)
	if err != nil {
		t.Fatalf("Could not unmarshal Transaction, err %s\n", err)
		return
	}

	sh, err := protoutil.GetSignatureHeader(tx2.Actions[0].Header)
	if err != nil {
		t.Fatalf("Could not unmarshal SignatureHeader, err %s\n", err)
		return
	}

	if bytes.Compare(sh.Creator, signerSerialized) != 0 {
		t.Fatalf("creator does not match")
		return
	}

	cap, err := protoutil.GetChaincodeActionPayload(tx2.Actions[0].Payload)
	if err != nil {
		t.Fatalf("Could not unmarshal ChaincodeActionPayload, err %s\n", err)
		return
	}
	assert.NotNil(t, cap)

	prp, err := protoutil.GetProposalResponsePayload(cap.Action.ProposalResponsePayload)
	if err != nil {
		t.Fatalf("Could not unmarshal ProposalResponsePayload, err %s\n", err)
		return
	}

	ca, err := protoutil.GetChaincodeAction(prp.Extension)
	if err != nil {
		t.Fatalf("Could not unmarshal ChaincodeAction, err %s\n", err)
		return
	}

	if ca.Response.Status != response.Status {
		t.Fatalf("response staus don't match")
		return
	}
	if bytes.Compare(ca.Response.Payload, response.Payload) != 0 {
		t.Fatalf("response payload don't match")
		return
	}

	if bytes.Compare(ca.Results, result) != 0 {
		t.Fatalf("results don't match")
		return
	}
}

func TestProposalTxID(t *testing.T) {
	nonce := []byte{1}
	creator := []byte{2}

	txid, err := protoutil.ComputeTxID(nonce, creator)
	assert.NotEmpty(t, txid, "TxID cannot be empty.")
	assert.NoError(t, err, "Failed computing txID")
	assert.Nil(t, protoutil.CheckTxID(txid, nonce, creator))
	assert.Error(t, protoutil.CheckTxID("", nonce, creator))

	txid, err = protoutil.ComputeTxID(nil, nil)
	assert.NotEmpty(t, txid, "TxID cannot be empty.")
	assert.NoError(t, err, "Failed computing txID")
}

func TestComputeProposalTxID(t *testing.T) {
	txid, err := protoutil.ComputeTxID([]byte{1}, []byte{1})
	assert.NoError(t, err, "Failed computing TxID")

	
	
	hf := sha256.New()
	hf.Write([]byte{1})
	hf.Write([]byte{1})
	hashOut := hf.Sum(nil)
	txid2 := hex.EncodeToString(hashOut)

	t.Logf("% x\n", hashOut)
	t.Logf("% s\n", txid)
	t.Logf("% s\n", txid2)

	assert.Equal(t, txid, txid2)
}

var signer msp.SigningIdentity
var signerSerialized []byte

func TestMain(m *testing.M) {
	
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not initialize msp")
		return
	}
	signer, err = mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not get signer")
		return
	}

	signerSerialized, err = signer.Serialize()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not serialize identity")
		return
	}

	os.Exit(m.Run())
}
