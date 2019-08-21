/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccprovider

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mcc-github/blockchain/bccsp/sw"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/stretchr/testify/assert"
)

func setupccdir() string {
	tempDir, err := ioutil.TempDir("/tmp", "ccprovidertest")
	if err != nil {
		panic(err)
	}
	SetChaincodesPath(tempDir)
	return tempDir
}

func processCDS(cds *pb.ChaincodeDeploymentSpec, tofs bool) (*CDSPackage, []byte, *ChaincodeData, error) {
	b := protoutil.MarshalOrPanic(cds)

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating bootBCCSP: %s", err)
	}
	ccpack := &CDSPackage{GetHasher: cryptoProvider}
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

func TestPutCDSCC(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, _, cd, err := processCDS(cds, true)
	if err != nil {
		t.Fatalf("error putting CDS to FS %s", err)
		return
	}

	if err = ccpack.ValidateCC(cd); err != nil {
		t.Fatalf("error validating package %s", err)
		return
	}
}

func TestPutCDSErrorPaths(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
		Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	ccpack, b, _, err := processCDS(cds, true)
	if err != nil {
		t.Fatalf("error putting CDS to FS %s", err)
		return
	}

	
	if err = ccpack.ValidateCC(&ChaincodeData{Name: "invalname", Version: "0"}); err == nil {
		t.Fatalf("expected error validating package")
		return
	}
	
	ccpack.buf = nil
	if err = ccpack.PutChaincodeToFS(); err == nil {
		t.Fatalf("expected error putting package on the FS")
		return
	}

	
	ccpack.buf = b
	savDepSpec := ccpack.depSpec
	ccpack.depSpec = nil
	if err = ccpack.PutChaincodeToFS(); err == nil {
		t.Fatalf("expected error putting package on the FS")
		return
	}

	
	ccpack.depSpec = savDepSpec

	
	os.RemoveAll(ccdir)
	if err = ccpack.PutChaincodeToFS(); err == nil {
		t.Fatalf("expected error putting package on the FS")
		return
	}
}

func TestCDSGetCCPackage(t *testing.T) {
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}

	b := protoutil.MarshalOrPanic(cds)

	ccpack, err := GetCCPackage(b)
	if err != nil {
		t.Fatalf("failed to get CDS CCPackage %s", err)
		return
	}

	cccdspack, ok := ccpack.(*CDSPackage)
	if !ok || cccdspack == nil {
		t.Fatalf("failed to get CDS CCPackage")
		return
	}

	cds2 := cccdspack.GetDepSpec()
	if cds2 == nil {
		t.Fatalf("nil dep spec in CDS CCPackage")
		return
	}

	if cds2.ChaincodeSpec.ChaincodeId.Name != cds.ChaincodeSpec.ChaincodeId.Name || cds2.ChaincodeSpec.ChaincodeId.Version != cds.ChaincodeSpec.ChaincodeId.Version {
		t.Fatalf("dep spec in CDS CCPackage does not match %v != %v", cds, cds2)
		return
	}
}


func TestCDSSwitchChaincodes(t *testing.T) {
	ccdir := setupccdir()
	defer os.RemoveAll(ccdir)

	
	cds := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("badcode")}

	
	badccpack, _, _, err := processCDS(cds, true)
	if err != nil {
		t.Fatalf("error putting CDS to FS %s", err)
		return
	}

	
	cds.CodePackage = []byte("goodcode")

	
	_, _, goodcd, err := processCDS(cds, false)
	if err != nil {
		t.Fatalf("error putting CDS to FS %s", err)
		return
	}

	if err = badccpack.ValidateCC(goodcd); err == nil {
		t.Fatalf("expected goodcd to fail against bad package but succeeded!")
		return
	}
}

func TestPutChaincodeToFSErrorPaths(t *testing.T) {
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	ccpack := &CDSPackage{GetHasher: cryptoProvider}
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uninitialized package", "Unexpected error returned")

	ccpack.buf = []byte("hello")
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id cannot be nil if buf is not nil", "Unexpected error returned")

	ccpack.id = []byte("cc123")
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "depspec cannot be nil if buf is not nil", "Unexpected error returned")

	ccpack.depSpec = &pb.ChaincodeDeploymentSpec{ChaincodeSpec: &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
		Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}, CodePackage: []byte("code")}
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil data", "Unexpected error returned")

	ccpack.data = &CDSData{}
	err = ccpack.PutChaincodeToFS()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil data bytes", "Unexpected error returned")
}

func TestValidateCCErrorPaths(t *testing.T) {
	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	assert.NoError(t, err)
	cpack := &CDSPackage{GetHasher: cryptoProvider}
	ccdata := &ChaincodeData{}
	err = cpack.ValidateCC(ccdata)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uninitialized package", "Unexpected error returned")

	cpack.depSpec = &pb.ChaincodeDeploymentSpec{
		CodePackage: []byte("code"),
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        1,
			ChaincodeId: &pb.ChaincodeID{Name: "testcc", Version: "0"},
			Input:       &pb.ChaincodeInput{Args: [][]byte{[]byte("")}},
		},
	}
	err = cpack.ValidateCC(ccdata)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil data", "Unexpected error returned")

	
	cpack = &CDSPackage{GetHasher: cryptoProvider}
	ccdata = &ChaincodeData{Name: "\027"}
	cpack.depSpec = &pb.ChaincodeDeploymentSpec{
		CodePackage: []byte("code"),
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: ccdata.Name, Version: "0"},
		},
	}
	cpack.data = &CDSData{}
	err = cpack.ValidateCC(ccdata)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `invalid chaincode name: "\x17"`)

	
	cpack = &CDSPackage{GetHasher: cryptoProvider}
	ccdata = &ChaincodeData{Name: "Tom"}
	cpack.depSpec = &pb.ChaincodeDeploymentSpec{
		CodePackage: []byte("code"),
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: "Jerry", Version: "0"},
		},
	}
	cpack.data = &CDSData{}
	err = cpack.ValidateCC(ccdata)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `invalid chaincode data name:"Tom"  (name:"Jerry" version:"0" )`)

	
	cpack = &CDSPackage{GetHasher: cryptoProvider}
	ccdata = &ChaincodeData{Name: "Tom", Version: "1"}
	cpack.depSpec = &pb.ChaincodeDeploymentSpec{
		CodePackage: []byte("code"),
		ChaincodeSpec: &pb.ChaincodeSpec{
			ChaincodeId: &pb.ChaincodeID{Name: ccdata.Name, Version: "0"},
		},
	}
	cpack.data = &CDSData{}
	err = cpack.ValidateCC(ccdata)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `invalid chaincode data name:"Tom" version:"1"  (name:"Tom" version:"0" )`)
}
