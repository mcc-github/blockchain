/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package lifecycle

import (
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/pkg/errors"
)


type ChaincodeStore interface {
	Save(name, version string, ccInstallPkg []byte) (hash []byte, err error)
}

type PackageParser interface {
	Parse(data []byte) (*persistence.ChaincodePackage, error)
}



type Lifecycle struct {
	ChaincodeStore ChaincodeStore
	PackageParser  PackageParser
}



func (l *Lifecycle) InstallChaincode(name, version string, chaincodeInstallPackage []byte) ([]byte, error) {
	
	_, err := l.PackageParser.Parse(chaincodeInstallPackage)
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse as a chaincode install package")
	}

	hash, err := l.ChaincodeStore.Save(name, version, chaincodeInstallPackage)
	if err != nil {
		return nil, errors.WithMessage(err, "could not save cc install package")
	}

	return hash, nil
}
