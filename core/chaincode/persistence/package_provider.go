/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"io/ioutil"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
)



type LegacyPackageProvider interface {
	GetChaincodeInstallPath() string
	ListInstalledChaincodes(dir string, de ccprovider.DirEnumerator, ce ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error)
}



type PackageProvider struct {
	LegacyPP LegacyPackageProvider
}



func (p *PackageProvider) ListInstalledChaincodesLegacy() ([]chaincode.InstalledChaincode, error) {
	installedChaincodesLegacy, err := p.LegacyPP.ListInstalledChaincodes(p.LegacyPP.GetChaincodeInstallPath(), ioutil.ReadDir, ccprovider.LoadPackage)
	if err != nil {
		return nil, err
	}

	return installedChaincodesLegacy, nil
}
