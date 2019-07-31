/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoutil

import (
	"bytes"
	"crypto/sha256"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)


func GetPayloads(txActions *peer.TransactionAction) (*peer.ChaincodeActionPayload, *peer.ChaincodeAction, error) {
	
	
	ccPayload, err := GetChaincodeActionPayload(txActions.Payload)
	if err != nil {
		return nil, nil, err
	}

	if ccPayload.Action == nil || ccPayload.Action.ProposalResponsePayload == nil {
		return nil, nil, errors.New("no payload in ChaincodeActionPayload")
	}
	pRespPayload, err := GetProposalResponsePayload(ccPayload.Action.ProposalResponsePayload)
	if err != nil {
		return nil, nil, err
	}

	if pRespPayload.Extension == nil {
		return nil, nil, errors.New("response payload is missing extension")
	}

	respPayload, err := GetChaincodeAction(pRespPayload.Extension)
	if err != nil {
		return ccPayload, nil, err
	}
	return ccPayload, respPayload, nil
}


func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	return env, nil
}



func CreateSignedEnvelope(
	txType common.HeaderType,
	channelID string,
	signer identity.SignerSerializer,
	dataMsg proto.Message,
	msgVersion int32,
	epoch uint64,
) (*common.Envelope, error) {
	return CreateSignedEnvelopeWithTLSBinding(txType, channelID, signer, dataMsg, msgVersion, epoch, nil)
}




func CreateSignedEnvelopeWithTLSBinding(
	txType common.HeaderType,
	channelID string,
	signer identity.SignerSerializer,
	dataMsg proto.Message,
	msgVersion int32,
	epoch uint64,
	tlsCertHash []byte,
) (*common.Envelope, error) {
	payloadChannelHeader := MakeChannelHeader(txType, msgVersion, channelID, epoch)
	payloadChannelHeader.TlsCertHash = tlsCertHash
	var err error
	payloadSignatureHeader := &common.SignatureHeader{}

	if signer != nil {
		payloadSignatureHeader, err = NewSignatureHeader(signer)
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(dataMsg)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling")
	}

	paylBytes := MarshalOrPanic(
		&common.Payload{
			Header: MakePayloadHeader(payloadChannelHeader, payloadSignatureHeader),
			Data:   data,
		},
	)

	var sig []byte
	if signer != nil {
		sig, err = signer.Sign(paylBytes)
		if err != nil {
			return nil, err
		}
	}

	env := &common.Envelope{
		Payload:   paylBytes,
		Signature: sig,
	}

	return env, nil
}


type Signer interface {
	Sign(msg []byte) ([]byte, error)
	Serialize() ([]byte, error)
}





func CreateSignedTx(
	proposal *peer.Proposal,
	signer Signer,
	resps ...*peer.ProposalResponse,
) (*common.Envelope, error) {
	if len(resps) == 0 {
		return nil, errors.New("at least one proposal response is required")
	}

	
	hdr, err := GetHeader(proposal.Header)
	if err != nil {
		return nil, err
	}

	
	pPayl, err := GetChaincodeProposalPayload(proposal.Payload)
	if err != nil {
		return nil, err
	}

	
	
	signerBytes, err := signer.Serialize()
	if err != nil {
		return nil, err
	}

	shdr, err := GetSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(signerBytes, shdr.Creator) != 0 {
		return nil, errors.New("signer must be the same as the one referenced in the header")
	}

	
	hdrExt, err := GetChaincodeHeaderExtension(hdr)
	if err != nil {
		return nil, err
	}

	
	var a1 []byte
	for n, r := range resps {
		if r.Response.Status < 200 || r.Response.Status >= 400 {
			return nil, errors.Errorf("proposal response was not successful, error code %d, msg %s", r.Response.Status, r.Response.Message)
		}

		if n == 0 {
			a1 = r.Payload
			continue
		}

		if bytes.Compare(a1, r.Payload) != 0 {
			return nil, errors.New("ProposalResponsePayloads do not match")
		}
	}

	
	endorsements := make([]*peer.Endorsement, len(resps))
	for n, r := range resps {
		endorsements[n] = r.Endorsement
	}

	
	cea := &peer.ChaincodeEndorsedAction{ProposalResponsePayload: resps[0].Payload, Endorsements: endorsements}

	
	propPayloadBytes, err := GetBytesProposalPayloadForTx(pPayl, hdrExt.PayloadVisibility)
	if err != nil {
		return nil, err
	}

	
	cap := &peer.ChaincodeActionPayload{ChaincodeProposalPayload: propPayloadBytes, Action: cea}
	capBytes, err := GetBytesChaincodeActionPayload(cap)
	if err != nil {
		return nil, err
	}

	
	taa := &peer.TransactionAction{Header: hdr.SignatureHeader, Payload: capBytes}
	taas := make([]*peer.TransactionAction, 1)
	taas[0] = taa
	tx := &peer.Transaction{Actions: taas}

	
	txBytes, err := GetBytesTransaction(tx)
	if err != nil {
		return nil, err
	}

	
	payl := &common.Payload{Header: hdr, Data: txBytes}
	paylBytes, err := GetBytesPayload(payl)
	if err != nil {
		return nil, err
	}

	
	sig, err := signer.Sign(paylBytes)
	if err != nil {
		return nil, err
	}

	
	return &common.Envelope{Payload: paylBytes, Signature: sig}, nil
}


func CreateProposalResponse(
	hdrbytes []byte,
	payl []byte,
	response *peer.Response,
	results []byte,
	events []byte,
	ccid *peer.ChaincodeID,
	visibility []byte,
	signingEndorser identity.SignerSerializer,
) (*peer.ProposalResponse, error) {
	hdr, err := GetHeader(hdrbytes)
	if err != nil {
		return nil, err
	}

	
	
	pHashBytes, err := GetProposalHash1(hdr, payl, visibility)
	if err != nil {
		return nil, errors.WithMessage(err, "error computing proposal hash")
	}

	
	prpBytes, err := GetBytesProposalResponsePayload(pHashBytes, response, results, events, ccid)
	if err != nil {
		return nil, err
	}

	
	endorser, err := signingEndorser.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "error serializing signing identity")
	}

	
	
	signature, err := signingEndorser.Sign(append(prpBytes, endorser...))
	if err != nil {
		return nil, errors.WithMessage(err, "could not sign the proposal response payload")
	}

	resp := &peer.ProposalResponse{
		
		Version: 1, 
		Endorsement: &peer.Endorsement{
			Signature: signature,
			Endorser:  endorser,
		},
		Payload: prpBytes,
		Response: &peer.Response{
			Status:  200,
			Message: "OK",
		},
	}

	return resp, nil
}




func CreateProposalResponseFailure(
	hdrbytes []byte,
	payl []byte,
	response *peer.Response,
	results []byte,
	events []byte,
	ccid *peer.ChaincodeID,
	visibility []byte,
) (*peer.ProposalResponse, error) {
	hdr, err := GetHeader(hdrbytes)
	if err != nil {
		return nil, err
	}

	
	pHashBytes, err := GetProposalHash1(hdr, payl, visibility)
	if err != nil {
		return nil, errors.WithMessage(err, "error computing proposal hash")
	}

	
	prpBytes, err := GetBytesProposalResponsePayload(pHashBytes, response, results, events, ccid)
	if err != nil {
		return nil, err
	}

	resp := &peer.ProposalResponse{
		
		Payload:  prpBytes,
		Response: response,
	}

	return resp, nil
}



func GetSignedProposal(prop *peer.Proposal, signer identity.SignerSerializer) (*peer.SignedProposal, error) {
	
	if prop == nil || signer == nil {
		return nil, errors.New("nil arguments")
	}

	propBytes, err := GetBytesProposal(prop)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(propBytes)
	if err != nil {
		return nil, err
	}

	return &peer.SignedProposal{ProposalBytes: propBytes, Signature: signature}, nil
}



func MockSignedEndorserProposalOrPanic(
	chainID string,
	cs *peer.ChaincodeSpec,
	creator,
	signature []byte,
) (*peer.SignedProposal, *peer.Proposal) {
	prop, _, err := CreateChaincodeProposal(
		common.HeaderType_ENDORSER_TRANSACTION,
		chainID,
		&peer.ChaincodeInvocationSpec{ChaincodeSpec: cs},
		creator)
	if err != nil {
		panic(err)
	}

	propBytes, err := GetBytesProposal(prop)
	if err != nil {
		panic(err)
	}

	return &peer.SignedProposal{ProposalBytes: propBytes, Signature: signature}, prop
}

func MockSignedEndorserProposal2OrPanic(
	chainID string,
	cs *peer.ChaincodeSpec,
	signer identity.SignerSerializer,
) (*peer.SignedProposal, *peer.Proposal) {
	serializedSigner, err := signer.Serialize()
	if err != nil {
		panic(err)
	}

	prop, _, err := CreateChaincodeProposal(
		common.HeaderType_ENDORSER_TRANSACTION,
		chainID,
		&peer.ChaincodeInvocationSpec{ChaincodeSpec: &peer.ChaincodeSpec{}},
		serializedSigner)
	if err != nil {
		panic(err)
	}

	sProp, err := GetSignedProposal(prop, signer)
	if err != nil {
		panic(err)
	}

	return sProp, prop
}



func GetBytesProposalPayloadForTx(
	payload *peer.ChaincodeProposalPayload,
	visibility []byte,
) ([]byte, error) {
	
	if payload == nil {
		return nil, errors.New("nil arguments")
	}

	
	
	cppNoTransient := &peer.ChaincodeProposalPayload{Input: payload.Input, TransientMap: nil}
	cppBytes, err := GetBytesChaincodeProposalPayload(cppNoTransient)
	if err != nil {
		return nil, err
	}

	
	
	
	
	
	
	

	return cppBytes, nil
}





func GetProposalHash2(header *common.Header, ccPropPayl []byte) ([]byte, error) {
	
	if header == nil ||
		header.ChannelHeader == nil ||
		header.SignatureHeader == nil ||
		ccPropPayl == nil {
		return nil, errors.New("nil arguments")
	}

	hash := sha256.New()
	
	hash.Write(header.ChannelHeader)
	
	hash.Write(header.SignatureHeader)
	
	hash.Write(ccPropPayl)
	return hash.Sum(nil), nil
}



func GetProposalHash1(header *common.Header, ccPropPayl []byte, visibility []byte) ([]byte, error) {
	
	if header == nil ||
		header.ChannelHeader == nil ||
		header.SignatureHeader == nil ||
		ccPropPayl == nil {
		return nil, errors.New("nil arguments")
	}

	
	cpp, err := GetChaincodeProposalPayload(ccPropPayl)
	if err != nil {
		return nil, err
	}

	ppBytes, err := GetBytesProposalPayloadForTx(cpp, visibility)
	if err != nil {
		return nil, err
	}

	hash2 := sha256.New()
	
	hash2.Write(header.ChannelHeader)
	
	hash2.Write(header.SignatureHeader)
	
	hash2.Write(ppBytes)
	return hash2.Sum(nil), nil
}




func GetOrComputeTxIDFromEnvelope(txEnvelopBytes []byte) (string, error) {
	txEnvelope, err := UnmarshalEnvelope(txEnvelopBytes)
	if err != nil {
		return "", errors.WithMessage(err, "error getting txID from envelope")
	}

	txPayload, err := UnmarshalPayload(txEnvelope.Payload)
	if err != nil {
		return "", errors.WithMessage(err, "error getting txID from payload")
	}

	if txPayload.Header == nil {
		return "", errors.New("error getting txID from header: payload header is nil")
	}

	chdr, err := UnmarshalChannelHeader(txPayload.Header.ChannelHeader)
	if err != nil {
		return "", errors.WithMessage(err, "error getting txID from channel header")
	}

	if chdr.TxId != "" {
		return chdr.TxId, nil
	}

	sighdr, err := UnmarshalSignatureHeader(txPayload.Header.SignatureHeader)
	if err != nil {
		return "", errors.WithMessage(err, "error getting nonce and creator for computing txID")
	}

	return ComputeTxID(sighdr.Nonce, sighdr.Creator)
}
