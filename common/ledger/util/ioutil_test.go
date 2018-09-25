/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dbPathTest = "/tmp/ledgertests/common/ledger/util"
var dbFileTest = dbPathTest + "/testFile"

func TestCreatingDBDirWithPathSeperator(t *testing.T) {

	
	dbPathTestWSeparator := dbPathTest + "/"
	cleanup(dbPathTestWSeparator) 
	defer cleanup(dbPathTestWSeparator)

	dirEmpty, err := CreateDirIfMissing(dbPathTestWSeparator)
	assert.NoError(t, err, "Error when trying to create a test db directory at [%s]", dbPathTestWSeparator)
	assert.True(t, dirEmpty) 
}

func TestCreatingDBDirWhenDirDoesAndDoesNotExists(t *testing.T) {

	cleanup(dbPathTest) 
	defer cleanup(dbPathTest)

	
	dirEmpty, err := CreateDirIfMissing(dbPathTest)
	assert.NoError(t, err, "Error when trying to create a test db directory at [%s]", dbPathTest)
	assert.True(t, dirEmpty)

	
	dirEmpty2, err2 := CreateDirIfMissing(dbPathTest)
	assert.NoError(t, err2, "Error not handling existing directory when trying to create a test db directory at [%s]", dbPathTest)
	assert.True(t, dirEmpty2)
}

func TestDirNotEmptyAndFileExists(t *testing.T) {

	cleanup(dbPathTest)
	defer cleanup(dbPathTest)

	
	dirEmpty, err := CreateDirIfMissing(dbPathTest)
	assert.NoError(t, err, "Error when trying to create a test db directory at [%s]", dbPathTest)
	assert.True(t, dirEmpty)

	
	exists2, size2, err2 := FileExists(dbFileTest)
	assert.NoError(t, err2, "Error when trying to determine if file exist when it does not at [%s]", dbFileTest)
	assert.Equal(t, int64(0), size2)
	assert.False(t, exists2) 

	
	testStr := "This is some test data in a file"
	sizeOfFileCreated, err3 := createAndWriteAFile(testStr)
	assert.NoError(t, err3, "Error when trying to create and write to file at [%s]", dbFileTest)
	assert.Equal(t, len(testStr), sizeOfFileCreated) 

	
	exists, size, err4 := FileExists(dbFileTest)
	assert.NoError(t, err4, "Error when trying to determine if file exist at [%s]", dbFileTest)
	assert.Equal(t, int64(sizeOfFileCreated), size)
	assert.True(t, exists) 

	
	dirEmpty5, err5 := DirEmpty(dbPathTest)
	assert.NoError(t, err5, "Error when detecting if empty at db directory [%s]", dbPathTest)
	assert.False(t, dirEmpty5) 
}

func TestListSubdirs(t *testing.T) {
	childFolders := []string{".childFolder1", "childFolder2", "childFolder3"}
	cleanup(dbPathTest)
	defer cleanup(dbPathTest)
	for _, folder := range childFolders {
		assert.NoError(t, os.MkdirAll(filepath.Join(dbPathTest, folder), 0755))
	}
	subFolders, err := ListSubdirs(dbPathTest)
	assert.NoError(t, err)
	assert.Equal(t, subFolders, childFolders)
}

func createAndWriteAFile(sentence string) (int, error) {
	
	f, err2 := os.Create(dbFileTest)
	if err2 != nil {
		return 0, err2
	}
	defer f.Close()

	
	return f.WriteString(sentence)
}

func cleanup(path string) {
	os.RemoveAll(path)
}
