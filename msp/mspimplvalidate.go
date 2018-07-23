/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"bytes"
	"crypto/x509"

	"time"

	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"reflect"

	"github.com/pkg/errors"
)

func (msp *bccspmsp) validateIdentity(id *identity) error {
	validationChain, err := msp.getCertificationChainForBCCSPIdentity(id)
	if err != nil {
		return errors.WithMessage(err, "could not obtain certification chain")
	}

	err = msp.validateIdentityAgainstChain(id, validationChain)
	if err != nil {
		return errors.WithMessage(err, "could not validate identity against certification chain")
	}

	err = msp.internalValidateIdentityOusFunc(id)
	if err != nil {
		return errors.WithMessage(err, "could not validate identity's OUs")
	}

	return nil
}

func (msp *bccspmsp) validateCAIdentity(id *identity) error {
	if !id.cert.IsCA {
		return errors.New("Only CA identities can be validated")
	}

	validationChain, err := msp.getUniqueValidationChain(id.cert, msp.getValidityOptsForCert(id.cert))
	if err != nil {
		return errors.WithMessage(err, "could not obtain certification chain")
	}
	if len(validationChain) == 1 {
		
		return nil
	}

	return msp.validateIdentityAgainstChain(id, validationChain)
}

func (msp *bccspmsp) validateTLSCAIdentity(cert *x509.Certificate, opts *x509.VerifyOptions) error {
	if !cert.IsCA {
		return errors.New("Only CA identities can be validated")
	}

	validationChain, err := msp.getUniqueValidationChain(cert, *opts)
	if err != nil {
		return errors.WithMessage(err, "could not obtain certification chain")
	}
	if len(validationChain) == 1 {
		
		return nil
	}

	return msp.validateCertAgainstChain(cert, validationChain)
}

func (msp *bccspmsp) validateIdentityAgainstChain(id *identity, validationChain []*x509.Certificate) error {
	return msp.validateCertAgainstChain(id.cert, validationChain)
}

func (msp *bccspmsp) validateCertAgainstChain(cert *x509.Certificate, validationChain []*x509.Certificate) error {
	

	
	SKI, err := getSubjectKeyIdentifierFromCert(validationChain[1])
	if err != nil {
		return errors.WithMessage(err, "could not obtain Subject Key Identifier for signer cert")
	}

	
	
	for _, crl := range msp.CRL {
		aki, err := getAuthorityKeyIdentifierFromCrl(crl)
		if err != nil {
			return errors.WithMessage(err, "could not obtain Authority Key Identifier for crl")
		}

		
		if bytes.Equal(aki, SKI) {
			
			for _, rc := range crl.TBSCertList.RevokedCertificates {
				if rc.SerialNumber.Cmp(cert.SerialNumber) == 0 {
					
					
					
					
					
					err = validationChain[1].CheckCRLSignature(crl)
					if err != nil {
						
						
						
						mspLogger.Warningf("Invalid signature over the identified CRL, error %+v", err)
						continue
					}

					
					
					
					
					
					
					return errors.New("The certificate has been revoked")
				}
			}
		}
	}

	return nil
}

func (msp *bccspmsp) validateIdentityOUsV1(id *identity) error {
	
	
	if len(msp.ouIdentifiers) > 0 {
		found := false

		for _, OU := range id.GetOrganizationalUnits() {
			certificationIDs, exists := msp.ouIdentifiers[OU.OrganizationalUnitIdentifier]

			if exists {
				for _, certificationID := range certificationIDs {
					if bytes.Equal(certificationID, OU.CertifiersIdentifier) {
						found = true
						break
					}
				}
			}
		}

		if !found {
			if len(id.GetOrganizationalUnits()) == 0 {
				return errors.New("the identity certificate does not contain an Organizational Unit (OU)")
			}
			return errors.Errorf("none of the identity's organizational units [%v] are in MSP %s", id.GetOrganizationalUnits(), msp.name)
		}
	}

	return nil
}

func (msp *bccspmsp) validateIdentityOUsV11(id *identity) error {
	
	err := msp.validateIdentityOUsV1(id)
	if err != nil {
		return err
	}

	
	
	
	if !msp.ouEnforcement {
		
		return nil
	}

	
	
	counter := 0
	for _, OU := range id.GetOrganizationalUnits() {
		
		var nodeOU *OUIdentifier
		switch OU.OrganizationalUnitIdentifier {
		case msp.clientOU.OrganizationalUnitIdentifier:
			nodeOU = msp.clientOU
		case msp.peerOU.OrganizationalUnitIdentifier:
			nodeOU = msp.peerOU
		default:
			continue
		}

		
		
		if len(nodeOU.CertifiersIdentifier) != 0 && !bytes.Equal(nodeOU.CertifiersIdentifier, OU.CertifiersIdentifier) {
			return errors.Errorf("certifiersIdentifier does not match: [%v], MSP: [%s]", id.GetOrganizationalUnits(), msp.name)
		}
		counter++
		if counter > 1 {
			break
		}
	}
	if counter != 1 {
		return errors.Errorf("the identity must be a client, a peer or an orderer identity to be valid, not a combination of them. OUs: [%v], MSP: [%s]", id.GetOrganizationalUnits(), msp.name)
	}

	return nil
}

func (msp *bccspmsp) getValidityOptsForCert(cert *x509.Certificate) x509.VerifyOptions {
	
	
	
	

	var tempOpts x509.VerifyOptions
	tempOpts.Roots = msp.opts.Roots
	tempOpts.DNSName = msp.opts.DNSName
	tempOpts.Intermediates = msp.opts.Intermediates
	tempOpts.KeyUsages = msp.opts.KeyUsages
	tempOpts.CurrentTime = cert.NotBefore.Add(time.Second)

	return tempOpts
}



type authorityKeyIdentifier struct {
	KeyIdentifier             []byte  `asn1:"optional,tag:0"`
	AuthorityCertIssuer       []byte  `asn1:"optional,tag:1"`
	AuthorityCertSerialNumber big.Int `asn1:"optional,tag:2"`
}




func getAuthorityKeyIdentifierFromCrl(crl *pkix.CertificateList) ([]byte, error) {
	aki := authorityKeyIdentifier{}

	for _, ext := range crl.TBSCertList.Extensions {
		
		
		if reflect.DeepEqual(ext.Id, asn1.ObjectIdentifier{2, 5, 29, 35}) {
			_, err := asn1.Unmarshal(ext.Value, &aki)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal AKI")
			}

			return aki.KeyIdentifier, nil
		}
	}

	return nil, errors.New("authorityKeyIdentifier not found in certificate")
}



func getSubjectKeyIdentifierFromCert(cert *x509.Certificate) ([]byte, error) {
	var SKI []byte

	for _, ext := range cert.Extensions {
		
		
		if reflect.DeepEqual(ext.Id, asn1.ObjectIdentifier{2, 5, 29, 14}) {
			_, err := asn1.Unmarshal(ext.Value, &SKI)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal Subject Key Identifier")
			}

			return SKI, nil
		}
	}

	return nil, errors.New("subjectKeyIdentifier not found in certificate")
}




func isCACert(cert *x509.Certificate) bool {
	_, err := getSubjectKeyIdentifierFromCert(cert)
	if err != nil {
		return false
	}

	if !cert.IsCA {
		return false
	}

	return true
}