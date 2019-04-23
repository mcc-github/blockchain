/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package manager

import (
	"github.com/mcc-github/blockchain/msp/mgmt"
	"github.com/mcc-github/blockchain/token/identity"
	"github.com/mcc-github/blockchain/token/tms/plain"
	"github.com/mcc-github/blockchain/token/transaction"
	"github.com/pkg/errors"
)





type FabricIdentityDeserializerManager struct {
}

func (*FabricIdentityDeserializerManager) Deserializer(channel string) (identity.Deserializer, error) {
	id, ok := mgmt.GetDeserializers()[channel]
	if !ok {
		return nil, errors.New("channel not found")
	}
	return id, nil
}


type Manager struct {
	IdentityDeserializerManager identity.DeserializerManager
}


func (m *Manager) GetTxProcessor(channel string) (transaction.TMSTxProcessor, error) {
	identityDeserializerManager, err := m.IdentityDeserializerManager.Deserializer(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting identity deserialiser manager for channel '%s'", channel)
	}

	return &plain.Verifier{
		IssuingValidator:    &AllIssuingValidator{Deserializer: identityDeserializerManager},
		TokenOwnerValidator: &FabricTokenOwnerValidator{Deserializer: identityDeserializerManager},
	}, nil
}
