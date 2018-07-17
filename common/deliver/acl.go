/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliver

import (
	"time"

	"github.com/mcc-github/blockchain/protos/common"
	"github.com/pkg/errors"
)


type ExpiresAtFunc func(identityBytes []byte) time.Time


type ConfigSequencer interface {
	Sequence() uint64
}



func NewSessionAC(chain ConfigSequencer, env *common.Envelope, policyChecker PolicyChecker, channelID string, expiresAt ExpiresAtFunc) (*SessionAccessControl, error) {
	signedData, err := env.AsSignedData()
	if err != nil {
		return nil, err
	}

	return &SessionAccessControl{
		envelope:       env,
		channelID:      channelID,
		sequencer:      chain,
		policyChecker:  policyChecker,
		sessionEndTime: expiresAt(signedData[0].Identity),
	}, nil
}




type SessionAccessControl struct {
	sequencer          ConfigSequencer
	policyChecker      PolicyChecker
	channelID          string
	envelope           *common.Envelope
	lastConfigSequence uint64
	sessionEndTime     time.Time
	usedAtLeastOnce    bool
}




func (ac *SessionAccessControl) Evaluate() error {
	if !ac.sessionEndTime.IsZero() && time.Now().After(ac.sessionEndTime) {
		return errors.Errorf("client identity expired %v before", time.Since(ac.sessionEndTime))
	}

	policyCheckNeeded := !ac.usedAtLeastOnce

	if currentConfigSequence := ac.sequencer.Sequence(); currentConfigSequence > ac.lastConfigSequence {
		ac.lastConfigSequence = currentConfigSequence
		policyCheckNeeded = true
	}

	if !policyCheckNeeded {
		return nil
	}

	ac.usedAtLeastOnce = true
	return ac.policyChecker.CheckPolicy(ac.envelope, ac.channelID)
}
