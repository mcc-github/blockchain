/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	validation "github.com/mcc-github/blockchain/core/handlers/validation/api"
	"github.com/mcc-github/blockchain/protoutil"
)


type PolicyEvaluator interface {
	validation.Dependency

	
	
	Evaluate(policyBytes []byte, signatureSet []*protoutil.SignedData) error
}


type SerializedPolicy interface {
	validation.ContextDatum

	
	Bytes() []byte
}
