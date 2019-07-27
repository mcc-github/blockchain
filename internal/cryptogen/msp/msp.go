/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package msp

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/mcc-github/blockchain/internal/cryptogen/ca"
	"github.com/mcc-github/blockchain/internal/cryptogen/csp"
	blockchainmsp "github.com/mcc-github/blockchain/msp"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	CLIENT = iota
	ORDERER
	PEER
	ADMIN
)

const (
	CLIENTOU  = "client"
	PEEROU    = "peer"
	ADMINOU   = "admin"
	ORDEREROU = "orderer"
)

var nodeOUMap = map[int]string{
	CLIENT:  CLIENTOU,
	PEER:    PEEROU,
	ADMIN:   ADMINOU,
	ORDERER: ORDEREROU,
}

func GenerateLocalMSP(
	baseDir,
	name string,
	sans []string,
	signCA *ca.CA,
	tlsCA *ca.CA,
	nodeType int,
	nodeOUs bool,
) error {

	
	mspDir := filepath.Join(baseDir, "msp")
	tlsDir := filepath.Join(baseDir, "tls")

	err := createFolderStructure(mspDir, true)
	if err != nil {
		return err
	}

	err = os.MkdirAll(tlsDir, 0755)
	if err != nil {
		return err
	}

	
	
	keystore := filepath.Join(mspDir, "keystore")

	
	priv, err := csp.GeneratePrivateKey(keystore)
	if err != nil {
		return err
	}

	
	var ous []string
	if nodeOUs {
		ous = []string{nodeOUMap[nodeType]}
	}
	cert, err := signCA.SignCertificate(
		filepath.Join(mspDir, "signcerts"),
		name,
		ous,
		nil,
		&priv.PublicKey,
		x509.KeyUsageDigitalSignature,
		[]x509.ExtKeyUsage{},
	)
	if err != nil {
		return err
	}

	

	
	err = x509Export(
		filepath.Join(mspDir, "cacerts", x509Filename(signCA.Name)),
		signCA.SignCert,
	)
	if err != nil {
		return err
	}
	
	err = x509Export(
		filepath.Join(mspDir, "tlscacerts", x509Filename(tlsCA.Name)),
		tlsCA.SignCert,
	)
	if err != nil {
		return err
	}

	
	if nodeOUs && (nodeType == PEER || nodeType == ORDERER) {

		exportConfig(mspDir, filepath.Join("cacerts", x509Filename(signCA.Name)), true)
	}

	
	
	
	
	
	
	
	err = x509Export(filepath.Join(mspDir, "admincerts", x509Filename(name)), cert)
	if err != nil {
		return err
	}

	

	
	tlsPrivKey, err := csp.GeneratePrivateKey(tlsDir)
	if err != nil {
		return err
	}

	
	_, err = tlsCA.SignCertificate(
		filepath.Join(tlsDir),
		name,
		nil,
		sans,
		&tlsPrivKey.PublicKey,
		x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment,
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth},
	)
	if err != nil {
		return err
	}
	err = x509Export(filepath.Join(tlsDir, "ca.crt"), tlsCA.SignCert)
	if err != nil {
		return err
	}

	
	tlsFilePrefix := "server"
	if nodeType == CLIENT {
		tlsFilePrefix = "client"
	}
	err = os.Rename(filepath.Join(tlsDir, x509Filename(name)),
		filepath.Join(tlsDir, tlsFilePrefix+".crt"))
	if err != nil {
		return err
	}

	err = keyExport(tlsDir, filepath.Join(tlsDir, tlsFilePrefix+".key"))
	if err != nil {
		return err
	}

	return nil
}

func GenerateVerifyingMSP(
	baseDir string,
	signCA,
	tlsCA *ca.CA,
	nodeOUs bool,
) error {

	
	err := createFolderStructure(baseDir, false)
	if err != nil {
		return err
	}
	
	err = x509Export(
		filepath.Join(baseDir, "cacerts", x509Filename(signCA.Name)),
		signCA.SignCert,
	)
	if err != nil {
		return err
	}
	
	err = x509Export(
		filepath.Join(baseDir, "tlscacerts", x509Filename(tlsCA.Name)),
		tlsCA.SignCert,
	)
	if err != nil {
		return err
	}

	
	if nodeOUs {
		exportConfig(baseDir, "cacerts/"+x509Filename(signCA.Name), true)
	}

	
	
	
	
	
	ksDir := filepath.Join(baseDir, "keystore")
	err = os.Mkdir(ksDir, 0755)
	defer os.RemoveAll(ksDir)
	if err != nil {
		return errors.WithMessage(err, "failed to create keystore directory")
	}
	priv, err := csp.GeneratePrivateKey(ksDir)
	if err != nil {
		return err
	}
	_, err = signCA.SignCertificate(
		filepath.Join(baseDir, "admincerts"),
		signCA.Name,
		nil,
		nil,
		&priv.PublicKey,
		x509.KeyUsageDigitalSignature,
		[]x509.ExtKeyUsage{},
	)
	if err != nil {
		return err
	}

	return nil
}

func createFolderStructure(rootDir string, local bool) error {

	var folders []string
	
	folders = []string{
		filepath.Join(rootDir, "admincerts"),
		filepath.Join(rootDir, "cacerts"),
		filepath.Join(rootDir, "tlscacerts"),
	}
	if local {
		folders = append(folders, filepath.Join(rootDir, "keystore"),
			filepath.Join(rootDir, "signcerts"))
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func x509Filename(name string) string {
	return name + "-cert.pem"
}

func x509Export(path string, cert *x509.Certificate) error {
	return pemExport(path, "CERTIFICATE", cert.Raw)
}

func keyExport(keystore, output string) error {
	return os.Rename(filepath.Join(keystore, "priv_sk"), output)
}

func pemExport(path, pemType string, bytes []byte) error {
	
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, &pem.Block{Type: pemType, Bytes: bytes})
}

func exportConfig(mspDir, caFile string, enable bool) error {
	var config = &blockchainmsp.Configuration{
		NodeOUs: &blockchainmsp.NodeOUs{
			Enable: enable,
			ClientOUIdentifier: &blockchainmsp.OrganizationalUnitIdentifiersConfiguration{
				Certificate:                  caFile,
				OrganizationalUnitIdentifier: CLIENTOU,
			},
			PeerOUIdentifier: &blockchainmsp.OrganizationalUnitIdentifiersConfiguration{
				Certificate:                  caFile,
				OrganizationalUnitIdentifier: PEEROU,
			},
			AdminOUIdentifier: &blockchainmsp.OrganizationalUnitIdentifiersConfiguration{
				Certificate:                  caFile,
				OrganizationalUnitIdentifier: ADMINOU,
			},
			OrdererOUIdentifier: &blockchainmsp.OrganizationalUnitIdentifiersConfiguration{
				Certificate:                  caFile,
				OrganizationalUnitIdentifier: ORDEREROU,
			},
		},
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(mspDir, "config.yaml"))
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.WriteString(string(configBytes))

	return err
}
