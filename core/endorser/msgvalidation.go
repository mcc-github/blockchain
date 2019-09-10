/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"crypto/sha256"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	cb "github.com/mcc-github/blockchain-protos-go/common"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)


type UnpackedProposal struct {
	ChaincodeName   string
	ChannelHeader   *cb.ChannelHeader
	Input           *pb.ChaincodeInput
	Proposal        *pb.Proposal
	SignatureHeader *cb.SignatureHeader
	SignedProposal  *pb.SignedProposal
	ProposalHash    []byte
}

func (up *UnpackedProposal) ChannelID() string {
	return up.ChannelHeader.ChannelId
}

func (up *UnpackedProposal) TxID() string {
	return up.ChannelHeader.TxId
}



func UnpackProposal(signedProp *pb.SignedProposal) (*UnpackedProposal, error) {
	prop, err := protoutil.UnmarshalProposal(signedProp.ProposalBytes)
	if err != nil {
		return nil, err
	}

	hdr, err := protoutil.UnmarshalHeader(prop.Header)
	if err != nil {
		return nil, err
	}

	chdr, err := protoutil.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, err
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, err
	}

	chaincodeHdrExt, err := protoutil.UnmarshalChaincodeHeaderExtension(chdr.Extension)
	if err != nil {
		return nil, err
	}

	if chaincodeHdrExt.ChaincodeId == nil {
		return nil, errors.Errorf("ChaincodeHeaderExtension.ChaincodeId is nil")
	}

	if chaincodeHdrExt.ChaincodeId.Name == "" {
		return nil, errors.Errorf("ChaincodeHeaderExtension.ChaincodeId.Name is empty")
	}

	cpp, err := protoutil.UnmarshalChaincodeProposalPayload(prop.Payload)
	if err != nil {
		return nil, err
	}

	cis, err := protoutil.UnmarshalChaincodeInvocationSpec(cpp.Input)
	if err != nil {
		return nil, err
	}

	if cis.ChaincodeSpec == nil {
		return nil, errors.Errorf("chaincode invocation spec did not contain chaincode spec")
	}

	if cis.ChaincodeSpec.Input == nil {
		return nil, errors.Errorf("chaincode input did not contain any input")
	}

	cppNoTransient := &pb.ChaincodeProposalPayload{Input: cpp.Input, TransientMap: nil}
	ppBytes, err := proto.Marshal(cppNoTransient)
	if err != nil {
		return nil, errors.WithMessage(err, "could not marshal non-transient portion of payload")
	}

	

	
	
	
	
	
	propHash := sha256.New()
	propHash.Write(hdr.ChannelHeader)
	propHash.Write(hdr.SignatureHeader)
	propHash.Write(ppBytes)

	return &UnpackedProposal{
		SignedProposal:  signedProp,
		Proposal:        prop,
		ChannelHeader:   chdr,
		SignatureHeader: shdr,
		ChaincodeName:   chaincodeHdrExt.ChaincodeId.Name,
		Input:           cis.ChaincodeSpec.Input,
		ProposalHash:    propHash.Sum(nil)[:],
	}, nil
}

func (up *UnpackedProposal) Validate(idDeserializer msp.IdentityDeserializer) error {
	
	switch common.HeaderType(up.ChannelHeader.Type) {
	case common.HeaderType_ENDORSER_TRANSACTION:
	case common.HeaderType_CONFIG:
		
		
		

	default:
		return errors.Errorf("invalid header type %s", common.HeaderType(up.ChannelHeader.Type))
	}

	
	if up.ChannelHeader.Epoch != 0 {
		return errors.Errorf("epoch is non-zero")
	}

	
	if len(up.SignatureHeader.Nonce) == 0 {
		return errors.Errorf("nonce is empty")
	}

	
	if len(up.SignatureHeader.Creator) == 0 {
		return errors.New("creator is empty")
	}

	expectedTxID := protoutil.ComputeTxID(up.SignatureHeader.Nonce, up.SignatureHeader.Creator)
	if up.TxID() != expectedTxID {
		return errors.Errorf("incorrectly computed txid '%s' -- expected '%s'", up.TxID(), expectedTxID)
	}

	if up.SignedProposal.ProposalBytes == nil {
		return errors.Errorf("empty proposal bytes")
	}

	if up.SignedProposal.Signature == nil {
		return errors.Errorf("empty signature bytes")
	}

	sId, err := protoutil.UnmarshalSerializedIdentity(up.SignatureHeader.Creator)
	if err != nil {
		return errors.Errorf("access denied: channel [%s] creator org unknown, creator is malformed", up.ChannelID())
	}

	genericAuthError := errors.Errorf("access denied: channel [%s] creator org [%s]", up.ChannelID(), sId.Mspid)

	
	creator, err := idDeserializer.DeserializeIdentity(up.SignatureHeader.Creator)
	if err != nil {
		endorserLogger.Warningf("%s access denied: channel [%s]: %s", up.TxID(), up.ChannelID(), err)
		return genericAuthError
	}

	endorserLogger.Debugf("%s creator is %s", up.TxID(), creator.GetIdentifier())

	
	err = creator.Validate()
	if err != nil {
		endorserLogger.Warningf("%s access denied: channel [%s]: identity is not valid: %s", up.TxID(), up.ChannelID(), err)
		return genericAuthError
	}

	endorserLogger.Debugf("%s creator is valid", up.TxID())

	
	err = creator.Verify(up.SignedProposal.ProposalBytes, up.SignedProposal.Signature)
	if err != nil {
		endorserLogger.Warningf("%s access denied: channel [%s]: creator's signature over the proposal is not valid: %s", up.TxID(), up.ChannelID(), err)
		return genericAuthError
	}

	endorserLogger.Debugf("%s signature is valid", up.TxID())

	return nil
}
