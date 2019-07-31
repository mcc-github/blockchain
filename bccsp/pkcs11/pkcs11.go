/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package pkcs11

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/miekg/pkcs11"
	"go.uber.org/zap/zapcore"
)

func loadLib(lib, pin, label string) (*pkcs11.Ctx, uint, *pkcs11.SessionHandle, error) {
	var slot uint
	logger.Debugf("Loading pkcs11 library [%s]\n", lib)
	if lib == "" {
		return nil, slot, nil, fmt.Errorf("No PKCS11 library default")
	}

	ctx := pkcs11.New(lib)
	if ctx == nil {
		return nil, slot, nil, fmt.Errorf("Instantiate failed [%s]", lib)
	}

	ctx.Initialize()
	slots, err := ctx.GetSlotList(true)
	if err != nil {
		return nil, slot, nil, fmt.Errorf("Could not get Slot List [%s]", err)
	}
	found := false
	for _, s := range slots {
		info, errToken := ctx.GetTokenInfo(s)
		if errToken != nil {
			continue
		}
		logger.Debugf("Looking for %s, found label %s\n", label, info.Label)
		if label == info.Label {
			found = true
			slot = s
			break
		}
	}
	if !found {
		return nil, slot, nil, fmt.Errorf("Could not find token with label %s", label)
	}

	var session pkcs11.SessionHandle
	for i := 0; i < 10; i++ {
		session, err = ctx.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
		if err != nil {
			logger.Warningf("OpenSession failed, retrying [%s]\n", err)
		} else {
			break
		}
	}
	if err != nil {
		logger.Fatalf("OpenSession [%s]\n", err)
	}
	logger.Debugf("Created new pkcs11 session %+v on slot %d\n", session, slot)

	if pin == "" {
		return nil, slot, nil, fmt.Errorf("No PIN set")
	}
	err = ctx.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		if err != pkcs11.Error(pkcs11.CKR_USER_ALREADY_LOGGED_IN) {
			return nil, slot, nil, fmt.Errorf("Login failed [%s]", err)
		}
	}

	return ctx, slot, &session, nil
}

func (csp *impl) getSession() (session pkcs11.SessionHandle) {
	select {
	case session = <-csp.sessions:
		logger.Debugf("Reusing existing pkcs11 session %+v on slot %d\n", session, csp.slot)

	default:
		
		var s pkcs11.SessionHandle
		var err error
		for i := 0; i < 10; i++ {
			s, err = csp.ctx.OpenSession(csp.slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
			if err != nil {
				logger.Warningf("OpenSession failed, retrying [%s]\n", err)
			} else {
				break
			}
		}
		if err != nil {
			panic(fmt.Errorf("OpenSession failed [%s]", err))
		}
		logger.Debugf("Created new pkcs11 session %+v on slot %d\n", s, csp.slot)
		session = s
	}
	return session
}

func (csp *impl) returnSession(session pkcs11.SessionHandle) {
	select {
	case csp.sessions <- session:
		
	default:
		
		csp.ctx.CloseSession(session)
	}
}



func (csp *impl) getECKey(ski []byte) (pubKey *ecdsa.PublicKey, isPriv bool, err error) {
	p11lib := csp.ctx
	session := csp.getSession()
	defer csp.returnSession(session)
	isPriv = true
	_, err = findKeyPairFromSKI(p11lib, session, ski, privateKeyType)
	if err != nil {
		isPriv = false
		logger.Debugf("Private key not found [%s] for SKI [%s], looking for Public key", err, hex.EncodeToString(ski))
	}

	publicKey, err := findKeyPairFromSKI(p11lib, session, ski, publicKeyType)
	if err != nil {
		return nil, false, fmt.Errorf("Public key not found [%s] for SKI [%s]", err, hex.EncodeToString(ski))
	}

	ecpt, marshaledOid, err := ecPoint(p11lib, session, *publicKey)
	if err != nil {
		return nil, false, fmt.Errorf("Public key not found [%s] for SKI [%s]", err, hex.EncodeToString(ski))
	}

	curveOid := new(asn1.ObjectIdentifier)
	_, err = asn1.Unmarshal(marshaledOid, curveOid)
	if err != nil {
		return nil, false, fmt.Errorf("Failed Unmarshaling Curve OID [%s]\n%s", err.Error(), hex.EncodeToString(marshaledOid))
	}

	curve := namedCurveFromOID(*curveOid)
	if curve == nil {
		return nil, false, fmt.Errorf("Cound not recognize Curve from OID")
	}
	x, y := elliptic.Unmarshal(curve, ecpt)
	if x == nil {
		return nil, false, fmt.Errorf("Failed Unmarshaling Public Key")
	}

	pubKey = &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return pubKey, isPriv, nil
}
















var (
	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
)

func namedCurveFromOID(oid asn1.ObjectIdentifier) elliptic.Curve {
	switch {
	case oid.Equal(oidNamedCurveP224):
		return elliptic.P224()
	case oid.Equal(oidNamedCurveP256):
		return elliptic.P256()
	case oid.Equal(oidNamedCurveP384):
		return elliptic.P384()
	case oid.Equal(oidNamedCurveP521):
		return elliptic.P521()
	}
	return nil
}

func (csp *impl) generateECKey(curve asn1.ObjectIdentifier, ephemeral bool) (ski []byte, pubKey *ecdsa.PublicKey, err error) {
	p11lib := csp.ctx
	session := csp.getSession()
	defer csp.returnSession(session)

	id := nextIDCtr()
	publabel := fmt.Sprintf("BCPUB%s", id.Text(16))
	prvlabel := fmt.Sprintf("BCPRV%s", id.Text(16))

	marshaledOID, err := asn1.Marshal(curve)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not marshal OID [%s]", err.Error())
	}

	pubkeyT := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, !ephemeral),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, marshaledOID),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),

		pkcs11.NewAttribute(pkcs11.CKA_ID, publabel),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, publabel),
	}

	prvkeyT := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, !ephemeral),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),

		pkcs11.NewAttribute(pkcs11.CKA_ID, prvlabel),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, prvlabel),

		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
	}

	pub, prv, err := p11lib.GenerateKeyPair(session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_EC_KEY_PAIR_GEN, nil)},
		pubkeyT, prvkeyT)

	if err != nil {
		return nil, nil, fmt.Errorf("P11: keypair generate failed [%s]", err)
	}

	ecpt, _, err := ecPoint(p11lib, session, pub)
	if err != nil {
		return nil, nil, fmt.Errorf("Error querying EC-point: [%s]", err)
	}
	hash := sha256.Sum256(ecpt)
	ski = hash[:]

	
	setskiT := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, ski),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, hex.EncodeToString(ski)),
	}

	logger.Infof("Generated new P11 key, SKI %x\n", ski)
	err = p11lib.SetAttributeValue(session, pub, setskiT)
	if err != nil {
		return nil, nil, fmt.Errorf("P11: set-ID-to-SKI[public] failed [%s]", err)
	}

	err = p11lib.SetAttributeValue(session, prv, setskiT)
	if err != nil {
		return nil, nil, fmt.Errorf("P11: set-ID-to-SKI[private] failed [%s]", err)
	}

	
	if csp.immutable {
		setCKAModifiable := []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_MODIFIABLE, false),
		}

		_, pubCopyerror := p11lib.CopyObject(session, pub, setCKAModifiable)
		if pubCopyerror != nil {
			return nil, nil, fmt.Errorf("P11: Public Key copy failed with error [%s] . Please contact your HSM vendor", pubCopyerror)
		}

		pubKeyDestroyError := p11lib.DestroyObject(session, pub)
		if pubKeyDestroyError != nil {
			return nil, nil, fmt.Errorf("P11: Public Key destroy failed with error [%s]. Please contact your HSM vendor", pubCopyerror)
		}

		_, prvCopyerror := p11lib.CopyObject(session, prv, setCKAModifiable)
		if prvCopyerror != nil {
			return nil, nil, fmt.Errorf("P11: Private Key copy failed with error [%s]. Please contact your HSM vendor", prvCopyerror)
		}
		prvKeyDestroyError := p11lib.DestroyObject(session, prv)
		if pubKeyDestroyError != nil {
			return nil, nil, fmt.Errorf("P11: Private Key destroy failed with error [%s]. Please contact your HSM vendor", prvKeyDestroyError)
		}
	}

	nistCurve := namedCurveFromOID(curve)
	if curve == nil {
		return nil, nil, fmt.Errorf("Cound not recognize Curve from OID")
	}
	x, y := elliptic.Unmarshal(nistCurve, ecpt)
	if x == nil {
		return nil, nil, fmt.Errorf("Failed Unmarshaling Public Key")
	}

	pubGoKey := &ecdsa.PublicKey{Curve: nistCurve, X: x, Y: y}

	if logger.IsEnabledFor(zapcore.DebugLevel) {
		listAttrs(p11lib, session, prv)
		listAttrs(p11lib, session, pub)
	}

	return ski, pubGoKey, nil
}

func (csp *impl) signP11ECDSA(ski []byte, msg []byte) (R, S *big.Int, err error) {
	p11lib := csp.ctx
	session := csp.getSession()
	defer csp.returnSession(session)

	privateKey, err := findKeyPairFromSKI(p11lib, session, ski, privateKeyType)
	if err != nil {
		return nil, nil, fmt.Errorf("Private key not found [%s]", err)
	}

	err = p11lib.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}, *privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Sign-initialize  failed [%s]", err)
	}

	var sig []byte

	sig, err = p11lib.Sign(session, msg)
	if err != nil {
		return nil, nil, fmt.Errorf("P11: sign failed [%s]", err)
	}

	R = new(big.Int)
	S = new(big.Int)
	R.SetBytes(sig[0 : len(sig)/2])
	S.SetBytes(sig[len(sig)/2:])

	return R, S, nil
}

func (csp *impl) verifyP11ECDSA(ski []byte, msg []byte, R, S *big.Int, byteSize int) (bool, error) {
	p11lib := csp.ctx
	session := csp.getSession()
	defer csp.returnSession(session)

	logger.Debugf("Verify ECDSA\n")

	publicKey, err := findKeyPairFromSKI(p11lib, session, ski, publicKeyType)
	if err != nil {
		return false, fmt.Errorf("Public key not found [%s]", err)
	}

	r := R.Bytes()
	s := S.Bytes()

	
	sig := make([]byte, 2*byteSize)
	copy(sig[byteSize-len(r):byteSize], r)
	copy(sig[2*byteSize-len(s):], s)

	err = p11lib.VerifyInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)},
		*publicKey)
	if err != nil {
		return false, fmt.Errorf("PKCS11: Verify-initialize [%s]", err)
	}
	err = p11lib.Verify(session, msg, sig)
	if err == pkcs11.Error(pkcs11.CKR_SIGNATURE_INVALID) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("PKCS11: Verify failed [%s]", err)
	}

	return true, nil
}

type keyType int8

const (
	publicKeyType keyType = iota
	privateKeyType
)

func findKeyPairFromSKI(mod *pkcs11.Ctx, session pkcs11.SessionHandle, ski []byte, keyType keyType) (*pkcs11.ObjectHandle, error) {
	ktype := pkcs11.CKO_PUBLIC_KEY
	if keyType == privateKeyType {
		ktype = pkcs11.CKO_PRIVATE_KEY
	}

	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, ktype),
		pkcs11.NewAttribute(pkcs11.CKA_ID, ski),
	}
	if err := mod.FindObjectsInit(session, template); err != nil {
		return nil, err
	}

	
	objs, _, err := mod.FindObjects(session, 1)
	if err != nil {
		return nil, err
	}
	if err = mod.FindObjectsFinal(session); err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		return nil, fmt.Errorf("Key not found [%s]", hex.Dump(ski))
	}

	return &objs[0], nil
}










































func ecPoint(p11lib *pkcs11.Ctx, session pkcs11.SessionHandle, key pkcs11.ObjectHandle) (ecpt, oid []byte, err error) {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, nil),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, nil),
	}

	attr, err := p11lib.GetAttributeValue(session, key, template)
	if err != nil {
		return nil, nil, fmt.Errorf("PKCS11: get(EC point) [%s]", err)
	}

	for _, a := range attr {
		if a.Type == pkcs11.CKA_EC_POINT {
			logger.Debugf("EC point: attr type %d/0x%x, len %d\n%s\n", a.Type, a.Type, len(a.Value), hex.Dump(a.Value))

			
			if ((len(a.Value) % 2) == 0) &&
				(byte(0x04) == a.Value[0]) &&
				(byte(0x04) == a.Value[len(a.Value)-1]) {
				logger.Debugf("Detected opencryptoki bug, trimming trailing 0x04")
				ecpt = a.Value[0 : len(a.Value)-1] 
			} else if byte(0x04) == a.Value[0] && byte(0x04) == a.Value[2] {
				logger.Debugf("Detected SoftHSM bug, trimming leading 0x04 0xXX")
				ecpt = a.Value[2:len(a.Value)]
			} else {
				ecpt = a.Value
			}
		} else if a.Type == pkcs11.CKA_EC_PARAMS {
			logger.Debugf("EC point: attr type %d/0x%x, len %d\n%s\n", a.Type, a.Type, len(a.Value), hex.Dump(a.Value))

			oid = a.Value
		}
	}
	if oid == nil || ecpt == nil {
		return nil, nil, fmt.Errorf("CKA_EC_POINT not found, perhaps not an EC Key?")
	}

	return ecpt, oid, nil
}

func listAttrs(p11lib *pkcs11.Ctx, session pkcs11.SessionHandle, obj pkcs11.ObjectHandle) {
	var cktype, ckclass uint
	var ckaid, cklabel []byte

	if p11lib == nil {
		return
	}

	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, ckclass),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, cktype),
		pkcs11.NewAttribute(pkcs11.CKA_ID, ckaid),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, cklabel),
	}

	
	attr, err := p11lib.GetAttributeValue(session, obj, template)
	if err != nil {
		logger.Debugf("P11: get(attrlist) [%s]\n", err)
	}

	for _, a := range attr {
		
		logger.Debugf("ListAttr: type %d/0x%x, length %d\n%s", a.Type, a.Type, len(a.Value), hex.Dump(a.Value))
	}
}

var (
	bigone  = new(big.Int).SetInt64(1)
	idCtr   = new(big.Int)
	idMutex sync.Mutex
)

func nextIDCtr() *big.Int {
	idMutex.Lock()
	idCtr = new(big.Int).Add(idCtr, bigone)
	idMutex.Unlock()
	return idCtr
}
