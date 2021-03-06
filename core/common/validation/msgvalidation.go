/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"bytes"

	"github.com/mcc-github/blockchain-protos-go/common"
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/flogging"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var putilsLogger = flogging.MustGetLogger("protoutils")




func checkSignatureFromCreator(creatorBytes []byte, sig []byte, msg []byte, ChainID string) error {
	putilsLogger.Debugf("begin")

	
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

	putilsLogger.Debugf("creator is %s", creator.GetIdentifier())

	
	err = creator.Validate()
	if err != nil {
		return errors.WithMessage(err, "creator certificate is not valid")
	}

	putilsLogger.Debugf("creator is valid")

	
	err = creator.Verify(msg, sig)
	if err != nil {
		return errors.WithMessage(err, "creator's signature over the proposal is not valid")
	}

	putilsLogger.Debugf("exits successfully")

	return nil
}


func validateSignatureHeader(sHdr *common.SignatureHeader) error {
	
	if sHdr == nil {
		return errors.New("nil SignatureHeader provided")
	}

	
	if sHdr.Nonce == nil || len(sHdr.Nonce) == 0 {
		return errors.New("invalid nonce specified in the header")
	}

	
	if sHdr.Creator == nil || len(sHdr.Creator) == 0 {
		return errors.New("invalid creator specified in the header")
	}

	return nil
}


func validateChannelHeader(cHdr *common.ChannelHeader) error {
	
	if cHdr == nil {
		return errors.New("nil ChannelHeader provided")
	}

	
	switch common.HeaderType(cHdr.Type) {
	case common.HeaderType_ENDORSER_TRANSACTION:
	case common.HeaderType_CONFIG_UPDATE:
	case common.HeaderType_CONFIG:
	default:
		return errors.Errorf("invalid header type %s", common.HeaderType(cHdr.Type))
	}

	putilsLogger.Debugf("validateChannelHeader info: header type %d", common.HeaderType(cHdr.Type))

	

	
	
	
	
	if cHdr.Epoch != 0 {
		return errors.Errorf("invalid Epoch in ChannelHeader. Expected 0, got [%d]", cHdr.Epoch)
	}

	

	return nil
}


func validateCommonHeader(hdr *common.Header) (*common.ChannelHeader, *common.SignatureHeader, error) {
	if hdr == nil {
		return nil, nil, errors.New("nil header")
	}

	chdr, err := protoutil.UnmarshalChannelHeader(hdr.ChannelHeader)
	if err != nil {
		return nil, nil, err
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, nil, err
	}

	err = validateChannelHeader(chdr)
	if err != nil {
		return nil, nil, err
	}

	err = validateSignatureHeader(shdr)
	if err != nil {
		return nil, nil, err
	}

	return chdr, shdr, nil
}



func validateConfigTransaction(data []byte, hdr *common.Header) error {
	putilsLogger.Debugf("validateConfigTransaction starts for data %p, header %s", data, hdr)

	
	if data == nil || hdr == nil {
		return errors.New("nil arguments")
	}

	

	return nil
}



func validateEndorserTransaction(data []byte, hdr *common.Header) error {
	putilsLogger.Debugf("validateEndorserTransaction starts for data %p, header %s", data, hdr)

	
	if data == nil || hdr == nil {
		return errors.New("nil arguments")
	}

	
	tx, err := protoutil.UnmarshalTransaction(data)
	if err != nil {
		return err
	}

	
	if tx == nil {
		return errors.New("nil transaction")
	}

	

	

	
	if len(tx.Actions) != 1 {
		return errors.Errorf("only one action per transaction is supported, tx contains %d", len(tx.Actions))
	}

	putilsLogger.Debugf("validateEndorserTransaction info: there are %d actions", len(tx.Actions))

	for _, act := range tx.Actions {
		
		if act == nil {
			return errors.New("nil action")
		}

		
		sHdr, err := protoutil.UnmarshalSignatureHeader(act.Header)
		if err != nil {
			return err
		}

		
		
		err = validateSignatureHeader(sHdr)
		if err != nil {
			return err
		}

		putilsLogger.Debugf("validateEndorserTransaction info: signature header is valid")

		
		ccActionPayload, err := protoutil.UnmarshalChaincodeActionPayload(act.Payload)
		if err != nil {
			return err
		}

		
		prp, err := protoutil.UnmarshalProposalResponsePayload(ccActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			return err
		}

		
		
		hdrOrig := &common.Header{ChannelHeader: hdr.ChannelHeader, SignatureHeader: act.Header}

		
		pHash, err := protoutil.GetProposalHash2(hdrOrig, ccActionPayload.ChaincodeProposalPayload)
		if err != nil {
			return err
		}

		
		if !bytes.Equal(pHash, prp.ProposalHash) {
			return errors.New("proposal hash does not match")
		}
	}

	return nil
}


func ValidateTransaction(e *common.Envelope) (*common.Payload, pb.TxValidationCode) {
	putilsLogger.Debugf("ValidateTransactionEnvelope starts for envelope %p", e)

	
	if e == nil {
		putilsLogger.Errorf("Error: nil envelope")
		return nil, pb.TxValidationCode_NIL_ENVELOPE
	}

	
	payload, err := protoutil.UnmarshalPayload(e.Payload)
	if err != nil {
		putilsLogger.Errorf("GetPayload returns err %s", err)
		return nil, pb.TxValidationCode_BAD_PAYLOAD
	}

	putilsLogger.Debugf("Header is %s", payload.Header)

	
	chdr, shdr, err := validateCommonHeader(payload.Header)
	if err != nil {
		putilsLogger.Errorf("validateCommonHeader returns err %s", err)
		return nil, pb.TxValidationCode_BAD_COMMON_HEADER
	}

	
	err = checkSignatureFromCreator(shdr.Creator, e.Signature, e.Payload, chdr.ChannelId)
	if err != nil {
		putilsLogger.Errorf("checkSignatureFromCreator returns err %s", err)
		return nil, pb.TxValidationCode_BAD_CREATOR_SIGNATURE
	}

	

	
	switch common.HeaderType(chdr.Type) {
	case common.HeaderType_ENDORSER_TRANSACTION:
		
		
		
		err = protoutil.CheckTxID(
			chdr.TxId,
			shdr.Nonce,
			shdr.Creator)

		if err != nil {
			putilsLogger.Errorf("CheckTxID returns err %s", err)
			return nil, pb.TxValidationCode_BAD_PROPOSAL_TXID
		}

		err = validateEndorserTransaction(payload.Data, payload.Header)
		putilsLogger.Debugf("ValidateTransactionEnvelope returns err %s", err)

		if err != nil {
			putilsLogger.Errorf("validateEndorserTransaction returns err %s", err)
			return payload, pb.TxValidationCode_INVALID_ENDORSER_TRANSACTION
		} else {
			return payload, pb.TxValidationCode_VALID
		}
	case common.HeaderType_CONFIG:
		
		

		err = validateConfigTransaction(payload.Data, payload.Header)

		if err != nil {
			putilsLogger.Errorf("validateConfigTransaction returns err %s", err)
			return payload, pb.TxValidationCode_INVALID_CONFIG_TRANSACTION
		} else {
			return payload, pb.TxValidationCode_VALID
		}
	default:
		return nil, pb.TxValidationCode_UNSUPPORTED_TX_PAYLOAD
	}
}
