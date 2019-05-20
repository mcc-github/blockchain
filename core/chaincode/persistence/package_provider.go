/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"io/ioutil"

	"github.com/mcc-github/blockchain/common/chaincode"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/pkg/errors"
)



type StorePackageProvider interface {
	GetChaincodeInstallPath() string
	ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error)
	Load(packageID persistence.PackageID) ([]byte, error)
}



type LegacyPackageProvider interface {
	GetChaincodeCodePackage(name, version string) (codePackage []byte, err error)
	ListInstalledChaincodes(dir string, de ccprovider.DirEnumerator, ce ccprovider.ChaincodeExtractor) ([]chaincode.InstalledChaincode, error)
}


type PackageParser interface {
	Parse(data []byte) (*ChaincodePackage, error)
}



type PackageProvider struct {
	Store    StorePackageProvider
	Parser   PackageParser
	LegacyPP LegacyPackageProvider
}





func (p *PackageProvider) GetChaincodeCodePackage(ccci *ccprovider.ChaincodeContainerInfo) ([]byte, error) {
	codePackage, err := p.getCodePackageFromStore(ccci.PackageID)
	if err == nil {
		return codePackage, nil
	}
	if _, ok := err.(*CodePackageNotFoundErr); !ok {
		
		
		return nil, err
	}

	codePackage, err = p.getCodePackageFromLegacyPP(ccci.Name, ccci.Version)
	if err != nil {
		logger.Debug(err.Error())
		err = errors.Errorf("code package not found for chaincode with name '%s', version '%s'", ccci.Name, ccci.Version)
		return nil, err
	}
	return codePackage, nil
}



func (p *PackageProvider) getCodePackageFromStore(packageID persistence.PackageID) ([]byte, error) {
	fsBytes, err := p.Store.Load(packageID)
	if _, ok := err.(*CodePackageNotFoundErr); ok {
		return nil, err
	}
	if err != nil {
		return nil, errors.WithMessage(err, "error loading code package from ChaincodeInstallPackage")
	}

	ccPackage, err := p.Parser.Parse(fsBytes)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing chaincode package")
	}

	return ccPackage.CodePackage, nil
}



func (p *PackageProvider) getCodePackageFromLegacyPP(name, version string) ([]byte, error) {
	codePackage, err := p.LegacyPP.GetChaincodeCodePackage(name, version)
	if err != nil {
		return nil, errors.Wrap(err, "error loading code package from ChaincodeDeploymentSpec")
	}
	return codePackage, nil
}



func (p *PackageProvider) ListInstalledChaincodesLegacy() ([]chaincode.InstalledChaincode, error) {
	installedChaincodesLegacy, err := p.LegacyPP.ListInstalledChaincodes(p.Store.GetChaincodeInstallPath(), ioutil.ReadDir, ccprovider.LoadPackage)
	if err != nil {
		return nil, err
	}

	return installedChaincodesLegacy, nil
}
