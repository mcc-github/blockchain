/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package couchdb

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/ledger/ledgerconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

var logger = flogging.MustGetLogger("couchdb")


const retryWaitTime = 125


type DBOperationResponse struct {
	Ok  bool
	id  string
	rev string
}


type DBInfo struct {
	DbName    string `json:"db_name"`
	UpdateSeq string `json:"update_seq"`
	Sizes     struct {
		File     int `json:"file"`
		External int `json:"external"`
		Active   int `json:"active"`
	} `json:"sizes"`
	PurgeSeq int `json:"purge_seq"`
	Other    struct {
		DataSize int `json:"data_size"`
	} `json:"other"`
	DocDelCount       int    `json:"doc_del_count"`
	DocCount          int    `json:"doc_count"`
	DiskSize          int    `json:"disk_size"`
	DiskFormatVersion int    `json:"disk_format_version"`
	DataSize          int    `json:"data_size"`
	CompactRunning    bool   `json:"compact_running"`
	InstanceStartTime string `json:"instance_start_time"`
}


type ConnectionInfo struct {
	Couchdb string `json:"couchdb"`
	Version string `json:"version"`
	Vendor  struct {
		Name string `json:"name"`
	} `json:"vendor"`
}


type RangeQueryResponse struct {
	TotalRows int32 `json:"total_rows"`
	Offset    int32 `json:"offset"`
	Rows      []struct {
		ID    string `json:"id"`
		Key   string `json:"key"`
		Value struct {
			Rev string `json:"rev"`
		} `json:"value"`
		Doc json.RawMessage `json:"doc"`
	} `json:"rows"`
}


type QueryResponse struct {
	Warning  string            `json:"warning"`
	Docs     []json.RawMessage `json:"docs"`
	Bookmark string            `json:"bookmark"`
}



type DocMetadata struct {
	ID              string          `json:"_id"`
	Rev             string          `json:"_rev"`
	Version         string          `json:"~version"`
	AttachmentsInfo json.RawMessage `json:"_attachments"`
}


type DocID struct {
	ID string `json:"_id"`
}


type QueryResult struct {
	ID          string
	Value       []byte
	Attachments []*AttachmentInfo
}


type CouchConnectionDef struct {
	URL                   string
	Username              string
	Password              string
	MaxRetries            int
	MaxRetriesOnStartup   int
	RequestTimeout        time.Duration
	CreateGlobalChangesDB bool
}


type CouchInstance struct {
	conf   CouchConnectionDef 
	client *http.Client       
	stats  *stats
}


type CouchDatabase struct {
	CouchInstance    *CouchInstance 
	DBName           string
	IndexWarmCounter int
}


type DBReturn struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
	Reason     string `json:"reason"`
}


type CreateIndexResponse struct {
	Result string `json:"result"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}


type AttachmentInfo struct {
	Name            string
	ContentType     string
	Length          uint64
	AttachmentBytes []byte
}


type FileDetails struct {
	Follows     bool   `json:"follows"`
	ContentType string `json:"content_type"`
	Length      int    `json:"length"`
}


type CouchDoc struct {
	JSONValue   []byte
	Attachments []*AttachmentInfo
}


type BatchRetrieveDocMetadataResponse struct {
	Rows []struct {
		ID          string `json:"id"`
		DocMetadata struct {
			ID      string `json:"_id"`
			Rev     string `json:"_rev"`
			Version string `json:"~version"`
		} `json:"doc"`
	} `json:"rows"`
}


type BatchUpdateResponse struct {
	ID     string `json:"id"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
	Ok     bool   `json:"ok"`
	Rev    string `json:"rev"`
}


type Base64Attachment struct {
	ContentType    string `json:"content_type"`
	AttachmentData string `json:"data"`
}


type IndexResult struct {
	DesignDocument string `json:"designdoc"`
	Name           string `json:"name"`
	Definition     string `json:"definition"`
}


type DatabaseSecurity struct {
	Admins struct {
		Names []string `json:"names"`
		Roles []string `json:"roles"`
	} `json:"admins"`
	Members struct {
		Names []string `json:"names"`
		Roles []string `json:"roles"`
	} `json:"members"`
}



func closeResponseBody(resp *http.Response) {
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body) 
		resp.Body.Close()
	}
}


func CreateConnectionDefinition(couchDBAddress, username, password string, maxRetries,
	maxRetriesOnStartup int, requestTimeout time.Duration, createGlobalChangesDB bool) (*CouchConnectionDef, error) {

	logger.Debugf("Entering CreateConnectionDefinition()")

	connectURL := &url.URL{
		Host:   couchDBAddress,
		Scheme: "http",
	}

	
	finalURL, err := url.Parse(connectURL.String())
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing connect URL: %s", connectURL)
	}

	logger.Debugf("Created database configuration  URL=[%s]", finalURL.String())
	logger.Debugf("Exiting CreateConnectionDefinition()")

	
	return &CouchConnectionDef{finalURL.String(), username, password, maxRetries,
		maxRetriesOnStartup, requestTimeout, createGlobalChangesDB}, nil

}


func (dbclient *CouchDatabase) CreateDatabaseIfNotExist() error {

	logger.Debugf("[%s] Entering CreateDatabaseIfNotExist()", dbclient.DBName)

	dbInfo, couchDBReturn, err := dbclient.GetDatabaseInfo()
	if err != nil {
		if couchDBReturn == nil || couchDBReturn.StatusCode != 404 {
			return err
		}
	}

	
	if dbInfo != nil && couchDBReturn.StatusCode == 200 {

		
		errSecurity := dbclient.applyDatabasePermissions()
		if errSecurity != nil {
			return errSecurity
		}

		logger.Debugf("[%s] Database already exists", dbclient.DBName)

		logger.Debugf("[%s] Exiting CreateDatabaseIfNotExist()", dbclient.DBName)

		return nil
	}

	logger.Debugf("[%s] Database does not exist.", dbclient.DBName)

	connectURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	
	resp, _, err := dbclient.handleRequest(http.MethodPut, "CreateDatabaseIfNotExist", connectURL, nil, "", "", maxRetries, true, nil)

	if err != nil {

		
		
		
		
		
		dbInfo, couchDBReturn, errDbInfo := dbclient.GetDatabaseInfo()
		
		if errDbInfo == nil && dbInfo != nil && couchDBReturn.StatusCode == 200 {

			errSecurity := dbclient.applyDatabasePermissions()
			if errSecurity != nil {
				return errSecurity
			}

			logger.Infof("[%s] Created state database", dbclient.DBName)
			logger.Debugf("[%s] Exiting CreateDatabaseIfNotExist()", dbclient.DBName)
			return nil
		}

		return err

	}
	defer closeResponseBody(resp)

	errSecurity := dbclient.applyDatabasePermissions()
	if errSecurity != nil {
		return errSecurity
	}

	logger.Infof("Created state database %s", dbclient.DBName)

	logger.Debugf("[%s] Exiting CreateDatabaseIfNotExist()", dbclient.DBName)

	return nil

}


func (dbclient *CouchDatabase) applyDatabasePermissions() error {

	
	if dbclient.CouchInstance.conf.Username == "" && dbclient.CouchInstance.conf.Password == "" {
		return nil
	}

	securityPermissions := &DatabaseSecurity{}

	securityPermissions.Admins.Names = append(securityPermissions.Admins.Names, dbclient.CouchInstance.conf.Username)
	securityPermissions.Members.Names = append(securityPermissions.Members.Names, dbclient.CouchInstance.conf.Username)

	err := dbclient.ApplyDatabaseSecurity(securityPermissions)
	if err != nil {
		return err
	}

	return nil
}


func (dbclient *CouchDatabase) GetDatabaseInfo() (*DBInfo, *DBReturn, error) {

	connectURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, couchDBReturn, err := dbclient.handleRequest(http.MethodGet, "GetDatabaseInfo", connectURL, nil, "", "", maxRetries, true, nil)
	if err != nil {
		return nil, couchDBReturn, err
	}
	defer closeResponseBody(resp)

	dbResponse := &DBInfo{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&dbResponse)
	if decodeErr != nil {
		return nil, nil, errors.Wrap(decodeErr, "error decoding response body")
	}

	
	logger.Debugw("GetDatabaseInfo()", "dbResponseJSON", dbResponse)

	return dbResponse, couchDBReturn, nil

}


func (couchInstance *CouchInstance) VerifyCouchConfig() (*ConnectionInfo, *DBReturn, error) {

	logger.Debugf("Entering VerifyCouchConfig()")
	defer logger.Debugf("Exiting VerifyCouchConfig()")

	connectURL, err := url.Parse(couchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, nil, errors.Wrapf(err, "error parsing couch instance URL: %s", couchInstance.conf.URL)
	}
	connectURL.Path = "/"

	
	maxRetriesOnStartup := couchInstance.conf.MaxRetriesOnStartup

	resp, couchDBReturn, err := couchInstance.handleRequest(http.MethodGet, "", "VerifyCouchConfig", connectURL, nil,
		couchInstance.conf.Username, couchInstance.conf.Password, maxRetriesOnStartup, true, nil)

	if err != nil {
		return nil, couchDBReturn, errors.WithMessage(err, "unable to connect to CouchDB, check the hostname and port")
	}
	defer closeResponseBody(resp)

	dbResponse := &ConnectionInfo{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&dbResponse)
	if decodeErr != nil {
		return nil, nil, errors.Wrap(decodeErr, "error decoding response body")
	}

	
	logger.Debugw("VerifyConnection() dbResponseJSON: %s", dbResponse)

	
	
	
	
	err = CreateSystemDatabasesIfNotExist(couchInstance)
	if err != nil {
		logger.Errorf("Unable to connect to CouchDB, error: %s. Check the admin username and password.", err)
		return nil, nil, errors.WithMessage(err, "unable to connect to CouchDB. Check the admin username and password")
	}

	return dbResponse, couchDBReturn, nil
}


func (dbclient *CouchDatabase) DropDatabase() (*DBOperationResponse, error) {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering DropDatabase()", dbName)

	connectURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodDelete, "DropDatabase", connectURL, nil, "", "", maxRetries, true, nil)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	dbResponse := &DBOperationResponse{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&dbResponse)
	if decodeErr != nil {
		return nil, errors.Wrap(decodeErr, "error decoding response body")
	}

	if dbResponse.Ok == true {
		logger.Debugf("[%s] Dropped database", dbclient.DBName)
	}

	logger.Debugf("[%s] Exiting DropDatabase()", dbclient.DBName)

	if dbResponse.Ok == true {

		return dbResponse, nil

	}

	return dbResponse, errors.New("error dropping database")

}


func (dbclient *CouchDatabase) EnsureFullCommit() (*DBOperationResponse, error) {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering EnsureFullCommit()", dbName)

	connectURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodPost, "EnsureFullCommit", connectURL, nil, "", "", maxRetries, true, nil, "_ensure_full_commit")
	if err != nil {
		logger.Errorf("Failed to invoke couchdb _ensure_full_commit. Error: %+v", err)
		return nil, err
	}
	defer closeResponseBody(resp)

	dbResponse := &DBOperationResponse{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&dbResponse)
	if decodeErr != nil {
		return nil, errors.Wrap(decodeErr, "error decoding response body")
	}

	
	
	
	
	
	if ledgerconfig.IsAutoWarmIndexesEnabled() {

		if dbclient.IndexWarmCounter >= ledgerconfig.GetWarmIndexesAfterNBlocks() {
			go dbclient.runWarmIndexAllIndexes()
			dbclient.IndexWarmCounter = 0
		}
		dbclient.IndexWarmCounter++

	}

	logger.Debugf("[%s] Exiting EnsureFullCommit()", dbclient.DBName)

	if dbResponse.Ok == true {

		return dbResponse, nil

	}

	return dbResponse, errors.New("error syncing database")
}


func (dbclient *CouchDatabase) SaveDoc(id string, rev string, couchDoc *CouchDoc) (string, error) {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering SaveDoc() id=[%s]", dbName, id)

	if !utf8.ValidString(id) {
		return "", errors.Errorf("doc id [%x] not a valid utf8 string", id)
	}

	saveURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return "", errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	data := []byte{}

	
	defaultBoundary := ""

	
	keepConnectionOpen := true

	
	if couchDoc.Attachments == nil {

		
		if IsJSON(string(couchDoc.JSONValue)) != true {
			return "", errors.New("JSON format is not valid")
		}

		
		data = couchDoc.JSONValue

	} else { 

		
		multipartData, multipartBoundary, err3 := createAttachmentPart(couchDoc, defaultBoundary)
		if err3 != nil {
			return "", err3
		}

		
		for _, attach := range couchDoc.Attachments {
			if attach.Length < 1 {
				keepConnectionOpen = false
			}
		}

		
		data = multipartData.Bytes()

		
		defaultBoundary = multipartBoundary

	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	
	resp, _, err := dbclient.handleRequestWithRevisionRetry(id, http.MethodPut, dbName, "SaveDoc", saveURL, data, rev, defaultBoundary, maxRetries, keepConnectionOpen, nil)

	if err != nil {
		return "", err
	}
	defer closeResponseBody(resp)

	
	revision, err := getRevisionHeader(resp)
	if err != nil {
		return "", err
	}

	logger.Debugf("[%s] Exiting SaveDoc()", dbclient.DBName)

	return revision, nil

}


func (dbclient *CouchDatabase) getDocumentRevision(id string) string {

	var rev = ""

	
	_, revdoc, err := dbclient.ReadDoc(id)
	if err == nil {
		
		rev = revdoc
	}
	return rev
}

func createAttachmentPart(couchDoc *CouchDoc, defaultBoundary string) (bytes.Buffer, string, error) {

	
	writeBuffer := new(bytes.Buffer)

	
	writer := multipart.NewWriter(writeBuffer)

	
	defaultBoundary = writer.Boundary()

	fileAttachments := map[string]FileDetails{}

	for _, attachment := range couchDoc.Attachments {
		fileAttachments[attachment.Name] = FileDetails{true, attachment.ContentType, len(attachment.AttachmentBytes)}
	}

	attachmentJSONMap := map[string]interface{}{
		"_attachments": fileAttachments}

	
	if couchDoc.JSONValue != nil {

		
		genericMap := make(map[string]interface{})

		
		decoder := json.NewDecoder(bytes.NewBuffer(couchDoc.JSONValue))
		decoder.UseNumber()
		decodeErr := decoder.Decode(&genericMap)
		if decodeErr != nil {
			return *writeBuffer, "", errors.Wrap(decodeErr, "error decoding json data")
		}

		
		for jsonKey, jsonValue := range genericMap {
			attachmentJSONMap[jsonKey] = jsonValue
		}

	}

	filesForUpload, err := json.Marshal(attachmentJSONMap)
	if err != nil {
		return *writeBuffer, "", errors.Wrap(err, "error marshalling json data")
	}

	logger.Debugf(string(filesForUpload))

	
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", "application/json")

	part, err := writer.CreatePart(header)
	if err != nil {
		return *writeBuffer, defaultBoundary, errors.Wrap(err, "error creating multipart")
	}

	part.Write(filesForUpload)

	for _, attachment := range couchDoc.Attachments {

		header := make(textproto.MIMEHeader)
		part, err2 := writer.CreatePart(header)
		if err2 != nil {
			return *writeBuffer, defaultBoundary, errors.Wrap(err2, "error creating multipart")
		}
		part.Write(attachment.AttachmentBytes)

	}

	err = writer.Close()
	if err != nil {
		return *writeBuffer, defaultBoundary, errors.Wrap(err, "error closing multipart writer")
	}

	return *writeBuffer, defaultBoundary, nil

}

func getRevisionHeader(resp *http.Response) (string, error) {

	if resp == nil {
		return "", errors.New("no response received from CouchDB")
	}

	revision := resp.Header.Get("Etag")

	if revision == "" {
		return "", errors.New("no revision tag detected")
	}

	reg := regexp.MustCompile(`"([^"]*)"`)
	revisionNoQuotes := reg.ReplaceAllString(revision, "${1}")
	return revisionNoQuotes, nil

}



func (dbclient *CouchDatabase) ReadDoc(id string) (*CouchDoc, string, error) {
	var couchDoc CouchDoc
	attachments := []*AttachmentInfo{}
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering ReadDoc()  id=[%s]", dbName, id)
	if !utf8.ValidString(id) {
		return nil, "", errors.Errorf("doc id [%x] not a valid utf8 string", id)
	}

	readURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, "", errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	query := readURL.Query()
	query.Add("attachments", "true")

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, couchDBReturn, err := dbclient.handleRequest(http.MethodGet, "ReadDoc", readURL, nil, "", "", maxRetries, true, &query, id)
	if err != nil {
		if couchDBReturn != nil && couchDBReturn.StatusCode == 404 {
			logger.Debugf("[%s] Document not found (404), returning nil value instead of 404 error", dbclient.DBName)
			
			
			return nil, "", nil
		}
		logger.Debugf("[%s] couchDBReturn=%v\n", dbclient.DBName, couchDBReturn)
		return nil, "", err
	}
	defer closeResponseBody(resp)

	
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	
	revision, err := getRevisionHeader(resp)
	if err != nil {
		return nil, "", err
	}

	
	if strings.HasPrefix(mediaType, "multipart/") {
		
		multipartReader := multipart.NewReader(resp.Body, params["boundary"])
		for {
			p, err := multipartReader.NextPart()
			if err == io.EOF {
				break 
			}
			if err != nil {
				return nil, "", errors.Wrap(err, "error reading next multipart")
			}

			defer p.Close()

			logger.Debugf("[%s] part header=%s", dbclient.DBName, p.Header)
			switch p.Header.Get("Content-Type") {
			case "application/json":
				partdata, err := ioutil.ReadAll(p)
				if err != nil {
					return nil, "", errors.Wrap(err, "error reading multipart data")
				}
				couchDoc.JSONValue = partdata
			default:

				
				attachment := &AttachmentInfo{}
				attachment.ContentType = p.Header.Get("Content-Type")
				contentDispositionParts := strings.Split(p.Header.Get("Content-Disposition"), ";")
				if strings.TrimSpace(contentDispositionParts[0]) == "attachment" {
					switch p.Header.Get("Content-Encoding") {
					case "gzip": 

						var respBody []byte

						gr, err := gzip.NewReader(p)
						if err != nil {
							return nil, "", errors.Wrap(err, "error creating gzip reader")
						}
						respBody, err = ioutil.ReadAll(gr)
						if err != nil {
							return nil, "", errors.Wrap(err, "error reading gzip data")
						}

						logger.Debugf("[%s] Retrieved attachment data", dbclient.DBName)
						attachment.AttachmentBytes = respBody
						attachment.Length = uint64(len(attachment.AttachmentBytes))
						attachment.Name = p.FileName()
						attachments = append(attachments, attachment)

					default:

						
						partdata, err := ioutil.ReadAll(p)
						if err != nil {
							return nil, "", errors.Wrap(err, "error reading multipart data")
						}
						logger.Debugf("[%s] Retrieved attachment data", dbclient.DBName)
						attachment.AttachmentBytes = partdata
						attachment.Length = uint64(len(attachment.AttachmentBytes))
						attachment.Name = p.FileName()
						attachments = append(attachments, attachment)

					} 
				} 
			} 
		} 

		couchDoc.Attachments = attachments

		return &couchDoc, revision, nil
	}

	
	couchDoc.JSONValue, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "error reading response body")
	}

	logger.Debugf("[%s] Exiting ReadDoc()", dbclient.DBName)
	return &couchDoc, revision, nil
}





func (dbclient *CouchDatabase) ReadDocRange(startKey, endKey string, limit int32) ([]*QueryResult, string, error) {
	dbName := dbclient.DBName
	logger.Debugf("[%s] Entering ReadDocRange()  startKey=%s, endKey=%s", dbName, startKey, endKey)

	var results []*QueryResult

	rangeURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, "", errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	queryParms := rangeURL.Query()
	
	queryParms.Set("limit", strconv.FormatInt(int64(limit+1), 10))
	queryParms.Add("include_docs", "true")
	queryParms.Add("inclusive_end", "false") 

	
	if startKey != "" {
		if startKey, err = encodeForJSON(startKey); err != nil {
			return nil, "", err
		}
		queryParms.Add("startkey", "\""+startKey+"\"")
	}

	
	if endKey != "" {
		var err error
		if endKey, err = encodeForJSON(endKey); err != nil {
			return nil, "", err
		}
		queryParms.Add("endkey", "\""+endKey+"\"")
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodGet, "RangeDocRange", rangeURL, nil, "", "", maxRetries, true, &queryParms, "_all_docs")
	if err != nil {
		return nil, "", err
	}
	defer closeResponseBody(resp)

	if logger.IsEnabledFor(zapcore.DebugLevel) {
		dump, err2 := httputil.DumpResponse(resp, true)
		if err2 != nil {
			log.Fatal(err2)
		}
		logger.Debugf("[%s] %s", dbclient.DBName, dump)
	}

	
	jsonResponseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "error reading response body")
	}

	var jsonResponse = &RangeQueryResponse{}
	err2 := json.Unmarshal(jsonResponseRaw, &jsonResponse)
	if err2 != nil {
		return nil, "", errors.Wrap(err2, "error unmarshalling json data")
	}

	
	
	if jsonResponse.TotalRows > limit {
		jsonResponse.TotalRows = limit
	}

	logger.Debugf("[%s] Total Rows: %d", dbclient.DBName, jsonResponse.TotalRows)

	
	nextStartKey := endKey

	for index, row := range jsonResponse.Rows {

		var docMetadata = &DocMetadata{}
		err3 := json.Unmarshal(row.Doc, &docMetadata)
		if err3 != nil {
			return nil, "", errors.Wrap(err3, "error unmarshalling json data")
		}

		
		
		if int32(index) >= jsonResponse.TotalRows {
			nextStartKey = docMetadata.ID
			continue
		}

		if docMetadata.AttachmentsInfo != nil {

			logger.Debugf("[%s] Adding JSON document and attachments for id: %s", dbclient.DBName, docMetadata.ID)

			couchDoc, _, err := dbclient.ReadDoc(docMetadata.ID)
			if err != nil {
				return nil, "", err
			}

			var addDocument = &QueryResult{docMetadata.ID, couchDoc.JSONValue, couchDoc.Attachments}
			results = append(results, addDocument)

		} else {

			logger.Debugf("[%s] Adding json docment for id: %s", dbclient.DBName, docMetadata.ID)

			var addDocument = &QueryResult{docMetadata.ID, row.Doc, nil}
			results = append(results, addDocument)

		}

	}

	logger.Debugf("[%s] Exiting ReadDocRange()", dbclient.DBName)

	return results, nextStartKey, nil

}


func (dbclient *CouchDatabase) DeleteDoc(id, rev string) error {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering DeleteDoc()  id=%s", dbName, id)

	deleteURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	
	resp, couchDBReturn, err := dbclient.handleRequestWithRevisionRetry(id, http.MethodDelete, dbName, "DeleteDoc",
		deleteURL, nil, "", "", maxRetries, true, nil)

	if err != nil {
		if couchDBReturn != nil && couchDBReturn.StatusCode == 404 {
			logger.Debugf("[%s] Document not found (404), returning nil value instead of 404 error", dbclient.DBName)
			
			
			return nil
		}
		return err
	}
	defer closeResponseBody(resp)

	logger.Debugf("[%s] Exiting DeleteDoc()", dbclient.DBName)

	return nil

}


func (dbclient *CouchDatabase) QueryDocuments(query string) ([]*QueryResult, string, error) {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering QueryDocuments()  query=%s", dbName, query)

	var results []*QueryResult

	queryURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, "", errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodPost, "QueryDocuments", queryURL, []byte(query), "", "", maxRetries, true, nil, "_find")
	if err != nil {
		return nil, "", err
	}
	defer closeResponseBody(resp)

	if logger.IsEnabledFor(zapcore.DebugLevel) {
		dump, err2 := httputil.DumpResponse(resp, true)
		if err2 != nil {
			log.Fatal(err2)
		}
		logger.Debugf("[%s] %s", dbclient.DBName, dump)
	}

	
	jsonResponseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "error reading response body")
	}

	var jsonResponse = &QueryResponse{}

	err2 := json.Unmarshal(jsonResponseRaw, &jsonResponse)
	if err2 != nil {
		return nil, "", errors.Wrap(err2, "error unmarshalling json data")
	}

	if jsonResponse.Warning != "" {
		logger.Warnf("The query [%s] caused the following warning: [%s]", query, jsonResponse.Warning)
	}

	for _, row := range jsonResponse.Docs {

		var docMetadata = &DocMetadata{}
		err3 := json.Unmarshal(row, &docMetadata)
		if err3 != nil {
			return nil, "", errors.Wrap(err3, "error unmarshalling json data")
		}

		if docMetadata.AttachmentsInfo != nil {

			logger.Debugf("[%s] Adding JSON docment and attachments for id: %s", dbclient.DBName, docMetadata.ID)

			couchDoc, _, err := dbclient.ReadDoc(docMetadata.ID)
			if err != nil {
				return nil, "", err
			}
			var addDocument = &QueryResult{ID: docMetadata.ID, Value: couchDoc.JSONValue, Attachments: couchDoc.Attachments}
			results = append(results, addDocument)

		} else {
			logger.Debugf("[%s] Adding json docment for id: %s", dbclient.DBName, docMetadata.ID)
			var addDocument = &QueryResult{ID: docMetadata.ID, Value: row, Attachments: nil}

			results = append(results, addDocument)

		}
	}

	logger.Debugf("[%s] Exiting QueryDocuments()", dbclient.DBName)

	return results, jsonResponse.Bookmark, nil

}


func (dbclient *CouchDatabase) ListIndex() ([]*IndexResult, error) {

	
	type indexDefinition struct {
		DesignDocument string          `json:"ddoc"`
		Name           string          `json:"name"`
		Type           string          `json:"type"`
		Definition     json.RawMessage `json:"def"`
	}

	
	type listIndexResponse struct {
		TotalRows int               `json:"total_rows"`
		Indexes   []indexDefinition `json:"indexes"`
	}

	dbName := dbclient.DBName
	logger.Debug("[%s] Entering ListIndex()", dbName)

	indexURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodGet, "ListIndex", indexURL, nil, "", "", maxRetries, true, nil, "_index")
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	
	jsonResponseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var jsonResponse = &listIndexResponse{}

	err2 := json.Unmarshal(jsonResponseRaw, jsonResponse)
	if err2 != nil {
		return nil, errors.Wrap(err2, "error unmarshalling json data")
	}

	var results []*IndexResult

	for _, row := range jsonResponse.Indexes {

		
		
		designDoc := row.DesignDocument
		s := strings.SplitAfterN(designDoc, "_design/", 2)
		if len(s) > 1 {
			designDoc = s[1]

			
			var addIndexResult = &IndexResult{DesignDocument: designDoc, Name: row.Name, Definition: fmt.Sprintf("%s", row.Definition)}
			results = append(results, addIndexResult)
		}

	}

	logger.Debugf("[%s] Exiting ListIndex()", dbclient.DBName)

	return results, nil

}


func (dbclient *CouchDatabase) CreateIndex(indexdefinition string) (*CreateIndexResponse, error) {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering CreateIndex()  indexdefinition=%s", dbName, indexdefinition)

	
	if IsJSON(indexdefinition) != true {
		return nil, errors.New("JSON format is not valid")
	}

	indexURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodPost, "CreateIndex", indexURL, []byte(indexdefinition), "", "", maxRetries, true, nil, "_index")
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	if resp == nil {
		return nil, errors.New("invalid response received from CouchDB")
	}

	
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	couchDBReturn := &CreateIndexResponse{}

	jsonBytes := []byte(respBody)

	
	err = json.Unmarshal(jsonBytes, &couchDBReturn)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling json data")
	}

	if couchDBReturn.Result == "created" {

		logger.Infof("Created CouchDB index [%s] in state database [%s] using design document [%s]", couchDBReturn.Name, dbclient.DBName, couchDBReturn.ID)

		return couchDBReturn, nil

	}

	logger.Infof("Updated CouchDB index [%s] in state database [%s] using design document [%s]", couchDBReturn.Name, dbclient.DBName, couchDBReturn.ID)

	return couchDBReturn, nil
}


func (dbclient *CouchDatabase) DeleteIndex(designdoc, indexname string) error {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering DeleteIndex()  designdoc=%s  indexname=%s", dbName, designdoc, indexname)

	indexURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodDelete, "DeleteIndex", indexURL, nil, "", "", maxRetries, true, nil, "_index", designdoc, "json", indexname)
	if err != nil {
		return err
	}
	defer closeResponseBody(resp)

	return nil

}


func (dbclient *CouchDatabase) WarmIndex(designdoc, indexname string) error {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering WarmIndex()  designdoc=%s  indexname=%s", dbName, designdoc, indexname)

	indexURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	queryParms := indexURL.Query()
	
	
	queryParms.Add("stale", "update_after")

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodGet, "WarmIndex", indexURL, nil, "", "", maxRetries, true, &queryParms, "_design", designdoc, "_view", indexname)
	if err != nil {
		return err
	}
	defer closeResponseBody(resp)

	return nil

}


func (dbclient *CouchDatabase) runWarmIndexAllIndexes() {

	err := dbclient.WarmIndexAllIndexes()
	if err != nil {
		logger.Errorf("Error detected during WarmIndexAllIndexes(): %+v", err)
	}

}


func (dbclient *CouchDatabase) WarmIndexAllIndexes() error {

	logger.Debugf("[%s] Entering WarmIndexAllIndexes()", dbclient.DBName)

	
	listResult, err := dbclient.ListIndex()
	if err != nil {
		return err
	}

	
	for _, elem := range listResult {

		err := dbclient.WarmIndex(elem.DesignDocument, elem.Name)
		if err != nil {
			return err
		}

	}

	logger.Debugf("[%s] Exiting WarmIndexAllIndexes()", dbclient.DBName)

	return nil

}


func (dbclient *CouchDatabase) GetDatabaseSecurity() (*DatabaseSecurity, error) {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering GetDatabaseSecurity()", dbName)

	securityURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodGet, "GetDatabaseSecurity", securityURL, nil, "", "", maxRetries, true, nil, "_security")

	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	
	jsonResponseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var jsonResponse = &DatabaseSecurity{}

	err2 := json.Unmarshal(jsonResponseRaw, jsonResponse)
	if err2 != nil {
		return nil, errors.Wrap(err2, "error unmarshalling json data")
	}

	logger.Debugf("[%s] Exiting GetDatabaseSecurity()", dbclient.DBName)

	return jsonResponse, nil

}


func (dbclient *CouchDatabase) ApplyDatabaseSecurity(databaseSecurity *DatabaseSecurity) error {
	dbName := dbclient.DBName

	logger.Debugf("[%s] Entering ApplyDatabaseSecurity()", dbName)

	securityURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	
	if databaseSecurity.Admins.Names == nil {
		databaseSecurity.Admins.Names = make([]string, 0)
	}
	if databaseSecurity.Admins.Roles == nil {
		databaseSecurity.Admins.Roles = make([]string, 0)
	}
	if databaseSecurity.Members.Names == nil {
		databaseSecurity.Members.Names = make([]string, 0)
	}
	if databaseSecurity.Members.Roles == nil {
		databaseSecurity.Members.Roles = make([]string, 0)
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	databaseSecurityJSON, err := json.Marshal(databaseSecurity)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling json data")
	}

	logger.Debugf("[%s] Applying security to database: %s", dbclient.DBName, string(databaseSecurityJSON))

	resp, _, err := dbclient.handleRequest(http.MethodPut, "ApplyDatabaseSecurity", securityURL, databaseSecurityJSON, "", "", maxRetries, true, nil, "_security")

	if err != nil {
		return err
	}
	defer closeResponseBody(resp)

	logger.Debugf("[%s] Exiting ApplyDatabaseSecurity()", dbclient.DBName)

	return nil

}



func (dbclient *CouchDatabase) BatchRetrieveDocumentMetadata(keys []string) ([]*DocMetadata, error) {

	logger.Debugf("[%s] Entering BatchRetrieveDocumentMetadata()  keys=%s", dbclient.DBName, keys)

	batchRetrieveURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	queryParms := batchRetrieveURL.Query()

	
	
	
	
	
	
	queryParms.Add("include_docs", "true")

	keymap := make(map[string]interface{})

	keymap["keys"] = keys

	jsonKeys, err := json.Marshal(keymap)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling json data")
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodPost, "BatchRetrieveDocumentMetadata", batchRetrieveURL, jsonKeys, "", "", maxRetries, true, &queryParms, "_all_docs")
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	if logger.IsEnabledFor(zapcore.DebugLevel) {
		dump, _ := httputil.DumpResponse(resp, false)
		
		logger.Debugf("[%s] HTTP Response: %s", dbclient.DBName, bytes.Replace(dump, []byte{0x0d, 0x0a}, []byte{0x20, 0x7c, 0x20}, -1))
	}

	
	jsonResponseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var jsonResponse = &BatchRetrieveDocMetadataResponse{}

	err2 := json.Unmarshal(jsonResponseRaw, &jsonResponse)
	if err2 != nil {
		return nil, errors.Wrap(err2, "error unmarshalling json data")
	}

	docMetadataArray := []*DocMetadata{}

	for _, row := range jsonResponse.Rows {
		docMetadata := &DocMetadata{ID: row.ID, Rev: row.DocMetadata.Rev, Version: row.DocMetadata.Version}
		docMetadataArray = append(docMetadataArray, docMetadata)
	}

	logger.Debugf("[%s] Exiting BatchRetrieveDocumentMetadata()", dbclient.DBName)

	return docMetadataArray, nil

}


func (dbclient *CouchDatabase) BatchUpdateDocuments(documents []*CouchDoc) ([]*BatchUpdateResponse, error) {
	dbName := dbclient.DBName

	if logger.IsEnabledFor(zapcore.DebugLevel) {
		documentIdsString, err := printDocumentIds(documents)
		if err == nil {
			logger.Debugf("[%s] Entering BatchUpdateDocuments()  document ids=[%s]", dbName, documentIdsString)
		} else {
			logger.Debugf("[%s] Entering BatchUpdateDocuments()  Could not print document ids due to error: %+v", dbName, err)
		}
	}

	batchUpdateURL, err := url.Parse(dbclient.CouchInstance.conf.URL)
	if err != nil {
		logger.Errorf("URL parse error: %s", err)
		return nil, errors.Wrapf(err, "error parsing CouchDB URL: %s", dbclient.CouchInstance.conf.URL)
	}

	documentMap := make(map[string]interface{})

	var jsonDocumentMap []interface{}

	for _, jsonDocument := range documents {

		
		var document = make(map[string]interface{})

		
		err = json.Unmarshal(jsonDocument.JSONValue, &document)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling json data")
		}

		
		if len(jsonDocument.Attachments) > 0 {

			
			fileAttachment := make(map[string]interface{})

			
			
			for _, attachment := range jsonDocument.Attachments {
				fileAttachment[attachment.Name] = Base64Attachment{attachment.ContentType,
					base64.StdEncoding.EncodeToString(attachment.AttachmentBytes)}
			}

			
			document["_attachments"] = fileAttachment

		}

		
		jsonDocumentMap = append(jsonDocumentMap, document)

	}

	
	documentMap["docs"] = jsonDocumentMap

	bulkDocsJSON, err := json.Marshal(documentMap)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling json data")
	}

	
	maxRetries := dbclient.CouchInstance.conf.MaxRetries

	resp, _, err := dbclient.handleRequest(http.MethodPost, "BatchUpdateDocuments", batchUpdateURL, bulkDocsJSON, "", "", maxRetries, true, nil, "_bulk_docs")
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp)

	if logger.IsEnabledFor(zapcore.DebugLevel) {
		dump, _ := httputil.DumpResponse(resp, false)
		
		logger.Debugf("[%s] HTTP Response: %s", dbclient.DBName, bytes.Replace(dump, []byte{0x0d, 0x0a}, []byte{0x20, 0x7c, 0x20}, -1))
	}

	
	jsonResponseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var jsonResponse = []*BatchUpdateResponse{}
	err2 := json.Unmarshal(jsonResponseRaw, &jsonResponse)
	if err2 != nil {
		return nil, errors.Wrap(err2, "error unmarshalling json data")
	}

	logger.Debugf("[%s] Exiting BatchUpdateDocuments() _bulk_docs response=[%s]", dbclient.DBName, string(jsonResponseRaw))

	return jsonResponse, nil

}





func (dbclient *CouchDatabase) handleRequestWithRevisionRetry(id, method, dbName, functionName string, connectURL *url.URL, data []byte, rev string,
	multipartBoundary string, maxRetries int, keepConnectionOpen bool, queryParms *url.Values) (*http.Response, *DBReturn, error) {

	
	revisionConflictDetected := false
	var resp *http.Response
	var couchDBReturn *DBReturn
	var errResp error

	
	
	
	for attempts := 0; attempts <= maxRetries; attempts++ {

		
		
		if rev == "" || revisionConflictDetected {
			rev = dbclient.getDocumentRevision(id)
		}

		
		resp, couchDBReturn, errResp = dbclient.CouchInstance.handleRequest(method, dbName, functionName, connectURL,
			data, rev, multipartBoundary, maxRetries, keepConnectionOpen, queryParms, id)

		
		
		if couchDBReturn != nil && couchDBReturn.StatusCode == 409 {
			logger.Warningf("CouchDB document revision conflict detected, retrying. Attempt:%v", attempts+1)
			revisionConflictDetected = true
		} else {
			break
		}
	}

	
	return resp, couchDBReturn, errResp
}

func (dbclient *CouchDatabase) handleRequest(method, functionName string, connectURL *url.URL, data []byte, rev, multipartBoundary string,
	maxRetries int, keepConnectionOpen bool, queryParms *url.Values, pathElements ...string) (*http.Response, *DBReturn, error) {

	return dbclient.CouchInstance.handleRequest(
		method, dbclient.DBName, functionName, connectURL, data, rev, multipartBoundary,
		maxRetries, keepConnectionOpen, queryParms, pathElements...,
	)
}





func (couchInstance *CouchInstance) handleRequest(method, dbName, functionName string, connectURL *url.URL, data []byte, rev string,
	multipartBoundary string, maxRetries int, keepConnectionOpen bool, queryParms *url.Values, pathElements ...string) (*http.Response, *DBReturn, error) {

	logger.Debugf("Entering handleRequest()  method=%s  url=%v  dbName=%s", method, connectURL, dbName)

	
	var resp *http.Response
	var errResp error
	couchDBReturn := &DBReturn{}
	defer couchInstance.recordMetric(time.Now(), dbName, functionName, couchDBReturn)

	
	waitDuration := retryWaitTime * time.Millisecond

	if maxRetries < 0 {
		return nil, nil, errors.New("number of retries must be zero or greater")
	}

	requestURL := constructCouchDBUrl(connectURL, dbName, pathElements...)

	if queryParms != nil {
		requestURL.RawQuery = queryParms.Encode()
	}

	logger.Debugf("Request URL: %s", requestURL)

	
	
	
	
	
	for attempts := 0; attempts <= maxRetries; attempts++ {

		
		payloadData := new(bytes.Buffer)

		payloadData.ReadFrom(bytes.NewReader(data))

		
		req, err := http.NewRequest(method, requestURL.String(), payloadData)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error creating http request")
		}

		
		
		
		if !keepConnectionOpen {
			req.Close = true
		}

		
		if method == http.MethodPut || method == http.MethodPost || method == http.MethodDelete {

			
			
			if multipartBoundary == "" {
				req.Header.Set("Content-Type", "application/json")
			} else {
				req.Header.Set("Content-Type", "multipart/related;boundary=\""+multipartBoundary+"\"")
			}

			
			if rev != "" {
				req.Header.Set("If-Match", rev)
			}
		}

		
		if method == http.MethodPut || method == http.MethodPost {
			req.Header.Set("Accept", "application/json")
		}

		
		if method == http.MethodGet {
			req.Header.Set("Accept", "multipart/related")
		}

		
		if couchInstance.conf.Username != "" && couchInstance.conf.Password != "" {
			
			req.SetBasicAuth(couchInstance.conf.Username, couchInstance.conf.Password)
		}

		if logger.IsEnabledFor(zapcore.DebugLevel) {
			dump, _ := httputil.DumpRequestOut(req, false)
			
			logger.Debugf("HTTP Request: %s", bytes.Replace(dump, []byte{0x0d, 0x0a}, []byte{0x20, 0x7c, 0x20}, -1))
		}

		
		resp, errResp = couchInstance.client.Do(req)

		
		if invalidCouchDBReturn(resp, errResp) {
			continue
		}

		
		if errResp == nil && resp != nil && resp.StatusCode < 500 {
			
			if resp.StatusCode >= 400 {
				
				jsonError, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, nil, errors.Wrap(err, "error reading response body")
				}
				defer closeResponseBody(resp)

				errorBytes := []byte(jsonError)
				
				err = json.Unmarshal(errorBytes, &couchDBReturn)
				if err != nil {
					return nil, nil, errors.Wrap(err, "error unmarshalling json data")
				}
			}

			break
		}

		
		if maxRetries > 0 {

			
			if errResp != nil {

				
				logger.Warningf("Retrying couchdb request in %s. Attempt:%v  Error:%v",
					waitDuration.String(), attempts+1, errResp.Error())

				
			} else {
				
				jsonError, err := ioutil.ReadAll(resp.Body)
				defer closeResponseBody(resp)
				if err != nil {
					return nil, nil, errors.Wrap(err, "error reading response body")
				}

				errorBytes := []byte(jsonError)
				
				err = json.Unmarshal(errorBytes, &couchDBReturn)
				if err != nil {
					return nil, nil, errors.Wrap(err, "error unmarshalling json data")
				}

				
				logger.Warningf("Retrying couchdb request in %s. Attempt:%v  Couch DB Error:%s,  Status Code:%v  Reason:%v",
					waitDuration.String(), attempts+1, couchDBReturn.Error, resp.Status, couchDBReturn.Reason)

			}
			
			time.Sleep(waitDuration)

			
			waitDuration *= 2

		}

	} 

	
	if errResp != nil {
		return nil, couchDBReturn, errResp
	}

	
	
	
	
	if invalidCouchDBReturn(resp, errResp) {
		return nil, nil, errors.New("unable to connect to CouchDB, check the hostname and port.")
	}

	
	couchDBReturn.StatusCode = resp.StatusCode

	
	
	
	if resp.StatusCode >= 400 {

		
		logger.Debugf("Error handling CouchDB request. Error:%s,  Status Code:%v,  Reason:%s",
			couchDBReturn.Error, resp.StatusCode, couchDBReturn.Reason)

		return nil, couchDBReturn, errors.Errorf("error handling CouchDB request. Error:%s,  Status Code:%v,  Reason:%s",
			couchDBReturn.Error, resp.StatusCode, couchDBReturn.Reason)

	}

	logger.Debugf("Exiting handleRequest()")

	
	return resp, couchDBReturn, nil
}

func (ci *CouchInstance) recordMetric(startTime time.Time, dbName, api string, couchDBReturn *DBReturn) {
	ci.stats.observeProcessingTime(startTime, dbName, api, strconv.Itoa(couchDBReturn.StatusCode))
}


func invalidCouchDBReturn(resp *http.Response, errResp error) bool {
	if resp == nil && errResp == nil {
		return true
	}
	return false
}


func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}





func encodePathElement(str string) string {

	u := &url.URL{}
	u.Path = str
	encodedStr := u.EscapedPath() 
	encodedStr = strings.Replace(encodedStr, "/", "%2F", -1)
	encodedStr = strings.Replace(encodedStr, "+", "%2B", -1)

	return encodedStr
}

func encodeForJSON(str string) (string, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(str); err != nil {
		return "", errors.Wrap(err, "error encoding json data")
	}
	
	buffer := buf.Bytes()
	return string(buffer[1 : len(buffer)-2]), nil
}



func printDocumentIds(documentPointers []*CouchDoc) (string, error) {

	documentIds := []string{}

	for _, documentPointer := range documentPointers {
		docMetadata := &DocMetadata{}
		err := json.Unmarshal(documentPointer.JSONValue, &docMetadata)
		if err != nil {
			return "", errors.Wrap(err, "error unmarshalling json data")
		}
		documentIds = append(documentIds, docMetadata.ID)
	}
	return strings.Join(documentIds, ","), nil
}
