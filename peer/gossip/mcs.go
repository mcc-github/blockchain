/*
Copyright IBM Corp. 2017 All Rights Reserved.

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

package gossip

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/gossip/api"
	"github.com/mcc-github/blockchain/gossip/common"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
	pcommon "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/utils"
	"github.com/pkg/errors"
)

var mcsLogger = flogging.MustGetLogger("peer/gossip/mcs")










type mspMessageCryptoService struct {
	channelPolicyManagerGetter policies.ChannelPolicyManagerGetter
	localSigner                crypto.LocalSigner
	deserializer               mgmt.DeserializersManager
}







func NewMCS(channelPolicyManagerGetter policies.ChannelPolicyManagerGetter, localSigner crypto.LocalSigner, deserializer mgmt.DeserializersManager) *mspMessageCryptoService {
	return &mspMessageCryptoService{channelPolicyManagerGetter: channelPolicyManagerGetter, localSigner: localSigner, deserializer: deserializer}
}




func (s *mspMessageCryptoService) ValidateIdentity(peerIdentity api.PeerIdentityType) error {
	
	
	

	_, _, err := s.getValidatedIdentity(peerIdentity)
	return err
}







func (s *mspMessageCryptoService) GetPKIidOfCert(peerIdentity api.PeerIdentityType) common.PKIidType {
	
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




func (s *mspMessageCryptoService) VerifyBlock(chainID common.ChainID, seqNum uint64, signedBlock []byte) error {
	
	block, err := utils.GetBlockFromBlockBytes(signedBlock)
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

	
	channelID, err := utils.GetChainIDFromBlock(block)
	if err != nil {
		return fmt.Errorf("Failed getting channel id from block with id [%d] on channel [%s]: [%s]", block.Header.Number, chainID, err)
	}

	if channelID != string(chainID) {
		return fmt.Errorf("Invalid block's channel id. Expected [%s]. Given [%s]", chainID, channelID)
	}

	
	if block.Metadata == nil || len(block.Metadata.Metadata) == 0 {
		return fmt.Errorf("Block with id [%d] on channel [%s] does not have metadata. Block not valid.", block.Header.Number, chainID)
	}

	metadata, err := utils.GetMetadataFromBlock(block, pcommon.BlockMetadataIndex_SIGNATURES)
	if err != nil {
		return fmt.Errorf("Failed unmarshalling medatata for signatures [%s]", err)
	}

	
	
	if !bytes.Equal(block.Data.Hash(), block.Header.DataHash) {
		return fmt.Errorf("Header.DataHash is different from Hash(block.Data) for block with id [%d] on channel [%s]", block.Header.Number, chainID)
	}

	

	
	cpm, ok := s.channelPolicyManagerGetter.Manager(channelID)
	if cpm == nil {
		return fmt.Errorf("Could not acquire policy manager for channel %s", channelID)
	}
	
	mcsLogger.Debugf("Got policy manager for channel [%s] with flag [%t]", channelID, ok)

	
	policy, ok := cpm.GetPolicy(policies.BlockValidation)
	
	mcsLogger.Debugf("Got block validation policy for channel [%s] with flag [%t]", channelID, ok)

	
	signatureSet := []*pcommon.SignedData{}
	for _, metadataSignature := range metadata.Signatures {
		shdr, err := utils.GetSignatureHeader(metadataSignature.SignatureHeader)
		if err != nil {
			return fmt.Errorf("Failed unmarshalling signature header for block with id [%d] on channel [%s]: [%s]", block.Header.Number, chainID, err)
		}
		signatureSet = append(
			signatureSet,
			&pcommon.SignedData{
				Identity:  shdr.Creator,
				Data:      util.ConcatenateBytes(metadata.Value, metadataSignature.SignatureHeader, block.Header.Bytes()),
				Signature: metadataSignature.Signature,
			},
		)
	}

	
	return policy.Evaluate(signatureSet)
}



func (s *mspMessageCryptoService) Sign(msg []byte) ([]byte, error) {
	return s.localSigner.Sign(msg)
}




func (s *mspMessageCryptoService) Verify(peerIdentity api.PeerIdentityType, signature, message []byte) error {
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





func (s *mspMessageCryptoService) VerifyByChannel(chainID common.ChainID, peerIdentity api.PeerIdentityType, signature, message []byte) error {
	
	if len(peerIdentity) == 0 {
		return errors.New("Invalid Peer Identity. It must be different from nil.")
	}

	
	cpm, flag := s.channelPolicyManagerGetter.Manager(string(chainID))
	if cpm == nil {
		return fmt.Errorf("Could not acquire policy manager for channel %s", string(chainID))
	}
	mcsLogger.Debugf("Got policy manager for channel [%s] with flag [%t]", string(chainID), flag)

	
	policy, flag := cpm.GetPolicy(policies.ChannelApplicationReaders)
	mcsLogger.Debugf("Got reader policy for channel [%s] with flag [%t]", string(chainID), flag)

	return policy.Evaluate(
		[]*pcommon.SignedData{{
			Data:      message,
			Identity:  []byte(peerIdentity),
			Signature: signature,
		}},
	)
}

func (s *mspMessageCryptoService) Expiration(peerIdentity api.PeerIdentityType) (time.Time, error) {
	id, _, err := s.getValidatedIdentity(peerIdentity)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Unable to extract msp.Identity from peer Identity")
	}
	return id.ExpiresAt(), nil

}

func (s *mspMessageCryptoService) getValidatedIdentity(peerIdentity api.PeerIdentityType) (msp.Identity, common.ChainID, error) {
	
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

		return identity, common.ChainID(chainID), nil
	}

	return nil, nil, fmt.Errorf("Peer Identity [% x] cannot be validated. No MSP found able to do that.", peerIdentity)
}
