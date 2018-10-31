/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/mcc-github/blockchain/token/tms/plain"
	"github.com/pkg/errors"
)



type Manager struct {
	LedgerManager ledger.LedgerManager
}




func (manager *Manager) GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error) {
	return &plain.Issuer{}, nil
}



func (manager *Manager) GetTransactor(channel string, privateCredential, publicCredential []byte) (Transactor, error) {
	ledger, err := manager.LedgerManager.GetLedgerReader(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting ledger for channel: %s", channel)
	}
	return &plain.Transactor{Ledger: ledger, PublicCredential: publicCredential}, nil
}
