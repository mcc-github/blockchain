/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/bccsp/signer"
	"github.com/mcc-github/blockchain/bccsp/sw"
	m "github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
)


type mspSetupFuncType func(config *m.FabricMSPConfig) error


type validateIdentityOUsFuncType func(id *identity) error


type satisfiesPrincipalInternalFuncType func(id Identity, principal *m.MSPPrincipal) error


type setupAdminInternalFuncType func(conf *m.FabricMSPConfig) error



type bccspmsp struct {
	
	version MSPVersion
	
	
	
	internalSetupFunc mspSetupFuncType

	
	internalValidateIdentityOusFunc validateIdentityOUsFuncType

	
	internalSatisfiesPrincipalInternalFunc satisfiesPrincipalInternalFuncType

	
	internalSetupAdmin setupAdminInternalFuncType

	
	rootCerts []Identity

	
	intermediateCerts []Identity

	
	tlsRootCerts [][]byte

	
	tlsIntermediateCerts [][]byte

	
	
	
	
	certificationTreeInternalNodesMap map[string]bool

	
	signer SigningIdentity

	
	admins []Identity

	
	bccsp bccsp.BCCSP

	
	name string

	
	opts *x509.VerifyOptions

	
	CRL []*pkix.CertificateList

	
	ouIdentifiers map[string][][]byte

	
	cryptoConfig *m.FabricCryptoConfig

	
	ouEnforcement bool
	
	
	clientOU, peerOU, adminOU, ordererOU *OUIdentifier
}





func newBccspMsp(version MSPVersion, defaultBCCSP bccsp.BCCSP) (MSP, error) {
	mspLogger.Debugf("Creating BCCSP-based MSP instance")

	theMsp := &bccspmsp{}
	theMsp.version = version
	theMsp.bccsp = defaultBCCSP
	switch version {
	case MSPv1_0:
		theMsp.internalSetupFunc = theMsp.setupV1
		theMsp.internalValidateIdentityOusFunc = theMsp.validateIdentityOUsV1
		theMsp.internalSatisfiesPrincipalInternalFunc = theMsp.satisfiesPrincipalInternalPreV13
		theMsp.internalSetupAdmin = theMsp.setupAdminsPreV142
	case MSPv1_1:
		theMsp.internalSetupFunc = theMsp.setupV11
		theMsp.internalValidateIdentityOusFunc = theMsp.validateIdentityOUsV11
		theMsp.internalSatisfiesPrincipalInternalFunc = theMsp.satisfiesPrincipalInternalPreV13
		theMsp.internalSetupAdmin = theMsp.setupAdminsPreV142
	case MSPv1_3:
		theMsp.internalSetupFunc = theMsp.setupV11
		theMsp.internalValidateIdentityOusFunc = theMsp.validateIdentityOUsV11
		theMsp.internalSatisfiesPrincipalInternalFunc = theMsp.satisfiesPrincipalInternalV13
		theMsp.internalSetupAdmin = theMsp.setupAdminsPreV142
	case MSPv1_4_2:
		theMsp.internalSetupFunc = theMsp.setupV142
		theMsp.internalValidateIdentityOusFunc = theMsp.validateIdentityOUsV142
		theMsp.internalSatisfiesPrincipalInternalFunc = theMsp.satisfiesPrincipalInternalV142
		theMsp.internalSetupAdmin = theMsp.setupAdminsV142
	default:
		return nil, errors.Errorf("Invalid MSP version [%v]", version)
	}

	return theMsp, nil
}



func NewBccspMspWithKeyStore(version MSPVersion, keyStore bccsp.KeyStore, bccsp bccsp.BCCSP) (MSP, error) {
	thisMSP, err := newBccspMsp(version, bccsp)
	if err != nil {
		return nil, err
	}

	csp, err := sw.NewWithParams(
		factory.GetDefaultOpts().SwOpts.SecLevel,
		factory.GetDefaultOpts().SwOpts.HashFamily,
		keyStore)
	if err != nil {
		return nil, err
	}
	thisMSP.(*bccspmsp).bccsp = csp

	return thisMSP, nil
}

func (msp *bccspmsp) getCertFromPem(idBytes []byte) (*x509.Certificate, error) {
	if idBytes == nil {
		return nil, errors.New("getCertFromPem error: nil idBytes")
	}

	
	pemCert, _ := pem.Decode(idBytes)
	if pemCert == nil {
		return nil, errors.Errorf("getCertFromPem error: could not decode pem bytes [%v]", idBytes)
	}

	
	var cert *x509.Certificate
	cert, err := x509.ParseCertificate(pemCert.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "getCertFromPem error: failed to parse x509 cert")
	}

	return cert, nil
}

func (msp *bccspmsp) getIdentityFromConf(idBytes []byte) (Identity, bccsp.Key, error) {
	
	cert, err := msp.getCertFromPem(idBytes)
	if err != nil {
		return nil, nil, err
	}

	
	certPubK, err := msp.bccsp.KeyImport(cert, &bccsp.X509PublicKeyImportOpts{Temporary: true})
	if err != nil {
		return nil, nil, err
	}

	mspId, err := newIdentity(cert, certPubK, msp)
	if err != nil {
		return nil, nil, err
	}

	return mspId, certPubK, nil
}

func (msp *bccspmsp) getSigningIdentityFromConf(sidInfo *m.SigningIdentityInfo) (SigningIdentity, error) {
	if sidInfo == nil {
		return nil, errors.New("getIdentityFromBytes error: nil sidInfo")
	}

	
	idPub, pubKey, err := msp.getIdentityFromConf(sidInfo.PublicSigner)
	if err != nil {
		return nil, err
	}

	
	privKey, err := msp.bccsp.GetKey(pubKey.SKI())
	
	if err != nil {
		mspLogger.Debugf("Could not find SKI [%s], trying KeyMaterial field: %+v\n", hex.EncodeToString(pubKey.SKI()), err)
		if sidInfo.PrivateSigner == nil || sidInfo.PrivateSigner.KeyMaterial == nil {
			return nil, errors.New("KeyMaterial not found in SigningIdentityInfo")
		}

		pemKey, _ := pem.Decode(sidInfo.PrivateSigner.KeyMaterial)
		if pemKey == nil {
			return nil, errors.Errorf("%s: wrong PEM encoding", sidInfo.PrivateSigner.KeyIdentifier)
		}
		privKey, err = msp.bccsp.KeyImport(pemKey.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: true})
		if err != nil {
			return nil, errors.WithMessage(err, "getIdentityFromBytes error: Failed to import EC private key")
		}
	}

	
	peerSigner, err := signer.New(msp.bccsp, privKey)
	if err != nil {
		return nil, errors.WithMessage(err, "getIdentityFromBytes error: Failed initializing bccspCryptoSigner")
	}

	return newSigningIdentity(idPub.(*identity).cert, idPub.(*identity).pk, peerSigner, msp)
}




func (msp *bccspmsp) Setup(conf1 *m.MSPConfig) error {
	if conf1 == nil {
		return errors.New("Setup error: nil conf reference")
	}

	
	conf := &m.FabricMSPConfig{}
	err := proto.Unmarshal(conf1.Config, conf)
	if err != nil {
		return errors.Wrap(err, "failed unmarshalling blockchain msp config")
	}

	
	msp.name = conf.Name
	mspLogger.Debugf("Setting up MSP instance %s", msp.name)

	
	return msp.internalSetupFunc(conf)
}


func (msp *bccspmsp) GetVersion() MSPVersion {
	return msp.version
}


func (msp *bccspmsp) GetType() ProviderType {
	return FABRIC
}


func (msp *bccspmsp) GetIdentifier() (string, error) {
	return msp.name, nil
}


func (msp *bccspmsp) GetTLSRootCerts() [][]byte {
	return msp.tlsRootCerts
}


func (msp *bccspmsp) GetTLSIntermediateCerts() [][]byte {
	return msp.tlsIntermediateCerts
}



func (msp *bccspmsp) GetDefaultSigningIdentity() (SigningIdentity, error) {
	mspLogger.Debugf("Obtaining default signing identity")

	if msp.signer == nil {
		return nil, errors.New("this MSP does not possess a valid default signing identity")
	}

	return msp.signer, nil
}



func (msp *bccspmsp) GetSigningIdentity(identifier *IdentityIdentifier) (SigningIdentity, error) {
	
	return nil, errors.Errorf("no signing identity for %#v", identifier)
}






func (msp *bccspmsp) Validate(id Identity) error {
	mspLogger.Debugf("MSP %s validating identity", msp.name)

	switch id := id.(type) {
	
	
	
	case *identity:
		return msp.validateIdentity(id)
	default:
		return errors.New("identity type not recognized")
	}
}





func (msp *bccspmsp) hasOURole(id Identity, mspRole m.MSPRole_MSPRoleType) error {
	
	if !msp.ouEnforcement {
		return errors.New("NodeOUs not activated. Cannot tell apart identities.")
	}

	mspLogger.Debugf("MSP %s checking if the identity is a client", msp.name)

	switch id := id.(type) {
	
	
	
	case *identity:
		return msp.hasOURoleInternal(id, mspRole)
	default:
		return errors.New("Identity type not recognized")
	}
}

func (msp *bccspmsp) hasOURoleInternal(id *identity, mspRole m.MSPRole_MSPRoleType) error {
	var nodeOU *OUIdentifier
	switch mspRole {
	case m.MSPRole_CLIENT:
		nodeOU = msp.clientOU
	case m.MSPRole_PEER:
		nodeOU = msp.peerOU
	case m.MSPRole_ADMIN:
		nodeOU = msp.adminOU
	case m.MSPRole_ORDERER:
		nodeOU = msp.ordererOU
	default:
		return errors.New("Invalid MSPRoleType. It must be CLIENT, PEER, ADMIN or ORDERER")
	}

	if nodeOU == nil {
		return errors.Errorf("cannot test for classification, node ou for type [%s], not defined, msp: [%s]", mspRole, msp.name)
	}

	for _, OU := range id.GetOrganizationalUnits() {
		if OU.OrganizationalUnitIdentifier == nodeOU.OrganizationalUnitIdentifier {
			return nil
		}
	}

	return errors.Errorf("The identity does not contain OU [%s], MSP: [%s]", mspRole, msp.name)
}



func (msp *bccspmsp) DeserializeIdentity(serializedID []byte) (Identity, error) {
	mspLogger.Debug("Obtaining identity")

	
	sId := &m.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sId)
	if err != nil {
		return nil, errors.Wrap(err, "could not deserialize a SerializedIdentity")
	}

	if sId.Mspid != msp.name {
		return nil, errors.Errorf("expected MSP ID %s, received %s", msp.name, sId.Mspid)
	}

	return msp.deserializeIdentityInternal(sId.IdBytes)
}


func (msp *bccspmsp) deserializeIdentityInternal(serializedIdentity []byte) (Identity, error) {
	
	bl, _ := pem.Decode(serializedIdentity)
	if bl == nil {
		return nil, errors.New("could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parseCertificate failed")
	}

	
	
	
	
	
	

	pub, err := msp.bccsp.KeyImport(cert, &bccsp.X509PublicKeyImportOpts{Temporary: true})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to import certificate's public key")
	}

	return newIdentity(cert, pub, msp)
}


func (msp *bccspmsp) SatisfiesPrincipal(id Identity, principal *m.MSPPrincipal) error {
	principals, err := collectPrincipals(principal, msp.GetVersion())
	if err != nil {
		return err
	}
	for _, principal := range principals {
		err = msp.internalSatisfiesPrincipalInternalFunc(id, principal)
		if err != nil {
			return err
		}
	}
	return nil
}


func collectPrincipals(principal *m.MSPPrincipal, mspVersion MSPVersion) ([]*m.MSPPrincipal, error) {
	switch principal.PrincipalClassification {
	case m.MSPPrincipal_COMBINED:
		
		if mspVersion <= MSPv1_1 {
			return nil, errors.Errorf("invalid principal type %d", int32(principal.PrincipalClassification))
		}
		
		principals := &m.CombinedPrincipal{}
		err := proto.Unmarshal(principal.Principal, principals)
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal CombinedPrincipal from principal")
		}
		
		if len(principals.Principals) == 0 {
			return nil, errors.New("No principals in CombinedPrincipal")
		}
		
		
		var principalsSlice []*m.MSPPrincipal
		for _, cp := range principals.Principals {
			internalSlice, err := collectPrincipals(cp, mspVersion)
			if err != nil {
				return nil, err
			}
			principalsSlice = append(principalsSlice, internalSlice...)
		}
		
		return principalsSlice, nil
	default:
		return []*m.MSPPrincipal{principal}, nil
	}
}




func (msp *bccspmsp) satisfiesPrincipalInternalPreV13(id Identity, principal *m.MSPPrincipal) error {
	switch principal.PrincipalClassification {
	
	
	case m.MSPPrincipal_ROLE:
		
		mspRole := &m.MSPRole{}
		err := proto.Unmarshal(principal.Principal, mspRole)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal MSPRole from principal")
		}

		
		
		if mspRole.MspIdentifier != msp.name {
			return errors.Errorf("the identity is a member of a different MSP (expected %s, got %s)", mspRole.MspIdentifier, id.GetMSPIdentifier())
		}

		
		switch mspRole.Role {
		case m.MSPRole_MEMBER:
			
			
			mspLogger.Debugf("Checking if identity satisfies MEMBER role for %s", msp.name)
			return msp.Validate(id)
		case m.MSPRole_ADMIN:
			mspLogger.Debugf("Checking if identity satisfies ADMIN role for %s", msp.name)
			
			
			if msp.isInAdmins(id.(*identity)) {
				return nil
			}
			return errors.New("This identity is not an admin")
		case m.MSPRole_CLIENT:
			fallthrough
		case m.MSPRole_PEER:
			mspLogger.Debugf("Checking if identity satisfies role [%s] for %s", m.MSPRole_MSPRoleType_name[int32(mspRole.Role)], msp.name)
			if err := msp.Validate(id); err != nil {
				return errors.Wrapf(err, "The identity is not valid under this MSP [%s]", msp.name)
			}

			if err := msp.hasOURole(id, mspRole.Role); err != nil {
				return errors.Wrapf(err, "The identity is not a [%s] under this MSP [%s]", m.MSPRole_MSPRoleType_name[int32(mspRole.Role)], msp.name)
			}
			return nil
		default:
			return errors.Errorf("invalid MSP role type %d", int32(mspRole.Role))
		}
	case m.MSPPrincipal_IDENTITY:
		
		
		principalId, err := msp.DeserializeIdentity(principal.Principal)
		if err != nil {
			return errors.WithMessage(err, "invalid identity principal, not a certificate")
		}

		if bytes.Equal(id.(*identity).cert.Raw, principalId.(*identity).cert.Raw) {
			return principalId.Validate()
		}

		return errors.New("The identities do not match")
	case m.MSPPrincipal_ORGANIZATION_UNIT:
		
		OU := &m.OrganizationUnit{}
		err := proto.Unmarshal(principal.Principal, OU)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal OrganizationUnit from principal")
		}

		
		
		if OU.MspIdentifier != msp.name {
			return errors.Errorf("the identity is a member of a different MSP (expected %s, got %s)", OU.MspIdentifier, id.GetMSPIdentifier())
		}

		
		
		err = msp.Validate(id)
		if err != nil {
			return err
		}

		
		for _, ou := range id.GetOrganizationalUnits() {
			if ou.OrganizationalUnitIdentifier == OU.OrganizationalUnitIdentifier &&
				bytes.Equal(ou.CertifiersIdentifier, OU.CertifiersIdentifier) {
				return nil
			}
		}

		
		return errors.New("The identities do not match")
	default:
		return errors.Errorf("invalid principal type %d", int32(principal.PrincipalClassification))
	}
}





func (msp *bccspmsp) satisfiesPrincipalInternalV13(id Identity, principal *m.MSPPrincipal) error {
	switch principal.PrincipalClassification {
	case m.MSPPrincipal_COMBINED:
		return errors.New("SatisfiesPrincipalInternal shall not be called with a CombinedPrincipal")
	case m.MSPPrincipal_ANONYMITY:
		anon := &m.MSPIdentityAnonymity{}
		err := proto.Unmarshal(principal.Principal, anon)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal MSPIdentityAnonymity from principal")
		}
		switch anon.AnonymityType {
		case m.MSPIdentityAnonymity_ANONYMOUS:
			return errors.New("Principal is anonymous, but X.509 MSP does not support anonymous identities")
		case m.MSPIdentityAnonymity_NOMINAL:
			return nil
		default:
			return errors.Errorf("Unknown principal anonymity type: %d", anon.AnonymityType)
		}

	default:
		
		return msp.satisfiesPrincipalInternalPreV13(id, principal)
	}
}





func (msp *bccspmsp) satisfiesPrincipalInternalV142(id Identity, principal *m.MSPPrincipal) error {
	_, okay := id.(*identity)
	if !okay {
		return errors.New("invalid identity type, expected *identity")
	}

	switch principal.PrincipalClassification {
	case m.MSPPrincipal_ROLE:
		if !msp.ouEnforcement {
			break
		}

		
		mspRole := &m.MSPRole{}
		err := proto.Unmarshal(principal.Principal, mspRole)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal MSPRole from principal")
		}

		
		
		if mspRole.MspIdentifier != msp.name {
			return errors.Errorf("the identity is a member of a different MSP (expected %s, got %s)", mspRole.MspIdentifier, id.GetMSPIdentifier())
		}

		
		switch mspRole.Role {
		case m.MSPRole_ADMIN:
			mspLogger.Debugf("Checking if identity has been named explicitly as an admin for %s", msp.name)
			
			
			if msp.isInAdmins(id.(*identity)) {
				return nil
			}

			
			mspLogger.Debugf("Checking if identity carries the admin ou for %s", msp.name)
			if err := msp.Validate(id); err != nil {
				return errors.Wrapf(err, "The identity is not valid under this MSP [%s]", msp.name)
			}

			if err := msp.hasOURole(id, m.MSPRole_ADMIN); err != nil {
				return errors.Wrapf(err, "The identity is not an admin under this MSP [%s]", msp.name)
			}

			return nil
		case m.MSPRole_ORDERER:
			mspLogger.Debugf("Checking if identity satisfies role [%s] for %s", m.MSPRole_MSPRoleType_name[int32(mspRole.Role)], msp.name)
			if err := msp.Validate(id); err != nil {
				return errors.Wrapf(err, "The identity is not valid under this MSP [%s]", msp.name)
			}

			if err := msp.hasOURole(id, mspRole.Role); err != nil {
				return errors.Wrapf(err, "The identity is not a [%s] under this MSP [%s]", m.MSPRole_MSPRoleType_name[int32(mspRole.Role)], msp.name)
			}
			return nil
		}
	}

	
	return msp.satisfiesPrincipalInternalV13(id, principal)
}

func (msp *bccspmsp) isInAdmins(id *identity) bool {
	for _, admincert := range msp.admins {
		if bytes.Equal(id.cert.Raw, admincert.(*identity).cert.Raw) {
			
			
			
			return true
		}
	}
	return false
}


func (msp *bccspmsp) getCertificationChain(id Identity) ([]*x509.Certificate, error) {
	mspLogger.Debugf("MSP %s getting certification chain", msp.name)

	switch id := id.(type) {
	
	
	
	case *identity:
		return msp.getCertificationChainForBCCSPIdentity(id)
	default:
		return nil, errors.New("identity type not recognized")
	}
}


func (msp *bccspmsp) getCertificationChainForBCCSPIdentity(id *identity) ([]*x509.Certificate, error) {
	if id == nil {
		return nil, errors.New("Invalid bccsp identity. Must be different from nil.")
	}

	
	if msp.opts == nil {
		return nil, errors.New("Invalid msp instance")
	}

	
	if id.cert.IsCA {
		return nil, errors.New("An X509 certificate with Basic Constraint: " +
			"Certificate Authority equals true cannot be used as an identity")
	}

	return msp.getValidationChain(id.cert, false)
}

func (msp *bccspmsp) getUniqueValidationChain(cert *x509.Certificate, opts x509.VerifyOptions) ([]*x509.Certificate, error) {
	
	if msp.opts == nil {
		return nil, errors.New("the supplied identity has no verify options")
	}
	validationChains, err := cert.Verify(opts)
	if err != nil {
		return nil, errors.WithMessage(err, "the supplied identity is not valid")
	}

	
	
	
	if len(validationChains) != 1 {
		return nil, errors.Errorf("this MSP only supports a single validation chain, got %d", len(validationChains))
	}

	return validationChains[0], nil
}

func (msp *bccspmsp) getValidationChain(cert *x509.Certificate, isIntermediateChain bool) ([]*x509.Certificate, error) {
	validationChain, err := msp.getUniqueValidationChain(cert, msp.getValidityOptsForCert(cert))
	if err != nil {
		return nil, errors.WithMessage(err, "failed getting validation chain")
	}

	
	if len(validationChain) < 2 {
		return nil, errors.Errorf("expected a chain of length at least 2, got %d", len(validationChain))
	}

	
	
	parentPosition := 1
	if isIntermediateChain {
		parentPosition = 0
	}
	if msp.certificationTreeInternalNodesMap[string(validationChain[parentPosition].Raw)] {
		return nil, errors.Errorf("invalid validation chain. Parent certificate should be a leaf of the certification tree [%v]", cert.Raw)
	}
	return validationChain, nil
}



func (msp *bccspmsp) getCertificationChainIdentifier(id Identity) ([]byte, error) {
	chain, err := msp.getCertificationChain(id)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed getting certification chain for [%v]", id)
	}

	
	
	return msp.getCertificationChainIdentifierFromChain(chain[1:])
}

func (msp *bccspmsp) getCertificationChainIdentifierFromChain(chain []*x509.Certificate) ([]byte, error) {
	
	
	hashOpt, err := bccsp.GetHashOpt(msp.cryptoConfig.IdentityIdentifierHashFunction)
	if err != nil {
		return nil, errors.WithMessage(err, "failed getting hash function options")
	}

	hf, err := msp.bccsp.GetHash(hashOpt)
	if err != nil {
		return nil, errors.WithMessage(err, "failed getting hash function when computing certification chain identifier")
	}
	for i := 0; i < len(chain); i++ {
		hf.Write(chain[i].Raw)
	}
	return hf.Sum(nil), nil
}




func (msp *bccspmsp) sanitizeCert(cert *x509.Certificate) (*x509.Certificate, error) {
	if isECDSASignedCert(cert) {
		
		var parentCert *x509.Certificate
		chain, err := msp.getUniqueValidationChain(cert, msp.getValidityOptsForCert(cert))
		if err != nil {
			return nil, err
		}

		
		
		if cert.IsCA && len(chain) == 1 {
			
			parentCert = cert
		} else {
			parentCert = chain[1]
		}

		
		cert, err = sanitizeECDSASignedCert(cert, parentCert)
		if err != nil {
			return nil, err
		}
	}
	return cert, nil
}




func (msp *bccspmsp) IsWellFormed(identity *m.SerializedIdentity) error {
	bl, _ := pem.Decode(identity.IdBytes)
	if bl == nil {
		return errors.New("PEM decoding resulted in an empty block")
	}
	
	
	
	
	
	if bl.Type != "CERTIFICATE" && bl.Type != "" {
		return errors.Errorf("pem type is %s, should be 'CERTIFICATE' or missing", bl.Type)
	}
	_, err := x509.ParseCertificate(bl.Bytes)
	return err
}
