/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp"
	idemixbccsp "github.com/mcc-github/blockchain/bccsp/idemix"
	"github.com/mcc-github/blockchain/bccsp/sw"
	m "github.com/mcc-github/blockchain/protos/msp"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const (
	
	AttributeIndexOU = iota

	
	AttributeIndexRole

	
	AttributeIndexEnrollmentId

	
	AttributeIndexRevocationHandle
)

const (
	
	AttributeNameOU = "OU"

	
	AttributeNameRole = "Role"

	
	AttributeNameEnrollmentId = "EnrollmentID"

	
	AttributeNameRevocationHandle = "RevocationHandle"
)


const rhIndex = 3




var discloseFlags = []byte{1, 1, 0, 0}

type idemixmsp struct {
	csp          bccsp.BCCSP
	version      MSPVersion
	ipk          bccsp.Key
	signer       *idemixSigningIdentity
	name         string
	revocationPK bccsp.Key
	epoch        int
}


func newIdemixMsp(version MSPVersion) (MSP, error) {
	mspLogger.Debugf("Creating Idemix-based MSP instance")

	csp, err := idemixbccsp.New(sw.NewDummyKeyStore())
	if err != nil {
		panic(fmt.Sprintf("unexpected condition, error received [%s]", err))
	}

	msp := idemixmsp{csp: csp}
	msp.version = version
	return &msp, nil
}

func (msp *idemixmsp) Setup(conf1 *m.MSPConfig) error {
	mspLogger.Debugf("Setting up Idemix-based MSP instance")

	if conf1 == nil {
		return errors.Errorf("setup error: nil conf reference")
	}

	if conf1.Type != int32(IDEMIX) {
		return errors.Errorf("setup error: config is not of type IDEMIX")
	}

	var conf m.IdemixMSPConfig
	err := proto.Unmarshal(conf1.Config, &conf)
	if err != nil {
		return errors.Wrap(err, "failed unmarshalling idemix msp config")
	}

	msp.name = conf.Name
	mspLogger.Debugf("Setting up Idemix MSP instance %s", msp.name)

	
	IssuerPublicKey, err := msp.csp.KeyImport(
		conf.Ipk,
		&bccsp.IdemixIssuerPublicKeyImportOpts{
			Temporary: true,
			AttributeNames: []string{
				AttributeNameOU,
				AttributeNameRole,
				AttributeNameEnrollmentId,
				AttributeNameRevocationHandle,
			},
		})
	if err != nil {
		importErr, ok := errors.Cause(err).(*bccsp.IdemixIssuerPublicKeyImporterError)
		if !ok {
			panic("unexpected condition, BCCSP did not return the expected *bccsp.IdemixIssuerPublicKeyImporterError")
		}
		switch importErr.Type {
		case bccsp.IdemixIssuerPublicKeyImporterUnmarshallingError:
			return errors.WithMessage(err, "failed to unmarshal ipk from idemix msp config")
		case bccsp.IdemixIssuerPublicKeyImporterHashError:
			return errors.WithMessage(err, "setting the hash of the issuer public key failed")
		case bccsp.IdemixIssuerPublicKeyImporterValidationError:
			return errors.WithMessage(err, "cannot setup idemix msp with invalid public key")
		case bccsp.IdemixIssuerPublicKeyImporterNumAttributesError:
			fallthrough
		case bccsp.IdemixIssuerPublicKeyImporterAttributeNameError:
			return errors.Errorf("issuer public key must have have attributes OU, Role, EnrollmentId, and RevocationHandle")
		default:
			panic(fmt.Sprintf("unexpected condtion, issuer public key import error not valid, got [%d]", importErr.Type))
		}
	}
	msp.ipk = IssuerPublicKey

	
	RevocationPublicKey, err := msp.csp.KeyImport(
		conf.RevocationPk,
		&bccsp.IdemixRevocationPublicKeyImportOpts{Temporary: true},
	)
	if err != nil {
		return errors.WithMessage(err, "failed to import revocation public key")
	}
	msp.revocationPK = RevocationPublicKey

	if conf.Signer == nil {
		
		mspLogger.Debug("idemix msp setup as verification only msp (no key material found)")
		return nil
	}

	

	
	UserKey, err := msp.csp.KeyImport(conf.Signer.Sk, &bccsp.IdemixUserSecretKeyImportOpts{Temporary: true})
	if err != nil {
		return errors.WithMessage(err, "failed importing signer secret key")
	}

	
	NymKey, err := msp.csp.KeyDeriv(UserKey, &bccsp.IdemixNymKeyDerivationOpts{Temporary: true, IssuerPK: IssuerPublicKey})
	if err != nil {
		return errors.WithMessage(err, "failed deriving nym")
	}
	NymPublicKey, err := NymKey.PublicKey()
	if err != nil {
		return errors.Wrapf(err, "failed getting public nym key")
	}

	role := &m.MSPRole{
		MspIdentifier: msp.name,
		Role:          m.MSPRole_MEMBER,
	}
	if checkRole(int(conf.Signer.Role), ADMIN) {
		role.Role = m.MSPRole_ADMIN
	}

	ou := &m.OrganizationUnit{
		MspIdentifier:                msp.name,
		OrganizationalUnitIdentifier: conf.Signer.OrganizationalUnitIdentifier,
		CertifiersIdentifier:         IssuerPublicKey.SKI(),
	}

	enrollmentId := conf.Signer.EnrollmentId

	
	valid, err := msp.csp.Verify(
		UserKey,
		conf.Signer.Cred,
		bccsp.IdemixEmptyDigest(),
		&bccsp.IdemixCredentialSignerOpts{
			IssuerPK: IssuerPublicKey,
			Attributes: []bccsp.IdemixAttribute{
				{Type: bccsp.IdemixBytesAttribute, Value: []byte(conf.Signer.OrganizationalUnitIdentifier)},
				{Type: bccsp.IdemixIntAttribute, Value: getIdemixRoleFromMSPRole(role)},
				{Type: bccsp.IdemixBytesAttribute, Value: []byte(enrollmentId)},
				{Type: bccsp.IdemixHiddenAttribute},
			},
		},
	)
	if err != nil || !valid {
		return errors.WithMessage(err, "Credential is not cryptographically valid")
	}

	
	proof, err := msp.csp.Sign(
		UserKey,
		bccsp.IdemixEmptyDigest(),
		&bccsp.IdemixSignerOpts{
			Credential: conf.Signer.Cred,
			Nym:        NymKey,
			IssuerPK:   IssuerPublicKey,
			Attributes: []bccsp.IdemixAttribute{
				{Type: bccsp.IdemixBytesAttribute},
				{Type: bccsp.IdemixIntAttribute},
				{Type: bccsp.IdemixHiddenAttribute},
				{Type: bccsp.IdemixHiddenAttribute},
			},
			RhIndex: rhIndex,
			CRI:     conf.Signer.CredentialRevocationInformation,
		},
	)
	if err != nil {
		return errors.WithMessage(err, "Failed to setup cryptographic proof of identity")
	}

	
	msp.signer = &idemixSigningIdentity{
		idemixidentity: newIdemixIdentity(msp, NymPublicKey, role, ou, proof),
		Cred:           conf.Signer.Cred,
		UserKey:        UserKey,
		NymKey:         NymKey,
		enrollmentId:   enrollmentId}

	return nil
}


func (msp *idemixmsp) GetVersion() MSPVersion {
	return msp.version
}

func (msp *idemixmsp) GetType() ProviderType {
	return IDEMIX
}

func (msp *idemixmsp) GetIdentifier() (string, error) {
	return msp.name, nil
}

func (msp *idemixmsp) GetSigningIdentity(identifier *IdentityIdentifier) (SigningIdentity, error) {
	return nil, errors.Errorf("GetSigningIdentity not implemented")
}

func (msp *idemixmsp) GetDefaultSigningIdentity() (SigningIdentity, error) {
	mspLogger.Debugf("Obtaining default idemix signing identity")

	if msp.signer == nil {
		return nil, errors.Errorf("no default signer setup")
	}
	return msp.signer, nil
}

func (msp *idemixmsp) DeserializeIdentity(serializedID []byte) (Identity, error) {
	sID := &m.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sID)
	if err != nil {
		return nil, errors.Wrap(err, "could not deserialize a SerializedIdentity")
	}

	if sID.Mspid != msp.name {
		return nil, errors.Errorf("expected MSP ID %s, received %s", msp.name, sID.Mspid)
	}

	return msp.deserializeIdentityInternal(sID.GetIdBytes())
}

func (msp *idemixmsp) deserializeIdentityInternal(serializedID []byte) (Identity, error) {
	mspLogger.Debug("idemixmsp: deserializing identity")
	serialized := new(m.SerializedIdemixIdentity)
	err := proto.Unmarshal(serializedID, serialized)
	if err != nil {
		return nil, errors.Wrap(err, "could not deserialize a SerializedIdemixIdentity")
	}
	if serialized.NymX == nil || serialized.NymY == nil {
		return nil, errors.Errorf("unable to deserialize idemix identity: pseudonym is invalid")
	}

	
	var rawNymPublicKey []byte
	rawNymPublicKey = append(rawNymPublicKey, serialized.NymX...)
	rawNymPublicKey = append(rawNymPublicKey, serialized.NymY...)
	NymPublicKey, err := msp.csp.KeyImport(
		rawNymPublicKey,
		&bccsp.IdemixNymPublicKeyImportOpts{Temporary: true},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to import nym public key")
	}

	
	ou := &m.OrganizationUnit{}
	err = proto.Unmarshal(serialized.Ou, ou)
	if err != nil {
		return nil, errors.Wrap(err, "cannot deserialize the OU of the identity")
	}

	
	role := &m.MSPRole{}
	err = proto.Unmarshal(serialized.Role, role)
	if err != nil {
		return nil, errors.Wrap(err, "cannot deserialize the role of the identity")
	}

	return newIdemixIdentity(msp, NymPublicKey, role, ou, serialized.Proof), nil
}

func (msp *idemixmsp) Validate(id Identity) error {
	var identity *idemixidentity
	switch t := id.(type) {
	case *idemixidentity:
		identity = id.(*idemixidentity)
	case *idemixSigningIdentity:
		identity = id.(*idemixSigningIdentity).idemixidentity
	default:
		return errors.Errorf("identity type %T is not recognized", t)
	}

	mspLogger.Debugf("Validating identity %+v", identity)
	if identity.GetMSPIdentifier() != msp.name {
		return errors.Errorf("the supplied identity does not belong to this msp")
	}
	return identity.verifyProof()
}

func (id *idemixidentity) verifyProof() error {
	
	valid, err := id.msp.csp.Verify(
		id.msp.ipk,
		id.associationProof,
		bccsp.IdemixEmptyDigest(),
		&bccsp.IdemixSignerOpts{
			RevocationPublicKey: id.msp.revocationPK,
			Attributes: []bccsp.IdemixAttribute{
				{Type: bccsp.IdemixBytesAttribute, Value: []byte(id.OU.OrganizationalUnitIdentifier)},
				{Type: bccsp.IdemixIntAttribute, Value: getIdemixRoleFromMSPRole(id.Role)},
				{Type: bccsp.IdemixHiddenAttribute},
				{Type: bccsp.IdemixHiddenAttribute},
			},
			RhIndex: rhIndex,
			Epoch:   id.msp.epoch,
		},
	)
	if err == nil && !valid {
		panic("unexpected condition, an error should be returned for an invalid signature")
	}

	return err
}

func (msp *idemixmsp) SatisfiesPrincipal(id Identity, principal *m.MSPPrincipal) error {
	err := msp.Validate(id)
	if err != nil {
		return errors.Wrap(err, "identity is not valid with respect to this MSP")
	}

	return msp.satisfiesPrincipalValidated(id, principal)
}



func (msp *idemixmsp) satisfiesPrincipalValidated(id Identity, principal *m.MSPPrincipal) error {
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
			return nil
		case m.MSPRole_ADMIN:
			mspLogger.Debugf("Checking if identity satisfies ADMIN role for %s", msp.name)
			if id.(*idemixidentity).Role.Role != m.MSPRole_ADMIN {
				return errors.Errorf("user is not an admin")
			}
			return nil
		case m.MSPRole_PEER:
			if msp.version >= MSPv1_3 {
				return errors.Errorf("idemixmsp only supports client use, so it cannot satisfy an MSPRole PEER principal")
			}
			fallthrough
		case m.MSPRole_CLIENT:
			if msp.version >= MSPv1_3 {
				return nil 
			}
			fallthrough
		default:
			return errors.Errorf("invalid MSP role type %d", int32(mspRole.Role))
		}
		
		
	case m.MSPPrincipal_IDENTITY:
		mspLogger.Debugf("Checking if identity satisfies IDENTITY principal")
		idBytes, err := id.Serialize()
		if err != nil {
			return errors.Wrap(err, "could not serialize this identity instance")
		}

		rv := bytes.Compare(idBytes, principal.Principal)
		if rv == 0 {
			return nil
		}
		return errors.Errorf("the identities do not match")

	case m.MSPPrincipal_ORGANIZATION_UNIT:
		ou := &m.OrganizationUnit{}
		err := proto.Unmarshal(principal.Principal, ou)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal OU from principal")
		}

		mspLogger.Debugf("Checking if identity is part of OU \"%s\" of mspid \"%s\"", ou.OrganizationalUnitIdentifier, ou.MspIdentifier)

		
		
		if ou.MspIdentifier != msp.name {
			return errors.Errorf("the identity is a member of a different MSP (expected %s, got %s)", ou.MspIdentifier, id.GetMSPIdentifier())
		}

		if ou.OrganizationalUnitIdentifier != id.(*idemixidentity).OU.OrganizationalUnitIdentifier {
			return errors.Errorf("user is not part of the desired organizational unit")
		}

		return nil
	case m.MSPPrincipal_COMBINED:
		if msp.version <= MSPv1_1 {
			return errors.Errorf("Combined MSP Principals are unsupported in MSPv1_1")
		}

		
		principals := &m.CombinedPrincipal{}
		err := proto.Unmarshal(principal.Principal, principals)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal CombinedPrincipal from principal")
		}
		
		if len(principals.Principals) == 0 {
			return errors.New("no principals in CombinedPrincipal")
		}
		
		
		for _, cp := range principals.Principals {
			err = msp.satisfiesPrincipalValidated(id, cp)
			if err != nil {
				return err
			}
		}
		
		return nil
	case m.MSPPrincipal_ANONYMITY:
		if msp.version <= MSPv1_1 {
			return errors.Errorf("Anonymity MSP Principals are unsupported in MSPv1_1")
		}

		anon := &m.MSPIdentityAnonymity{}
		err := proto.Unmarshal(principal.Principal, anon)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal MSPIdentityAnonymity from principal")
		}
		switch anon.AnonymityType {
		case m.MSPIdentityAnonymity_ANONYMOUS:
			return nil
		case m.MSPIdentityAnonymity_NOMINAL:
			return errors.New("principal is nominal, but idemix MSP is anonymous")
		default:
			return errors.Errorf("unknown principal anonymity type: %d", anon.AnonymityType)
		}
	default:
		return errors.Errorf("invalid principal type %d", int32(principal.PrincipalClassification))
	}
}




func (id *idemixmsp) IsWellFormed(identity *m.SerializedIdentity) error {
	sId := new(m.SerializedIdemixIdentity)
	err := proto.Unmarshal(identity.IdBytes, sId)
	if err != nil {
		return errors.Wrap(err, "not an idemix identity")
	}
	return nil
}

func (msp *idemixmsp) GetTLSRootCerts() [][]byte {
	
	return nil
}

func (msp *idemixmsp) GetTLSIntermediateCerts() [][]byte {
	
	return nil
}

type idemixidentity struct {
	NymPublicKey bccsp.Key
	msp          *idemixmsp
	id           *IdentityIdentifier
	Role         *m.MSPRole
	OU           *m.OrganizationUnit
	
	
	
	associationProof []byte
}

func (id *idemixidentity) Anonymous() bool {
	return true
}

func newIdemixIdentity(msp *idemixmsp, NymPublicKey bccsp.Key, role *m.MSPRole, ou *m.OrganizationUnit, proof []byte) *idemixidentity {
	id := &idemixidentity{}
	id.NymPublicKey = NymPublicKey
	id.msp = msp
	id.Role = role
	id.OU = ou
	id.associationProof = proof

	raw, err := NymPublicKey.Bytes()
	if err != nil {
		panic(fmt.Sprintf("unexpected condition, failed marshalling nym public key [%s]", err))
	}
	id.id = &IdentityIdentifier{
		Mspid: msp.name,
		Id:    bytes.NewBuffer(raw).String(),
	}

	return id
}

func (id *idemixidentity) ExpiresAt() time.Time {
	
	
	return time.Time{}
}

func (id *idemixidentity) GetIdentifier() *IdentityIdentifier {
	return id.id
}

func (id *idemixidentity) GetMSPIdentifier() string {
	mspid, _ := id.msp.GetIdentifier()
	return mspid
}

func (id *idemixidentity) GetOrganizationalUnits() []*OUIdentifier {
	
	certifiersIdentifier, err := id.msp.ipk.Bytes()
	if err != nil {
		mspIdentityLogger.Errorf("Failed to marshal ipk in GetOrganizationalUnits: %s", err)
		return nil
	}

	return []*OUIdentifier{{certifiersIdentifier, id.OU.OrganizationalUnitIdentifier}}
}

func (id *idemixidentity) Validate() error {
	return id.msp.Validate(id)
}

func (id *idemixidentity) Verify(msg []byte, sig []byte) error {
	if mspIdentityLogger.IsEnabledFor(zapcore.DebugLevel) {
		mspIdentityLogger.Debugf("Verify Idemix sig: msg = %s", hex.Dump(msg))
		mspIdentityLogger.Debugf("Verify Idemix sig: sig = %s", hex.Dump(sig))
	}

	_, err := id.msp.csp.Verify(
		id.NymPublicKey,
		sig,
		msg,
		&bccsp.IdemixNymSignerOpts{
			IssuerPK: id.msp.ipk,
		},
	)
	return err
}

func (id *idemixidentity) SatisfiesPrincipal(principal *m.MSPPrincipal) error {
	return id.msp.SatisfiesPrincipal(id, principal)
}

func (id *idemixidentity) Serialize() ([]byte, error) {
	serialized := &m.SerializedIdemixIdentity{}

	raw, err := id.NymPublicKey.Bytes()
	if err != nil {
		return nil, errors.Wrapf(err, "could not serialize nym of identity %s", id.id)
	}
	
	
	serialized.NymX = raw[:len(raw)/2]
	serialized.NymY = raw[len(raw)/2:]
	ouBytes, err := proto.Marshal(id.OU)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal OU of identity %s", id.id)
	}

	roleBytes, err := proto.Marshal(id.Role)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal role of identity %s", id.id)
	}

	serialized.Ou = ouBytes
	serialized.Role = roleBytes
	serialized.Proof = id.associationProof

	idemixIDBytes, err := proto.Marshal(serialized)
	if err != nil {
		return nil, err
	}

	sID := &m.SerializedIdentity{Mspid: id.GetMSPIdentifier(), IdBytes: idemixIDBytes}
	idBytes, err := proto.Marshal(sID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal a SerializedIdentity structure for identity %s", id.id)
	}

	return idBytes, nil
}

type idemixSigningIdentity struct {
	*idemixidentity
	Cred         []byte
	UserKey      bccsp.Key
	NymKey       bccsp.Key
	enrollmentId string
}

func (id *idemixSigningIdentity) Sign(msg []byte) ([]byte, error) {
	mspLogger.Debugf("Idemix identity %s is signing", id.GetIdentifier())

	sig, err := id.msp.csp.Sign(
		id.UserKey,
		msg,
		&bccsp.IdemixNymSignerOpts{
			Nym:      id.NymKey,
			IssuerPK: id.msp.ipk,
		},
	)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (id *idemixSigningIdentity) GetPublicVersion() Identity {
	return id.idemixidentity
}
