/*
Copyright IBM Corp. 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccpackage

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain-protos-go/common"
	mspprotos "github.com/mcc-github/blockchain-protos-go/msp"
	"github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/cauthdsl"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/protoutil"
)

func ownerCreateCCDepSpec(codepackage []byte, sigpolicy *common.SignaturePolicyEnvelope, owner msp.SigningIdentity) (*common.Envelope, error) {
	cds := &peer.ChaincodeDeploymentSpec{CodePackage: codepackage}
	return OwnerCreateSignedCCDepSpec(cds, sigpolicy, owner)
}


func createInstantiationPolicy(mspid string, role mspprotos.MSPRole_MSPRoleType) *common.SignaturePolicyEnvelope {
	principals := []*mspprotos.MSPPrincipal{{
		PrincipalClassification: mspprotos.MSPPrincipal_ROLE,
		Principal:               protoutil.MarshalOrPanic(&mspprotos.MSPRole{Role: role, MspIdentifier: mspid})}}
	sigspolicy := []*common.SignaturePolicy{cauthdsl.SignedBy(int32(0))}

	
	p := &common.SignaturePolicyEnvelope{
		Version:    0,
		Rule:       cauthdsl.NOutOf(1, sigspolicy),
		Identities: principals,
	}

	return p
}

func TestOwnerCreateSignedCCDepSpec(t *testing.T) {
	mspid, _ := localmsp.GetIdentifier()
	sigpolicy := createInstantiationPolicy(mspid, mspprotos.MSPRole_ADMIN)
	env, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, signer)
	if err != nil || env == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}
}

func TestAddSignature(t *testing.T) {
	mspid, _ := localmsp.GetIdentifier()
	sigpolicy := createInstantiationPolicy(mspid, mspprotos.MSPRole_ADMIN)
	env, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, signer)
	if err != nil || env == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}
	
	env, err = SignExistingPackage(env, signer)
	if err != nil || env == nil {
		t.Fatalf("error signing existing package %s", err)
		return
	}
	
	env, err = SignExistingPackage(env, signer)
	if err != nil || env == nil {
		t.Fatalf("error signing existing package %s", err)
		return
	}

	p := &common.Payload{}
	if err = proto.Unmarshal(env.Payload, p); err != nil {
		t.Fatalf("fatal error unmarshal payload")
		return
	}

	sigdepspec := &peer.SignedChaincodeDeploymentSpec{}
	if err = proto.Unmarshal(p.Data, sigdepspec); err != nil || sigdepspec == nil {
		t.Fatalf("fatal error unmarshal sigdepspec")
		return
	}

	if len(sigdepspec.OwnerEndorsements) != 3 {
		t.Fatalf("invalid number of endorsements %d", len(sigdepspec.OwnerEndorsements))
		return
	}
}

func TestMissingSigaturePolicy(t *testing.T) {
	env, err := ownerCreateCCDepSpec([]byte("codepackage"), nil, signer)
	if err == nil || env != nil {
		t.Fatalf("expected error on missing signature policy")
		return
	}
}

func TestCreateSignedCCDepSpecForInstall(t *testing.T) {
	mspid, _ := localmsp.GetIdentifier()
	sigpolicy := createInstantiationPolicy(mspid, mspprotos.MSPRole_ADMIN)
	env1, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, nil)
	if err != nil || env1 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	env2, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, nil)
	if err != nil || env2 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	pack := []*common.Envelope{env1, env2}
	env, err := CreateSignedCCDepSpecForInstall(pack)
	if err != nil || env == nil {
		t.Fatalf("error creating install package %s", err)
		return
	}

	p := &common.Payload{}
	if err = proto.Unmarshal(env.Payload, p); err != nil {
		t.Fatalf("fatal error unmarshal payload")
		return
	}

	cip2 := &peer.SignedChaincodeDeploymentSpec{}
	if err = proto.Unmarshal(p.Data, cip2); err != nil {
		t.Fatalf("fatal error unmarshal cip")
		return
	}

	p = &common.Payload{}
	if err = proto.Unmarshal(env1.Payload, p); err != nil {
		t.Fatalf("fatal error unmarshal payload")
		return
	}

	cip1 := &peer.SignedChaincodeDeploymentSpec{}
	if err = proto.Unmarshal(p.Data, cip1); err != nil {
		t.Fatalf("fatal error unmarshal cip")
		return
	}

	if err = ValidateCip(cip1, cip2); err != nil {
		t.Fatalf("fatal error validating cip1 (%v) against cip2(%v)", cip1, cip2)
		return
	}
}

func TestCreateSignedCCDepSpecForInstallWithEndorsements(t *testing.T) {
	mspid, _ := localmsp.GetIdentifier()
	sigpolicy := createInstantiationPolicy(mspid, mspprotos.MSPRole_ADMIN)
	env1, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, signer)
	if err != nil || env1 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	env2, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, signer)
	if err != nil || env2 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	pack := []*common.Envelope{env1, env2}
	env, err := CreateSignedCCDepSpecForInstall(pack)
	if err != nil || env == nil {
		t.Fatalf("error creating install package %s", err)
		return
	}

	p := &common.Payload{}
	if err = proto.Unmarshal(env.Payload, p); err != nil {
		t.Fatalf("fatal error unmarshal payload")
		return
	}

	cip2 := &peer.SignedChaincodeDeploymentSpec{}
	if err = proto.Unmarshal(p.Data, cip2); err != nil {
		t.Fatalf("fatal error unmarshal cip")
		return
	}

	if len(cip2.OwnerEndorsements) != 2 {
		t.Fatalf("invalid number of endorsements %d", len(cip2.OwnerEndorsements))
		return
	}

	p = &common.Payload{}
	if err = proto.Unmarshal(env1.Payload, p); err != nil {
		t.Fatalf("fatal error unmarshal payload")
		return
	}

	cip1 := &peer.SignedChaincodeDeploymentSpec{}
	if err = proto.Unmarshal(p.Data, cip1); err != nil {
		t.Fatalf("fatal error unmarshal cip")
		return
	}

	if len(cip1.OwnerEndorsements) != 1 {
		t.Fatalf("invalid number of endorsements %d", len(cip1.OwnerEndorsements))
		return
	}
}

func TestMismatchedCodePackages(t *testing.T) {
	mspid, _ := localmsp.GetIdentifier()
	sigpolicy := createInstantiationPolicy(mspid, mspprotos.MSPRole_ADMIN)
	env1, err := ownerCreateCCDepSpec([]byte("codepackage1"), sigpolicy, nil)
	if err != nil || env1 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	env2, err := ownerCreateCCDepSpec([]byte("codepackage2"), sigpolicy, nil)
	if err != nil || env2 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}
	pack := []*common.Envelope{env1, env2}
	env, err := CreateSignedCCDepSpecForInstall(pack)
	if err == nil || env != nil {
		t.Fatalf("expected error creating install from mismatched code package but succeeded")
		return
	}
}

func TestMismatchedEndorsements(t *testing.T) {
	mspid, _ := localmsp.GetIdentifier()
	sigpolicy := createInstantiationPolicy(mspid, mspprotos.MSPRole_ADMIN)
	env1, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, signer)
	if err != nil || env1 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	env2, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy, nil)
	if err != nil || env2 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}
	pack := []*common.Envelope{env1, env2}
	env, err := CreateSignedCCDepSpecForInstall(pack)
	if err == nil || env != nil {
		t.Fatalf("expected error creating install from mismatched endorsed package but succeeded")
		return
	}
}

func TestMismatchedSigPolicy(t *testing.T) {
	sigpolicy1 := createInstantiationPolicy("mspid1", mspprotos.MSPRole_ADMIN)
	env1, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy1, signer)
	if err != nil || env1 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}

	sigpolicy2 := createInstantiationPolicy("mspid2", mspprotos.MSPRole_ADMIN)
	env2, err := ownerCreateCCDepSpec([]byte("codepackage"), sigpolicy2, signer)
	if err != nil || env2 == nil {
		t.Fatalf("error owner creating package %s", err)
		return
	}
	pack := []*common.Envelope{env1, env2}
	env, err := CreateSignedCCDepSpecForInstall(pack)
	if err == nil || env != nil {
		t.Fatalf("expected error creating install from mismatched signature policies but succeeded")
		return
	}
}

var localmsp msp.MSP
var signer msp.SigningIdentity
var signerSerialized []byte

func TestMain(m *testing.M) {
	
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not initialize msp")
		return
	}
	localmsp = mspmgmt.GetLocalMSP()
	if localmsp == nil {
		os.Exit(-1)
		fmt.Printf("Could not get msp")
		return
	}
	signer, err = localmsp.GetDefaultSigningIdentity()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not get signer")
		return
	}

	signerSerialized, err = signer.Serialize()
	if err != nil {
		os.Exit(-1)
		fmt.Printf("Could not serialize identity")
		return
	}

	os.Exit(m.Run())
}
