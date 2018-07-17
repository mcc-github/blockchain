/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package entities

import (
	"encoding/json"

	"github.com/pkg/errors"
)




type SignedMessage struct {
	
	ID []byte `json:"id"`

	
	Payload []byte `json:"payload"`

	
	Sig []byte `json:"sig"`
}


func (m *SignedMessage) Sign(signer Signer) error {
	if signer == nil {
		return errors.New("nil signer")
	}

	m.Sig = nil
	bytes, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "sign error: json.Marshal returned")
	}
	sig, err := signer.Sign(bytes)
	if err != nil {
		return errors.WithMessage(err, "sign error: signer.Sign returned")
	}
	m.Sig = sig

	return nil
}


func (m *SignedMessage) Verify(verifier Signer) (bool, error) {
	if verifier == nil {
		return false, errors.New("nil verifier")
	}

	sig := m.Sig
	m.Sig = nil
	defer func() {
		m.Sig = sig
	}()

	bytes, err := json.Marshal(m)
	if err != nil {
		return false, errors.Wrap(err, "sign error: json.Marshal returned")
	}

	return verifier.Verify(sig, bytes)
}


func (m *SignedMessage) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}


func (m *SignedMessage) FromBytes(d []byte) error {
	return json.Unmarshal(d, m)
}
