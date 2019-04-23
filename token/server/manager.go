/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/mcc-github/blockchain/token/ledger"
	"github.com/mcc-github/blockchain/token/tms/plain"
	"github.com/pkg/errors"
)



type Manager struct {
	LedgerManager              ledger.LedgerManager
	TokenOwnerValidatorManager identity.TokenOwnerValidatorManager
}




func (m *Manager) GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error) {
	tokenOwnerValidator, err := m.TokenOwnerValidatorManager.Get(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting token owner validator for channel: %s", channel)
	}

	return &plain.Issuer{TokenOwnerValidator: tokenOwnerValidator}, nil
}



func (m *Manager) GetTransactor(channel string, privateCredential, publicCredential []byte) (Transactor, error) {
	ledger, err := m.LedgerManager.GetLedgerReader(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting ledger for channel: %s", channel)
	}

	tokenOwnerValidator, err := m.TokenOwnerValidatorManager.Get(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting token owner validator for channel: %s", channel)
	}

	return &plain.Transactor{
		Ledger:              ledger,
		PublicCredential:    publicCredential,
		TokenOwnerValidator: tokenOwnerValidator}, nil
}
