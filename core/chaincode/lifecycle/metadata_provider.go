/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)



type LegacyMetadataProvider interface {
	Metadata(channel string, cc string, includeCollections bool) *chaincode.Metadata
}



type ChaincodeInfoProvider interface {
	
	
	ChaincodeInfo(channelID, name string) (*LocalChaincodeInfo, error)
}


type ChannelPolicyReferenceProvider interface {
	
	NewPolicy(channelID, channelConfigPolicyReference string) (policies.Policy, error)
}

type channelPolicyReferenceProviderImpl struct {
	ChannelConfigSource
}


func (c *channelPolicyReferenceProviderImpl) NewPolicy(channelID, channelConfigPolicyReference string) (policies.Policy, error) {
	cc := c.GetStableChannelConfig(channelID)
	pm := cc.PolicyManager()
	p, ok := pm.GetPolicy(channelConfigPolicyReference)
	if !ok {
		return nil, errors.Errorf("could not retrieve policy for reference '%s' on channel '%s'", channelConfigPolicyReference, channelID)
	}

	return p, nil
}


func NewMetadataProvider(cip ChaincodeInfoProvider, lmp LegacyMetadataProvider, ccs ChannelConfigSource) *MetadataProvider {
	return &MetadataProvider{
		ChaincodeInfoProvider:  cip,
		LegacyMetadataProvider: lmp,
		ChannelPolicyReferenceProvider: &channelPolicyReferenceProviderImpl{
			ChannelConfigSource: ccs,
		},
	}
}

type MetadataProvider struct {
	ChaincodeInfoProvider          ChaincodeInfoProvider
	LegacyMetadataProvider         LegacyMetadataProvider
	ChannelPolicyReferenceProvider ChannelPolicyReferenceProvider
}

func (mp *MetadataProvider) toSignaturePolicyEnvelope(channelID string, policyBytes []byte) ([]byte, error) {
	p := &peer.ApplicationPolicy{}
	err := proto.Unmarshal(policyBytes, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal ApplicationPolicy bytes")
	}

	switch policy := p.Type.(type) {
	case *peer.ApplicationPolicy_SignaturePolicy:
		return protoutil.MarshalOrPanic(policy.SignaturePolicy), nil
	case *peer.ApplicationPolicy_ChannelConfigPolicyReference:
		p, err := mp.ChannelPolicyReferenceProvider.NewPolicy(channelID, policy.ChannelConfigPolicyReference)
		if err != nil {
			return nil, errors.WithMessagef(err, "could not retrieve policy for reference '%s' on channel '%s'", policy.ChannelConfigPolicyReference, channelID)
		}

		cp, ok := p.(policies.Converter)
		if !ok {
			return nil, errors.Errorf("policy with reference '%s' on channel '%s' is not convertible to SignaturePolicyEnvelope", policy.ChannelConfigPolicyReference, channelID)
		}

		spe, err := cp.Convert()
		if err != nil {
			return nil, errors.WithMessagef(err, "error converting policy with reference '%s' on channel '%s' to SignaturePolicyEnvelope", policy.ChannelConfigPolicyReference, channelID)
		}

		return proto.Marshal(spe)
	default:
		
		return nil, errors.Errorf("unsupported policy type %T on channel '%s'", policy, channelID)
	}
}


func (mp *MetadataProvider) Metadata(channel string, ccName string, includeCollections bool) *chaincode.Metadata {
	ccInfo, err := mp.ChaincodeInfoProvider.ChaincodeInfo(channel, ccName)
	if err != nil {
		logger.Debugf("chaincode '%s' on channel '%s' not defined in _lifecycle. requesting metadata from lscc", ccName, channel)
		
		return mp.LegacyMetadataProvider.Metadata(channel, ccName, includeCollections)
	}

	spe, err := mp.toSignaturePolicyEnvelope(channel, ccInfo.Definition.ValidationInfo.ValidationParameter)
	if err != nil {
		logger.Errorf("could not convert policy for chaincode '%s' on channel '%s', err '%s'", ccName, channel, err)
		spe = cauthdsl.MarshaledRejectAllPolicy
	}

	
	
	
	ccMetadata := &chaincode.Metadata{
		Name:              ccName,
		Version:           strconv.FormatInt(ccInfo.Definition.Sequence, 10),
		Policy:            spe,
		CollectionsConfig: ccInfo.Definition.Collections,
	}

	return ccMetadata
}
