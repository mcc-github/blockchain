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

package common

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/common/util"
)





type SignedData struct {
	Data      []byte
	Identity  []byte
	Signature []byte
}


type Signable interface {
	
	AsSignedData() ([]*SignedData, error)
}


func (ce *ConfigUpdateEnvelope) AsSignedData() ([]*SignedData, error) {
	if ce == nil {
		return nil, fmt.Errorf("No signatures for nil SignedConfigItem")
	}

	result := make([]*SignedData, len(ce.Signatures))
	for i, configSig := range ce.Signatures {
		sigHeader := &SignatureHeader{}
		err := proto.Unmarshal(configSig.SignatureHeader, sigHeader)
		if err != nil {
			return nil, err
		}

		result[i] = &SignedData{
			Data:      util.ConcatenateBytes(configSig.SignatureHeader, ce.ConfigUpdate),
			Identity:  sigHeader.Creator,
			Signature: configSig.Signature,
		}

	}

	return result, nil
}


func (env *Envelope) AsSignedData() ([]*SignedData, error) {
	if env == nil {
		return nil, fmt.Errorf("No signatures for nil Envelope")
	}

	payload := &Payload{}
	err := proto.Unmarshal(env.Payload, payload)
	if err != nil {
		return nil, err
	}

	if payload.Header == nil  {
		return nil, fmt.Errorf("Missing Header")
	}

	shdr := &SignatureHeader{}
	err = proto.Unmarshal(payload.Header.SignatureHeader, shdr)
	if err != nil {
		return nil, fmt.Errorf("GetSignatureHeaderFromBytes failed, err %s", err)
	}

	return []*SignedData{{
		Data:      env.Payload,
		Identity:  shdr.Creator,
		Signature: env.Signature,
	}}, nil
}
