/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/mcc-github/blockchain/core/handlers/endorsement/api"
	"github.com/mcc-github/blockchain/core/handlers/endorsement/builtin"
	"github.com/mcc-github/blockchain/integration/pluggable"
)




func NewPluginFactory() endorsement.PluginFactory {
	pluggable.PublishEndorsementPluginActivation()
	return &builtin.DefaultEndorsementFactory{}
}
