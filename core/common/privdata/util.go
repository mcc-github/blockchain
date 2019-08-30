/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/mcc-github/blockchain-protos-go/common"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/pkg/errors"
)




func getPolicy(collectionPolicyConfig *common.CollectionPolicyConfig, deserializer msp.IdentityDeserializer) (policies.Policy, error) {
	if collectionPolicyConfig == nil {
		return nil, errors.New("collection policy config is nil")
	}
	accessPolicyEnvelope := collectionPolicyConfig.GetSignaturePolicy()
	if accessPolicyEnvelope == nil {
		return nil, errors.New("collection config access policy is nil")
	}
	

	pp := cauthdsl.EnvelopeBasedPolicyProvider{Deserializer: deserializer}
	accessPolicy, err := pp.NewPolicy(accessPolicyEnvelope)
	if err != nil {
		return nil, errors.WithMessage(err, "failed constructing policy object out of collection policy config")
	}

	return accessPolicy, nil
}
