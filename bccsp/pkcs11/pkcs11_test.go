

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
package pkcs11

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"testing"

	"github.com/mcc-github/blockchain/bccsp"
	"github.com/miekg/pkcs11"
	"github.com/stretchr/testify/assert"
)

func TestKeyGenFailures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestKeyGenFailures")
	}
	var testOpts bccsp.KeyGenOpts
	ki := currentBCCSP
	_, err := ki.KeyGen(testOpts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Opts parameter. It must not be nil.")
}

func TestLoadLib(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestLoadLib")
	}
	
	lib, pin, label := FindPKCS11Lib()

	
	_, _, _, err := loadLib("", pin, label)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No PKCS11 library default")

	
	_, _, _, err = loadLib("badLib", pin, label)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Instantiate failed")

	
	_, _, _, err = loadLib(lib, pin, "badLabel")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Could not find token with label")

	
	_, _, _, err = loadLib(lib, "", label)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No PIN set")
}

func TestOIDFromNamedCurve(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestOIDFromNamedCurve")
	}
	
	testOID, boolValue := oidFromNamedCurve(elliptic.P224())
	assert.Equal(t, oidNamedCurveP224, testOID, "Did not receive expected OID for elliptic.P224")
	assert.Equal(t, true, boolValue, "Did not receive a true value when acquiring OID for elliptic.P224")

	
	testOID, boolValue = oidFromNamedCurve(elliptic.P256())
	assert.Equal(t, oidNamedCurveP256, testOID, "Did not receive expected OID for elliptic.P256")
	assert.Equal(t, true, boolValue, "Did not receive a true value when acquiring OID for elliptic.P256")

	
	testOID, boolValue = oidFromNamedCurve(elliptic.P384())
	assert.Equal(t, oidNamedCurveP384, testOID, "Did not receive expected OID for elliptic.P384")
	assert.Equal(t, true, boolValue, "Did not receive a true value when acquiring OID for elliptic.P384")

	
	testOID, boolValue = oidFromNamedCurve(elliptic.P521())
	assert.Equal(t, oidNamedCurveP521, testOID, "Did not receive expected OID for elliptic.P521")
	assert.Equal(t, true, boolValue, "Did not receive a true value when acquiring OID for elliptic.P521")

	var testCurve elliptic.Curve
	testOID, _ = oidFromNamedCurve(testCurve)
	if testOID != nil {
		t.Fatal("Expected nil to be returned.")
	}
}

func TestNamedCurveFromOID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestNamedCurveFromOID")
	}
	
	namedCurve := namedCurveFromOID(oidNamedCurveP224)
	assert.Equal(t, elliptic.P224(), namedCurve, "Did not receive expected named curve for oidNamedCurveP224")

	
	namedCurve = namedCurveFromOID(oidNamedCurveP256)
	assert.Equal(t, elliptic.P256(), namedCurve, "Did not receive expected named curve for oidNamedCurveP256")

	
	namedCurve = namedCurveFromOID(oidNamedCurveP384)
	assert.Equal(t, elliptic.P384(), namedCurve, "Did not receive expected named curve for oidNamedCurveP384")

	
	namedCurve = namedCurveFromOID(oidNamedCurveP521)
	assert.Equal(t, elliptic.P521(), namedCurve, "Did not receive expected named curved for oidNamedCurveP521")

	testAsn1Value := asn1.ObjectIdentifier{4, 9, 15, 1}
	namedCurve = namedCurveFromOID(testAsn1Value)
	if namedCurve != nil {
		t.Fatal("Expected nil to be returned.")
	}
}

func TestPKCS11GetSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestPKCS11GetSession")
	}
	var sessions []pkcs11.SessionHandle
	for i := 0; i < 3*sessionCacheSize; i++ {
		sessions = append(sessions, currentBCCSP.(*impl).getSession())
	}

	
	for _, session := range sessions {
		currentBCCSP.(*impl).returnSession(session)
	}
	sessions = nil

	
	oldSlot := currentBCCSP.(*impl).slot
	currentBCCSP.(*impl).slot = ^uint(0)

	
	for i := 0; i < sessionCacheSize; i++ {
		sessions = append(sessions, currentBCCSP.(*impl).getSession())
	}

	
	assert.Panics(t, func() {
		currentBCCSP.(*impl).getSession()
	}, "Should not been able to create another session")

	
	for _, session := range sessions {
		currentBCCSP.(*impl).returnSession(session)
	}
	currentBCCSP.(*impl).slot = oldSlot
}

func TestPKCS11ECKeySignVerify(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestPKCS11ECKeySignVerify")
	}

	msg1 := []byte("This is my very authentic message")
	msg2 := []byte("This is my very unauthentic message")
	hash1, _ := currentBCCSP.Hash(msg1, &bccsp.SHAOpts{})
	hash2, _ := currentBCCSP.Hash(msg2, &bccsp.SHAOpts{})

	var oid asn1.ObjectIdentifier
	if currentTestConfig.securityLevel == 256 {
		oid = oidNamedCurveP256
	} else if currentTestConfig.securityLevel == 384 {
		oid = oidNamedCurveP384
	}

	key, pubKey, err := currentBCCSP.(*impl).generateECKey(oid, true)
	if err != nil {
		t.Fatalf("Failed generating Key [%s]", err)
	}

	R, S, err := currentBCCSP.(*impl).signP11ECDSA(key, hash1)

	if err != nil {
		t.Fatalf("Failed signing message [%s]", err)
	}

	_, _, err = currentBCCSP.(*impl).signP11ECDSA(nil, hash1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Private key not found")

	pass, err := currentBCCSP.(*impl).verifyP11ECDSA(key, hash1, R, S, currentTestConfig.securityLevel/8)
	if err != nil {
		t.Fatalf("Error verifying message 1 [%s]", err)
	}
	if pass == false {
		t.Fatal("Signature should match!")
	}

	pass = ecdsa.Verify(pubKey, hash1, R, S)
	if pass == false {
		t.Fatal("Signature should match with software verification!")
	}

	pass, err = currentBCCSP.(*impl).verifyP11ECDSA(key, hash2, R, S, currentTestConfig.securityLevel/8)
	if err != nil {
		t.Fatalf("Error verifying message 2 [%s]", err)
	}

	if pass != false {
		t.Fatal("Signature should not match!")
	}

	pass = ecdsa.Verify(pubKey, hash2, R, S)
	if pass != false {
		t.Fatal("Signature should not match with software verification!")
	}
}
