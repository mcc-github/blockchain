/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/common/policies/inquire"
	common2 "github.com/mcc-github/blockchain/protos/common"
)

var logger = flogging.MustGetLogger("discovery/DiscoverySupport")

type MetadataRetriever interface {
	Metadata(channel string, cc string, loadCollections bool) *chaincode.Metadata
}



type DiscoverySupport struct {
	ci MetadataRetriever
}


func NewDiscoverySupport(ci MetadataRetriever) *DiscoverySupport {
	s := &DiscoverySupport{
		ci: ci,
	}
	return s
}

func (s *DiscoverySupport) PolicyByChaincode(channel string, cc string) policies.InquireablePolicy {
	chaincodeData := s.ci.Metadata(channel, cc, false)
	if chaincodeData == nil {
		logger.Info("Chaincode", cc, "wasn't found")
		return nil
	}
	pol := &common2.SignaturePolicyEnvelope{}
	if err := proto.Unmarshal(chaincodeData.Policy, pol); err != nil {
		logger.Warning("Failed unmarshaling policy for chaincode", cc, ":", err)
		return nil
	}
	if len(pol.Identities) == 0 || pol.Rule == nil {
		logger.Warningf("Invalid policy, either Identities(%v) or Rule(%v) are empty:", pol.Identities, pol.Rule)
		return nil
	}
	return inquire.NewInquireableSignaturePolicy(pol)
}
