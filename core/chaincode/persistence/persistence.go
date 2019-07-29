/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mcc-github/blockchain/common/chaincode"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	persistence "github.com/mcc-github/blockchain/core/chaincode/persistence/intf"

	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("chaincode.persistence")



type IOReadWriter interface {
	ReadDir(string) ([]os.FileInfo, error)
	ReadFile(string) ([]byte, error)
	Remove(name string) error
	WriteFile(string, string, []byte) error
	MakeDir(string, os.FileMode) error
	Exists(path string) (bool, error)
}


type FilesystemIO struct {
}




func (f *FilesystemIO) WriteFile(path, name string, data []byte) error {
	if path == "" {
		return errors.New("empty path not allowed")
	}
	tmpFile, err := ioutil.TempFile(path, ".ccpackage.")
	if err != nil {
		return errors.Wrapf(err, "error creating temp file in directory '%s'", path)
	}
	defer os.Remove(tmpFile.Name())

	if n, err := tmpFile.Write(data); err != nil || n != len(data) {
		if err == nil {
			err = errors.Errorf(
				"failed to write the entire content of the file, expected %d, wrote %d",
				len(data), n)
		}
		return errors.Wrapf(err, "error writing to temp file '%s'", tmpFile.Name())
	}

	if err := tmpFile.Close(); err != nil {
		return errors.Wrapf(err, "error closing temp file '%s'", tmpFile.Name())
	}

	if err := os.Rename(tmpFile.Name(), filepath.Join(path, name)); err != nil {
		return errors.Wrapf(err, "error renaming temp file '%s'", tmpFile.Name())
	}

	return nil
}



func (f *FilesystemIO) Remove(name string) error {
	return os.Remove(name)
}


func (f *FilesystemIO) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}


func (f *FilesystemIO) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}



func (f *FilesystemIO) MakeDir(dirname string, mode os.FileMode) error {
	return os.MkdirAll(dirname, mode)
}


func (*FilesystemIO) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, errors.Wrapf(err, "could not determine whether file '%s' exists", path)
}


type Store struct {
	Path       string
	ReadWriter IOReadWriter
}



func NewStore(path string) *Store {
	store := &Store{
		Path:       path,
		ReadWriter: &FilesystemIO{},
	}
	store.Initialize()
	return store
}



func (s *Store) Initialize() {
	var (
		exists bool
		err    error
	)
	if exists, err = s.ReadWriter.Exists(s.Path); exists {
		return
	}
	if err != nil {
		panic(fmt.Sprintf("Initialization of chaincode store failed: %s", err))
	}
	if err = s.ReadWriter.MakeDir(s.Path, 0750); err != nil {
		panic(fmt.Sprintf("Could not create _lifecycle chaincodes install path: %s", err))
	}
}



func (s *Store) Save(label string, ccInstallPkg []byte) (persistence.PackageID, error) {
	hash := util.ComputeSHA256(ccInstallPkg)
	packageID := packageID(label, hash)

	ccInstallPkgFileName := packageID.String() + ".bin"
	ccInstallPkgFilePath := filepath.Join(s.Path, ccInstallPkgFileName)

	if exists, _ := s.ReadWriter.Exists(ccInstallPkgFilePath); exists {
		
		return packageID, nil
	}

	if err := s.ReadWriter.WriteFile(s.Path, ccInstallPkgFileName, ccInstallPkg); err != nil {
		err = errors.Wrapf(err, "error writing chaincode install package to %s", ccInstallPkgFilePath)
		logger.Error(err.Error())
		return "", err
	}

	return packageID, nil
}



func (s *Store) Load(packageID persistence.PackageID) ([]byte, error) {
	ccInstallPkgPath := filepath.Join(s.Path, packageID.String()+".bin")

	exists, err := s.ReadWriter.Exists(ccInstallPkgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not determine whether chaincode install package '%s' exists", packageID)
	}
	if !exists {
		return nil, &CodePackageNotFoundErr{
			PackageID: packageID,
		}
	}

	ccInstallPkg, err := s.ReadWriter.ReadFile(ccInstallPkgPath)
	if err != nil {
		err = errors.Wrapf(err, "error reading chaincode install package at %s", ccInstallPkgPath)
		return nil, err
	}

	return ccInstallPkg, nil
}



type CodePackageNotFoundErr struct {
	PackageID persistence.PackageID
}

func (e CodePackageNotFoundErr) Error() string {
	return fmt.Sprintf("chaincode install package '%s' not found", e.PackageID)
}



func (s *Store) ListInstalledChaincodes() ([]chaincode.InstalledChaincode, error) {
	files, err := s.ReadWriter.ReadDir(s.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading chaincode directory at %s", s.Path)
	}

	installedChaincodes := []chaincode.InstalledChaincode{}
	for _, file := range files {
		if instCC, isInstCC := installedChaincodeFromFilename(file.Name()); isInstCC {
			installedChaincodes = append(installedChaincodes, instCC)
		}
	}
	return installedChaincodes, nil
}



func (s *Store) GetChaincodeInstallPath() string {
	return s.Path
}

func packageID(label string, hash []byte) persistence.PackageID {
	return persistence.PackageID(fmt.Sprintf("%s:%x", label, hash))
}

var packageFileMatcher = regexp.MustCompile("^(.+):([0-9abcdef]+)[.]bin$")

func installedChaincodeFromFilename(fileName string) (chaincode.InstalledChaincode, bool) {
	matches := packageFileMatcher.FindStringSubmatch(fileName)
	if len(matches) == 3 {
		label := matches[1]
		hash, _ := hex.DecodeString(matches[2])
		packageID := packageID(label, hash)

		return chaincode.InstalledChaincode{
			Label:     label,
			Hash:      hash,
			PackageID: packageID,
		}, true
	}

	return chaincode.InstalledChaincode{}, false
}
