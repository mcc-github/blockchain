/*
Copyright London Stock Exchange 2016 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package util

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WriteFileToPackage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "utiltest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)
	err = WriteFileToPackage("blah", "", tw)
	assert.Error(t, err, "Expected error writing non existent file to package")

	
	filename := "test.txt"
	filecontent := "hello"
	filePath := filepath.Join(tempDir, filename)
	err = ioutil.WriteFile(filePath, bytes.NewBufferString(filecontent).Bytes(), 0600)
	require.NoError(t, err, "Error creating file %s", filePath)

	err = WriteFileToPackage(filePath, filename, tw)
	assert.NoError(t, err, "Error returned by WriteFileToPackage while writing existing file")
	tw.Close()
	gw.Close()

	
	
	err = WriteFileToPackage(filePath, "", tw)
	assert.Error(t, err, "Expected error writing using a closed writer")

	
	r := bytes.NewReader(buf.Bytes())
	gr, err := gzip.NewReader(r)
	require.NoError(t, err, "Error creating a gzip reader")
	defer gr.Close()

	tr := tar.NewReader(gr)
	header, err := tr.Next()
	require.NoError(t, err, "Error getting the file from the tar")
	assert.Equal(t, filename, header.Name, "filename read from archive does not match what was added")

	b := make([]byte, 5)
	n, err := tr.Read(b)
	assert.Equal(t, 5, n)
	assert.True(t, err == nil || err == io.EOF, "Error reading file from the archive") 
	assert.Equal(t, filecontent, string(b), "file content from archive does not equal original content")
}

func Test_WriteStreamToPackage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "utiltest")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tarw := tar.NewWriter(bytes.NewBuffer(nil))
	input := bytes.NewReader([]byte("hello"))

	
	err = WriteStreamToPackage(nil, "nonexistentpath", "", tarw)
	assert.Error(t, err, "Expected error getting info of non existent file")

	
	err = WriteStreamToPackage(input, tempDir, "", tarw)
	assert.Error(t, err, "Expected error copying to the tar writer (tarw)")

	tarw.Close()

	
	err = WriteStreamToPackage(input, tempDir, "", tarw)
	assert.Error(t, err, "Expected error copying to closed tar writer (tarw)")

	
	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	
	filename := "test.txt"
	filecontent := "hello"
	filePath := filepath.Join(tempDir, filename)
	err = ioutil.WriteFile(filePath, bytes.NewBufferString(filecontent).Bytes(), 0600)
	assert.NoError(t, err, "Error creating file %s", filePath)

	
	var b []byte
	b, err = ioutil.ReadFile(filePath)
	assert.NoError(t, err, "Error reading file %s", filePath)
	is := bytes.NewReader(b)

	
	err = WriteStreamToPackage(is, filePath, filename, tw)
	assert.NoError(t, err, "Error copying file to the tar writer (tw)")

	
	tw.Close()
	gw.Close()

	
	br := bytes.NewReader(buf.Bytes())
	gr, err := gzip.NewReader(br)
	defer gr.Close()
	assert.NoError(t, err, "Error creating a gzip reader")
	tr := tar.NewReader(gr)
	header, err := tr.Next()
	assert.NoError(t, err, "Error getting the file from the tar")
	assert.Equal(t, filename, header.Name, "filename read from archive does not match what was added")

	b1 := make([]byte, 5)
	n, err := tr.Read(b1)
	assert.Equal(t, 5, n)
	assert.True(t, err == nil || err == io.EOF, "Error reading file from the archive") 
	assert.Equal(t, filecontent, string(b), "file content from archive does not equal original content")
}


func Test_WriteFolderToTarPackage1(t *testing.T) {
	srcPath := filepath.Join("testdata", "sourcefiles")
	filePath := "src/src/Hello.java"
	includeFileTypes := map[string]bool{
		".java": true,
	}
	excludeFileTypes := map[string]bool{}

	tarBytes := createTestTar(t, srcPath, []string{}, includeFileTypes, excludeFileTypes)

	
	entries := tarContents(t, tarBytes)
	assert.ElementsMatch(t, []string{filePath}, entries, "archive should only contain one file")
}


func Test_WriteFolderToTarPackage2(t *testing.T) {
	srcPath := filepath.Join("testdata", "sourcefiles")
	tarBytes := createTestTar(t, srcPath, []string{"src"}, nil, nil)

	entries := tarContents(t, tarBytes)
	assert.ElementsMatch(t, []string{"src/artifact.xml", "META-INF/statedb/couchdb/indexes/indexOwner.json"}, entries)
}


func Test_WriteFolderToTarPackage3(t *testing.T) {
	srcPath := filepath.Join("testdata", "sourcefiles")
	filePath := "META-INF/statedb/couchdb/indexes/indexOwner.json"

	tarBytes := createTestTar(t, srcPath, []string{}, nil, nil)

	
	entries := tarContents(t, tarBytes)
	assert.Contains(t, entries, filePath, "should have found statedb index artifact in META-INF directory")
}


func Test_WriteFolderToTarPackage4(t *testing.T) {
	srcPath := filepath.Join("testdata", "sourcefiles") + string(filepath.Separator)
	filePath := "META-INF/statedb/couchdb/indexes/indexOwner.json"

	tarBytes := createTestTar(t, srcPath, []string{}, nil, nil)

	
	entries := tarContents(t, tarBytes)
	assert.Contains(t, entries, filePath, "should have found statedb index artifact in META-INF directory")
}


func Test_WriteFolderToTarPackage5(t *testing.T) {
	srcPath := filepath.Join("testdata", "sourcefiles")
	filePath := "META-INF/.hiddenfile"

	assert.FileExists(t, filepath.Join(srcPath, "META-INF", ".hiddenfile"))

	tarBytes := createTestTar(t, srcPath, []string{}, nil, nil)

	
	entries := tarContents(t, tarBytes)
	assert.NotContains(t, entries, filePath, "should not contain .hiddenfile in META-INF directory")
}


func Test_WriteFolderToTarPackageFailure1(t *testing.T) {
	srcPath, err := ioutil.TempDir("", "utiltest")
	require.NoError(t, err)
	defer os.RemoveAll(srcPath)

	tw := tar.NewWriter(bytes.NewBuffer(nil))
	defer tw.Close()

	err = WriteFolderToTarPackage(tw, srcPath, []string{}, nil, nil)
	assert.Contains(t, err.Error(), "no source files found")
}


func Test_WriteFolderToTarPackageFailure2(t *testing.T) {
	srcPath := filepath.Join("testdata", "BadMetadataInvalidIndex")
	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	err := WriteFolderToTarPackage(tw, srcPath, []string{}, nil, nil)
	assert.Error(t, err, "Should have received error writing folder to package")
	assert.Contains(t, err.Error(), "Index metadata file [META-INF/statedb/couchdb/indexes/bad.json] is not a valid JSON")

	tw.Close()
	gw.Close()
}


func Test_WriteFolderToTarPackageFailure3(t *testing.T) {
	srcPath := filepath.Join("testdata", "BadMetadataUnexpectedFolderContent")
	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	err := WriteFolderToTarPackage(tw, srcPath, []string{}, nil, nil)
	assert.Error(t, err, "Should have received error writing folder to package")
	assert.Contains(t, err.Error(), "metadata file path must begin with META-INF/statedb")

	tw.Close()
	gw.Close()
}

func Test_WriteJavaProjectToPackage(t *testing.T) {
	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tw := tar.NewWriter(gw)

	srcPath := filepath.Join("testdata", "sourcefiles")
	assert.FileExists(t, filepath.Join(srcPath, "src", "Hello.class"))

	err := WriteJavaProjectToPackage(tw, srcPath)
	assert.NoError(t, err, "Error writing java project to package")

	
	
	tw.Close()
	gw.Close()

	entries := tarContents(t, inputbuf.Bytes())
	assert.Contains(t, entries, "META-INF/statedb/couchdb/indexes/indexOwner.json")
	assert.Contains(t, entries, "src/artifact.xml")
	assert.Contains(t, entries, "src/src/Hello.java")
	assert.NotContains(t, entries, "src/src/Hello.class")

	err = WriteJavaProjectToPackage(tw, srcPath)
	assert.Error(t, err, "WriteJavaProjectToPackage was called with closed writer, should have failed")
}

func Test_WriteBytesToPackage(t *testing.T) {
	inputbuf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(inputbuf)
	defer tw.Close()
	err := WriteBytesToPackage("foo", []byte("blah"), tw)
	assert.NoError(t, err, "Error writing bytes to package")
}

func createTestTar(t *testing.T, srcPath string, excludeDir []string, includeFileTypeMap map[string]bool, excludeFileTypeMap map[string]bool) []byte {
	buf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	err := WriteFolderToTarPackage(tw, srcPath, excludeDir, includeFileTypeMap, excludeFileTypeMap)
	assert.NoError(t, err, "Error writing folder to package")

	tw.Close()
	gw.Close()
	return buf.Bytes()
}

func tarContents(t *testing.T, buf []byte) []string {
	br := bytes.NewReader(buf)
	gr, err := gzip.NewReader(br)
	require.NoError(t, err)
	defer gr.Close()

	tr := tar.NewReader(gr)

	var entries []string
	for {
		header, err := tr.Next()
		if err == io.EOF { 
			break
		}
		require.NoError(t, err, "failed to get next entry")
		entries = append(entries, header.Name)
	}

	return entries
}
