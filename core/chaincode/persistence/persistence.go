/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package persistence

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("chaincode/persistence")



type IOReadWriter interface {
	ReadFile(string) ([]byte, error)
	Remove(name string) error
	Stat(string) (os.FileInfo, error)
	WriteFile(string, []byte, os.FileMode) error
}


type FilesystemIO struct {
}


func (f *FilesystemIO) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}


func (f *FilesystemIO) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}



func (f *FilesystemIO) Remove(name string) error {
	return os.Remove(name)
}


func (f *FilesystemIO) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}


type Store struct {
	Path       string
	ReadWriter IOReadWriter
}



func (s *Store) Save(name, version string, ccInstallPkg []byte) error {
	metadataJSON, err := toJSON(name, version)
	if err != nil {
		return err
	}

	hashString := hex.EncodeToString(util.ComputeSHA256(ccInstallPkg))
	metadataPath := filepath.Join(s.Path, hashString+".json")
	if _, err := s.ReadWriter.Stat(metadataPath); err == nil {
		return errors.Errorf("chaincode metadata already exists at %s", metadataPath)
	}

	ccInstallPkgPath := filepath.Join(s.Path, hashString+".bin")
	if _, err := s.ReadWriter.Stat(ccInstallPkgPath); err == nil {
		return errors.Errorf("ChaincodeInstallPackage already exists at %s", ccInstallPkgPath)
	}

	if err := s.ReadWriter.WriteFile(metadataPath, metadataJSON, 0600); err != nil {
		return errors.Wrapf(err, "error writing metadata file to %s", metadataPath)
	}

	if err := s.ReadWriter.WriteFile(ccInstallPkgPath, ccInstallPkg, 0600); err != nil {
		err = errors.Wrapf(err, "error writing chaincode install package to %s", ccInstallPkgPath)
		logger.Error(err.Error())

		
		if err2 := s.ReadWriter.Remove(metadataPath); err2 != nil {
			logger.Errorf("error removing metadata file at %s: %s", metadataPath, err2)
		}
		return err
	}

	return nil
}



func (s *Store) Load(hash []byte) (ccInstallPkg []byte, name, version string, err error) {
	hashString := hex.EncodeToString(hash)
	ccInstallPkgPath := filepath.Join(s.Path, hashString+".bin")
	ccInstallPkg, err = s.ReadWriter.ReadFile(ccInstallPkgPath)
	if err != nil {
		err = errors.Wrapf(err, "error reading chaincode install package at %s", ccInstallPkgPath)
		return nil, "", "", err
	}

	metadataPath := filepath.Join(s.Path, hashString+".json")
	metadataBytes, err := s.ReadWriter.ReadFile(metadataPath)
	if err != nil {
		err = errors.Wrapf(err, "error reading metadata at %s", metadataPath)
		return nil, "", "", err
	}
	ccMetadata := &ChaincodeMetadata{}
	err = json.Unmarshal(metadataBytes, ccMetadata)
	if err != nil {
		err = errors.Wrapf(err, "error unmarshaling metadata at %s", metadataPath)
		return nil, "", "", err
	}

	name = ccMetadata.Name
	version = ccMetadata.Version

	return ccInstallPkg, name, version, nil
}


type ChaincodeMetadata struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

func toJSON(name, version string) ([]byte, error) {
	metadata := &ChaincodeMetadata{
		Name:    name,
		Version: version,
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling name and version into JSON")
	}

	return metadataBytes, nil
}
