/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mgmt

import (
	"reflect"
	"sync"

	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/cache"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)


func LoadLocalMspWithType(dir string, bccspConfig *factory.FactoryOpts, mspID, mspType string) error {
	if mspID == "" {
		return errors.New("the local MSP must have an ID")
	}

	conf, err := msp.GetLocalMspConfigWithType(dir, bccspConfig, mspID, mspType)
	if err != nil {
		return err
	}

	return GetLocalMSP().Setup(conf)
}


func LoadLocalMsp(dir string, bccspConfig *factory.FactoryOpts, mspID string) error {
	if mspID == "" {
		return errors.New("the local MSP must have an ID")
	}

	conf, err := msp.GetLocalMspConfig(dir, bccspConfig, mspID)
	if err != nil {
		return err
	}

	return GetLocalMSP().Setup(conf)
}






var m sync.Mutex
var localMsp msp.MSP
var mspMap map[string]msp.MSPManager = make(map[string]msp.MSPManager)
var mspLogger = flogging.MustGetLogger("msp")




type mspMgmtMgr struct {
	msp.MSPManager
	
	up bool
}

func (mgr *mspMgmtMgr) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	if !mgr.up {
		return nil, errors.New("channel doesn't exist")
	}
	return mgr.MSPManager.DeserializeIdentity(serializedIdentity)
}

func (mgr *mspMgmtMgr) Setup(msps []msp.MSP) error {
	err := mgr.MSPManager.Setup(msps)
	if err == nil {
		mgr.up = true
	}
	return err
}



func GetManagerForChain(chainID string) msp.MSPManager {
	m.Lock()
	defer m.Unlock()

	mspMgr, ok := mspMap[chainID]
	if !ok {
		mspLogger.Debugf("Created new msp manager for channel `%s`", chainID)
		mspMgmtMgr := &mspMgmtMgr{msp.NewMSPManager(), false}
		mspMap[chainID] = mspMgmtMgr
		mspMgr = mspMgmtMgr
	} else {
		
		
		
		
		if !(reflect.TypeOf(mspMgr).Elem().Name() == "mspManagerImpl" || reflect.TypeOf(mspMgr).Elem().Name() == "mspMgmtMgr") {
			panic("Found unexpected MSPManager type.")
		}
		mspLogger.Debugf("Returning existing manager for channel '%s'", chainID)
	}
	return mspMgr
}


func GetDeserializers() map[string]msp.IdentityDeserializer {
	m.Lock()
	defer m.Unlock()

	clone := make(map[string]msp.IdentityDeserializer)

	for key, mspManager := range mspMap {
		clone[key] = mspManager
	}

	return clone
}




func XXXSetMSPManager(chainID string, manager msp.MSPManager) {
	m.Lock()
	defer m.Unlock()

	mspMap[chainID] = &mspMgmtMgr{manager, true}
}


func GetLocalMSP() msp.MSP {
	m.Lock()
	defer m.Unlock()

	if localMsp != nil {
		return localMsp
	}

	localMsp = loadLocaMSP()

	return localMsp
}

func loadLocaMSP() msp.MSP {
	
	mspType := viper.GetString("peer.localMspType")
	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}

	var mspOpts = map[string]msp.NewOpts{
		msp.ProviderTypeToString(msp.FABRIC): &msp.BCCSPNewOpts{NewBaseOpts: msp.NewBaseOpts{Version: msp.MSPv1_0}},
		msp.ProviderTypeToString(msp.IDEMIX): &msp.IdemixNewOpts{msp.NewBaseOpts{Version: msp.MSPv1_1}},
	}
	newOpts, found := mspOpts[mspType]
	if !found {
		mspLogger.Panicf("msp type " + mspType + " unknown")
	}

	mspInst, err := msp.New(newOpts)
	if err != nil {
		mspLogger.Fatalf("Failed to initialize local MSP, received err %+v", err)
	}
	switch mspType {
	case msp.ProviderTypeToString(msp.FABRIC):
		mspInst, err = cache.New(mspInst)
		if err != nil {
			mspLogger.Fatalf("Failed to initialize local MSP, received err %+v", err)
		}
	case msp.ProviderTypeToString(msp.IDEMIX):
		
	default:
		panic("msp type " + mspType + " unknown")
	}

	mspLogger.Debugf("Created new local MSP")

	return mspInst
}


func GetIdentityDeserializer(chainID string) msp.IdentityDeserializer {
	if chainID == "" {
		return GetLocalMSP()
	}

	return GetManagerForChain(chainID)
}



func GetLocalSigningIdentityOrPanic() msp.SigningIdentity {
	id, err := GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		mspLogger.Panicf("Failed getting local signing identity [%+v]", err)
	}
	return id
}
