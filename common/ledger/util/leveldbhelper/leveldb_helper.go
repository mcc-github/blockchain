/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package leveldbhelper

import (
	"fmt"
	"sync"
	"syscall"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/ledger/util"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	goleveldbutil "github.com/syndtr/goleveldb/leveldb/util"
)

var logger = flogging.MustGetLogger("leveldbhelper")

type dbState int32

const (
	closed dbState = iota
	opened
)


type DB struct {
	conf    *Conf
	db      *leveldb.DB
	dbState dbState
	mutex   sync.RWMutex

	readOpts        *opt.ReadOptions
	writeOptsNoSync *opt.WriteOptions
	writeOptsSync   *opt.WriteOptions
}


func CreateDB(conf *Conf) *DB {
	readOpts := &opt.ReadOptions{}
	writeOptsNoSync := &opt.WriteOptions{}
	writeOptsSync := &opt.WriteOptions{}
	writeOptsSync.Sync = true

	return &DB{
		conf:            conf,
		dbState:         closed,
		readOpts:        readOpts,
		writeOptsNoSync: writeOptsNoSync,
		writeOptsSync:   writeOptsSync}
}


func (dbInst *DB) Open() {
	dbInst.mutex.Lock()
	defer dbInst.mutex.Unlock()
	if dbInst.dbState == opened {
		return
	}
	dbOpts := &opt.Options{}
	dbPath := dbInst.conf.DBPath
	var err error
	var dirEmpty bool
	if dirEmpty, err = util.CreateDirIfMissing(dbPath); err != nil {
		panic(fmt.Sprintf("Error creating dir if missing: %s", err))
	}
	dbOpts.ErrorIfMissing = !dirEmpty
	if dbInst.db, err = leveldb.OpenFile(dbPath, dbOpts); err != nil {
		panic(fmt.Sprintf("Error opening leveldb: %s", err))
	}
	dbInst.dbState = opened
}

func (dbInst *DB) isEmpty() (bool, error) {
	itr := dbInst.db.NewIterator(&goleveldbutil.Range{}, dbInst.readOpts)
	defer itr.Release()
	hasItems := itr.Next()
	return !hasItems,
		errors.Wrapf(itr.Error(), "error while trying to see if the leveldb at path [%s] is empty", dbInst.conf.DBPath)
}


func (dbInst *DB) Close() {
	dbInst.mutex.Lock()
	defer dbInst.mutex.Unlock()
	if dbInst.dbState == closed {
		return
	}
	if err := dbInst.db.Close(); err != nil {
		logger.Errorf("Error closing leveldb: %s", err)
	}
	dbInst.dbState = closed
}


func (dbInst *DB) Get(key []byte) ([]byte, error) {
	dbInst.mutex.RLock()
	defer dbInst.mutex.RUnlock()
	value, err := dbInst.db.Get(key, dbInst.readOpts)
	if err == leveldb.ErrNotFound {
		value = nil
		err = nil
	}
	if err != nil {
		logger.Errorf("Error retrieving leveldb key [%#v]: %s", key, err)
		return nil, errors.Wrapf(err, "error retrieving leveldb key [%#v]", key)
	}
	return value, nil
}


func (dbInst *DB) Put(key []byte, value []byte, sync bool) error {
	dbInst.mutex.RLock()
	defer dbInst.mutex.RUnlock()
	wo := dbInst.writeOptsNoSync
	if sync {
		wo = dbInst.writeOptsSync
	}
	err := dbInst.db.Put(key, value, wo)
	if err != nil {
		logger.Errorf("Error writing leveldb key [%#v]", key)
		return errors.Wrapf(err, "error writing leveldb key [%#v]", key)
	}
	return nil
}


func (dbInst *DB) Delete(key []byte, sync bool) error {
	dbInst.mutex.RLock()
	defer dbInst.mutex.RUnlock()
	wo := dbInst.writeOptsNoSync
	if sync {
		wo = dbInst.writeOptsSync
	}
	err := dbInst.db.Delete(key, wo)
	if err != nil {
		logger.Errorf("Error deleting leveldb key [%#v]", key)
		return errors.Wrapf(err, "error deleting leveldb key [%#v]", key)
	}
	return nil
}




func (dbInst *DB) GetIterator(startKey []byte, endKey []byte) iterator.Iterator {
	dbInst.mutex.RLock()
	defer dbInst.mutex.RUnlock()
	return dbInst.db.NewIterator(&goleveldbutil.Range{Start: startKey, Limit: endKey}, dbInst.readOpts)
}


func (dbInst *DB) WriteBatch(batch *leveldb.Batch, sync bool) error {
	dbInst.mutex.RLock()
	defer dbInst.mutex.RUnlock()
	wo := dbInst.writeOptsNoSync
	if sync {
		wo = dbInst.writeOptsSync
	}
	if err := dbInst.db.Write(batch, wo); err != nil {
		return errors.Wrap(err, "error writing batch to leveldb")
	}
	return nil
}





type FileLock struct {
	db       *leveldb.DB
	filePath string
}


func NewFileLock(filePath string) *FileLock {
	return &FileLock{
		filePath: filePath,
	}
}









func (f *FileLock) Lock() error {
	dbOpts := &opt.Options{}
	var err error
	var dirEmpty bool
	if dirEmpty, err = util.CreateDirIfMissing(f.filePath); err != nil {
		panic(fmt.Sprintf("Error creating dir if missing: %s", err))
	}
	dbOpts.ErrorIfMissing = !dirEmpty
	f.db, err = leveldb.OpenFile(f.filePath, dbOpts)
	if err != nil && err == syscall.EAGAIN {
		return errors.Errorf("lock is already acquired on file %s", f.filePath)
	}
	if err != nil {
		panic(fmt.Sprintf("Error acquiring lock on file %s: %s", f.filePath, err))
	}
	return nil
}



func (f *FileLock) Unlock() {
	if f.db == nil {
		return
	}
	if err := f.db.Close(); err != nil {
		logger.Warningf("unable to release the lock on file %s: %s", f.filePath, err)
		return
	}
	f.db = nil
}
