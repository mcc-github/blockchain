/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	m "github.com/mcc-github/blockchain/protos/msp"
	errors "github.com/pkg/errors"
)

func (msp *bccspmsp) getCertifiersIdentifier(certRaw []byte) ([]byte, error) {
	
	cert, err := msp.getCertFromPem(certRaw)
	if err != nil {
		return nil, fmt.Errorf("Failed getting certificate for [%v]: [%s]", certRaw, err)
	}

	
	cert, err = msp.sanitizeCert(cert)
	if err != nil {
		return nil, fmt.Errorf("sanitizeCert failed %s", err)
	}

	found := false
	root := false
	
	for _, v := range msp.rootCerts {
		if v.(*identity).cert.Equal(cert) {
			found = true
			root = true
			break
		}
	}
	if !found {
		
		for _, v := range msp.intermediateCerts {
			if v.(*identity).cert.Equal(cert) {
				found = true
				break
			}
		}
	}
	if !found {
		
		return nil, fmt.Errorf("Failed adding OU. Certificate [%v] not in root or intermediate certs.", cert)
	}

	
	var certifiersIdentifier []byte
	var chain []*x509.Certificate
	if root {
		chain = []*x509.Certificate{cert}
	} else {
		chain, err = msp.getValidationChain(cert, true)
		if err != nil {
			return nil, fmt.Errorf("Failed computing validation chain for [%v]. [%s]", cert, err)
		}
	}

	
	certifiersIdentifier, err = msp.getCertificationChainIdentifierFromChain(chain)
	if err != nil {
		return nil, fmt.Errorf("Failed computing Certifiers Identifier for [%v]. [%s]", certRaw, err)
	}

	return certifiersIdentifier, nil

}

func (msp *bccspmsp) setupCrypto(conf *m.FabricMSPConfig) error {
	msp.cryptoConfig = conf.CryptoConfig
	if msp.cryptoConfig == nil {
		
		msp.cryptoConfig = &m.FabricCryptoConfig{
			SignatureHashFamily:            bccsp.SHA2,
			IdentityIdentifierHashFunction: bccsp.SHA256,
		}
		mspLogger.Debugf("CryptoConfig was nil. Move to defaults.")
	}
	if msp.cryptoConfig.SignatureHashFamily == "" {
		msp.cryptoConfig.SignatureHashFamily = bccsp.SHA2
		mspLogger.Debugf("CryptoConfig.SignatureHashFamily was nil. Move to defaults.")
	}
	if msp.cryptoConfig.IdentityIdentifierHashFunction == "" {
		msp.cryptoConfig.IdentityIdentifierHashFunction = bccsp.SHA256
		mspLogger.Debugf("CryptoConfig.IdentityIdentifierHashFunction was nil. Move to defaults.")
	}

	return nil
}

func (msp *bccspmsp) setupCAs(conf *m.FabricMSPConfig) error {
	
	if len(conf.RootCerts) == 0 {
		return errors.New("expected at least one CA certificate")
	}

	
	
	
	
	
	msp.opts = &x509.VerifyOptions{Roots: x509.NewCertPool(), Intermediates: x509.NewCertPool()}
	for _, v := range conf.RootCerts {
		cert, err := msp.getCertFromPem(v)
		if err != nil {
			return err
		}
		msp.opts.Roots.AddCert(cert)
	}
	for _, v := range conf.IntermediateCerts {
		cert, err := msp.getCertFromPem(v)
		if err != nil {
			return err
		}
		msp.opts.Intermediates.AddCert(cert)
	}

	
	
	msp.rootCerts = make([]Identity, len(conf.RootCerts))
	for i, trustedCert := range conf.RootCerts {
		id, _, err := msp.getIdentityFromConf(trustedCert)
		if err != nil {
			return err
		}

		msp.rootCerts[i] = id
	}

	
	msp.intermediateCerts = make([]Identity, len(conf.IntermediateCerts))
	for i, trustedCert := range conf.IntermediateCerts {
		id, _, err := msp.getIdentityFromConf(trustedCert)
		if err != nil {
			return err
		}

		msp.intermediateCerts[i] = id
	}

	
	msp.opts = &x509.VerifyOptions{Roots: x509.NewCertPool(), Intermediates: x509.NewCertPool()}
	for _, id := range msp.rootCerts {
		msp.opts.Roots.AddCert(id.(*identity).cert)
	}
	for _, id := range msp.intermediateCerts {
		msp.opts.Intermediates.AddCert(id.(*identity).cert)
	}

	return nil
}

func (msp *bccspmsp) setupAdmins(conf *m.FabricMSPConfig) error {
	return msp.internalSetupAdmin(conf)
}

func (msp *bccspmsp) setupAdminsPreV142(conf *m.FabricMSPConfig) error {
	
	msp.admins = make([]Identity, len(conf.Admins))
	for i, admCert := range conf.Admins {
		id, _, err := msp.getIdentityFromConf(admCert)
		if err != nil {
			return err
		}

		msp.admins[i] = id
	}

	return nil
}

func (msp *bccspmsp) setupAdminsV142(conf *m.FabricMSPConfig) error {
	
	if err := msp.setupAdminsPreV142(conf); err != nil {
		return err
	}

	if len(msp.admins) == 0 && (!msp.ouEnforcement || msp.adminOU == nil) {
		return errors.New("administrators must be declared when no admin ou classification is set")
	}

	return nil
}

func (msp *bccspmsp) setupCRLs(conf *m.FabricMSPConfig) error {
	
	msp.CRL = make([]*pkix.CertificateList, len(conf.RevocationList))
	for i, crlbytes := range conf.RevocationList {
		crl, err := x509.ParseCRL(crlbytes)
		if err != nil {
			return errors.Wrap(err, "could not parse RevocationList")
		}

		
		
		
		

		msp.CRL[i] = crl
	}

	return nil
}

func (msp *bccspmsp) finalizeSetupCAs() error {
	
	for _, id := range append(append([]Identity{}, msp.rootCerts...), msp.intermediateCerts...) {
		if !id.(*identity).cert.IsCA {
			return errors.Errorf("CA Certificate did not have the CA attribute, (SN: %x)", id.(*identity).cert.SerialNumber)
		}
		if _, err := getSubjectKeyIdentifierFromCert(id.(*identity).cert); err != nil {
			return errors.WithMessagef(err, "CA Certificate problem with Subject Key Identifier extension, (SN: %x)", id.(*identity).cert.SerialNumber)
		}

		if err := msp.validateCAIdentity(id.(*identity)); err != nil {
			return errors.WithMessagef(err, "CA Certificate is not valid, (SN: %s)", id.(*identity).cert.SerialNumber)
		}
	}

	
	
	msp.certificationTreeInternalNodesMap = make(map[string]bool)
	for _, id := range append([]Identity{}, msp.intermediateCerts...) {
		chain, err := msp.getUniqueValidationChain(id.(*identity).cert, msp.getValidityOptsForCert(id.(*identity).cert))
		if err != nil {
			return errors.WithMessagef(err, "failed getting validation chain, (SN: %s)", id.(*identity).cert.SerialNumber)
		}

		
		for i := 1; i < len(chain); i++ {
			msp.certificationTreeInternalNodesMap[string(chain[i].Raw)] = true
		}
	}

	return nil
}

func (msp *bccspmsp) setupNodeOUs(config *m.FabricMSPConfig) error {
	if config.FabricNodeOus != nil {

		msp.ouEnforcement = config.FabricNodeOus.Enable

		if config.FabricNodeOus.ClientOuIdentifier == nil || len(config.FabricNodeOus.ClientOuIdentifier.OrganizationalUnitIdentifier) == 0 {
			return errors.New("Failed setting up NodeOUs. ClientOU must be different from nil.")
		}

		if config.FabricNodeOus.PeerOuIdentifier == nil || len(config.FabricNodeOus.PeerOuIdentifier.OrganizationalUnitIdentifier) == 0 {
			return errors.New("Failed setting up NodeOUs. PeerOU must be different from nil.")
		}

		
		msp.clientOU = &OUIdentifier{OrganizationalUnitIdentifier: config.FabricNodeOus.ClientOuIdentifier.OrganizationalUnitIdentifier}
		if len(config.FabricNodeOus.ClientOuIdentifier.Certificate) != 0 {
			certifiersIdentifier, err := msp.getCertifiersIdentifier(config.FabricNodeOus.ClientOuIdentifier.Certificate)
			if err != nil {
				return err
			}
			msp.clientOU.CertifiersIdentifier = certifiersIdentifier
		}

		
		msp.peerOU = &OUIdentifier{OrganizationalUnitIdentifier: config.FabricNodeOus.PeerOuIdentifier.OrganizationalUnitIdentifier}
		if len(config.FabricNodeOus.PeerOuIdentifier.Certificate) != 0 {
			certifiersIdentifier, err := msp.getCertifiersIdentifier(config.FabricNodeOus.PeerOuIdentifier.Certificate)
			if err != nil {
				return err
			}
			msp.peerOU.CertifiersIdentifier = certifiersIdentifier
		}

	} else {
		msp.ouEnforcement = false
	}

	return nil
}

func (msp *bccspmsp) setupNodeOUsV142(config *m.FabricMSPConfig) error {
	if config.FabricNodeOus == nil {
		msp.ouEnforcement = false
		return nil
	}

	msp.ouEnforcement = config.FabricNodeOus.Enable

	counter := 0
	
	if config.FabricNodeOus.ClientOuIdentifier != nil {
		msp.clientOU = &OUIdentifier{OrganizationalUnitIdentifier: config.FabricNodeOus.ClientOuIdentifier.OrganizationalUnitIdentifier}
		if len(config.FabricNodeOus.ClientOuIdentifier.Certificate) != 0 {
			certifiersIdentifier, err := msp.getCertifiersIdentifier(config.FabricNodeOus.ClientOuIdentifier.Certificate)
			if err != nil {
				return err
			}
			msp.clientOU.CertifiersIdentifier = certifiersIdentifier
		}
		counter++
	} else {
		msp.clientOU = nil
	}

	
	if config.FabricNodeOus.PeerOuIdentifier != nil {
		msp.peerOU = &OUIdentifier{OrganizationalUnitIdentifier: config.FabricNodeOus.PeerOuIdentifier.OrganizationalUnitIdentifier}
		if len(config.FabricNodeOus.PeerOuIdentifier.Certificate) != 0 {
			certifiersIdentifier, err := msp.getCertifiersIdentifier(config.FabricNodeOus.PeerOuIdentifier.Certificate)
			if err != nil {
				return err
			}
			msp.peerOU.CertifiersIdentifier = certifiersIdentifier
		}
		counter++
	} else {
		msp.peerOU = nil
	}

	
	if config.FabricNodeOus.AdminOuIdentifier != nil {
		msp.adminOU = &OUIdentifier{OrganizationalUnitIdentifier: config.FabricNodeOus.AdminOuIdentifier.OrganizationalUnitIdentifier}
		if len(config.FabricNodeOus.AdminOuIdentifier.Certificate) != 0 {
			certifiersIdentifier, err := msp.getCertifiersIdentifier(config.FabricNodeOus.AdminOuIdentifier.Certificate)
			if err != nil {
				return err
			}
			msp.adminOU.CertifiersIdentifier = certifiersIdentifier
		}
		counter++
	} else {
		msp.adminOU = nil
	}

	
	if config.FabricNodeOus.OrdererOuIdentifier != nil {
		msp.ordererOU = &OUIdentifier{OrganizationalUnitIdentifier: config.FabricNodeOus.OrdererOuIdentifier.OrganizationalUnitIdentifier}
		if len(config.FabricNodeOus.OrdererOuIdentifier.Certificate) != 0 {
			certifiersIdentifier, err := msp.getCertifiersIdentifier(config.FabricNodeOus.OrdererOuIdentifier.Certificate)
			if err != nil {
				return err
			}
			msp.ordererOU.CertifiersIdentifier = certifiersIdentifier
		}
		counter++
	} else {
		msp.ordererOU = nil
	}

	if counter == 0 {
		
		msp.ouEnforcement = false
	}

	return nil
}

func (msp *bccspmsp) setupSigningIdentity(conf *m.FabricMSPConfig) error {
	if conf.SigningIdentity != nil {
		sid, err := msp.getSigningIdentityFromConf(conf.SigningIdentity)
		if err != nil {
			return err
		}

		expirationTime := sid.ExpiresAt()
		now := time.Now()
		if expirationTime.After(now) {
			mspLogger.Debug("Signing identity expires at", expirationTime)
		} else if expirationTime.IsZero() {
			mspLogger.Debug("Signing identity has no known expiration time")
		} else {
			return errors.Errorf("signing identity expired %v ago", now.Sub(expirationTime))
		}

		msp.signer = sid
	}

	return nil
}

func (msp *bccspmsp) setupOUs(conf *m.FabricMSPConfig) error {
	msp.ouIdentifiers = make(map[string][][]byte)
	for _, ou := range conf.OrganizationalUnitIdentifiers {

		certifiersIdentifier, err := msp.getCertifiersIdentifier(ou.Certificate)
		if err != nil {
			return errors.WithMessagef(err, "failed getting certificate for [%v]", ou)
		}

		
		found := false
		for _, id := range msp.ouIdentifiers[ou.OrganizationalUnitIdentifier] {
			if bytes.Equal(id, certifiersIdentifier) {
				mspLogger.Warningf("Duplicate found in ou identifiers [%s, %v]", ou.OrganizationalUnitIdentifier, id)
				found = true
				break
			}
		}

		if !found {
			
			msp.ouIdentifiers[ou.OrganizationalUnitIdentifier] = append(
				msp.ouIdentifiers[ou.OrganizationalUnitIdentifier],
				certifiersIdentifier,
			)
		}
	}

	return nil
}

func (msp *bccspmsp) setupTLSCAs(conf *m.FabricMSPConfig) error {

	opts := &x509.VerifyOptions{Roots: x509.NewCertPool(), Intermediates: x509.NewCertPool()}

	
	msp.tlsRootCerts = make([][]byte, len(conf.TlsRootCerts))
	rootCerts := make([]*x509.Certificate, len(conf.TlsRootCerts))
	for i, trustedCert := range conf.TlsRootCerts {
		cert, err := msp.getCertFromPem(trustedCert)
		if err != nil {
			return err
		}

		rootCerts[i] = cert
		msp.tlsRootCerts[i] = trustedCert
		opts.Roots.AddCert(cert)
	}

	
	msp.tlsIntermediateCerts = make([][]byte, len(conf.TlsIntermediateCerts))
	intermediateCerts := make([]*x509.Certificate, len(conf.TlsIntermediateCerts))
	for i, trustedCert := range conf.TlsIntermediateCerts {
		cert, err := msp.getCertFromPem(trustedCert)
		if err != nil {
			return err
		}

		intermediateCerts[i] = cert
		msp.tlsIntermediateCerts[i] = trustedCert
		opts.Intermediates.AddCert(cert)
	}

	
	for _, cert := range append(append([]*x509.Certificate{}, rootCerts...), intermediateCerts...) {
		if cert == nil {
			continue
		}

		if !cert.IsCA {
			return errors.Errorf("CA Certificate did not have the CA attribute, (SN: %x)", cert.SerialNumber)
		}
		if _, err := getSubjectKeyIdentifierFromCert(cert); err != nil {
			return errors.WithMessagef(err, "CA Certificate problem with Subject Key Identifier extension, (SN: %x)", cert.SerialNumber)
		}

		if err := msp.validateTLSCAIdentity(cert, opts); err != nil {
			return errors.WithMessagef(err, "CA Certificate is not valid, (SN: %s)", cert.SerialNumber)
		}
	}

	return nil
}

func (msp *bccspmsp) setupV1(conf1 *m.FabricMSPConfig) error {
	err := msp.preSetupV1(conf1)
	if err != nil {
		return err
	}

	err = msp.postSetupV1(conf1)
	if err != nil {
		return err
	}

	return nil
}

func (msp *bccspmsp) preSetupV1(conf *m.FabricMSPConfig) error {
	
	if err := msp.setupCrypto(conf); err != nil {
		return err
	}

	
	if err := msp.setupCAs(conf); err != nil {
		return err
	}

	
	if err := msp.setupAdmins(conf); err != nil {
		return err
	}

	
	if err := msp.setupCRLs(conf); err != nil {
		return err
	}

	
	if err := msp.finalizeSetupCAs(); err != nil {
		return err
	}

	
	if err := msp.setupSigningIdentity(conf); err != nil {
		return err
	}

	
	if err := msp.setupTLSCAs(conf); err != nil {
		return err
	}

	
	if err := msp.setupOUs(conf); err != nil {
		return err
	}

	return nil
}

func (msp *bccspmsp) preSetupV142(conf *m.FabricMSPConfig) error {
	
	if err := msp.setupCrypto(conf); err != nil {
		return err
	}

	
	if err := msp.setupCAs(conf); err != nil {
		return err
	}

	
	if err := msp.setupCRLs(conf); err != nil {
		return err
	}

	
	if err := msp.finalizeSetupCAs(); err != nil {
		return err
	}

	
	if err := msp.setupSigningIdentity(conf); err != nil {
		return err
	}

	
	if err := msp.setupTLSCAs(conf); err != nil {
		return err
	}

	
	if err := msp.setupOUs(conf); err != nil {
		return err
	}

	
	if err := msp.setupNodeOUsV142(conf); err != nil {
		return err
	}

	
	if err := msp.setupAdmins(conf); err != nil {
		return err
	}

	return nil
}

func (msp *bccspmsp) postSetupV1(conf *m.FabricMSPConfig) error {
	
	
	
	for i, admin := range msp.admins {
		err := admin.Validate()
		if err != nil {
			return errors.WithMessagef(err, "admin %d is invalid", i)
		}
	}

	return nil
}

func (msp *bccspmsp) setupV11(conf *m.FabricMSPConfig) error {
	err := msp.preSetupV1(conf)
	if err != nil {
		return err
	}

	
	if err := msp.setupNodeOUs(conf); err != nil {
		return err
	}

	err = msp.postSetupV11(conf)
	if err != nil {
		return err
	}

	return nil
}

func (msp *bccspmsp) setupV142(conf *m.FabricMSPConfig) error {
	err := msp.preSetupV142(conf)
	if err != nil {
		return err
	}

	err = msp.postSetupV142(conf)
	if err != nil {
		return err
	}

	return nil
}

func (msp *bccspmsp) postSetupV11(conf *m.FabricMSPConfig) error {
	
	if !msp.ouEnforcement {
		
		return msp.postSetupV1(conf)
	}

	
	principalBytes, err := proto.Marshal(&m.MSPRole{Role: m.MSPRole_CLIENT, MspIdentifier: msp.name})
	if err != nil {
		return errors.Wrapf(err, "failed creating MSPRole_CLIENT")
	}
	principal := &m.MSPPrincipal{
		PrincipalClassification: m.MSPPrincipal_ROLE,
		Principal:               principalBytes}
	for i, admin := range msp.admins {
		err = admin.SatisfiesPrincipal(principal)
		if err != nil {
			return errors.WithMessagef(err, "admin %d is invalid", i)
		}
	}

	return nil
}

func (msp *bccspmsp) postSetupV142(conf *m.FabricMSPConfig) error {
	
	if !msp.ouEnforcement {
		
		return msp.postSetupV1(conf)
	}

	
	for i, admin := range msp.admins {
		err1 := msp.hasOURole(admin, m.MSPRole_CLIENT)
		err2 := msp.hasOURole(admin, m.MSPRole_ADMIN)
		if err1 != nil && err2 != nil {
			return errors.Errorf("admin %d is invalid [%s,%s]", i, err1, err2)
		}
	}

	return nil
}
