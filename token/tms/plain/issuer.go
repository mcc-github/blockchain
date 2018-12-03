/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"github.com/mcc-github/blockchain/protos/token"
)


type Issuer struct{}


func (i *Issuer) RequestImport(tokensToIssue []*token.TokenToIssue) (*token.TokenTransaction, error) {
	var outputs []*token.PlainOutput
	for _, tti := range tokensToIssue {
		outputs = append(outputs, &token.PlainOutput{
			Owner:    tti.Recipient,
			Type:     tti.Type,
			Quantity: tti.Quantity,
		})
	}

	return &token.TokenTransaction{
		Action: &token.TokenTransaction_PlainAction{
			PlainAction: &token.PlainTokenAction{
				Data: &token.PlainTokenAction_PlainImport{
					PlainImport: &token.PlainImport{
						Outputs: outputs,
					},
				},
			},
		},
	}, nil
}



func (i *Issuer) RequestExpectation(request *token.ExpectationRequest) (*token.TokenTransaction, error) {
	panic("not implemented yet")
}
