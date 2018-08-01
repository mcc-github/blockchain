/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"github.com/pkg/errors"
)



type StorePackageProvider interface {
	Load(hash []byte) (codePackage []byte, name, version string, err error)
	RetrieveHash(name, version string) (hash []byte, err error)
}



type LegacyPackageProvider interface {
	GetChaincodeCodePackage(name, version string) (codePackage []byte, err error)
}



type PackageProvider struct {
	Store    StorePackageProvider
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

	codePackage, _, _, err := p.Store.Load(hash)
	if err != nil {
		return nil, errors.WithMessage(err, "error loading code package from ChaincodeInstallPackage")
	}
	return codePackage, nil
}



func (p *PackageProvider) getCodePackageFromLegacyPP(name, version string) ([]byte, error) {
	codePackage, err := p.LegacyPP.GetChaincodeCodePackage(name, version)
	if err != nil {
		return nil, errors.Wrap(err, "error loading code package from ChaincodeDeploymentSpec")
	}
	return codePackage, nil
}
