/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tms

import "github.com/mcc-github/blockchain/protos/token"






type TransactionData struct {
	Tx *token.TokenTransaction
	
	TxID string
}
