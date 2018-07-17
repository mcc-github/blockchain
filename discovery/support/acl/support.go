/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package acl

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)

var (
	logger = flogging.MustGetLogger("discovery/acl")
)


type ChannelConfigGetter interface {
	
	GetChannelConfig(cid string) channelconfig.Resources
}


type ChannelConfigGetterFunc func(cid string) channelconfig.Resources


func (f ChannelConfigGetterFunc) GetChannelConfig(cid string) channelconfig.Resources {
	return f(cid)
}


type Verifier interface {
	
	
	
	
	VerifyByChannel(channel string, sd *cb.SignedData) error
}



type Evaluator interface {
	
	Evaluate(signatureSet []*cb.SignedData) error
}



type DiscoverySupport struct {
	ChannelConfigGetter
	Verifier
	Evaluator
}


func NewDiscoverySupport(v Verifier, e Evaluator, chanConf ChannelConfigGetter) *DiscoverySupport {
	return &DiscoverySupport{Verifier: v, Evaluator: e, ChannelConfigGetter: chanConf}
}



func (s *DiscoverySupport) EligibleForService(channel string, data cb.SignedData) error {
	if channel == "" {
		return s.Evaluate([]*cb.SignedData{&data})
	}
	return s.VerifyByChannel(channel, &data)
}


func (s *DiscoverySupport) ConfigSequence(channel string) uint64 {
	
	if channel == "" {
		return 0
	}
	conf := s.GetChannelConfig(channel)
	if conf == nil {
		logger.Panic("Failed obtaining channel config for channel", channel)
	}
	v := conf.ConfigtxValidator()
	if v == nil {
		logger.Panic("ConfigtxValidator for channel", channel, "is nil")
	}
	return v.Sequence()
}

func (s *DiscoverySupport) SatisfiesPrincipal(channel string, rawIdentity []byte, principal *msp.MSPPrincipal) error {
	conf := s.GetChannelConfig(channel)
	if conf == nil {
		return errors.Errorf("channel %s doesn't exist", channel)
	}
	mspMgr := conf.MSPManager()
	if mspMgr == nil {
		return errors.Errorf("could not find MSP manager for channel %s", channel)
	}
	identity, err := mspMgr.DeserializeIdentity(rawIdentity)
	if err != nil {
		return errors.Wrap(err, "failed deserializing identity")
	}
	return identity.SatisfiesPrincipal(principal)
}


func (s *DiscoverySupport) MSPOfPrincipal(principal *msp.MSPPrincipal) string {
	if principal == nil {
		return ""
	}
	switch principal.PrincipalClassification {
	case msp.MSPPrincipal_ROLE:
		
		mspRole := &msp.MSPRole{}
		err := proto.Unmarshal(principal.Principal, mspRole)
		if err != nil {
			logger.Warning("Failed unmarshaling principal:", err)
			return ""
		}
		return mspRole.MspIdentifier
	case msp.MSPPrincipal_IDENTITY:
		sId := &msp.SerializedIdentity{}
		err := proto.Unmarshal(principal.Principal, sId)
		if err != nil {
			logger.Warning("Failed unmarshaling principal:", err)
			return ""
		}
		return sId.Mspid
	case msp.MSPPrincipal_ORGANIZATION_UNIT:
		ouRole := &msp.OrganizationUnit{}
		err := proto.Unmarshal(principal.Principal, ouRole)
		if err != nil {
			logger.Warning("Failed unmarshaling principal:", err)
			return ""
		}
		return ouRole.MspIdentifier
	}
	logger.Warning("Received principal of unknown classification:", principal)
	return ""
}





type ChannelPolicyManagerGetter interface {
	
	
	Manager(channelID string) (policies.Manager, bool)
}


func NewChannelVerifier(policy string, polMgr policies.ChannelPolicyManagerGetter) *ChannelVerifier {
	return &ChannelVerifier{
		Policy: policy,
		ChannelPolicyManagerGetter: polMgr,
	}
}


type ChannelVerifier struct {
	policies.ChannelPolicyManagerGetter
	Policy string
}





func (cv *ChannelVerifier) VerifyByChannel(channel string, sd *cb.SignedData) error {
	mgr, _ := cv.Manager(channel)
	if mgr == nil {
		return errors.Errorf("policy manager for channel %s doesn't exist", channel)
	}
	pol, _ := mgr.GetPolicy(cv.Policy)
	if pol == nil {
		return errors.New("failed obtaining channel application writers policy")
	}
	return pol.Evaluate([]*cb.SignedData{sd})
}
