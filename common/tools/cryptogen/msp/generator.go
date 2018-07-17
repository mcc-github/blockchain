/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package msp

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/tools/cryptogen/ca"
	"github.com/mcc-github/blockchain/common/tools/cryptogen/csp"
	blockchainmsp "github.com/mcc-github/blockchain/msp"
)

const (
	CLIENT = iota
	ORDERER
	PEER
)

const (
	CLIENTOU = "client"
	PEEROU   = "peer"
)

var nodeOUMap = map[int]string{
	CLIENT: CLIENTOU,
	PEER:   PEEROU,
}

func GenerateLocalMSP(baseDir, name string, sans []string, signCA *ca.CA,
	tlsCA *ca.CA, nodeType int, nodeOUs bool) error {

	
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

	
	priv, _, err := csp.GeneratePrivateKey(keystore)
	if err != nil {
		return err
	}

	
	ecPubKey, err := csp.GetECPublicKey(priv)
	if err != nil {
		return err
	}
	
	var ous []string
	if nodeOUs {
		ous = []string{nodeOUMap[nodeType]}
	}
	cert, err := signCA.SignCertificate(filepath.Join(mspDir, "signcerts"),
		name, ous, nil, ecPubKey, x509.KeyUsageDigitalSignature, []x509.ExtKeyUsage{})
	if err != nil {
		return err
	}

	

	
	err = x509Export(filepath.Join(mspDir, "cacerts", x509Filename(signCA.Name)), signCA.SignCert)
	if err != nil {
		return err
	}
	
	err = x509Export(filepath.Join(mspDir, "tlscacerts", x509Filename(tlsCA.Name)), tlsCA.SignCert)
	if err != nil {
		return err
	}

	
	if nodeOUs && nodeType == PEER {
		exportConfig(mspDir, "cacerts/"+x509Filename(signCA.Name), true)
	}

	
	
	
	
	
	
	
	err = x509Export(filepath.Join(mspDir, "admincerts", x509Filename(name)), cert)
	if err != nil {
		return err
	}

	

	
	tlsPrivKey, _, err := csp.GeneratePrivateKey(tlsDir)
	if err != nil {
		return err
	}
	
	tlsPubKey, err := csp.GetECPublicKey(tlsPrivKey)
	if err != nil {
		return err
	}
	
	_, err = tlsCA.SignCertificate(filepath.Join(tlsDir),
		name, nil, sans, tlsPubKey, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment,
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth})
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

	err = keyExport(tlsDir, filepath.Join(tlsDir, tlsFilePrefix+".key"), tlsPrivKey)
	if err != nil {
		return err
	}

	return nil
}

func GenerateVerifyingMSP(baseDir string, signCA *ca.CA, tlsCA *ca.CA, nodeOUs bool) error {

	
	err := createFolderStructure(baseDir, false)
	if err == nil {
		
		err = x509Export(filepath.Join(baseDir, "cacerts", x509Filename(signCA.Name)), signCA.SignCert)
		if err != nil {
			return err
		}
		
		err = x509Export(filepath.Join(baseDir, "tlscacerts", x509Filename(tlsCA.Name)), tlsCA.SignCert)
		if err != nil {
			return err
		}
	}

	
	if nodeOUs {
		exportConfig(baseDir, "cacerts/"+x509Filename(signCA.Name), true)
	}

	
	
	
	
	
	factory.InitFactories(nil)
	bcsp := factory.GetDefault()
	priv, err := bcsp.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: true})
	ecPubKey, err := csp.GetECPublicKey(priv)
	if err != nil {
		return err
	}
	_, err = signCA.SignCertificate(filepath.Join(baseDir, "admincerts"), signCA.Name,
		nil, nil, ecPubKey, x509.KeyUsageDigitalSignature, []x509.ExtKeyUsage{})
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

func keyExport(keystore, output string, key bccsp.Key) error {
	id := hex.EncodeToString(key.SKI())

	return os.Rename(filepath.Join(keystore, id+"_sk"), output)
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
