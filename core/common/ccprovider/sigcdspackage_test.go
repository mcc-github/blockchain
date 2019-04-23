/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider

import (
	"fmt"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/core/common/ccpackage"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

func processSignedCDS(cds *pb.ChaincodeDeploymentSpec, policy *common.SignaturePolicyEnvelope, tofs bool) (*SignedCDSPackage, []byte, *ChaincodeData, error) {
	env, err := ccpackage.OwnerCreateSignedCCDepSpec(cds, policy, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not create package %s", err)
	}

	b := protoutil.MarshalOrPanic(env)

	ccpack := &SignedCDSPackage{}
	cd, err := ccpack.InitFromBuffer(b)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error owner creating package %s", err)
	}

	if tofs {
		if err = ccpack.PutChaincodeToFS(); err != nil {
			return nil, nil, nil, fmt.Errorf("error putting package on the FS %s", err)
		}
	}

	return ccpack, b, cd, nil
}

func TestPutSigCDSCC(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, _, cd, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, true)
	if err != nil {
		t.Fatalf("cannot create package %s", err)
		return
	}

	if err = ccpack.ValidateCC(cd); err != nil {
		t.Fatalf("error validating package %s", err)
		return
	}
}

func TestPutSignedCDSErrorPaths(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
		Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, b, _, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, true)
	if err != nil {
		t.Fatalf("cannot create package %s", err)
		return
	}

	
	ccpack.buf = nil
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uninitialized package", "Unexpected error putting package on the FS")

	
	ccpack.buf = b
	id := ccpack.id
	ccpack.id = nil 
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil if buf is not nil", "Unexpected error putting package on the FS")

	assert.Panics(t, func() {
		ccpack.GetId()
	}, "GetId should have paniced if chaincode package ID is nil")

	
	ccpack.id = id
	id1 := ccpack.GetId()
	assert.Equal(t, id, id1)

	savDepSpec := ccpack.sDepSpec
	ccpack.sDepSpec = nil 
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "depspec cannot be nil if buf is not nil", "Unexpected error putting package on the FS")
	assert.Panics(t, func() {
		ccpack.GetInstantiationPolicy()
	}, "GetChaincodeData should have paniced if signed chaincode deployment spec is nil")
	assert.Panics(t, func() {
		ccpack.GetDepSpecBytes()
	}, "GetDepSpecBytes should have paniced if signed chaincode deployment spec is nil")
	ccpack.sDepSpec = savDepSpec 
	sdepspec1 := ccpack.GetInstantiationPolicy()
	assert.NotNil(t, sdepspec1)
	depspecBytes := ccpack.GetDepSpecBytes()
	assert.NotNil(t, depspecBytes)

	
	depSpec := ccpack.depSpec
	ccpack.depSpec = nil 
	assert.Panics(t, func() {
		ccpack.GetDepSpec()
	}, "GetDepSec should have paniced if chaincode deployment spec is nil")
	assert.Panics(t, func() {
		ccpack.GetChaincodeData()
	}, "GetChaincodeData should have paniced if chaincode deployment spec is nil")
	ccpack.depSpec = depSpec 
	depSpec1 := ccpack.GetDepSpec()
	assert.NotNil(t, depSpec1)

	env := ccpack.env
	ccpack.env = nil 
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "env cannot be nil if buf and depspec are not nil", "Unexpected error putting package on the FS")
	ccpack.env = env 
	env1 := ccpack.GetPackageObject()
	assert.Equal(t, env, env1)

	data := ccpack.data
	ccpack.data = nil 
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil data", "Unexpected error putting package on the FS")
	ccpack.data = data 

	datab := ccpack.datab
	ccpack.datab = nil 
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil data bytes", "Unexpected error putting package on the FS")
	ccpack.datab = datab 

	
	os.RemoveAll(ccdir)
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err, "Expected error putting package on the FS")
}

func TestGetCDSDataErrorPaths(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
		Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, _, _, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, true)
	if err != nil {
		t.Fatalf("cannot create package %s", err)
		return
	}

	
	assert.Panics(t, func() {
		_, _, _, err = ccpack.getCDSData(nil)
	}, "getCDSData should have paniced when called with nil signed chaincode deployment spec")

	
	scdp := &pb.SignedChaincodeDeploymentSpec{ChaincodeDeploymentSpec: []byte("bad spec")}
	_, _, _, err = ccpack.getCDSData(scdp)
	assert.Error(t, err)

	
	instPolicy := ccpack.sDepSpec.InstantiationPolicy
	ccpack.sDepSpec.InstantiationPolicy = nil
	_, _, _, err = ccpack.getCDSData(ccpack.sDepSpec)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instantiation policy cannot be nil for chaincode", "Unexpected error returned by getCDSData")
	ccpack.sDepSpec.InstantiationPolicy = instPolicy

	ccpack.sDepSpec.OwnerEndorsements = make([]*pb.Endorsement, 1)
	ccpack.sDepSpec.OwnerEndorsements[0] = &pb.Endorsement{}
	_, _, _, err = ccpack.getCDSData(ccpack.sDepSpec)
	assert.NoError(t, err)
}

func TestInitFromBufferErrorPaths(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
		Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, _, _, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, true)
	if err != nil {
		t.Fatalf("cannot create package %s", err)
		return
	}

	_, err = ccpack.InitFromBuffer([]byte("bad buffer"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal envelope from bytes", "Unexpected error returned by InitFromBuffer")
}

func TestValidateSignedCCErrorPaths(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
		Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, _, _, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, true)
	if err != nil {
		t.Fatalf("cannot create package %s", err)
		return
	}

	
	cd := &ChaincodeData{Name: "invalname", Version: "0"}
	err = ccpack.ValidateCC(cd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid chaincode data", "Unexpected error validating package")

	savDepSpec := ccpack.sDepSpec
	ccpack.sDepSpec = nil
	err = ccpack.ValidateCC(cd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uninitialized package", "Unexpected error validating package")
	ccpack.sDepSpec = savDepSpec

	cdspec := ccpack.sDepSpec.ChaincodeDeploymentSpec
	ccpack.sDepSpec.ChaincodeDeploymentSpec = nil
	err = ccpack.ValidateCC(cd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signed chaincode deployment spec cannot be nil in a package", "Unexpected error validating package")
	ccpack.sDepSpec.ChaincodeDeploymentSpec = cdspec

	depspec := ccpack.depSpec
	ccpack.depSpec = nil
	err = ccpack.ValidateCC(cd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chaincode deployment spec cannot be nil in a package", "Unexpected error validating package")
	ccpack.depSpec = depspec

	cd = &ChaincodeData{Name: "\027", Version: "0"}
	err = ccpack.ValidateCC(cd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `invalid chaincode name: "\x17"`)
}

func TestSigCDSGetCCPackage(t *testing.T) {
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	env, err := ccpackage.OwnerCreateSignedCCDepSpec(cds, &common.SignaturePolicyEnvelope{Version: 1}, nil)
	if err != nil {
		t.Fatalf("cannot create package")
		return
	}

	b := protoutil.MarshalOrPanic(env)

	ccpack, err := GetCCPackage(b)
	if err != nil {
		t.Fatalf("failed to get CCPackage %s", err)
		return
	}

	cccdspack, ok := ccpack.(*CDSPackage)
	if ok || cccdspack != nil {
		t.Fatalf("expected CDSPackage type cast to fail but succeeded")
		return
	}

	ccsignedcdspack, ok := ccpack.(*SignedCDSPackage)
	if !ok || ccsignedcdspack == nil {
		t.Fatalf("failed to get Signed CDS CCPackage")
		return
	}

	cds2 := ccsignedcdspack.GetDepSpec()
	if cds2 == nil {
		t.Fatalf("nil dep spec in Signed CDS CCPackage")
		return
	}

	if cds2.ChaincodeSpec.ChaincodeId.Name != cds.ChaincodeSpec.ChaincodeId.Name || cds2.ChaincodeSpec.ChaincodeId.Version != cds.ChaincodeSpec.ChaincodeId.Version {
		t.Fatalf("dep spec in Signed CDS CCPackage does not match %v != %v", cds, cds2)
		return
	}
}

func TestInvalidSigCDSGetCCPackage(t *testing.T) {
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	b := protoutil.MarshalOrPanic(cds)
	ccpack, err := GetCCPackage(b)
	if err != nil {
		t.Fatalf("failed to get CCPackage %s", err)
	}

	ccsignedcdspack, ok := ccpack.(*SignedCDSPackage)
	if ok || ccsignedcdspack != nil {
		t.Fatalf("expected failure to get Signed CDS CCPackage but succeeded")
	}
}


func TestSignedCDSSwitchChaincodes(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("badcode")}

	
	badccpack, _, _, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, true)
	if err != nil {
		t.Fatalf("error putting CDS to FS %s", err)
		return
	}

	
	cds.CodePackage = []byte("goodcode")

	
	_, _, goodcd, err := processSignedCDS(cds, &common.SignaturePolicyEnvelope{Version: 1}, false)
	if err != nil {
		t.Fatalf("error putting CDS to FS %s", err)
		return
	}

	if err = badccpack.ValidateCC(goodcd); err == nil {
		t.Fatalf("expected goodcd to fail against bad package but succeeded!")
		return
	}
}
