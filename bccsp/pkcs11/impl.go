/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pkcs11

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"os"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/miekg/pkcs11"
	"github.com/pkg/errors"
)

var (
	logger           = flogging.MustGetLogger("bccsp_p11")
	sessionCacheSize = 10
)



func New(opts PKCS11Opts, keyStore bccsp.KeyStore) (bccsp.BCCSP, error) {
	
	conf := &config{}
	err := conf.setSecurityLevel(opts.SecLevel, opts.HashFamily)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed initializing configuration")
	}

	swCSP, err := sw.NewWithParams(opts.SecLevel, opts.HashFamily, keyStore)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed initializing fallback SW BCCSP")
	}

	
	if keyStore == nil {
		return nil, errors.New("Invalid bccsp.KeyStore instance. It must be different from nil")
	}

	lib := opts.Library
	pin := opts.Pin
	label := opts.Label
	ctx, slot, session, err := loadLib(lib, pin, label)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed initializing PKCS11 library %s %s",
			lib, label)
	}

	sessions := make(chan pkcs11.SessionHandle, sessionCacheSize)
	csp := &impl{swCSP, conf, keyStore, ctx, sessions, slot, lib, opts.SoftVerify, opts.Immutable}
	csp.returnSession(*session)
	return csp, nil
}

type impl struct {
	bccsp.BCCSP

	conf *config
	ks   bccsp.KeyStore

	ctx      *pkcs11.Ctx
	sessions chan pkcs11.SessionHandle
	slot     uint

	lib        string
	softVerify bool
	
	immutable bool
}


func (csp *impl) KeyGen(opts bccsp.KeyGenOpts) (k bccsp.Key, err error) {
	
	if opts == nil {
		return nil, errors.New("Invalid Opts parameter. It must not be nil")
	}

	
	switch opts.(type) {
	case *bccsp.ECDSAKeyGenOpts:
		ski, pub, err := csp.generateECKey(csp.conf.ellipticCurve, opts.Ephemeral())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed generating ECDSA key")
		}
		k = &ecdsaPrivateKey{ski, ecdsaPublicKey{ski, pub}}

	case *bccsp.ECDSAP256KeyGenOpts:
		ski, pub, err := csp.generateECKey(oidNamedCurveP256, opts.Ephemeral())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed generating ECDSA P256 key")
		}

		k = &ecdsaPrivateKey{ski, ecdsaPublicKey{ski, pub}}

	case *bccsp.ECDSAP384KeyGenOpts:
		ski, pub, err := csp.generateECKey(oidNamedCurveP384, opts.Ephemeral())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed generating ECDSA P384 key")
		}

		k = &ecdsaPrivateKey{ski, ecdsaPublicKey{ski, pub}}

	default:
		return csp.BCCSP.KeyGen(opts)
	}

	return k, nil
}



func (csp *impl) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	
	if raw == nil {
		return nil, errors.New("Invalid raw. Cannot be nil")
	}

	if opts == nil {
		return nil, errors.New("Invalid Opts parameter. It must not be nil")
	}

	switch opts.(type) {

	case *bccsp.X509PublicKeyImportOpts:
		x509Cert, ok := raw.(*x509.Certificate)
		if !ok {
			return nil, errors.New("[X509PublicKeyImportOpts] Invalid raw material. Expected *x509.Certificate")
		}

		pk := x509Cert.PublicKey

		switch pk.(type) {
		case *ecdsa.PublicKey:
			return csp.KeyImport(pk, &bccsp.ECDSAGoPublicKeyImportOpts{Temporary: opts.Ephemeral()})
		case *rsa.PublicKey:
			return csp.KeyImport(pk, &bccsp.RSAGoPublicKeyImportOpts{Temporary: opts.Ephemeral()})
		default:
			return nil, errors.New("Certificate's public key type not recognized. Supported keys: [ECDSA, RSA]")
		}

	default:
		return csp.BCCSP.KeyImport(raw, opts)

	}
}



func (csp *impl) GetKey(ski []byte) (bccsp.Key, error) {
	pubKey, isPriv, err := csp.getECKey(ski)
	if err == nil {
		if isPriv {
			return &ecdsaPrivateKey{ski, ecdsaPublicKey{ski, pubKey}}, nil
		}
		return &ecdsaPublicKey{ski, pubKey}, nil
	}
	return csp.BCCSP.GetKey(ski)
}







func (csp *impl) Sign(k bccsp.Key, digest []byte, opts bccsp.SignerOpts) ([]byte, error) {
	
	if k == nil {
		return nil, errors.New("Invalid Key. It must not be nil")
	}
	if len(digest) == 0 {
		return nil, errors.New("Invalid digest. Cannot be empty")
	}

	
	switch key := k.(type) {
	case *ecdsaPrivateKey:
		return csp.signECDSA(*key, digest, opts)
	default:
		return csp.BCCSP.Sign(key, digest, opts)
	}
}


func (csp *impl) Verify(k bccsp.Key, signature, digest []byte, opts bccsp.SignerOpts) (bool, error) {
	
	if k == nil {
		return false, errors.New("Invalid Key. It must not be nil")
	}
	if len(signature) == 0 {
		return false, errors.New("Invalid signature. Cannot be empty")
	}
	if len(digest) == 0 {
		return false, errors.New("Invalid digest. Cannot be empty")
	}

	
	switch key := k.(type) {
	case *ecdsaPrivateKey:
		return csp.verifyECDSA(key.pub, signature, digest, opts)
	case *ecdsaPublicKey:
		return csp.verifyECDSA(*key, signature, digest, opts)
	default:
		return csp.BCCSP.Verify(k, signature, digest, opts)
	}
}



func (csp *impl) Encrypt(k bccsp.Key, plaintext []byte, opts bccsp.EncrypterOpts) ([]byte, error) {
	
	return csp.BCCSP.Encrypt(k, plaintext, opts)
}



func (csp *impl) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) ([]byte, error) {
	return csp.BCCSP.Decrypt(k, ciphertext, opts)
}




func FindPKCS11Lib() (lib, pin, label string) {
	
	lib = os.Getenv("PKCS11_LIB")
	if lib == "" {
		pin = "98765432"
		label = "ForFabric"
		possibilities := []string{
			"/usr/lib/softhsm/libsofthsm2.so",                            
			"/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so",           
			"/usr/lib/s390x-linux-gnu/softhsm/libsofthsm2.so",            
			"/usr/lib/powerpc64le-linux-gnu/softhsm/libsofthsm2.so",      
			"/usr/local/Cellar/softhsm/2.5.0/lib/softhsm/libsofthsm2.so", 
		}
		for _, path := range possibilities {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				lib = path
				break
			}
		}
	} else {
		pin = os.Getenv("PKCS11_PIN")
		label = os.Getenv("PKCS11_LABEL")
	}
	return lib, pin, label
}
