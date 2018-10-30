/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/mcc-github/blockchain/token/tms/plain"
)



type Manager struct {
}




func (t *Manager) GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error) {
	return &plain.Issuer{}, nil
}
