/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channelconfig

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/cache"
	mspprotos "github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)

type pendingMSPConfig struct {
	mspConfig *mspprotos.MSPConfig
	msp       msp.MSP
}


type MSPConfigHandler struct {
	version msp.MSPVersion
	idMap   map[string]*pendingMSPConfig
}

func NewMSPConfigHandler(mspVersion msp.MSPVersion) *MSPConfigHandler {
	return &MSPConfigHandler{
		version: mspVersion,
		idMap:   make(map[string]*pendingMSPConfig),
	}
}


func (bh *MSPConfigHandler) ProposeMSP(mspConfig *mspprotos.MSPConfig) (msp.MSP, error) {
	var theMsp msp.MSP
	var err error

	switch mspConfig.Type {
	case int32(msp.FABRIC):
		
		mspInst, err := msp.New(&msp.BCCSPNewOpts{NewBaseOpts: msp.NewBaseOpts{Version: bh.version}})
		if err != nil {
			return nil, errors.WithMessage(err, "creating the MSP manager failed")
		}

		
		theMsp, err = cache.New(mspInst)
		if err != nil {
			return nil, errors.WithMessage(err, "creating the MSP cache failed")
		}
	case int32(msp.IDEMIX):
		
		theMsp, err = msp.New(&msp.IdemixNewOpts{msp.NewBaseOpts{Version: bh.version}})
		if err != nil {
			return nil, errors.WithMessage(err, "creating the MSP manager failed")
		}
	default:
		return nil, errors.New(fmt.Sprintf("Setup error: unsupported msp type %d", mspConfig.Type))
	}

	
	err = theMsp.Setup(mspConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "setting up the MSP manager failed")
	}

	
	mspID, _ := theMsp.GetIdentifier()

	existingPendingMSPConfig, ok := bh.idMap[mspID]
	if ok && !proto.Equal(existingPendingMSPConfig.mspConfig, mspConfig) {
		return nil, errors.New(fmt.Sprintf("Attempted to define two different versions of MSP: %s", mspID))
	}

	if !ok {
		bh.idMap[mspID] = &pendingMSPConfig{
			mspConfig: mspConfig,
			msp:       theMsp,
		}
	}

	return theMsp, nil
}

func (bh *MSPConfigHandler) CreateMSPManager() (msp.MSPManager, error) {
	mspList := make([]msp.MSP, len(bh.idMap))
	i := 0
	for _, pendingMSP := range bh.idMap {
		mspList[i] = pendingMSP.msp
		i++
	}

	manager := msp.NewMSPManager()
	err := manager.Setup(mspList)
	return manager, err
}
