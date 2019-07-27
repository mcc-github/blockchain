/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package csp

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)



func LoadPrivateKey(keystorePath string) (*ecdsa.PrivateKey, error) {
	var priv *ecdsa.PrivateKey

	walkFunc := func(path string, info os.FileInfo, pathErr error) error {

		if !strings.HasSuffix(path, "_sk") {
			return nil
		}

		rawKey, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		priv, err = parsePrivateKeyPEM(rawKey)
		if err != nil {
			return errors.WithMessage(err, path)
		}

		return nil
	}

	err := filepath.Walk(keystorePath, walkFunc)
	if err != nil {
		return nil, err
	}

	return priv, err
}

func parsePrivateKeyPEM(rawKey []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(rawKey)
	if block == nil {
		return nil, errors.New("bytes are not PEM encoded")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.WithMessage(err, "pem bytes are not PKCS8 encoded ")
	}

	priv, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("pem bytes do not contain an EC private key")
	}
	return priv, nil
}



func GeneratePrivateKey(keystorePath string) (*ecdsa.PrivateKey, error) {

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate private key")
	}

	pkcs8Encoded, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal private key")
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8Encoded})

	keyFile := filepath.Join(keystorePath, "priv_sk")
	err = ioutil.WriteFile(keyFile, pemEncoded, 0600)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to save private key to file %s", keyFile)
	}

	return priv, err
}


type ECDSASigner struct {
	PrivateKey *ecdsa.PrivateKey
}


func (e *ECDSASigner) Public() crypto.PublicKey {
	return &e.PrivateKey.PublicKey
}


func (e *ECDSASigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand, e.PrivateKey, digest)
	if err != nil {
		return nil, err
	}

	
	sig := toLowS(
		e.PrivateKey.PublicKey,
		ECDSASignature{
			R: r,
			S: s,
		},
	)

	
	return asn1.Marshal(sig)
}


func toLowS(key ecdsa.PublicKey, sig ECDSASignature) ECDSASignature {
	
	halfOrder := new(big.Int).Div(key.Curve.Params().N, big.NewInt(2))
	
	if sig.S.Cmp(halfOrder) == 1 {
		
		sig.S.Sub(key.Params().N, sig.S)
	}
	return sig
}

type ECDSASignature struct {
	R, S *big.Int
}
