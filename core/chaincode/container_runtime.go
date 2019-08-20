/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	"github.com/pkg/errors"
)


type CertGenerator interface {
	
	
	Generate(ccName string) (*accesscontrol.CertAndPrivKeyPair, error)
}






type ContainerRouter interface {
	Build(ccid string) error
	Start(ccid string, peerConnection *ccintf.PeerConnection) error
	Stop(ccid string) error
	Wait(ccid string) (int, error)
}


type ContainerRuntime struct {
	CertGenerator   CertGenerator
	ContainerRouter ContainerRouter
	CACert          []byte
	PeerAddress     string
}


func (c *ContainerRuntime) Start(ccid string) error {
	var tlsConfig *ccintf.TLSConfig
	if c.CertGenerator != nil {
		certKeyPair, err := c.CertGenerator.Generate(string(ccid))
		if err != nil {
			return errors.WithMessagef(err, "failed to generate TLS certificates for %s", ccid)
		}

		tlsConfig = &ccintf.TLSConfig{
			ClientCert: certKeyPair.Cert,
			ClientKey:  certKeyPair.Key,
			RootCert:   c.CACert,
		}
	}

	if err := c.ContainerRouter.Build(ccid); err != nil {
		return errors.WithMessage(err, "error building image")
	}

	chaincodeLogger.Debugf("start container: %s", ccid)

	if err := c.ContainerRouter.Start(
		ccid,
		&ccintf.PeerConnection{
			Address:   c.PeerAddress,
			TLSConfig: tlsConfig,
		},
	); err != nil {
		return errors.WithMessage(err, "error starting container")
	}

	return nil
}


func (c *ContainerRuntime) Stop(ccid string) error {
	if err := c.ContainerRouter.Stop(ccid); err != nil {
		return errors.WithMessage(err, "error stopping container")
	}

	return nil
}


func (c *ContainerRuntime) Wait(ccid string) (int, error) {
	return c.ContainerRouter.Wait(ccid)
}
