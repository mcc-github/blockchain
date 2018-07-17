/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package msp

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)

var mspLogger = flogging.MustGetLogger("msp")

type mspManagerImpl struct {
	
	mspsMap map[string]MSP

	
	mspsByProviders map[ProviderType][]MSP

	
	up bool
}




func NewMSPManager() MSPManager {
	return &mspManagerImpl{}
}


func (mgr *mspManagerImpl) Setup(msps []MSP) error {
	if mgr.up {
		mspLogger.Infof("MSP manager already up")
		return nil
	}

	mspLogger.Debugf("Setting up the MSP manager (%d msps)", len(msps))

	
	mgr.mspsMap = make(map[string]MSP)

	
	mgr.mspsByProviders = make(map[ProviderType][]MSP)

	for _, msp := range msps {
		
		mspID, err := msp.GetIdentifier()
		if err != nil {
			return errors.WithMessage(err, "could not extract msp identifier")
		}
		mgr.mspsMap[mspID] = msp
		providerType := msp.GetType()
		mgr.mspsByProviders[providerType] = append(mgr.mspsByProviders[providerType], msp)
	}

	mgr.up = true

	mspLogger.Debugf("MSP manager setup complete, setup %d msps", len(msps))

	return nil
}


func (mgr *mspManagerImpl) GetMSPs() (map[string]MSP, error) {
	return mgr.mspsMap, nil
}


func (mgr *mspManagerImpl) DeserializeIdentity(serializedID []byte) (Identity, error) {
	
	sId := &msp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sId)
	if err != nil {
		return nil, errors.Wrap(err, "could not deserialize a SerializedIdentity")
	}

	
	msp := mgr.mspsMap[sId.Mspid]
	if msp == nil {
		return nil, errors.Errorf("MSP %s is unknown", sId.Mspid)
	}

	switch t := msp.(type) {
	case *bccspmsp:
		return t.deserializeIdentityInternal(sId.IdBytes)
	case *idemixmsp:
		return t.deserializeIdentityInternal(sId.IdBytes)
	default:
		return t.DeserializeIdentity(serializedID)
	}
}

func (mgr *mspManagerImpl) IsWellFormed(identity *msp.SerializedIdentity) error {
	
	
	for _, mspList := range mgr.mspsByProviders {
		
		msp := mspList[0]
		if err := msp.IsWellFormed(identity); err == nil {
			return nil
		}
	}
	return errors.New("no MSP provider recognizes the identity")
}
