/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package testutil

import (
	"archive/tar"
	"bytes"
	"crypto/rand"
	"io"
	"os"
	"path/filepath"
	"testing"
)


func ConstructRandomBytes(t testing.TB, size int) []byte {
	value := make([]byte, size)
	_, err := rand.Read(value)
	if err != nil {
		t.Fatalf("Error while generating random bytes: %s", err)
	}
	return value
}


type TarFileEntry struct {
	Name, Body string
}


func CreateTarBytesForTest(testFiles []*TarFileEntry) []byte {
	
	buffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buffer)

	for _, file := range testFiles {
		tarHeader := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		err := tarWriter.WriteHeader(tarHeader)
		if err != nil {
			return nil
		}
		_, err = tarWriter.Write([]byte(file.Body))
		if err != nil {
			return nil
		}
	}
	
	tarWriter.Close()
	return buffer.Bytes()
}


func CopyDir(srcroot, destroot string, copyOnlySubdirs bool) error {
	if !copyOnlySubdirs {
		_, lastSegment := filepath.Split(srcroot)
		destroot = filepath.Join(destroot, lastSegment)
	}

	walkFunc := func(srcpath string, info os.FileInfo, err error) error {
		srcsubpath, err := filepath.Rel(srcroot, srcpath)
		if err != nil {
			return err
		}
		destpath := filepath.Join(destroot, srcsubpath)

		if info.IsDir() { 
			if err = os.MkdirAll(destpath, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		
		if err = copyFile(srcpath, destpath); err != nil {
			return err
		}
		return nil
	}

	return filepath.Walk(srcroot, walkFunc)
}

func copyFile(srcpath, destpath string) error {
	var srcFile, destFile *os.File
	var err error
	if srcFile, err = os.Open(srcpath); err != nil {
		return err
	}
	if destFile, err = os.Create(destpath); err != nil {
		return err
	}
	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}
	if err = srcFile.Close(); err != nil {
		return err
	}
	if err = destFile.Close(); err != nil {
		return err
	}
	return nil
}
