/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/pkg/errors"
)


type CertGenerator interface {
	
	
	Generate(ccName string) (*accesscontrol.CertAndPrivKeyPair, error)
}






type ContainerRouter interface {
	Build(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error
	Start(ccid ccintf.CCID, peerConnection *ccintf.PeerConnection) error
	Stop(ccid ccintf.CCID) error
	Wait(ccid ccintf.CCID) (int, error)
}


type ContainerRuntime struct {
	CertGenerator   CertGenerator
	ContainerRouter ContainerRouter
	CACert          []byte
	PeerAddress     string
}


func (c *ContainerRuntime) Start(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error {
	packageID := ccci.PackageID.String()

	var tlsConfig *ccintf.TLSConfig
	if c.CertGenerator != nil {
		certKeyPair, err := c.CertGenerator.Generate(packageID)
		if err != nil {
			return errors.WithMessagef(err, "failed to generate TLS certificates for %s", packageID)
		}

		tlsConfig = &ccintf.TLSConfig{
			ClientCert: []byte(certKeyPair.Cert),
			ClientKey:  []byte(certKeyPair.Key),
			RootCert:   c.CACert,
		}
	}

	if err := c.ContainerRouter.Build(ccci, codePackage); err != nil {
		return errors.WithMessage(err, "error building image")
	}

	chaincodeLogger.Debugf("start container: %s", packageID)

	if err := c.ContainerRouter.Start(
		ccintf.New(ccci.PackageID),
		&ccintf.PeerConnection{
			Address:   c.PeerAddress,
			TLSConfig: tlsConfig,
		},
	); err != nil {
		return errors.WithMessage(err, "error starting container")
	}

	return nil
}


func (c *ContainerRuntime) Stop(ccci *ccprovider.ChaincodeContainerInfo) error {
	if err := c.ContainerRouter.Stop(ccintf.New(ccci.PackageID)); err != nil {
		return errors.WithMessage(err, "error stopping container")
	}

	return nil
}


func (c *ContainerRuntime) Wait(ccci *ccprovider.ChaincodeContainerInfo) (int, error) {
	return c.ContainerRouter.Wait(ccintf.New(ccci.PackageID))
}
