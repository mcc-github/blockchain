/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtxgentest

import (
	"fmt"

	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/internal/configtxgen/localconfig"
)

func Load(profile string) *localconfig.Profile {
	devConfigDir, err := configtest.GetDevConfigDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get dev config dir: %s", err))
	}
	return localconfig.Load(profile, devConfigDir)
}

func LoadTopLevel() *localconfig.TopLevel {
	devConfigDir, err := configtest.GetDevConfigDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get dev config dir: %s", err))
	}
	return localconfig.LoadTopLevel(devConfigDir)
}
