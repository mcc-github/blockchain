/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package leveldbhelper

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)



const internalDBName = "_"

var (
	dbNameKeySep     = []byte{0x00}
	lastKeyIndicator = byte(0x01)
	formatVersionKey = []byte{'f'} 
)








type Conf struct {
	DBPath                string
	ExpectedFormatVersion string
}



type ErrFormatVersionMismatch struct {
	DBPath                string
	ExpectedFormatVersion string
	DataFormatVersion     string
}

func (e *ErrFormatVersionMismatch) Error() string {
	return fmt.Sprintf("unexpected format. db path = [%s], data format = [%s], expected format = [%s]",
		e.DBPath, e.DataFormatVersion, e.ExpectedFormatVersion,
	)
}


func IsFormatVersionMismatch(err error) bool {
	_, ok := err.(*ErrFormatVersionMismatch)
	return ok
}


type Provider struct {
	db *DB

	mux       sync.Mutex
	dbHandles map[string]*DBHandle
}


func NewProvider(conf *Conf) (*Provider, error) {
	db, err := openDBAndCheckFormat(conf)
	if err != nil {
		return nil, err
	}
	return &Provider{
		db:        db,
		dbHandles: make(map[string]*DBHandle),
	}, nil
}

func openDBAndCheckFormat(conf *Conf) (d *DB, e error) {
	db := CreateDB(conf)
	db.Open()

	defer func() {
		if e != nil {
			db.Close()
		}
	}()

	internalDB := &DBHandle{
		db:     db,
		dbName: internalDBName,
	}

	dbEmpty, err := db.isEmpty()
	if err != nil {
		return nil, err
	}

	if dbEmpty && conf.ExpectedFormatVersion != "" {
		logger.Infof("DB is empty Setting db format as %s", conf.ExpectedFormatVersion)
		if err := internalDB.Put(formatVersionKey, []byte(conf.ExpectedFormatVersion), true); err != nil {
			return nil, err
		}
		return db, nil
	}

	formatVersion, err := internalDB.Get(formatVersionKey)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Checking for db format at path [%s]", conf.DBPath)

	if !bytes.Equal(formatVersion, []byte(conf.ExpectedFormatVersion)) {
		logger.Errorf("the db at path [%s] contains data in unexpected format. expected data format = [%s], data format = [%s]",
			conf.DBPath, conf.ExpectedFormatVersion, formatVersion)
		return nil, &ErrFormatVersionMismatch{ExpectedFormatVersion: conf.ExpectedFormatVersion, DataFormatVersion: string(formatVersion), DBPath: conf.DBPath}
	}

	logger.Debug("format is latest, nothing to do")
	return db, nil
}


func (p *Provider) GetDataFormat() (string, error) {
	f, err := p.GetDBHandle(internalDBName).Get(formatVersionKey)
	return string(f), err
}


func (p *Provider) GetDBHandle(dbName string) *DBHandle {
	p.mux.Lock()
	defer p.mux.Unlock()
	dbHandle := p.dbHandles[dbName]
	if dbHandle == nil {
		dbHandle = &DBHandle{dbName, p.db}
		p.dbHandles[dbName] = dbHandle
	}
	return dbHandle
}


func (p *Provider) Close() {
	p.db.Close()
}


type DBHandle struct {
	dbName string
	db     *DB
}


func (h *DBHandle) Get(key []byte) ([]byte, error) {
	return h.db.Get(constructLevelKey(h.dbName, key))
}


func (h *DBHandle) Put(key []byte, value []byte, sync bool) error {
	return h.db.Put(constructLevelKey(h.dbName, key), value, sync)
}


func (h *DBHandle) Delete(key []byte, sync bool) error {
	return h.db.Delete(constructLevelKey(h.dbName, key), sync)
}


func (h *DBHandle) WriteBatch(batch *UpdateBatch, sync bool) error {
	if len(batch.KVs) == 0 {
		return nil
	}
	levelBatch := &leveldb.Batch{}
	for k, v := range batch.KVs {
		key := constructLevelKey(h.dbName, []byte(k))
		if v == nil {
			levelBatch.Delete(key)
		} else {
			levelBatch.Put(key, v)
		}
	}
	if err := h.db.WriteBatch(levelBatch, sync); err != nil {
		return err
	}
	return nil
}




func (h *DBHandle) GetIterator(startKey []byte, endKey []byte) *Iterator {
	sKey := constructLevelKey(h.dbName, startKey)
	eKey := constructLevelKey(h.dbName, endKey)
	if endKey == nil {
		
		eKey[len(eKey)-1] = lastKeyIndicator
	}
	logger.Debugf("Getting iterator for range [%#v] - [%#v]", sKey, eKey)
	return &Iterator{h.db.GetIterator(sKey, eKey)}
}


type UpdateBatch struct {
	KVs map[string][]byte
}


func NewUpdateBatch() *UpdateBatch {
	return &UpdateBatch{make(map[string][]byte)}
}


func (batch *UpdateBatch) Put(key []byte, value []byte) {
	if value == nil {
		panic("Nil value not allowed")
	}
	batch.KVs[string(key)] = value
}


func (batch *UpdateBatch) Delete(key []byte) {
	batch.KVs[string(key)] = nil
}


func (batch *UpdateBatch) Len() int {
	return len(batch.KVs)
}


type Iterator struct {
	iterator.Iterator
}


func (itr *Iterator) Key() []byte {
	return retrieveAppKey(itr.Iterator.Key())
}

func constructLevelKey(dbName string, key []byte) []byte {
	return append(append([]byte(dbName), dbNameKeySep...), key...)
}

func retrieveAppKey(levelKey []byte) []byte {
	return bytes.SplitN(levelKey, dbNameKeySep, 2)[1]
}
