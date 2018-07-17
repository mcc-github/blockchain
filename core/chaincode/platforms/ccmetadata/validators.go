/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ccmetadata

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)


type fileValidator func(fileName string, fileBytes []byte) error


const AllowedCharsCollectionName = "[A-Za-z0-9_-]+"


var fileValidators = map[*regexp.Regexp]fileValidator{
	regexp.MustCompile("^META-INF/statedb/couchdb/indexes/.*[.]json"):                                                couchdbIndexFileValidator,
	regexp.MustCompile("^META-INF/statedb/couchdb/collections/" + AllowedCharsCollectionName + "/indexes/.*[.]json"): couchdbIndexFileValidator,
}

var collectionNameValid = regexp.MustCompile("^" + AllowedCharsCollectionName)

var fileNameValid = regexp.MustCompile("^.*[.]json")

var validDatabases = []string{"couchdb"}


type UnhandledDirectoryError struct {
	err string
}

func (e *UnhandledDirectoryError) Error() string {
	return e.err
}


type InvalidIndexContentError struct {
	err string
}

func (e *InvalidIndexContentError) Error() string {
	return e.err
}



func ValidateMetadataFile(filePathName string, fileBytes []byte) error {
	
	fileValidator := selectFileValidator(filePathName)

	
	if fileValidator == nil {
		return &UnhandledDirectoryError{buildMetadataFileErrorMessage(filePathName)}
	}

	
	err := fileValidator(filePathName, fileBytes)
	if err != nil {
		return err
	}

	
	return nil
}

func buildMetadataFileErrorMessage(filePathName string) string {

	dir, filename := filepath.Split(filePathName)

	if !strings.HasPrefix(filePathName, "META-INF/statedb") {
		return fmt.Sprintf("metadata file path must begin with META-INF/statedb, found: %s", dir)
	}
	directoryArray := strings.Split(filepath.Clean(dir), "/")
	
	if len(directoryArray) < 4 {
		return fmt.Sprintf("metadata file path must include a database and index directory: %s", dir)
	}
	
	if !contains(validDatabases, directoryArray[2]) {
		return fmt.Sprintf("database name [%s] is not supported, valid options: %s", directoryArray[2], validDatabases)
	}
	
	if len(directoryArray) == 4 && directoryArray[3] != "indexes" {
		return fmt.Sprintf("metadata file path does not have an indexes directory: %s", dir)
	}
	
	if len(directoryArray) != 6 {
		return fmt.Sprintf("metadata file path for collections must include a collections and index directory: %s", dir)
	}
	
	if directoryArray[3] != "collections" || directoryArray[5] != "indexes" {
		return fmt.Sprintf("metadata file path for collections must have a collections and indexes directory: %s", dir)
	}
	
	if !collectionNameValid.MatchString(directoryArray[4]) {
		return fmt.Sprintf("collection name is not valid: %s", directoryArray[4])
	}

	
	if !fileNameValid.MatchString(filename) {
		return fmt.Sprintf("artifact file name is not valid: %s", filename)
	}

	return fmt.Sprintf("metadata file path or name is not supported: %s", dir)

}

func contains(validStrings []string, target string) bool {
	for _, str := range validStrings {
		if str == target {
			return true
		}
	}
	return false
}

func selectFileValidator(filePathName string) fileValidator {
	for validateExp, fileValidator := range fileValidators {
		isValid := validateExp.MatchString(filePathName)
		if isValid {
			return fileValidator
		}
	}
	return nil
}


func couchdbIndexFileValidator(fileName string, fileBytes []byte) error {

	
	boolIsJSON, indexDefinition := isJSON(fileBytes)
	if !boolIsJSON {
		return &InvalidIndexContentError{fmt.Sprintf("Index metadata file [%s] is not a valid JSON", fileName)}
	}

	
	err := validateIndexJSON(indexDefinition)
	if err != nil {
		return &InvalidIndexContentError{fmt.Sprintf("Index metadata file [%s] is not a valid index definition: %s", fileName, err)}
	}

	return nil

}


func isJSON(s []byte) (bool, map[string]interface{}) {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil, js
}

func validateIndexJSON(indexDefinition map[string]interface{}) error {

	
	indexIncluded := false

	
	for jsonKey, jsonValue := range indexDefinition {

		
		switch jsonKey {

		case "index":

			if reflect.TypeOf(jsonValue).Kind() != reflect.Map {
				return fmt.Errorf("Invalid entry, \"index\" must be a JSON")
			}

			err := processIndexMap(jsonValue.(map[string]interface{}))
			if err != nil {
				return err
			}

			indexIncluded = true

		case "ddoc":

			
			if reflect.TypeOf(jsonValue).Kind() != reflect.String {
				return fmt.Errorf("Invalid entry, \"ddoc\" must be a string")
			}

			logger.Debugf("Found index object: \"%s\":\"%s\"", jsonKey, jsonValue)

		case "name":

			
			if reflect.TypeOf(jsonValue).Kind() != reflect.String {
				return fmt.Errorf("Invalid entry, \"name\" must be a string")
			}

			logger.Debugf("Found index object: \"%s\":\"%s\"", jsonKey, jsonValue)

		case "type":

			if jsonValue != "json" {
				return fmt.Errorf("Index type must be json")
			}

			logger.Debugf("Found index object: \"%s\":\"%s\"", jsonKey, jsonValue)

		default:

			return fmt.Errorf("Invalid Entry.  Entry %s", jsonKey)

		}

	}

	if !indexIncluded {
		return fmt.Errorf("Index definition must include a \"fields\" definition")
	}

	return nil

}



func processIndexMap(jsonFragment map[string]interface{}) error {

	
	for jsonKey, jsonValue := range jsonFragment {

		switch jsonKey {

		case "fields":

			switch jsonValueType := jsonValue.(type) {

			case []interface{}:

				
				for _, itemValue := range jsonValueType {

					switch reflect.TypeOf(itemValue).Kind() {

					case reflect.String:
						
						logger.Debugf("Found index field name: \"%s\"", itemValue)

					case reflect.Map:
						
						err := validateFieldMap(itemValue.(map[string]interface{}))
						if err != nil {
							return err
						}

					}
				}

			default:
				return fmt.Errorf("Expecting a JSON array of fields")
			}

		case "partial_filter_selector":

			
			

		default:

			
			
			return fmt.Errorf("Invalid Entry.  Entry %s", jsonKey)

		}

	}

	return nil

}


func validateFieldMap(jsonFragment map[string]interface{}) error {

	
	for jsonKey, jsonValue := range jsonFragment {

		switch jsonValue.(type) {

		case string:
			
			if !(strings.ToLower(jsonValue.(string)) == "asc" || strings.ToLower(jsonValue.(string)) == "desc") {
				return fmt.Errorf("Sort must be either \"asc\" or \"desc\".  \"%s\" was found.", jsonValue)
			}
			logger.Debugf("Found index field name: \"%s\":\"%s\"", jsonKey, jsonValue)

		default:
			return fmt.Errorf("Invalid field definition, fields must be in the form \"fieldname\":\"sort\"")

		}
	}

	return nil

}
