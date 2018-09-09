/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plain

import (
	"sync"

	"github.com/mcc-github/blockchain/protos/token"
	"github.com/mcc-github/blockchain/token/tms"
	"github.com/pkg/errors"
)




type Credential interface {
	Public() []byte
	Private() []byte
}




type PolicyValidator interface {
	
	IsIssuer(creator Credential, tokenType string) error
}




type Pool interface {
	CommitUpdate(transactionData []tms.TransactionData) error
}


type Verifier struct {
	Pool            Pool
	PolicyValidator PolicyValidator

	mutex sync.Mutex
}



func (v *Verifier) Validate(creator Credential, data tms.TransactionData) error {
	plainAction := data.Tx.GetPlainAction()
	if plainAction == nil {
		return errors.Errorf("validation failed: unknown action: %T", data.Tx.GetAction())
	}

	switch action := plainAction.Data.(type) {
	case *token.PlainTokenAction_PlainImport:
		return nil
	default:
		return errors.Errorf("validation failed: unknown plain token action: %T", action)
	}
}



func (v *Verifier) Commit(creator Credential, transactionData []tms.TransactionData) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	for _, data := range transactionData {
		plainAction := data.Tx.GetPlainAction()
		if plainAction == nil {
			return errors.Errorf("commit failed: unknown action: %T", data.Tx.GetAction())
		}

		switch action := plainAction.Data.(type) {
		case *token.PlainTokenAction_PlainImport:
			err := v.commitCheckImport(creator, data.TxID, action.PlainImport)
			if err != nil {
				return errors.WithMessage(err, "commit failed")
			}
		default:
			return errors.Errorf("commit failed: unknown plain token action: %T", action)
		}
	}

	err := v.Pool.CommitUpdate(transactionData)
	if err != nil {
		return errors.Wrap(err, "commit failed")
	}
	return nil
}

func (v *Verifier) commitCheckImport(creator Credential, txID string, importData *token.PlainImport) error {
	for _, output := range importData.Outputs {
		err := v.PolicyValidator.IsIssuer(creator, output.Type)
		if err != nil {
			return err
		}
	}
	return nil
}
