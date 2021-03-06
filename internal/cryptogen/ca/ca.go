/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/internal/cryptogen/csp"
	"github.com/pkg/errors"
)

type CA struct {
	Name               string
	Country            string
	Province           string
	Locality           string
	OrganizationalUnit string
	StreetAddress      string
	PostalCode         string
	Signer             crypto.Signer
	SignCert           *x509.Certificate
}



func NewCA(
	baseDir,
	org,
	name,
	country,
	province,
	locality,
	orgUnit,
	streetAddress,
	postalCode string,
) (*CA, error) {

	var ca *CA

	err := os.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, err
	}

	priv, err := csp.GeneratePrivateKey(baseDir)
	if err != nil {
		return nil, err
	}

	template := x509Template()
	
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageDigitalSignature |
		x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign
	template.ExtKeyUsage = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageServerAuth,
	}

	
	subject := subjectTemplateAdditional(country, province, locality, orgUnit, streetAddress, postalCode)
	subject.Organization = []string{org}
	subject.CommonName = name

	template.Subject = subject
	template.SubjectKeyId = computeSKI(priv)

	x509Cert, err := genCertificateECDSA(
		baseDir,
		name,
		&template,
		&template,
		&priv.PublicKey,
		priv,
	)
	if err != nil {
		return nil, err
	}
	ca = &CA{
		Name: name,
		Signer: &csp.ECDSASigner{
			PrivateKey: priv,
		},
		SignCert:           x509Cert,
		Country:            country,
		Province:           province,
		Locality:           locality,
		OrganizationalUnit: orgUnit,
		StreetAddress:      streetAddress,
		PostalCode:         postalCode,
	}

	return ca, err
}



func (ca *CA) SignCertificate(
	baseDir,
	name string,
	orgUnits,
	alternateNames []string,
	pub *ecdsa.PublicKey,
	ku x509.KeyUsage,
	eku []x509.ExtKeyUsage,
) (*x509.Certificate, error) {

	template := x509Template()
	template.KeyUsage = ku
	template.ExtKeyUsage = eku

	
	subject := subjectTemplateAdditional(
		ca.Country,
		ca.Province,
		ca.Locality,
		ca.OrganizationalUnit,
		ca.StreetAddress,
		ca.PostalCode,
	)
	subject.CommonName = name

	subject.OrganizationalUnit = append(subject.OrganizationalUnit, orgUnits...)

	template.Subject = subject
	for _, san := range alternateNames {
		
		ip := net.ParseIP(san)
		if ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, alternateNames...)
		}
	}

	cert, err := genCertificateECDSA(
		baseDir,
		name,
		&template,
		ca.SignCert,
		pub,
		ca.Signer,
	)

	if err != nil {
		return nil, err
	}

	return cert, nil
}


func computeSKI(privKey *ecdsa.PrivateKey) []byte {
	
	raw := elliptic.Marshal(privKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)

	
	hash := sha256.Sum256(raw)
	return hash[:]
}


func subjectTemplate() pkix.Name {
	return pkix.Name{
		Country:  []string{"US"},
		Locality: []string{"San Francisco"},
		Province: []string{"California"},
	}
}


func subjectTemplateAdditional(
	country,
	province,
	locality,
	orgUnit,
	streetAddress,
	postalCode string,
) pkix.Name {
	name := subjectTemplate()
	if len(country) >= 1 {
		name.Country = []string{country}
	}
	if len(province) >= 1 {
		name.Province = []string{province}
	}

	if len(locality) >= 1 {
		name.Locality = []string{locality}
	}
	if len(orgUnit) >= 1 {
		name.OrganizationalUnit = []string{orgUnit}
	}
	if len(streetAddress) >= 1 {
		name.StreetAddress = []string{streetAddress}
	}
	if len(postalCode) >= 1 {
		name.PostalCode = []string{postalCode}
	}
	return name
}


func x509Template() x509.Certificate {

	
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	
	expiry := 3650 * 24 * time.Hour
	
	notBefore := time.Now().Round(time.Minute).Add(-5 * time.Minute).UTC()

	
	x509 := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		BasicConstraintsValid: true,
	}
	return x509

}


func genCertificateECDSA(
	baseDir,
	name string,
	template,
	parent *x509.Certificate,
	pub *ecdsa.PublicKey,
	priv interface{},
) (*x509.Certificate, error) {

	
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	if err != nil {
		return nil, err
	}

	
	fileName := filepath.Join(baseDir, name+"-cert.pem")
	certFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certFile.Close()
	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}
	return x509Cert, nil
}


func LoadCertificateECDSA(certPath string) (*x509.Certificate, error) {
	var cert *x509.Certificate
	var err error

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".pem") {
			rawCert, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			block, _ := pem.Decode(rawCert)
			if block == nil || block.Type != "CERTIFICATE" {
				return errors.Errorf("%s: wrong PEM encoding", path)
			}
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return errors.Errorf("%s: wrong DER encoding", path)
			}
		}
		return nil
	}

	err = filepath.Walk(certPath, walkFunc)
	if err != nil {
		return nil, err
	}

	return cert, err
}
