/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/protos/common"
)


type PolicyEvaluator interface {
	validation.Dependency

	
	
	Evaluate(policyBytes []byte, signatureSet []*common.SignedData) error
}


type SerializedPolicy interface {
	validation.ContextDatum

	
	Bytes() []byte
}
