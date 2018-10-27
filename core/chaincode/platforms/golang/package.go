/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package golang

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	ccutil "github.com/mcc-github/blockchain/core/chaincode/platforms/util"
)

var includeFileTypes = map[string]bool{
	".c":    true,
	".h":    true,
	".s":    true,
	".go":   true,
	".yaml": true,
	".json": true,
}

var logger = flogging.MustGetLogger("golang-platform")

func getCodeFromFS(path string) (codegopath string, err error) {
	logger.Debugf("getCodeFromFS %s", path)
	gopath, err := getGopath()
	if err != nil {
		return "", err
	}

	tmppath := filepath.Join(gopath, "src", path)
	if err := ccutil.IsCodeExist(tmppath); err != nil {
		return "", fmt.Errorf("code does not exist %s", err)
	}

	return gopath, nil
}

type CodeDescriptor struct {
	Gopath, Pkg string
	Cleanup     func()
}






func getCode(path string) (*CodeDescriptor, error) {
	if path == "" {
		return nil, errors.New("Cannot collect files from empty chaincode path")
	}

	
	var gopath string
	gopath, err := getCodeFromFS(path)
	if err != nil {
		return nil, fmt.Errorf("Error getting code %s", err)
	}

	return &CodeDescriptor{Gopath: gopath, Pkg: path, Cleanup: nil}, nil
}

type SourceDescriptor struct {
	Name, Path string
	IsMetadata bool
	Info       os.FileInfo
}
type SourceMap map[string]SourceDescriptor

type Sources []SourceDescriptor

func (s Sources) Len() int {
	return len(s)
}

func (s Sources) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sources) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) < 0
}

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
			return fmt.Errorf("error obtaining relative path for %s: %s", path, err)
		}

		sources[name] = SourceDescriptor{Name: name, Path: path, IsMetadata: isMetadataDir(path, tld), Info: info}

		return nil
	}

	if err := filepath.Walk(tld, walkFn); err != nil {
		return nil, fmt.Errorf("Error walking directory: %s", err)
	}

	return sources, nil
}


func isMetadataDir(path, tld string) bool {
	return strings.HasPrefix(path, filepath.Join(tld, "META-INF"))
}
