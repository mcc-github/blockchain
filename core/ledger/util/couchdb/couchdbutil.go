/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package couchdb

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/common/util"
)

var expectedDatabaseNamePattern = `[a-z][a-z0-9.$_()-]*`
var maxLength = 238






var chainNameAllowedLength = 50
var namespaceNameAllowedLength = 50
var collectionNameAllowedLength = 50


func CreateCouchInstance(couchDBConnectURL, id, pw string, maxRetries,
	maxRetriesOnStartup int, connectionTimeout time.Duration) (*CouchInstance, error) {

	couchConf, err := CreateConnectionDefinition(couchDBConnectURL,
		id, pw, maxRetries, maxRetriesOnStartup, connectionTimeout)
	if err != nil {
		logger.Errorf("Error during CouchDB CreateConnectionDefinition(): %s\n", err.Error())
		return nil, err
	}

	
	
	
	client := &http.Client{Timeout: couchConf.RequestTimeout}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	transport.DisableCompression = false
	client.Transport = transport

	
	couchInstance := &CouchInstance{conf: *couchConf, client: client}

	connectInfo, retVal, verifyErr := couchInstance.VerifyCouchConfig()
	if verifyErr != nil {
		return nil, verifyErr
	}

	
	if retVal.StatusCode != 200 {
		return nil, fmt.Errorf("CouchDB connection error, expecting return code of 200, received %v", retVal.StatusCode)
	}

	
	errVersion := checkCouchDBVersion(connectInfo.Version)
	if errVersion != nil {
		return nil, errVersion
	}

	return couchInstance, nil
}


func checkCouchDBVersion(version string) error {

	
	majorVersion := strings.Split(version, ".")

	
	majorVersionInt, _ := strconv.Atoi(majorVersion[0])
	if majorVersionInt < 2 {
		return fmt.Errorf("CouchDB must be at least version 2.0.0.  Detected version %s", version)
	}

	return nil
}


func CreateCouchDatabase(couchInstance *CouchInstance, dbName string) (*CouchDatabase, error) {

	databaseName, err := mapAndValidateDatabaseName(dbName)
	if err != nil {
		logger.Errorf("Error during CouchDB CreateDatabaseIfNotExist() for dbName: %s  error: %s\n", dbName, err.Error())
		return nil, err
	}

	couchDBDatabase := CouchDatabase{CouchInstance: couchInstance, DBName: databaseName, IndexWarmCounter: 1}

	
	err = couchDBDatabase.CreateDatabaseIfNotExist()
	if err != nil {
		logger.Errorf("Error during CouchDB CreateDatabaseIfNotExist() for dbName: %s  error: %s\n", dbName, err.Error())
		return nil, err
	}

	return &couchDBDatabase, nil
}


func CreateSystemDatabasesIfNotExist(couchInstance *CouchInstance) error {

	dbName := "_users"
	systemCouchDBDatabase := CouchDatabase{CouchInstance: couchInstance, DBName: dbName, IndexWarmCounter: 1}
	err := systemCouchDBDatabase.CreateDatabaseIfNotExist()
	if err != nil {
		logger.Errorf("Error during CouchDB CreateDatabaseIfNotExist() for system dbName: %s  error: %s\n", dbName, err.Error())
		return err
	}

	dbName = "_replicator"
	systemCouchDBDatabase = CouchDatabase{CouchInstance: couchInstance, DBName: dbName, IndexWarmCounter: 1}
	err = systemCouchDBDatabase.CreateDatabaseIfNotExist()
	if err != nil {
		logger.Errorf("Error during CouchDB CreateDatabaseIfNotExist() for system dbName: %s  error: %s\n", dbName, err.Error())
		return err
	}

	dbName = "_global_changes"
	systemCouchDBDatabase = CouchDatabase{CouchInstance: couchInstance, DBName: dbName, IndexWarmCounter: 1}
	err = systemCouchDBDatabase.CreateDatabaseIfNotExist()
	if err != nil {
		logger.Errorf("Error during CouchDB CreateDatabaseIfNotExist() for system dbName: %s  error: %s\n", dbName, err.Error())
		return err
	}

	return nil

}



func ConstructMetadataDBName(dbName string) string {
	if len(dbName) > maxLength {
		untruncatedDBName := dbName
		
		
		dbName = dbName[:chainNameAllowedLength]
		
		
		dbName = dbName + "(" + hex.EncodeToString(util.ComputeSHA256([]byte(untruncatedDBName))) + ")"
		
	}
	return dbName + "_"
}



func ConstructNamespaceDBName(chainName, namespace string) string {
	
	escapedNamespace := escapeUpperCase(namespace)
	namespaceDBName := chainName + "_" + escapedNamespace

	
	
	
	
	
	
	
	
	

	if len(namespaceDBName) > maxLength {
		
		
		hashOfNamespaceDBName := hex.EncodeToString(util.ComputeSHA256([]byte(chainName + "_" + namespace)))

		
		
		if len(chainName) > chainNameAllowedLength {
			
			chainName = chainName[0:chainNameAllowedLength]
		}
		
		
		
		
		names := strings.Split(escapedNamespace, "$$")
		namespace := names[0]
		if len(namespace) > namespaceNameAllowedLength {
			
			namespace = namespace[0:namespaceNameAllowedLength]
		}

		escapedNamespace = namespace

		
		if len(names) == 2 {
			collection := names[1]
			if len(collection) > collectionNameAllowedLength {
				
				collection = collection[0:collectionNameAllowedLength]
			}
			
			escapedNamespace = escapedNamespace + "$$" + collection
		}
		
		
		
		return chainName + "_" + escapedNamespace + "(" + hashOfNamespaceDBName + ")"
	}
	return namespaceDBName
}











func mapAndValidateDatabaseName(databaseName string) (string, error) {
	
	if len(databaseName) <= 0 {
		return "", fmt.Errorf("Database name is illegal, cannot be empty")
	}
	if len(databaseName) > maxLength {
		return "", fmt.Errorf("Database name is illegal, cannot be longer than %d", maxLength)
	}
	re, err := regexp.Compile(expectedDatabaseNamePattern)
	if err != nil {
		return "", err
	}
	matched := re.FindString(databaseName)
	if len(matched) != len(databaseName) {
		return "", fmt.Errorf("databaseName '%s' does not matches pattern '%s'", databaseName, expectedDatabaseNamePattern)
	}
	
	
	databaseName = strings.Replace(databaseName, ".", "$", -1)
	return databaseName, nil
}



func escapeUpperCase(dbName string) string {
	re := regexp.MustCompile(`([A-Z])`)
	dbName = re.ReplaceAllString(dbName, "$$"+"$1")
	return strings.ToLower(dbName)
}
