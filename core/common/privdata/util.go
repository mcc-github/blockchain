/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package privdata

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/common/policies"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)

func getPolicy(collectionPolicyConfig *common.CollectionPolicyConfig, deserializer msp.IdentityDeserializer) (policies.Policy, error) {
	if collectionPolicyConfig == nil {
		return nil, errors.New("Collection policy config is nil")
	}
	accessPolicyEnvelope := collectionPolicyConfig.GetSignaturePolicy()
	if accessPolicyEnvelope == nil {
		return nil, errors.New("Collection config access policy is nil")
	}
	
	npp := cauthdsl.NewPolicyProvider(deserializer)
	polBytes, err := proto.Marshal(accessPolicyEnvelope)
	if err != nil {
		return nil, err
	}
	accessPolicy, _, err := npp.NewPolicy(polBytes)
	if err != nil {
		return nil, err
	}
	return accessPolicy, nil
}
