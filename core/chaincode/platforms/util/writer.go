/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/internal/ccmetadata"
	"github.com/pkg/errors"
)

var vmLogger = flogging.MustGetLogger("container")




func WriteFolderToTarPackage(tw *tar.Writer, srcPath string, excludeDirs []string, includeFileTypeMap map[string]bool, excludeFileTypeMap map[string]bool) error {
	fileCount := 0
	rootDirectory := srcPath

	
	if rootDirectory[len(rootDirectory)-1] == '/' {
		rootDirectory = rootDirectory[:len(rootDirectory)-1]
	}

	vmLogger.Debugf("rootDirectory = %s", rootDirectory)

	
	updatedExcludeDirs := make([]string, 0)
	for _, excludeDir := range excludeDirs {
		if excludeDir != "" && strings.LastIndex(excludeDir, "/") < len(excludeDir)-1 {
			excludeDir = excludeDir + "/"
			updatedExcludeDirs = append(updatedExcludeDirs, excludeDir)
		}
	}

	rootDirLen := len(rootDirectory)
	walkFn := func(localpath string, info os.FileInfo, err error) error {

		
		if strings.Contains(localpath, ".git") {
			return nil
		}

		if info.Mode().IsDir() {
			return nil
		}

		
		for _, excludeDir := range updatedExcludeDirs {
			if strings.Index(localpath, excludeDir) == rootDirLen+1 {
				return nil
			}
		}
		
		if len(localpath[rootDirLen:]) == 0 {
			return nil
		}
		ext := filepath.Ext(localpath)

		if includeFileTypeMap != nil {
			
			if _, ok := includeFileTypeMap[ext]; ok != true {
				return nil
			}
		}

		
		if excludeFileTypeMap != nil {
			if exclude, ok := excludeFileTypeMap[ext]; ok && exclude {
				return nil
			}
		}

		var packagepath string

		
		
		if strings.HasPrefix(localpath, filepath.Join(rootDirectory, "META-INF")) {
			packagepath = localpath[rootDirLen+1:]

			
			_, filename := filepath.Split(packagepath)

			
			
			if strings.HasPrefix(filename, ".") {
				vmLogger.Warningf("Ignoring hidden file in metadata directory: %s", packagepath)
				return nil
			}

			fileBytes, errRead := ioutil.ReadFile(localpath)
			if errRead != nil {
				return errRead
			}

			
			
			err = ccmetadata.ValidateMetadataFile(packagepath, fileBytes)
			if err != nil {
				return err
			}

		} else { 
			packagepath = fmt.Sprintf("src%s", localpath[rootDirLen:])
		}

		err = WriteFileToPackage(localpath, packagepath, tw)
		if err != nil {
			return fmt.Errorf("Error writing file to package: %s", err)
		}
		fileCount++

		return nil
	}

	if err := filepath.Walk(rootDirectory, walkFn); err != nil {
		vmLogger.Infof("Error walking rootDirectory: %s", err)
		return err
	}
	
	if fileCount == 0 {
		return errors.Errorf("no source files found in '%s'", srcPath)
	}
	return nil
}


func WriteFileToPackage(localpath string, packagepath string, tw *tar.Writer) error {
	vmLogger.Debug("Writing file to tarball:", packagepath)
	fd, err := os.Open(localpath)
	if err != nil {
		return fmt.Errorf("%s: %s", localpath, err)
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("%s: %s", localpath, err)
	}

	header, err := tar.FileInfoHeader(fi, localpath)
	if err != nil {
		return fmt.Errorf("failed calculating FileInfoHeader: %s", err)
	}

	
	var zeroTime time.Time
	header.AccessTime = zeroTime
	header.ModTime = zeroTime
	header.ChangeTime = zeroTime
	header.Name = packagepath
	header.Mode = 0100644
	header.Uid = 500
	header.Gid = 500

	err = tw.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("failed to write header for %s: %s", localpath, err)
	}

	is := bufio.NewReader(fd)
	_, err = io.Copy(tw, is)
	if err != nil {
		return fmt.Errorf("failed to write %s as %s: %s", localpath, packagepath, err)
	}

	return nil
}
