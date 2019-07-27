/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gossip

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/internal/pkg/identity"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	pcommon "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

var mcsLogger = flogging.MustGetLogger("peer.gossip.mcs")










type MSPMessageCryptoService struct {
	channelPolicyManagerGetter policies.ChannelPolicyManagerGetter
	localSigner                identity.SignerSerializer
	deserializer               mgmt.DeserializersManager
}







func NewMCS(
	channelPolicyManagerGetter policies.ChannelPolicyManagerGetter,
	localSigner identity.SignerSerializer,
	deserializer mgmt.DeserializersManager,
) *MSPMessageCryptoService {
	return &MSPMessageCryptoService{
		channelPolicyManagerGetter: channelPolicyManagerGetter,
		localSigner:                localSigner,
		deserializer:               deserializer,
	}
}




func (s *MSPMessageCryptoService) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	
	
	

	_, _, err := s.getValidatedIdentity(peerIdentity)
	return err
}







func (s *MSPMessageCryptoService) GetPKIidOfCert(peerIdentity api.PeerIdentityType) common.PKIidType {
	
	if len(peerIdentity) == 0 {
		mcsLogger.Error("Invalid Peer Identity. It must be different from nil.")

		return nil
	}

	sid, err := s.deserializer.Deserialize(peerIdentity)
	if err != nil {
		mcsLogger.Errorf("Failed getting validated identity from peer identity [% x]: [%s]", peerIdentity, err)

		return nil
	}

	
	
	

	mspIdRaw := []byte(sid.Mspid)
	raw := append(mspIdRaw, sid.IdBytes...)

	
	digest, err := factory.GetDefault().Hash(raw, &bccsp.SHA256Opts{})
	if err != nil {
		mcsLogger.Errorf("Failed computing digest of serialized identity [% x]: [%s]", peerIdentity, err)

		return nil
	}

	return digest
}




func (s *MSPMessageCryptoService) VerifyBlock(chainID common.ChannelID, seqNum uint64, signedBlock []byte) error {
	
	block, err := protoutil.GetBlockFromBlockBytes(signedBlock)
	if err != nil {
		return fmt.Errorf("Failed unmarshalling block bytes on channel [%s]: [%s]", chainID, err)
	}

	if block.Header == nil {
		return fmt.Errorf("Invalid Block on channel [%s]. Header must be different from nil.", chainID)
	}

	blockSeqNum := block.Header.Number
	if seqNum != blockSeqNum {
		return fmt.Errorf("Claimed seqNum is [%d] but actual seqNum inside block is [%d]", seqNum, blockSeqNum)
	}

	
	channelID, err := protoutil.GetChainIDFromBlock(block)
	if err != nil {
		return fmt.Errorf("Failed getting channel id from block with id [%d] on channel [%s]: [%s]", block.Header.Number, chainID, err)
	}

	if channelID != string(chainID) {
		return fmt.Errorf("Invalid block's channel id. Expected [%s]. Given [%s]", chainID, channelID)
	}

	
	if block.Metadata == nil || len(block.Metadata.Metadata) == 0 {
		return fmt.Errorf("Block with id [%d] on channel [%s] does not have metadata. Block not valid.", block.Header.Number, chainID)
	}

	metadata, err := protoutil.GetMetadataFromBlock(block, pcommon.BlockMetadataIndex_SIGNATURES)
	if err != nil {
		return fmt.Errorf("Failed unmarshalling medatata for signatures [%s]", err)
	}

	
	
	if !bytes.Equal(protoutil.BlockDataHash(block.Data), block.Header.DataHash) {
		return fmt.Errorf("Header.DataHash is different from Hash(block.Data) for block with id [%d] on channel [%s]", block.Header.Number, chainID)
	}

	

	
	cpm := s.channelPolicyManagerGetter.Manager(channelID)
	if cpm == nil {
		return fmt.Errorf("Could not acquire policy manager for channel %s", channelID)
	}
	mcsLogger.Debugf("Got policy manager for channel [%s]", channelID)

	
	policy, ok := cpm.GetPolicy(policies.BlockValidation)
	
	mcsLogger.Debugf("Got block validation policy for channel [%s] with flag [%t]", channelID, ok)

	
	signatureSet := []*protoutil.SignedData{}
	for _, metadataSignature := range metadata.Signatures {
		shdr, err := protoutil.GetSignatureHeader(metadataSignature.SignatureHeader)
		if err != nil {
			return fmt.Errorf("Failed unmarshalling signature header for block with id [%d] on channel [%s]: [%s]", block.Header.Number, chainID, err)
		}
		signatureSet = append(
			signatureSet,
			&protoutil.SignedData{
				Identity:  shdr.Creator,
				Data:      util.ConcatenateBytes(metadata.Value, metadataSignature.SignatureHeader, protoutil.BlockHeaderBytes(block.Header)),
				Signature: metadataSignature.Signature,
			},
		)
	}

	
	return policy.Evaluate(signatureSet)
}



func (s *MSPMessageCryptoService) Sign(msg []byte) ([]byte, error) {
	return s.localSigner.Sign(msg)
}




func (s *MSPMessageCryptoService) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
	identity, chainID, err := s.getValidatedIdentity(peerIdentity)
	if err != nil {
		mcsLogger.Errorf("Failed getting validated identity from peer identity [%s]", err)

		return err
	}

	if len(chainID) == 0 {
		
		
		
		return identity.Verify(message, signature)
	}

	
	
	

	return s.VerifyByChannel(chainID, peerIdentity, signature, message)
}





func (s *MSPMessageCryptoService) VerifyByChannel(chainID common.ChannelID, peerIdentity api.PeerIdentityType, signature, message []byte) error {
	
	if len(peerIdentity) == 0 {
		return errors.New("Invalid Peer Identity. It must be different from nil.")
	}

	
	cpm := s.channelPolicyManagerGetter.Manager(string(chainID))
	if cpm == nil {
		return fmt.Errorf("Could not acquire policy manager for channel %s", string(chainID))
	}
	mcsLogger.Debugf("Got policy manager for channel [%s]", string(chainID))

	
	policy, flag := cpm.GetPolicy(policies.ChannelApplicationReaders)
	mcsLogger.Debugf("Got reader policy for channel [%s] with flag [%t]", string(chainID), flag)

	return policy.Evaluate(
		[]*protoutil.SignedData{{
			Data:      message,
			Identity:  []byte(peerIdentity),
			Signature: signature,
		}},
	)
}

func (s *MSPMessageCryptoService) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	id, _, err := s.getValidatedIdentity(peerIdentity)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Unable to extract msp.Identity from peer Identity")
	}
	return id.ExpiresAt(), nil

}

func (s *MSPMessageCryptoService) getValidatedIdentity(peerIdentity api.PeerIdentityType) (msp.Identity, common.ChannelID, error) {
	
	if len(peerIdentity) == 0 {
		return nil, nil, errors.New("Invalid Peer Identity. It must be different from nil.")
	}

	sId, err := s.deserializer.Deserialize(peerIdentity)
	if err != nil {
		mcsLogger.Error("failed deserializing identity", err)
		return nil, nil, err
	}

	
	

	
	
	
	
	lDes := s.deserializer.GetLocalDeserializer()
	identity, err := lDes.DeserializeIdentity([]byte(peerIdentity))
	if err == nil {
		
		
		if err := lDes.IsWellFormed(sId); err != nil {
			return nil, nil, errors.Wrap(err, "identity is not well formed")
		}
		
		
		
		
		
		
		
		
		if identity.GetMSPIdentifier() == s.deserializer.GetLocalMSPIdentifier() {
			

			
			
			
			return identity, nil, identity.Validate()
		}
	}

	
	for chainID, mspManager := range s.deserializer.GetChannelDeserializers() {
		
		identity, err := mspManager.DeserializeIdentity([]byte(peerIdentity))
		if err != nil {
			mcsLogger.Debugf("Failed deserialization identity [% x] on [%s]: [%s]", peerIdentity, chainID, err)
			continue
		}

		
		if err := mspManager.IsWellFormed(sId); err != nil {
			return nil, nil, errors.Wrap(err, "identity is not well formed")
		}

		
		
		
		

		if err := identity.Validate(); err != nil {
			mcsLogger.Debugf("Failed validating identity [% x] on [%s]: [%s]", peerIdentity, chainID, err)
			continue
		}

		mcsLogger.Debugf("Validation succeeded [% x] on [%s]", peerIdentity, chainID)

		return identity, common.ChannelID(chainID), nil
	}

	return nil, nil, fmt.Errorf("Peer Identity [% x] cannot be validated. No MSP found able to do that.", peerIdentity)
}
