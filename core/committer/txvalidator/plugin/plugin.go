/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package plugin

import validation "github.com/mcc-github/blockchain/core/handlers/validation/api"


type Name string



type Mapper interface {
	FactoryByName(name Name) validation.PluginFactory
}


type MapBasedMapper map[string]validation.PluginFactory


func (m MapBasedMapper) FactoryByName(name Name) validation.PluginFactory {
	return m[string(name)]
}


type SerializedPolicy []byte


func (sp SerializedPolicy) Bytes() []byte {
	return sp
}
