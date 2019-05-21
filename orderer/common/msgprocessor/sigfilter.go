/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protos/orderer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)


type SigFilterSupport interface {
	
	PolicyManager() policies.Manager
	
	OrdererConfig() (channelconfig.Orderer, bool)
}



type SigFilter struct {
	normalPolicyName      string
	maintenancePolicyName string
	support               SigFilterSupport
}






func NewSigFilter(normalPolicyName, maintenancePolicyName string, support SigFilterSupport) *SigFilter {
	return &SigFilter{
		normalPolicyName:      normalPolicyName,
		maintenancePolicyName: maintenancePolicyName,
		support:               support,
	}
}


func (sf *SigFilter) Apply(message *cb.Envelope) error {
	ordererConf, ok := sf.support.OrdererConfig()
	if !ok {
		logger.Panic("Programming error: orderer config not found")
	}

	signedData, err := protoutil.EnvelopeAsSignedData(message)

	if err != nil {
		return fmt.Errorf("could not convert message to signedData: %s", err)
	}

	
	
	
	var policyName = sf.normalPolicyName
	if ordererConf.ConsensusState() == orderer.ConsensusType_STATE_MAINTENANCE {
		policyName = sf.maintenancePolicyName
	}

	policy, ok := sf.support.PolicyManager().GetPolicy(policyName)
	if !ok {
		return fmt.Errorf("could not find policy %s", policyName)
	}

	err = policy.Evaluate(signedData)
	if err != nil {
		logger.Debugf("SigFilter evaluation failed: %s, policyName: %s, ConsensusState: %s", err.Error(), policyName, ordererConf.ConsensusState())
		return errors.Wrap(errors.WithStack(ErrPermissionDenied), err.Error())
	}
	return nil
}
