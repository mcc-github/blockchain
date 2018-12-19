/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/mocks/config"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/stretchr/testify/assert"
)

func createTestTransactionEnvelope(channel string, response *peer.Response, simRes []byte) (*common.Envelope, error) {
	prop, sProp, err := createTestProposalAndSignedProposal(channel)
	if err != nil {
		return nil, fmt.Errorf("failed to create test proposal and signed proposal, err %s", err)
	}

	
	_, _, _, err = ValidateProposalMessage(sProp)
	if err != nil {
		return nil, fmt.Errorf("ValidateProposalMessage failed, err %s", err)
	}

	
	presp, err := utils.CreateProposalResponse(prop.Header, prop.Payload, response, simRes, nil, getChaincodeID(), nil, signer)
	if err != nil {
		return nil, fmt.Errorf("CreateProposalResponse failed, err %s", err)
	}

	
	tx, err := utils.CreateSignedTx(prop, signer, presp)
	if err != nil {
		return nil, fmt.Errorf("CreateSignedTx failed, err %s", err)
	}

	return tx, nil
}

func createTestProposalAndSignedProposal(channel string) (*peer.Proposal, *peer.SignedProposal, error) {
	
	prop, err := getProposal(channel)
	if err != nil {
		return nil, nil, fmt.Errorf("getProposal failed, err %s", err)
	}

	
	sProp, err := utils.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, nil, fmt.Errorf("GetSignedProposal failed, err %s", err)
	}
	return prop, sProp, nil
}

func protoMarshal(t *testing.T, m proto.Message) []byte {
	bytes, err := proto.Marshal(m)
	assert.NoError(t, err)
	return bytes
}


func getTokenTransaction() *token.TokenTransaction {
	return &token.TokenTransaction{
		Action: &token.TokenTransaction_PlainAction{
			PlainAction: &token.PlainTokenAction{
				Data: &token.PlainTokenAction_PlainImport{
					PlainImport: &token.PlainImport{
						Outputs: []*token.PlainOutput{{
							Owner:    []byte("token-owner"),
							Type:     "PDQ",
							Quantity: 777,
						}},
					},
				},
			},
		},
	}
}



func createTestHeader(t *testing.T, txType common.HeaderType, channelId string, creator []byte, useGoodTxid bool) (*common.Header, error) {
	nonce := []byte("nonce-abc-12345")

	
	txid := "bad"
	if useGoodTxid {
		var err error
		txid, err = utils.ComputeTxID(nonce, creator)
		assert.NoError(t, err)
	}

	chdr := &common.ChannelHeader{
		Type:      int32(txType),
		ChannelId: channelId,
		TxId:      txid,
		Epoch:     uint64(0),
	}

	shdr := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}

	return &common.Header{
		ChannelHeader:   protoMarshal(t, chdr),
		SignatureHeader: protoMarshal(t, shdr),
	}, nil
}

func createTestEnvelope(t *testing.T, data []byte, header *common.Header, signer msp.SigningIdentity) (*common.Envelope, error) {
	payload := &common.Payload{
		Header: header,
		Data:   data,
	}
	payloadBytes := protoMarshal(t, payload)

	signature, err := signer.Sign(payloadBytes)
	assert.NoError(t, err)

	return &common.Envelope{
		Payload:   payloadBytes,
		Signature: signature,
	}, nil
}

func TestCheckSignatureFromCreator(t *testing.T) {
	response := &peer.Response{Status: 200}
	simRes := []byte("simulation_result")

	env, err := createTestTransactionEnvelope(util.GetTestChainID(), response, simRes)
	assert.Nil(t, err, "failed to create test transaction: %s", err)
	assert.NotNil(t, env)

	
	payload, err := utils.GetPayload(env)
	assert.NoError(t, err, "GetPayload returns err %s", err)

	
	chdr, shdr, err := validateCommonHeader(payload.Header)
	assert.NoError(t, err, "validateCommonHeader returns err %s", err)

	
	err = checkSignatureFromCreator(shdr.Creator, env.Signature, env.Payload, chdr.ChannelId)
	assert.NoError(t, err, "checkSignatureFromCreator returns err %s", err)

	
	err = checkSignatureFromCreator([]byte("junk"), env.Signature, env.Payload, chdr.ChannelId)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MSP error: could not deserialize")

	
	err = checkSignatureFromCreator(shdr.Creator, env.Signature, env.Payload, "junkchannel")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MSP error: channel doesn't exist")
}

func TestValidateProposalMessage(t *testing.T) {
	
	fakeChannel := "fakechannel"
	_, sProp, err := createTestProposalAndSignedProposal(fakeChannel)
	assert.NoError(t, err)
	
	_, _, _, err = ValidateProposalMessage(sProp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("access denied: channel [%s] creator org [%s]", fakeChannel, signerMSPId))

	
	_, sProp, err = createTestProposalAndSignedProposal(util.GetTestChainID())
	assert.NoError(t, err)
	sigCopy := make([]byte, len(sProp.Signature))
	copy(sigCopy, sProp.Signature)
	for i := 0; i < len(sProp.Signature); i++ {
		sigCopy[i] = byte(int(sigCopy[i]+1) % 255)
	}
	
	_, _, _, err = ValidateProposalMessage(&peer.SignedProposal{ProposalBytes: sProp.ProposalBytes, Signature: sigCopy})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("access denied: channel [%s] creator org [%s]", util.GetTestChainID(), signerMSPId))
}

func TestValidateTokenTransaction(t *testing.T) {
	tokenTx := getTokenTransaction()
	txBytes := protoMarshal(t, tokenTx)
	err := validateTokenTransaction(txBytes)
	assert.NoError(t, err)
}

func TestValidateTokenTransactionBadData(t *testing.T) {
	txBytes := []byte("bad-data")
	err := validateTokenTransaction(txBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error unmarshaling the token Transaction")
}

func TestValidateTransactionGoodTokenTx(t *testing.T) {
	tokenTx := getTokenTransaction()
	txBytes := protoMarshal(t, tokenTx)
	header, err := createTestHeader(t, common.HeaderType_TOKEN_TRANSACTION, util.GetTestChainID(), signerSerialized, true)
	assert.NoError(t, err)
	envelope, err := createTestEnvelope(t, txBytes, header, signer)
	assert.NoError(t, err)
	payload, code := ValidateTransaction(envelope, &config.MockApplicationCapabilities{})
	assert.Equal(t, code, peer.TxValidationCode_VALID)
	assert.Equal(t, payload.Data, txBytes)
}

func TestValidateTransactionBadTokenTxData(t *testing.T) {
	txBytes := []byte("bad-data")
	header, err := createTestHeader(t, common.HeaderType_TOKEN_TRANSACTION, util.GetTestChainID(), signerSerialized, true)
	assert.NoError(t, err)
	envelope, err := createTestEnvelope(t, txBytes, header, signer)
	assert.NoError(t, err)
	payload, code := ValidateTransaction(envelope, &config.MockApplicationCapabilities{})
	assert.Equal(t, code, peer.TxValidationCode_BAD_PAYLOAD)
	assert.Equal(t, payload.Data, txBytes)
}

func TestValidateTransactionBadTokenTxID(t *testing.T) {
	tokenTx := getTokenTransaction()
	txBytes := protoMarshal(t, tokenTx)
	header, err := createTestHeader(t, common.HeaderType_TOKEN_TRANSACTION, util.GetTestChainID(), signerSerialized, false)
	assert.NoError(t, err)
	envelope, err := createTestEnvelope(t, txBytes, header, signer)
	assert.NoError(t, err)
	payload, code := ValidateTransaction(envelope, &config.MockApplicationCapabilities{})
	assert.Equal(t, code, peer.TxValidationCode_BAD_PROPOSAL_TXID)
	assert.Nil(t, payload)
}
