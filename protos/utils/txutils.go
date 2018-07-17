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

package utils

import (
	"errors"
	"fmt"

	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/peer"
)


func GetPayloads(txActions *peer.TransactionAction) (*peer.ChaincodeActionPayload, *peer.ChaincodeAction, error) {
	
	ccPayload := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(txActions.Payload, ccPayload)
	if err != nil {
		return nil, nil, err
	}

	if ccPayload.Action == nil || ccPayload.Action.ProposalResponsePayload == nil {
		return nil, nil, fmt.Errorf("no payload in ChaincodeActionPayload")
	}
	pRespPayload := &peer.ProposalResponsePayload{}
	err = proto.Unmarshal(ccPayload.Action.ProposalResponsePayload, pRespPayload)
	if err != nil {
		return nil, nil, err
	}

	if pRespPayload.Extension == nil {
		return nil, nil, fmt.Errorf("response payload is missing extension")
	}

	respPayload := &peer.ChaincodeAction{}
	err = proto.Unmarshal(pRespPayload.Extension, respPayload)
	if err != nil {
		return ccPayload, nil, err
	}
	return ccPayload, respPayload, nil
}


func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, fmt.Errorf("Error getting envelope(%s)", err)
	}

	return env, nil
}


func CreateSignedEnvelope(txType common.HeaderType, channelID string, signer crypto.LocalSigner, dataMsg proto.Message, msgVersion int32, epoch uint64) (*common.Envelope, error) {
	return CreateSignedEnvelopeWithTLSBinding(txType, channelID, signer, dataMsg, msgVersion, epoch, nil)
}



func CreateSignedEnvelopeWithTLSBinding(txType common.HeaderType, channelID string, signer crypto.LocalSigner, dataMsg proto.Message, msgVersion int32, epoch uint64, tlsCertHash []byte) (*common.Envelope, error) {
	payloadChannelHeader := MakeChannelHeader(txType, msgVersion, channelID, epoch)
	payloadChannelHeader.TlsCertHash = tlsCertHash
	var err error
	payloadSignatureHeader := &common.SignatureHeader{}

	if signer != nil {
		payloadSignatureHeader, err = signer.NewSignatureHeader()
		if err != nil {
			return nil, err
		}
	}

	data, err := proto.Marshal(dataMsg)
	if err != nil {
		return nil, err
	}

	paylBytes := MarshalOrPanic(&common.Payload{
		Header: MakePayloadHeader(payloadChannelHeader, payloadSignatureHeader),
		Data:   data,
	})

	var sig []byte
	if signer != nil {
		sig, err = signer.Sign(paylBytes)
		if err != nil {
			return nil, err
		}
	}

	return &common.Envelope{Payload: paylBytes, Signature: sig}, nil
}




func CreateSignedTx(proposal *peer.Proposal, signer msp.SigningIdentity, resps ...*peer.ProposalResponse) (*common.Envelope, error) {
	if len(resps) == 0 {
		return nil, fmt.Errorf("At least one proposal response is necessary")
	}

	
	hdr, err := GetHeader(proposal.Header)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal header")
	}

	
	pPayl, err := GetChaincodeProposalPayload(proposal.Payload)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal payload")
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
		return nil, fmt.Errorf("The signer needs to be the same as the one referenced in the header")
	}

	
	hdrExt, err := GetChaincodeHeaderExtension(hdr)
	if err != nil {
		return nil, err
	}

	
	var a1 []byte
	for n, r := range resps {
		if n == 0 {
			a1 = r.Payload
			if r.Response.Status != 200 {
				return nil, fmt.Errorf("Proposal response was not successful, error code %d, msg %s", r.Response.Status, r.Response.Message)
			}
			continue
		}

		if bytes.Compare(a1, r.Payload) != 0 {
			return nil, fmt.Errorf("ProposalResponsePayloads do not match")
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


func CreateProposalResponse(hdrbytes []byte, payl []byte, response *peer.Response, results []byte, events []byte, ccid *peer.ChaincodeID, visibility []byte, signingEndorser msp.SigningIdentity) (*peer.ProposalResponse, error) {
	hdr, err := GetHeader(hdrbytes)
	if err != nil {
		return nil, err
	}

	
	pHashBytes, err := GetProposalHash1(hdr, payl, visibility)
	if err != nil {
		return nil, fmt.Errorf("Could not compute proposal hash: err %s", err)
	}

	
	prpBytes, err := GetBytesProposalResponsePayload(pHashBytes, response, results, events, ccid)
	if err != nil {
		return nil, errors.New("Failure while marshaling the ProposalResponsePayload")
	}

	
	endorser, err := signingEndorser.Serialize()
	if err != nil {
		return nil, fmt.Errorf("Could not serialize the signing identity for %s, err %s", signingEndorser.GetIdentifier(), err)
	}

	
	signature, err := signingEndorser.Sign(append(prpBytes, endorser...))
	if err != nil {
		return nil, fmt.Errorf("Could not sign the proposal response payload, err %s", err)
	}

	resp := &peer.ProposalResponse{
		
		Version:     1, 
		Endorsement: &peer.Endorsement{Signature: signature, Endorser: endorser},
		Payload:     prpBytes,
		Response:    &peer.Response{Status: 200, Message: "OK"}}

	return resp, nil
}




func CreateProposalResponseFailure(hdrbytes []byte, payl []byte, response *peer.Response, results []byte, events []byte, ccid *peer.ChaincodeID, visibility []byte) (*peer.ProposalResponse, error) {
	hdr, err := GetHeader(hdrbytes)
	if err != nil {
		return nil, err
	}

	
	pHashBytes, err := GetProposalHash1(hdr, payl, visibility)
	if err != nil {
		return nil, fmt.Errorf("Could not compute proposal hash: err %s", err)
	}

	
	prpBytes, err := GetBytesProposalResponsePayload(pHashBytes, response, results, events, ccid)
	if err != nil {
		return nil, errors.New("Failure while marshaling the ProposalResponsePayload")
	}

	resp := &peer.ProposalResponse{
		
		Payload:  prpBytes,
		Response: response}

	return resp, nil
}


func GetSignedProposal(prop *peer.Proposal, signer msp.SigningIdentity) (*peer.SignedProposal, error) {
	
	if prop == nil || signer == nil {
		return nil, fmt.Errorf("Nil arguments")
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


func GetSignedEvent(evt *peer.Event, signer msp.SigningIdentity) (*peer.SignedEvent, error) {
	
	if evt == nil || signer == nil {
		return nil, errors.New("nil arguments")
	}

	evtBytes, err := proto.Marshal(evt)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(evtBytes)
	if err != nil {
		return nil, err
	}

	return &peer.SignedEvent{EventBytes: evtBytes, Signature: signature}, nil
}


func MockSignedEndorserProposalOrPanic(chainID string, cs *peer.ChaincodeSpec, creator, signature []byte) (*peer.SignedProposal, *peer.Proposal) {
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

func MockSignedEndorserProposal2OrPanic(chainID string, cs *peer.ChaincodeSpec, signer msp.SigningIdentity) (*peer.SignedProposal, *peer.Proposal) {
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



func GetBytesProposalPayloadForTx(payload *peer.ChaincodeProposalPayload, visibility []byte) ([]byte, error) {
	
	if payload == nil  {
		return nil, fmt.Errorf("Nil arguments")
	}

	
	cppNoTransient := &peer.ChaincodeProposalPayload{Input: payload.Input, TransientMap: nil}
	cppBytes, err := GetBytesChaincodeProposalPayload(cppNoTransient)
	if err != nil {
		return nil, errors.New("Failure while marshalling the ChaincodeProposalPayload!")
	}

	
	
	
	
	
	
	

	return cppBytes, nil
}





func GetProposalHash2(header *common.Header, ccPropPayl []byte) ([]byte, error) {
	
	if header == nil ||
		header.ChannelHeader == nil ||
		header.SignatureHeader == nil ||
		ccPropPayl == nil {
		return nil, fmt.Errorf("Nil arguments")
	}

	hash, err := factory.GetDefault().GetHash(&bccsp.SHA256Opts{})
	if err != nil {
		return nil, fmt.Errorf("Failed instantiating hash function [%s]", err)
	}
	hash.Write(header.ChannelHeader)   
	hash.Write(header.SignatureHeader) 
	hash.Write(ccPropPayl)             

	return hash.Sum(nil), nil
}



func GetProposalHash1(header *common.Header, ccPropPayl []byte, visibility []byte) ([]byte, error) {
	
	if header == nil ||
		header.ChannelHeader == nil ||
		header.SignatureHeader == nil ||
		ccPropPayl == nil  {
		return nil, fmt.Errorf("Nil arguments")
	}

	
	cpp := &peer.ChaincodeProposalPayload{}
	err := proto.Unmarshal(ccPropPayl, cpp)
	if err != nil {
		return nil, errors.New("Failure while unmarshalling the ChaincodeProposalPayload!")
	}

	ppBytes, err := GetBytesProposalPayloadForTx(cpp, visibility)
	if err != nil {
		return nil, err
	}

	hash2, err := factory.GetDefault().GetHash(&bccsp.SHA256Opts{})
	if err != nil {
		return nil, fmt.Errorf("Failed instantiating hash function [%s]", err)
	}
	hash2.Write(header.ChannelHeader)   
	hash2.Write(header.SignatureHeader) 
	hash2.Write(ppBytes)                

	return hash2.Sum(nil), nil
}
