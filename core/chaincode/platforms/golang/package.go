/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package golang

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)

var includeFileTypes = map[string]bool{
	".c":    true,
	".h":    true,
	".s":    true,
	".go":   true,
	".yaml": true,
	".json": true,
}

var logger = flogging.MustGetLogger("chaincode.platform.golang")

func getCodeFromFS(path string) (codegopath string, err error) {
	logger.Debugf("getCodeFromFS %s", path)
	gopath, err := getGopath()
	if err != nil {
		return "", err
	}

	tmppath := filepath.Join(gopath, "src", path)
	if err := isCodeExist(tmppath); err != nil {
		return "", errors.Wrap(err, "code does not exist")
	}

	return gopath, nil
}


func isCodeExist(tmppath string) error {
	file, err := os.Open(tmppath)
	if err != nil {
		return errors.Wrap(err, "open failed")
	}

	fi, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, "stat failed")
	}

	if !fi.IsDir() {
		return errors.Errorf("%s is not a directory", file.Name())
	}

	return nil
}

type CodeDescriptor struct {
	Gopath, Pkg string
	Cleanup     func()
}






func getCode(path string) (*CodeDescriptor, error) {
	if path == "" {
		return nil, errors.New("cannot collect files from empty chaincode path")
	}

	
	var gopath string
	gopath, err := getCodeFromFS(path)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get code")
	}

	return &CodeDescriptor{Gopath: gopath, Pkg: path, Cleanup: nil}, nil
}

type SourceDescriptor struct {
	Name, Path string
	IsMetadata bool
	Info       os.FileInfo
}

type SourceMap map[string]SourceDescriptor

func findSource(gopath, pkg string) (SourceMap, error) {
	sources := make(SourceMap)
	tld := filepath.Join(gopath, "src", pkg)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			
			if path == tld {
				return nil
			}

			
			
			if isMetadataDir(path, tld) {
				logger.Debug("Files in META-INF directory will be included in code package tar:", path)
				return nil
			}

			
			logger.Debugf("skipping dir: %s", path)
			return filepath.SkipDir
		}

		ext := filepath.Ext(path)
		
		if _, ok := includeFileTypes[ext]; ok != true {
			return nil
		}

		name, err := filepath.Rel(gopath, path)
		if err != nil {
			return errors.Wrapf(err, "failed to calculate relative path for %s", path)
		}

		sources[name] = SourceDescriptor{Name: name, Path: path, IsMetadata: isMetadataDir(path, tld), Info: info}

		return nil
	}

	if err := filepath.Walk(tld, walkFn); err != nil {
		return nil, errors.Wrap(err, "walk failed")
	}

	return sources, nil
}


func isMetadataDir(path, tld string) bool {
	return strings.HasPrefix(path, filepath.Join(tld, "META-INF"))
}
