/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"fmt"

	"github.com/mcc-github/blockchain/common/policies"
	cb "github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)


type SigFilterSupport interface {
	
	PolicyManager() policies.Manager
}



type SigFilter struct {
	policyName string
	support    SigFilterSupport
}



func NewSigFilter(policyName string, support SigFilterSupport) *SigFilter {
	return &SigFilter{
		policyName: policyName,
		support:    support,
	}
}


func (sf *SigFilter) Apply(message *cb.Envelope) error {
	signedData, err := message.AsSignedData()

	if err != nil {
		return fmt.Errorf("could not convert message to signedData: %s", err)
	}

	policy, ok := sf.support.PolicyManager().GetPolicy(sf.policyName)
	if !ok {
		return fmt.Errorf("could not find policy %s", sf.policyName)
	}

	err = policy.Evaluate(signedData)
	if err != nil {
		return errors.Wrap(errors.WithStack(ErrPermissionDenied), err.Error())
	}
	return nil
}
