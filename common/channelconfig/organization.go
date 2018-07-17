/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"

	"github.com/mcc-github/blockchain/msp"
	cb "github.com/mcc-github/blockchain/protos/common"
	mspprotos "github.com/mcc-github/blockchain/protos/msp"

	"github.com/pkg/errors"
)

const (
	
	MSPKey = "MSP"
)


type OrganizationProtos struct {
	MSP *mspprotos.MSPConfig
}


type OrganizationConfig struct {
	protos *OrganizationProtos

	mspConfigHandler *MSPConfigHandler
	msp              msp.MSP
	mspID            string
	name             string
}


func NewOrganizationConfig(name string, orgGroup *cb.ConfigGroup, mspConfigHandler *MSPConfigHandler) (*OrganizationConfig, error) {
	if len(orgGroup.Groups) > 0 {
		return nil, fmt.Errorf("organizations do not support sub-groups")
	}

	oc := &OrganizationConfig{
		protos:           &OrganizationProtos{},
		name:             name,
		mspConfigHandler: mspConfigHandler,
	}

	if err := DeserializeProtoValuesFromGroup(orgGroup, oc.protos); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize values")
	}

	if err := oc.Validate(); err != nil {
		return nil, err
	}

	return oc, nil
}


func (oc *OrganizationConfig) Name() string {
	return oc.name
}


func (oc *OrganizationConfig) MSPID() string {
	return oc.mspID
}


func (oc *OrganizationConfig) Validate() error {
	return oc.validateMSP()
}

func (oc *OrganizationConfig) validateMSP() error {
	var err error

	logger.Debugf("Setting up MSP for org %s", oc.name)
	oc.msp, err = oc.mspConfigHandler.ProposeMSP(oc.protos.MSP)
	if err != nil {
		return err
	}

	oc.mspID, _ = oc.msp.GetIdentifier()

	if oc.mspID == "" {
		return fmt.Errorf("MSP for org %s has empty MSP ID", oc.name)
	}

	return nil
}
