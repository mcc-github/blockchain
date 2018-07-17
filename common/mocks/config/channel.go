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

package config

import (
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/msp"
)

func nearIdentityHash(input []byte) []byte {
	return util.ConcatenateBytes([]byte("FakeHash("), input, []byte(""))
}


type Channel struct {
	
	HashingAlgorithmVal func([]byte) []byte
	
	BlockDataHashingStructureWidthVal uint32
	
	OrdererAddressesVal []string
	
	CapabilitiesVal channelconfig.ChannelCapabilities
}


func (scm *Channel) HashingAlgorithm() func([]byte) []byte {
	if scm.HashingAlgorithmVal == nil {
		return nearIdentityHash
	}
	return scm.HashingAlgorithmVal
}


func (scm *Channel) BlockDataHashingStructureWidth() uint32 {
	return scm.BlockDataHashingStructureWidthVal
}


func (scm *Channel) OrdererAddresses() []string {
	return scm.OrdererAddressesVal
}


func (scm *Channel) Capabilities() channelconfig.ChannelCapabilities {
	return scm.CapabilitiesVal
}


type ChannelCapabilities struct {
	
	SupportedErr error

	
	MSPVersionVal msp.MSPVersion
}


func (cc *ChannelCapabilities) Supported() error {
	return cc.SupportedErr
}


func (cc *ChannelCapabilities) MSPVersion() msp.MSPVersion {
	return cc.MSPVersionVal
}
