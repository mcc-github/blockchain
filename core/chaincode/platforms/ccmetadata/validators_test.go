/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccmetadata

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var packageTestDir = filepath.Join(os.TempDir(), "ccmetadata-validator-test")

func TestGoodIndexJSON(t *testing.T) {
	testDir := filepath.Join(packageTestDir, "GoodIndexJSON")
	cleanupDir(testDir)
	defer cleanupDir(testDir)

	fileName := "META-INF/statedb/couchdb/indexes/myIndex.json"
	fileBytes := []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err := ValidateMetadataFile(fileName, fileBytes)
	assert.NoError(t, err, "Error validating a good index")
}

func TestBadIndexJSON(t *testing.T) {
	testDir := filepath.Join(packageTestDir, "BadIndexJSON")
	cleanupDir(testDir)
	defer cleanupDir(testDir)

	fileName := "META-INF/statedb/couchdb/indexes/myIndex.json"
	fileBytes := []byte("invalid json")

	err := ValidateMetadataFile(fileName, fileBytes)

	assert.Error(t, err, "Should have received an InvalidIndexContentError")

	
	_, ok := err.(*InvalidIndexContentError)
	assert.True(t, ok, "Should have received an InvalidIndexContentError")

	t.Log("SAMPLE ERROR STRING:", err.Error())
}

func TestIndexWrongLocation(t *testing.T) {
	testDir := filepath.Join(packageTestDir, "IndexWrongLocation")
	cleanupDir(testDir)
	defer cleanupDir(testDir)

	fileName := "META-INF/statedb/couchdb/myIndex.json"
	fileBytes := []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err := ValidateMetadataFile(fileName, fileBytes)
	assert.Error(t, err, "Should have received an UnhandledDirectoryError")

	
	_, ok := err.(*UnhandledDirectoryError)
	assert.True(t, ok, "Should have received an UnhandledDirectoryError")

	t.Log("SAMPLE ERROR STRING:", err.Error())
}

func TestInvalidMetadataType(t *testing.T) {
	testDir := filepath.Join(packageTestDir, "InvalidMetadataType")
	cleanupDir(testDir)
	defer cleanupDir(testDir)

	fileName := "myIndex.json"
	fileBytes := []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err := ValidateMetadataFile(fileName, fileBytes)
	assert.Error(t, err, "Should have received an UnhandledDirectoryError")

	
	_, ok := err.(*UnhandledDirectoryError)
	assert.True(t, ok, "Should have received an UnhandledDirectoryError")
}

func TestBadMetadataExtension(t *testing.T) {
	testDir := filepath.Join(packageTestDir, "BadMetadataExtension")
	cleanupDir(testDir)
	defer cleanupDir(testDir)

	fileName := "myIndex.go"
	fileBytes := []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err := ValidateMetadataFile(fileName, fileBytes)
	assert.Error(t, err, "Should have received an error")

}

func TestBadFilePaths(t *testing.T) {
	testDir := filepath.Join(packageTestDir, "BadMetadataExtension")
	cleanupDir(testDir)
	defer cleanupDir(testDir)

	
	fileName := "META-INF1/statedb/couchdb/indexes/test1.json"
	fileBytes := []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err := ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for bad META-INF directory")

	
	fileName = "META-INF/statedb/test1.json"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for bad length")

	
	fileName = "META-INF/statedb/goleveldb/indexes/test1.json"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for invalid database")

	
	fileName = "META-INF/statedb/couchdb/index/test1.json"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for invalid indexes directory")

	
	fileName = "META-INF/statedb/couchdb/collection/testcoll/indexes/test1.json"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for invalid collections directory")

	
	fileName = "META-INF/statedb/couchdb/collections/testcoll/indexes/test1.json"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.NoError(t, err, "Error should not have been thrown for a valid collection name")

	
	fileName = "META-INF/statedb/couchdb/collections/#testcoll/indexes/test1.json"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for an invalid collection name")

	
	fileName = "META-INF/statedb/couchdb/collections/testcoll/indexes/test1.txt"
	fileBytes = []byte(`{"index":{"fields":["data.docType","data.owner"]},"name":"indexOwner","type":"json"}`)

	err = ValidateMetadataFile(fileName, fileBytes)
	fmt.Println(err)
	assert.Error(t, err, "Should have received an error for an invalid file name")

}

func TestIndexValidation(t *testing.T) {

	
	indexDef := []byte(`{"index":{"fields":[{"size":"desc"}, {"color":"desc"}]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition := isJSON(indexDef)
	err := validateIndexJSON(indexDefinition)
	assert.NoError(t, err)

	
	indexDef = []byte(`{"index":{"fields":["size","color"]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.NoError(t, err)

	
	indexDef = []byte(`{"index":{"fields":["size","color"]}}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.NoError(t, err)

	
	indexDef = []byte(`{
		  "index": {
		    "partial_filter_selector": {
		      "status": {
		        "$ne": "archived"
		      }
		    },
		    "fields": ["type"]
		  },
		  "ddoc" : "type-not-archived",
		  "type" : "json"
		}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.NoError(t, err)

}

func TestIndexValidationInvalidParameters(t *testing.T) {

	
	indexDef := []byte(`{"index":{"fields":[{"size":"desc"}, {"color":"desc"}]},"ddoc":1, "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition := isJSON(indexDef)
	err := validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for numeric design doc")

	
	indexDef = []byte(`{"index":{"fields":["size","color"]},"ddoc1":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid design doc parameter")

	
	indexDef = []byte(`{"index":{"fields":["size","color"]},"ddoc":"indexSizeSortName", "name1":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid name parameter")

	
	indexDef = []byte(`{"index":{"fields":["size","color"]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":1}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for numeric type parameter")

	
	indexDef = []byte(`{"index":{"fields":["size","color"]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"text"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid type parameter")

	
	indexDef = []byte(`{"index1":{"fields":["size","color"]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid index parameter")

	
	indexDef = []byte(`{"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for missing index parameter")

}

func TestIndexValidationInvalidFields(t *testing.T) {

	
	indexDef := []byte(`{"index":{"fields1":[{"size":"desc"}, {"color":"desc"}]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition := isJSON(indexDef)
	err := validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid fields parameter")

	
	indexDef = []byte(`{"index":{"fields":["size", 1]},"ddoc1":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for field name defined as numeric")

	
	indexDef = []byte(`{"index":{"fields":[{"size":"desc1"}, {"color":"desc"}]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid field sort")

	
	indexDef = []byte(`{"index":{"fields":[{"size":1}, {"color":"desc"}]},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for a numeric in field sort")

	
	indexDef = []byte(`{"index":{"fields":"size"},"ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for invalid field json")

	
	indexDef = []byte(`{"index":"fields","ddoc":"indexSizeSortName", "name":"indexSizeSortName","type":"json"}`)
	_, indexDefinition = isJSON(indexDef)
	err = validateIndexJSON(indexDefinition)
	assert.Error(t, err, "Error should have been thrown for missing JSON for fields")

}

func cleanupDir(dir string) error {
	
	err := os.RemoveAll(dir)
	if err != nil {
		return nil
	}
	return os.Mkdir(dir, os.ModePerm)
}

func writeToFile(filename string, bytes []byte) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}
