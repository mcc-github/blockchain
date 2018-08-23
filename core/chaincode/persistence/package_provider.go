/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"io/ioutil"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/pkg/errors"
)



type StorePackageProvider interface {
	GetChaincodeInstallPath() string
	ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error)
	Load(hash []byte) (codePackage []byte, name, version string, err error)
	RetrieveHash(name, version string) (hash []byte, err error)
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





func (p *PackageProvider) GetChaincodeCodePackage(name, version string) ([]byte, error) {
	codePackage, err := p.getCodePackageFromStore(name, version)
	if err == nil {
		return codePackage, nil
	}
	if _, ok := err.(*CodePackageNotFoundErr); !ok {
		
		
		return nil, err
	}

	codePackage, err = p.getCodePackageFromLegacyPP(name, version)
	if err != nil {
		logger.Debug(err.Error())
		err = errors.Errorf("code package not found for chaincode with name '%s', version '%s'", name, version)
		return nil, err
	}
	return codePackage, nil
}



func (p *PackageProvider) getCodePackageFromStore(name, version string) ([]byte, error) {
	hash, err := p.Store.RetrieveHash(name, version)
	if _, ok := err.(*CodePackageNotFoundErr); ok {
		return nil, err
	}
	if err != nil {
		return nil, errors.WithMessage(err, "error retrieving hash")
	}

	fsBytes, _, _, err := p.Store.Load(hash)
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



func (p *PackageProvider) ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error) {
	
	installedChaincodes, err := p.Store.ListInstalledChaincodes()

	if err != nil {
		
		logger.Debugf("error getting installed chaincodes from persistence store: %s", err)
	}

	
	installedChaincodesLegacy, err := p.LegacyPP.ListInstalledChaincodes(p.Store.GetChaincodeInstallPath(), ioutil.ReadDir, ccprovider.LoadPackage)

	if err != nil {
		
		logger.Debugf("error getting installed chaincodes from ccprovider: %s", err)
	}

	for _, cc := range installedChaincodesLegacy {
		installedChaincodes = append(installedChaincodes, cc)
	}

	return installedChaincodes, nil
}
