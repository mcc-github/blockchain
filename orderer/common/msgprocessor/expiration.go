/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"time"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/crypto"
	"github.com/mcc-github/blockchain/protos/common"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/pkg/errors"
)

type resources interface {
	
	
	OrdererConfig() (channelconfig.Orderer, bool)
}



func NewExpirationRejectRule(filterSupport resources) Rule {
	return &expirationRejectRule{filterSupport: filterSupport}
}

type expirationRejectRule struct {
	filterSupport resources
}


func (exp *expirationRejectRule) Apply(message *common.Envelope) error {
	ordererConf, ok := exp.filterSupport.OrdererConfig()
	if !ok {
		logger.Panic("Programming error: orderer config not found")
	}
	if !ordererConf.Capabilities().ExpirationCheck() {
		return nil
	}
	signedData, err := protoutil.EnvelopeAsSignedData(message)

	if err != nil {
		return errors.Errorf("could not convert message to signedData: %s", err)
	}
	expirationTime := crypto.ExpiresAt(signedData[0].Identity)
	
	if expirationTime.IsZero() || time.Now().Before(expirationTime) {
		return nil
	}
	return errors.New("identity expired")
}
