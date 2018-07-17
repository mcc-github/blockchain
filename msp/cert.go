/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package msp

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/mcc-github/blockchain/bccsp/utils"
	"github.com/pkg/errors"
)

type validity struct {
	NotBefore, NotAfter time.Time
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

type certificate struct {
	Raw                asn1.RawContent
	TBSCertificate     tbsCertificate
	SignatureAlgorithm pkix.AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type tbsCertificate struct {
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	UniqueId           asn1.BitString   `asn1:"optional,tag:1"`
	SubjectUniqueId    asn1.BitString   `asn1:"optional,tag:2"`
	Extensions         []pkix.Extension `asn1:"optional,explicit,tag:3"`
}

func isECDSASignedCert(cert *x509.Certificate) bool {
	return cert.SignatureAlgorithm == x509.ECDSAWithSHA1 ||
		cert.SignatureAlgorithm == x509.ECDSAWithSHA256 ||
		cert.SignatureAlgorithm == x509.ECDSAWithSHA384 ||
		cert.SignatureAlgorithm == x509.ECDSAWithSHA512
}





func sanitizeECDSASignedCert(cert *x509.Certificate, parentCert *x509.Certificate) (*x509.Certificate, error) {
	if cert == nil {
		return nil, errors.New("certificate must be different from nil")
	}
	if parentCert == nil {
		return nil, errors.New("parent certificate must be different from nil")
	}

	expectedSig, err := utils.SignatureToLowS(parentCert.PublicKey.(*ecdsa.PublicKey), cert.Signature)
	if err != nil {
		return nil, err
	}

	
	if bytes.Equal(cert.Signature, expectedSig) {
		return cert, nil
	}
	

	
	
	
	var newCert certificate
	newCert, err = certFromX509Cert(cert)
	if err != nil {
		return nil, err
	}

	
	newCert.SignatureValue = asn1.BitString{Bytes: expectedSig, BitLength: len(expectedSig) * 8}

	
	newCert.Raw = nil
	newRaw, err := asn1.Marshal(newCert)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling of the certificate failed")
	}

	
	return x509.ParseCertificate(newRaw)
}

func certFromX509Cert(cert *x509.Certificate) (certificate, error) {
	var newCert certificate
	_, err := asn1.Unmarshal(cert.Raw, &newCert)
	if err != nil {
		return certificate{}, errors.Wrap(err, "unmarshalling of the certificate failed")
	}
	return newCert, nil
}


func (c certificate) String() string {
	b, err := asn1.Marshal(c)
	if err != nil {
		return fmt.Sprintf("Failed marshaling cert: %v", err)
	}
	block := &pem.Block{
		Bytes: b,
		Type:  "CERTIFICATE",
	}
	b = pem.EncodeToMemory(block)
	return string(b)
}



func certToPEM(certificate *x509.Certificate) string {
	cert, err := certFromX509Cert(certificate)
	if err != nil {
		mspIdentityLogger.Warning("Failed converting certificate to asn1", err)
		return ""
	}
	return cert.String()
}
