/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"crypto/sha256"

	"github.com/golang/protobuf/proto"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protos/common"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
	pb "github.com/mcc-github/blockchain/protos/peer"
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
		return nil, errors.WithMessage(err, "could not unmarshal proposal bytes")
	}

	
	hdr, err := protoutil.UnmarshalHeader(prop.Header)
	if err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal header")
	}

	
	chdr, err := protoutil.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal channel header")
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal signature header")
	}

	
	chaincodeHdrExt, err := protoutil.UnmarshalChaincodeHeaderExtension(chdr.Extension)
	if err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal header extension")
	}

	if chaincodeHdrExt.ChaincodeId == nil {
		return nil, errors.Errorf("ChaincodeHeaderExtension.ChaincodeId is nil")
	}

	if chaincodeHdrExt.ChaincodeId.Name == "" {
		return nil, errors.New("ChaincodeHeaderExtension.ChaincodeId.Name is empty")
	}

	
	
	
	
	
	
	
	
	if chaincodeHdrExt.PayloadVisibility != nil {
		return nil, errors.New("invalid payload visibility field")
	}

	cis, err := protoutil.GetChaincodeInvocationSpec(prop)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get invocation spec")
	}

	if cis.ChaincodeSpec == nil {
		return nil, errors.Errorf("chaincode invocation spec did not contain chaincode spec")
	}

	if cis.ChaincodeSpec.Input == nil {
		return nil, errors.Errorf("chaincode input did not contain any input")
	}

	
	cpp, err := protoutil.UnmarshalChaincodeProposalPayload(prop.Payload)
	if err != nil {
		return nil, err
	}

	ppBytes, err := protoutil.GetBytesProposalPayloadForTx(cpp)
	if err != nil {
		return nil, err
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




func ValidateUnpackedProposal(unpackedProp *UnpackedProposal) error {
	if err := validateChannelHeader(unpackedProp.ChannelHeader); err != nil {
		return err
	}

	if err := validateSignatureHeader(unpackedProp.SignatureHeader); err != nil {
		return err
	}

	
	err := checkSignatureFromCreator(
		unpackedProp.SignatureHeader.Creator,
		unpackedProp.SignedProposal.Signature,
		unpackedProp.SignedProposal.ProposalBytes,
		unpackedProp.ChannelHeader.ChannelId,
	)
	if err != nil {
		
		
		endorserLogger.Warningf("channel [%s]: %s", unpackedProp.ChannelHeader.ChannelId, err)
		sId := &msp.SerializedIdentity{}
		err := proto.Unmarshal(unpackedProp.SignatureHeader.Creator, sId)
		if err != nil {
			
			err = errors.Wrap(err, "could not deserialize a SerializedIdentity")
			endorserLogger.Warningf("channel [%s]: %s", unpackedProp.ChannelHeader.ChannelId, err)
		}
		return errors.Errorf("access denied: channel [%s] creator org [%s]", unpackedProp.ChannelHeader.ChannelId, sId.Mspid)
	}

	
	
	
	err = protoutil.CheckTxID(
		unpackedProp.ChannelHeader.TxId,
		unpackedProp.SignatureHeader.Nonce,
		unpackedProp.SignatureHeader.Creator,
	)
	if err != nil {
		return err
	}

	return nil
}




func checkSignatureFromCreator(creatorBytes []byte, sig []byte, msg []byte, ChainID string) error {
	
	if creatorBytes == nil || sig == nil || msg == nil {
		return errors.New("nil arguments")
	}

	mspObj := mspmgmt.GetIdentityDeserializer(ChainID)
	if mspObj == nil {
		return errors.Errorf("could not get msp for channel [%s]", ChainID)
	}

	
	creator, err := mspObj.DeserializeIdentity(creatorBytes)
	if err != nil {
		return errors.WithMessage(err, "MSP error")
	}

	endorserLogger.Debugf("creator is %s", creator.GetIdentifier())

	
	err = creator.Validate()
	if err != nil {
		return errors.WithMessage(err, "creator certificate is not valid")
	}

	endorserLogger.Debugf("creator is valid")

	
	err = creator.Verify(msg, sig)
	if err != nil {
		return errors.WithMessage(err, "creator's signature over the proposal is not valid")
	}

	endorserLogger.Debugf("exits successfully")

	return nil
}


func validateSignatureHeader(sHdr *common.SignatureHeader) error {
	
	if sHdr.Nonce == nil || len(sHdr.Nonce) == 0 {
		return errors.New("invalid nonce specified in the header")
	}

	
	if sHdr.Creator == nil || len(sHdr.Creator) == 0 {
		return errors.New("invalid creator specified in the header")
	}

	return nil
}


func validateChannelHeader(cHdr *common.ChannelHeader) error {
	
	switch common.HeaderType(cHdr.Type) {
	case common.HeaderType_ENDORSER_TRANSACTION:
	case common.HeaderType_CONFIG:
		
		
		

	default:
		return errors.Errorf("invalid header type %s", common.HeaderType(cHdr.Type))
	}

	return nil
}
