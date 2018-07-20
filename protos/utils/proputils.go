/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"errors"
	"fmt"

	"encoding/binary"

	"encoding/hex"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)


func GetChaincodeInvocationSpec(prop *peer.Proposal) (*peer.ChaincodeInvocationSpec, error) {
	if prop == nil {
		return nil, fmt.Errorf("Proposal is nil")
	}
	_, err := GetHeader(prop.Header)
	if err != nil {
		return nil, err
	}
	ccPropPayload := &peer.ChaincodeProposalPayload{}
	err = proto.Unmarshal(prop.Payload, ccPropPayload)
	if err != nil {
		return nil, err
	}
	cis := &peer.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(ccPropPayload.Input, cis)
	return cis, err
}


func GetChaincodeProposalContext(prop *peer.Proposal) ([]byte, map[string][]byte, error) {
	if prop == nil {
		return nil, nil, fmt.Errorf("Proposal is nil")
	}
	if len(prop.Header) == 0 {
		return nil, nil, fmt.Errorf("Proposal's header is nil")
	}
	if len(prop.Payload) == 0 {
		return nil, nil, fmt.Errorf("Proposal's payload is nil")
	}

	
	hdr, err := GetHeader(prop.Header)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not extract the header from the proposal: %s", err)
	}
	if hdr == nil {
		return nil, nil, fmt.Errorf("Unmarshalled header is nil")
	}

	chdr, err := UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not extract the channel header from the proposal: %s", err)
	}

	if common.HeaderType(chdr.Type) != common.HeaderType_ENDORSER_TRANSACTION &&
		common.HeaderType(chdr.Type) != common.HeaderType_CONFIG {
		return nil, nil, fmt.Errorf("Invalid proposal type expected ENDORSER_TRANSACTION or CONFIG. Was: %d", chdr.Type)
	}

	shdr, err := GetSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not extract the signature header from the proposal: %s", err)
	}

	ccPropPayload := &peer.ChaincodeProposalPayload{}
	err = proto.Unmarshal(prop.Payload, ccPropPayload)
	if err != nil {
		return nil, nil, err
	}

	return shdr.Creator, ccPropPayload.TransientMap, nil
}


func GetHeader(bytes []byte) (*common.Header, error) {
	hdr := &common.Header{}
	err := proto.Unmarshal(bytes, hdr)
	return hdr, err
}


func GetNonce(prop *peer.Proposal) ([]byte, error) {
	if prop == nil {
		return nil, fmt.Errorf("Proposal is nil")
	}
	
	hdr, err := GetHeader(prop.Header)
	if err != nil {
		return nil, fmt.Errorf("Could not extract the header from the proposal: %s", err)
	}

	chdr, err := UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not extract the channel header from the proposal: %s", err)
	}

	if common.HeaderType(chdr.Type) != common.HeaderType_ENDORSER_TRANSACTION &&
		common.HeaderType(chdr.Type) != common.HeaderType_CONFIG {
		return nil, fmt.Errorf("Invalid proposal type expected ENDORSER_TRANSACTION or CONFIG. Was: %d", chdr.Type)
	}

	shdr, err := GetSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not extract the signature header from the proposal: %s", err)
	}

	if hdr.SignatureHeader == nil {
		return nil, errors.New("Invalid signature header. It must be different from nil.")
	}

	return shdr.Nonce, nil
}


func GetChaincodeHeaderExtension(hdr *common.Header) (*peer.ChaincodeHeaderExtension, error) {
	chdr, err := UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, err
	}

	chaincodeHdrExt := &peer.ChaincodeHeaderExtension{}
	err = proto.Unmarshal(chdr.Extension, chaincodeHdrExt)
	return chaincodeHdrExt, err
}


func GetProposalResponse(prBytes []byte) (*peer.ProposalResponse, error) {
	proposalResponse := &peer.ProposalResponse{}
	err := proto.Unmarshal(prBytes, proposalResponse)
	return proposalResponse, err
}


func GetChaincodeDeploymentSpec(code []byte, pr *platforms.Registry) (*peer.ChaincodeDeploymentSpec, error) {
	cds := &peer.ChaincodeDeploymentSpec{}
	err := proto.Unmarshal(code, cds)
	if err != nil {
		return nil, err
	}

	
	return cds, pr.ValidateDeploymentSpec(cds.CCType(), cds.Bytes())
}


func GetChaincodeAction(caBytes []byte) (*peer.ChaincodeAction, error) {
	chaincodeAction := &peer.ChaincodeAction{}
	err := proto.Unmarshal(caBytes, chaincodeAction)
	return chaincodeAction, err
}


func GetResponse(resBytes []byte) (*peer.Response, error) {
	response := &peer.Response{}
	err := proto.Unmarshal(resBytes, response)
	return response, err
}


func GetChaincodeEvents(eBytes []byte) (*peer.ChaincodeEvent, error) {
	chaincodeEvent := &peer.ChaincodeEvent{}
	err := proto.Unmarshal(eBytes, chaincodeEvent)
	return chaincodeEvent, err
}


func GetProposalResponsePayload(prpBytes []byte) (*peer.ProposalResponsePayload, error) {
	prp := &peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpBytes, prp)
	return prp, err
}


func GetProposal(propBytes []byte) (*peer.Proposal, error) {
	prop := &peer.Proposal{}
	err := proto.Unmarshal(propBytes, prop)
	return prop, err
}


func GetPayload(e *common.Envelope) (*common.Payload, error) {
	payload := &common.Payload{}
	err := proto.Unmarshal(e.Payload, payload)
	return payload, err
}


func GetTransaction(txBytes []byte) (*peer.Transaction, error) {
	tx := &peer.Transaction{}
	err := proto.Unmarshal(txBytes, tx)
	return tx, err
}


func GetChaincodeActionPayload(capBytes []byte) (*peer.ChaincodeActionPayload, error) {
	cap := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, cap)
	return cap, err
}


func GetChaincodeProposalPayload(bytes []byte) (*peer.ChaincodeProposalPayload, error) {
	cpp := &peer.ChaincodeProposalPayload{}
	err := proto.Unmarshal(bytes, cpp)
	return cpp, err
}


func GetSignatureHeader(bytes []byte) (*common.SignatureHeader, error) {
	sh := &common.SignatureHeader{}
	err := proto.Unmarshal(bytes, sh)
	return sh, err
}



func CreateChaincodeProposal(typ common.HeaderType, chainID string, cis *peer.ChaincodeInvocationSpec, creator []byte) (*peer.Proposal, string, error) {
	return CreateChaincodeProposalWithTransient(typ, chainID, cis, creator, nil)
}



func CreateChaincodeProposalWithTransient(typ common.HeaderType, chainID string, cis *peer.ChaincodeInvocationSpec, creator []byte, transientMap map[string][]byte) (*peer.Proposal, string, error) {
	
	nonce, err := crypto.GetRandomNonce()
	if err != nil {
		return nil, "", err
	}

	
	txid, err := ComputeProposalTxID(nonce, creator)
	if err != nil {
		return nil, "", err
	}

	return CreateChaincodeProposalWithTxIDNonceAndTransient(txid, typ, chainID, cis, nonce, creator, transientMap)
}



func CreateChaincodeProposalWithTxIDAndTransient(typ common.HeaderType, chainID string, cis *peer.ChaincodeInvocationSpec, creator []byte, txid string, transientMap map[string][]byte) (*peer.Proposal, string, error) {
	
	nonce, err := crypto.GetRandomNonce()
	if err != nil {
		return nil, "", err
	}

	
	if txid == "" {
		txid, err = ComputeProposalTxID(nonce, creator)
		if err != nil {
			return nil, "", err
		}
	}

	return CreateChaincodeProposalWithTxIDNonceAndTransient(txid, typ, chainID, cis, nonce, creator, transientMap)
}


func CreateChaincodeProposalWithTxIDNonceAndTransient(txid string, typ common.HeaderType, chainID string, cis *peer.ChaincodeInvocationSpec, nonce, creator []byte, transientMap map[string][]byte) (*peer.Proposal, string, error) {
	ccHdrExt := &peer.ChaincodeHeaderExtension{ChaincodeId: cis.ChaincodeSpec.ChaincodeId}
	ccHdrExtBytes, err := proto.Marshal(ccHdrExt)
	if err != nil {
		return nil, "", err
	}

	cisBytes, err := proto.Marshal(cis)
	if err != nil {
		return nil, "", err
	}

	ccPropPayload := &peer.ChaincodeProposalPayload{Input: cisBytes, TransientMap: transientMap}
	ccPropPayloadBytes, err := proto.Marshal(ccPropPayload)
	if err != nil {
		return nil, "", err
	}

	
	
	var epoch uint64 = 0

	timestamp := util.CreateUtcTimestamp()

	hdr := &common.Header{ChannelHeader: MarshalOrPanic(&common.ChannelHeader{
		Type:      int32(typ),
		TxId:      txid,
		Timestamp: timestamp,
		ChannelId: chainID,
		Extension: ccHdrExtBytes,
		Epoch:     epoch}),
		SignatureHeader: MarshalOrPanic(&common.SignatureHeader{Nonce: nonce, Creator: creator})}

	hdrBytes, err := proto.Marshal(hdr)
	if err != nil {
		return nil, "", err
	}

	return &peer.Proposal{Header: hdrBytes, Payload: ccPropPayloadBytes}, txid, nil
}


func GetBytesProposalResponsePayload(hash []byte, response *peer.Response, result []byte, event []byte, ccid *peer.ChaincodeID) ([]byte, error) {
	cAct := &peer.ChaincodeAction{Events: event, Results: result, Response: response, ChaincodeId: ccid}
	cActBytes, err := proto.Marshal(cAct)
	if err != nil {
		return nil, err
	}

	prp := &peer.ProposalResponsePayload{Extension: cActBytes, ProposalHash: hash}
	prpBytes, err := proto.Marshal(prp)
	return prpBytes, err
}


func GetBytesChaincodeProposalPayload(cpp *peer.ChaincodeProposalPayload) ([]byte, error) {
	cppBytes, err := proto.Marshal(cpp)
	return cppBytes, err
}


func GetBytesResponse(res *peer.Response) ([]byte, error) {
	resBytes, err := proto.Marshal(res)
	return resBytes, err
}


func GetBytesChaincodeEvent(event *peer.ChaincodeEvent) ([]byte, error) {
	eventBytes, err := proto.Marshal(event)
	return eventBytes, err
}


func GetBytesChaincodeActionPayload(cap *peer.ChaincodeActionPayload) ([]byte, error) {
	capBytes, err := proto.Marshal(cap)
	return capBytes, err
}


func GetBytesProposalResponse(pr *peer.ProposalResponse) ([]byte, error) {
	respBytes, err := proto.Marshal(pr)
	return respBytes, err
}


func GetBytesProposal(prop *peer.Proposal) ([]byte, error) {
	propBytes, err := proto.Marshal(prop)
	return propBytes, err
}


func GetBytesHeader(hdr *common.Header) ([]byte, error) {
	bytes, err := proto.Marshal(hdr)
	return bytes, err
}


func GetBytesSignatureHeader(hdr *common.SignatureHeader) ([]byte, error) {
	bytes, err := proto.Marshal(hdr)
	return bytes, err
}


func GetBytesTransaction(tx *peer.Transaction) ([]byte, error) {
	bytes, err := proto.Marshal(tx)
	return bytes, err
}


func GetBytesPayload(payl *common.Payload) ([]byte, error) {
	bytes, err := proto.Marshal(payl)
	return bytes, err
}


func GetBytesEnvelope(env *common.Envelope) ([]byte, error) {
	bytes, err := proto.Marshal(env)
	return bytes, err
}


func GetActionFromEnvelope(envBytes []byte) (*peer.ChaincodeAction, error) {
	env, err := GetEnvelopeFromBlock(envBytes)
	if err != nil {
		return nil, err
	}

	payl, err := GetPayload(env)
	if err != nil {
		return nil, err
	}

	tx, err := GetTransaction(payl.Data)
	if err != nil {
		return nil, err
	}

	if len(tx.Actions) == 0 {
		return nil, fmt.Errorf("At least one TransactionAction is required")
	}

	_, respPayload, err := GetPayloads(tx.Actions[0])
	return respPayload, err
}


func CreateProposalFromCISAndTxid(txid string, typ common.HeaderType, chainID string, cis *peer.ChaincodeInvocationSpec, creator []byte) (*peer.Proposal, string, error) {
	nonce, err := crypto.GetRandomNonce()
	if err != nil {
		return nil, "", err
	}
	return CreateChaincodeProposalWithTxIDNonceAndTransient(txid, typ, chainID, cis, nonce, creator, nil)
}


func CreateProposalFromCIS(typ common.HeaderType, chainID string, cis *peer.ChaincodeInvocationSpec, creator []byte) (*peer.Proposal, string, error) {
	return CreateChaincodeProposal(typ, chainID, cis, creator)
}


func CreateGetChaincodesProposal(chainID string, creator []byte) (*peer.Proposal, string, error) {
	ccinp := &peer.ChaincodeInput{Args: [][]byte{[]byte("getchaincodes")}}
	lsccSpec := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			Type:        peer.ChaincodeSpec_GOLANG,
			ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
			Input:       ccinp},
	}
	return CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, chainID, lsccSpec, creator)
}


func CreateGetInstalledChaincodesProposal(creator []byte) (*peer.Proposal, string, error) {
	ccinp := &peer.ChaincodeInput{Args: [][]byte{[]byte("getinstalledchaincodes")}}
	lsccSpec := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			Type:        peer.ChaincodeSpec_GOLANG,
			ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
			Input:       ccinp},
	}
	return CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, "", lsccSpec, creator)
}


func CreateInstallProposalFromCDS(ccpack proto.Message, creator []byte) (*peer.Proposal, string, error) {
	return createProposalFromCDS("", ccpack, creator, "install")
}


func CreateDeployProposalFromCDS(
	chainID string,
	cds *peer.ChaincodeDeploymentSpec,
	creator []byte,
	policy []byte,
	escc []byte,
	vscc []byte,
	collectionConfig []byte) (*peer.Proposal, string, error) {
	if collectionConfig == nil {
		return createProposalFromCDS(chainID, cds, creator, "deploy", policy, escc, vscc)
	} else {
		return createProposalFromCDS(chainID, cds, creator, "deploy", policy, escc, vscc, collectionConfig)
	}
}


func CreateUpgradeProposalFromCDS(
	chainID string,
	cds *peer.ChaincodeDeploymentSpec,
	creator []byte,
	policy []byte,
	escc []byte,
	vscc []byte,
	collectionConfig []byte) (*peer.Proposal, string, error) {
	if collectionConfig == nil {
		return createProposalFromCDS(chainID, cds, creator, "upgrade", policy, escc, vscc)
	} else {
		return createProposalFromCDS(chainID, cds, creator, "upgrade", policy, escc, vscc, collectionConfig)
	}
}


func createProposalFromCDS(chainID string, msg proto.Message, creator []byte, propType string, args ...[]byte) (*peer.Proposal, string, error) {
	
	var ccinp *peer.ChaincodeInput
	var b []byte
	var err error
	if msg != nil {
		b, err = proto.Marshal(msg)
		if err != nil {
			return nil, "", err
		}
	}
	switch propType {
	case "deploy":
		fallthrough
	case "upgrade":
		cds, ok := msg.(*peer.ChaincodeDeploymentSpec)
		if !ok || cds == nil {
			return nil, "", fmt.Errorf("invalid message for creating lifecycle chaincode proposal from")
		}
		Args := [][]byte{[]byte(propType), []byte(chainID), b}
		Args = append(Args, args...)

		ccinp = &peer.ChaincodeInput{Args: Args}
	case "install":
		ccinp = &peer.ChaincodeInput{Args: [][]byte{[]byte(propType), b}}
	}

	
	lsccSpec := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			Type:        peer.ChaincodeSpec_GOLANG,
			ChaincodeId: &peer.ChaincodeID{Name: "lscc"},
			Input:       ccinp}}

	
	return CreateProposalFromCIS(common.HeaderType_ENDORSER_TRANSACTION, chainID, lsccSpec, creator)
}



func ComputeProposalTxID(nonce, creator []byte) (string, error) {
	
	
	digest, err := factory.GetDefault().Hash(
		append(nonce, creator...),
		&bccsp.SHA256Opts{})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(digest), nil
}



func CheckProposalTxID(txid string, nonce, creator []byte) error {
	computedTxID, err := ComputeProposalTxID(nonce, creator)
	if err != nil {
		return fmt.Errorf("Failed computing target TXID for comparison [%s]", err)
	}

	if txid != computedTxID {
		return fmt.Errorf("Transaction is not valid. Got [%s], expected [%s]", txid, computedTxID)
	}

	return nil
}


func ComputeProposalBinding(proposal *peer.Proposal) ([]byte, error) {
	if proposal == nil {
		return nil, fmt.Errorf("Porposal is nil")
	}
	if len(proposal.Header) == 0 {
		return nil, fmt.Errorf("Proposal's Header is nil")
	}

	h, err := GetHeader(proposal.Header)
	if err != nil {
		return nil, err
	}

	chdr, err := UnmarshalChannelHeader(h.ChannelHeader)
	if err != nil {
		return nil, err
	}
	shdr, err := GetSignatureHeader(h.SignatureHeader)
	if err != nil {
		return nil, err
	}

	return computeProposalBindingInternal(shdr.Nonce, shdr.Creator, chdr.Epoch)
}

func computeProposalBindingInternal(nonce, creator []byte, epoch uint64) ([]byte, error) {
	epochBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(epochBytes, epoch)

	
	return factory.GetDefault().Hash(
		append(append(nonce, creator...), epochBytes...),
		&bccsp.SHA256Opts{})
}
