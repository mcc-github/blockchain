/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)




func ConvertCDSToChaincodeInstallPackage(cds *peer.ChaincodeDeploymentSpec) (name, version string, cip *peer.ChaincodeInstallPackage, err error) {
	if cds.ChaincodeSpec == nil {
		return "", "", nil, errors.New("nil ChaincodeSpec")
	}

	if cds.ChaincodeSpec.ChaincodeId == nil {
		return "", "", nil, errors.New("nil ChaincodeId")
	}

	name = cds.ChaincodeSpec.ChaincodeId.Name
	version = cds.ChaincodeSpec.ChaincodeId.Version

	cip = &peer.ChaincodeInstallPackage{
		CodePackage: cds.CodePackage,
		Path:        cds.ChaincodeSpec.ChaincodeId.Path,
		Type:        cds.ChaincodeSpec.Type.String(),
	}

	return name, version, cip, nil
}




func ConvertSignedCDSToChaincodeInstallPackage(scds *peer.SignedChaincodeDeploymentSpec) (name, version string, cip *peer.ChaincodeInstallPackage, err error) {
	cds, err := UnmarshalChaincodeDeploymentSpec(scds.ChaincodeDeploymentSpec)
	if err != nil {
		return "", "", nil, err
	}
	return ConvertCDSToChaincodeInstallPackage(cds)
}



func UnmarshalChaincodeDeploymentSpec(cdsBytes []byte) (*peer.ChaincodeDeploymentSpec, error) {
	cds := &peer.ChaincodeDeploymentSpec{}
	err := proto.Unmarshal(cdsBytes, cds)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeDeploymentSpec")
	}

	return cds, nil
}
