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

package ccprovider

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/bccsp/sw"
	"github.com/mcc-github/blockchain/protos/peer"
	"github.com/stretchr/testify/assert"
)

func getDepSpec(name string, path string, version string, initArgs [][]byte) (*peer.ChaincodeDeploymentSpec, error) {
	spec := &peer.ChaincodeSpec{Type: 1, ChaincodeId: &peer.ChaincodeID{Name: name, Path: path, Version: version}, Input: &peer.ChaincodeInput{Args: initArgs}}

	codePackageBytes := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(codePackageBytes)
	tw := tar.NewWriter(gz)

	payload := []byte(name + path + version)
	err := tw.WriteHeader(&tar.Header{
		Name: "src/garbage.go",
		Size: int64(len(payload)),
		Mode: 0100644,
	})
	if err != nil {
		return nil, err
	}

	_, err = tw.Write(payload)
	if err != nil {
		return nil, err
	}

	tw.Close()
	gz.Close()

	return &peer.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes.Bytes()}, nil
}

func buildPackage(name string, path string, version string, initArgs [][]byte) (CCPackage, error) {
	depSpec, err := getDepSpec(name, path, version, initArgs)
	if err != nil {
		return nil, err
	}

	buf, err := proto.Marshal(depSpec)
	if err != nil {
		return nil, err
	}

	cryptoProvider, err := sw.NewDefaultSecurityLevelWithKeystore(sw.NewDummyKeyStore())
	if err != nil {
		return nil, err
	}
	cccdspack := &CDSPackage{GetHasher: cryptoProvider}
	if _, err := cccdspack.InitFromBuffer(buf); err != nil {
		return nil, err
	}

	return cccdspack, nil
}

type mockCCInfoFSStorageMgrImpl struct {
	CCMap map[string]CCPackage
}

func (m *mockCCInfoFSStorageMgrImpl) GetChaincode(ccNameVersion string) (CCPackage, error) {
	return m.CCMap[ccNameVersion], nil
}


func TestCCInfoCache(t *testing.T) {
	ccinfoFs := &mockCCInfoFSStorageMgrImpl{CCMap: map[string]CCPackage{}}
	cccache := NewCCInfoCache(ccinfoFs)

	

	
	_, err := cccache.GetChaincodeData("foo:1.0")
	assert.Error(t, err)

	
	pack, err := buildPackage("foo", "mychaincode", "1.0", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)
	ccinfoFs.CCMap["foo:1.0"] = pack

	
	cd1, err := cccache.GetChaincodeData("foo:1.0")
	assert.NoError(t, err)

	
	cd2, err := cccache.GetChaincodeData("foo:1.0")
	assert.NoError(t, err)

	
	assert.NotNil(t, cd1)
	assert.NotNil(t, cd2)

	
	pack, err = buildPackage("foo", "mychaincode", "2.0", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)
	ccinfoFs.CCMap["foo:2.0"] = pack

	
	_, err = getDepSpec("foo", "mychaincode", "2.0", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	
	cd1, err = cccache.GetChaincodeData("foo:2.0")
	assert.NoError(t, err)

	
	cd2, err = cccache.GetChaincodeData("foo:2.0")
	assert.NoError(t, err)

	
	assert.NotNil(t, cd1)
	assert.NotNil(t, cd2)
}

func TestPutChaincode(t *testing.T) {
	ccname := ""
	ccver := "1.0"
	ccpath := "mychaincode"

	ccinfoFs := &mockCCInfoFSStorageMgrImpl{CCMap: map[string]CCPackage{}}
	NewCCInfoCache(ccinfoFs)

	
	
	_, err := getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	
	ccname = "foo"
	ccver = ""
	_, err = getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	
	ccinfoFs = &mockCCInfoFSStorageMgrImpl{CCMap: map[string]CCPackage{}}
	NewCCInfoCache(ccinfoFs)

	ccname = "foo"
	ccver = "1.0"
	_, err = getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)
}


func TestCCInfoFSPeerInstance(t *testing.T) {
	ccname := "bar"
	ccver := "1.0"
	ccpath := "mychaincode"

	
	_, err := GetChaincodeFromFS("bar:1.0")
	assert.Error(t, err)

	
	ds, err := getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	
	err = PutChaincodeIntoFS(ds)
	assert.NoError(t, err)

	
	resp, err := GetInstalledChaincodes()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, len(resp.Chaincodes), "GetInstalledChaincodes should not have returned 0 chaincodes")

	
	_, err = GetChaincodeData("bar:1.0")
	assert.NoError(t, err)
}

func TestGetInstalledChaincodesErrorPaths(t *testing.T) {
	
	
	cip := chaincodeInstallPath
	defer SetChaincodesPath(cip)

	
	dir, err := ioutil.TempDir(os.TempDir(), "chaincodes")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	
	SetChaincodesPath(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "idontexist.1.0"), []byte("test"), 0777)
	assert.NoError(t, err)
	resp, err := GetInstalledChaincodes()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp.Chaincodes),
		"Expected 0 chaincodes but GetInstalledChaincodes returned %s chaincodes", len(resp.Chaincodes))
}

func TestChaincodePackageExists(t *testing.T) {
	_, err := ChaincodePackageExists("foo1", "1.0")
	assert.Error(t, err)
}

func TestSetChaincodesPath(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "setchaincodes")
	if err != nil {
		assert.Fail(t, err.Error(), "Unable to create temp dir")
	}
	defer os.RemoveAll(dir)
	t.Logf("created temp dir %s", dir)

	
	
	cip := chaincodeInstallPath
	defer SetChaincodesPath(cip)

	f, err := ioutil.TempFile(dir, "chaincodes")
	assert.NoError(t, err)
	assert.Panics(t, func() {
		SetChaincodesPath(f.Name())
	}, "SetChaincodesPath should have paniced if a file is passed to it")

	
	
	
	
	
	
	
	

	
	
	
	
	
	
}

var ccinfocachetestpath = "/tmp/ccinfocachetest"

func TestMain(m *testing.M) {
	os.RemoveAll(ccinfocachetestpath)

	SetChaincodesPath(ccinfocachetestpath)
	rc := m.Run()
	os.RemoveAll(ccinfocachetestpath)
	os.Exit(rc)
}
