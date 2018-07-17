/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msptesttools

import (
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/config/configtest"
	"github.com/mcc-github/blockchain/msp"
	"github.com/mcc-github/blockchain/msp/mgmt"
)



func LoadMSPSetupForTesting() error {
	dir, err := configtest.GetDevMspDir()
	if err != nil {
		return err
	}
	conf, err := msp.GetLocalMspConfig(dir, nil, "SampleOrg")
	if err != nil {
		return err
	}

	err = mgmt.GetLocalMSP().Setup(conf)
	if err != nil {
		return err
	}

	err = mgmt.GetManagerForChain(util.GetTestChainID()).Setup([]msp.MSP{mgmt.GetLocalMSP()})
	if err != nil {
		return err
	}

	return nil
}


func LoadDevMsp() error {
	mspDir, err := configtest.GetDevMspDir()
	if err != nil {
		return err
	}

	return mgmt.LoadLocalMsp(mspDir, nil, "SampleOrg")
}
