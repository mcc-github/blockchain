/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package handlers

import (
	"reflect"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/pkg/errors"
)


type CredentialRequestSigner struct {
	
	CredRequest CredRequest
}

func (c *CredentialRequestSigner) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) ([]byte, error) {
	userSecretKey, ok := k.(*userSecretKey)
	if !ok {
		return nil, errors.New("invalid key, expected *userSecretKey")
	}
	credentialRequestSignerOpts, ok := opts.(*bccsp.IdemixCredentialRequestSignerOpts)
	if !ok {
		return nil, errors.New("invalid options, expected *IdemixCredentialRequestSignerOpts")
	}
	if credentialRequestSignerOpts.IssuerPK == nil {
		return nil, errors.New("invalid options, missing issuer public key")
	}
	issuerPK, ok := credentialRequestSignerOpts.IssuerPK.(*issuerPublicKey)
	if !ok {
		return nil, errors.New("invalid options, expected IssuerPK as *issuerPublicKey")
	}
	if !reflect.DeepEqual(digest, bccsp.IdemixEmptyDigest()) {
		return nil, errors.New("invalid digest, the idemix empty digest is expected")
	}

	return c.CredRequest.Sign(userSecretKey.sk, issuerPK.pk)
}


type CredentialRequestVerifier struct {
	
	CredRequest CredRequest
}

func (c *CredentialRequestVerifier) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (bool, error) {
	issuerPublicKey, ok := k.(*issuerPublicKey)
	if !ok {
		return false, errors.New("invalid key, expected *issuerPublicKey")
	}
	if !reflect.DeepEqual(digest, bccsp.IdemixEmptyDigest()) {
		return false, errors.New("invalid digest, the idemix empty digest is expected")
	}

	err := c.CredRequest.Verify(signature, issuerPublicKey.pk)
	if err != nil {
		return false, err
	}

	return true, nil
}

type CredentialSigner struct {
	Credential Credential
}

func (s *CredentialSigner) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	issuerSecretKey, ok := k.(*issuerSecretKey)
	if !ok {
		return nil, errors.New("invalid key, expected *issuerSecretKey")
	}
	credOpts, ok := opts.(*bccsp.IdemixCredentialSignerOpts)
	if !ok {
		return nil, errors.New("invalid options, expected *IdemixCredentialSignerOpts")
	}

	signature, err = s.Credential.Sign(issuerSecretKey.sk, digest, credOpts.Attributes)
	if err != nil {
		return nil, err
	}

	return
}

type CredentialVerifier struct {
	Credential Credential
}

func (v *CredentialVerifier) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (valid bool, err error) {
	userSecretKey, ok := k.(*userSecretKey)
	if !ok {
		return false, errors.New("invalid key, expected *userSecretKey")
	}
	credOpts, ok := opts.(*bccsp.IdemixCredentialSignerOpts)
	if !ok {
		return false, errors.New("invalid options, expected *IdemixCredentialSignerOpts")
	}
	if credOpts.IssuerPK == nil {
		return false, errors.New("invalid options, missing issuer public key")
	}
	ipk, ok := credOpts.IssuerPK.(*issuerPublicKey)
	if !ok {
		return false, errors.New("invalid issuer public key, expected *issuerPublicKey")
	}
	if !reflect.DeepEqual(digest, bccsp.IdemixEmptyDigest()) {
		return false, errors.New("invalid digest, the idemix empty digest is expected")
	}
	if len(signature) == 0 {
		return false, errors.New("invalid signature, it must not be empty")
	}

	err = v.Credential.Verify(userSecretKey.sk, ipk.pk, signature, credOpts.Attributes)
	if err != nil {
		return false, err
	}

	return true, nil
}
