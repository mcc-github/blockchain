/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package leveldbhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestLevelDBHelperWriteWithoutOpen(t *testing.T) {
	env := newTestDBEnv(t, testDBPath)
	defer env.cleanup()
	db := env.db
	defer func() {
		if recover() == nil {
			t.Fatalf("A panic is expected when writing to db before opening")
		}
	}()
	db.Put([]byte("key"), []byte("value"), false)
}

func TestLevelDBHelperReadWithoutOpen(t *testing.T) {
	env := newTestDBEnv(t, testDBPath)
	defer env.cleanup()
	db := env.db
	defer func() {
		if recover() == nil {
			t.Fatalf("A panic is expected when writing to db before opening")
		}
	}()
	db.Get([]byte("key"))
}

func TestLevelDBHelper(t *testing.T) {
	env := newTestDBEnv(t, testDBPath)
	
	db := env.db

	db.Open()
	
	db.Open()
	isEmpty, err := db.isEmpty()
	assert.NoError(t, err)
	assert.True(t, isEmpty)
	db.Put([]byte("key1"), []byte("value1"), false)
	db.Put([]byte("key2"), []byte("value2"), true)
	db.Put([]byte("key3"), []byte("value3"), true)

	val, _ := db.Get([]byte("key2"))
	assert.Equal(t, "value2", string(val))

	db.Delete([]byte("key1"), false)
	db.Delete([]byte("key2"), true)

	val1, err1 := db.Get([]byte("key1"))
	assert.NoError(t, err1, "")
	assert.Equal(t, "", string(val1))

	val2, err2 := db.Get([]byte("key2"))
	assert.NoError(t, err2, "")
	assert.Equal(t, "", string(val2))

	db.Close()
	
	db.Close()

	_, err = db.isEmpty()
	require.Error(t, err)

	val3, err3 := db.Get([]byte("key3"))
	assert.Error(t, err3)

	db.Open()
	isEmpty, err = db.isEmpty()
	assert.NoError(t, err)
	assert.False(t, isEmpty)

	batch := &leveldb.Batch{}
	batch.Put([]byte("key1"), []byte("value1"))
	batch.Put([]byte("key2"), []byte("value2"))
	batch.Delete([]byte("key3"))
	db.WriteBatch(batch, true)

	val1, err1 = db.Get([]byte("key1"))
	assert.NoError(t, err1, "")
	assert.Equal(t, "value1", string(val1))

	val2, err2 = db.Get([]byte("key2"))
	assert.NoError(t, err2, "")
	assert.Equal(t, "value2", string(val2))

	val3, err3 = db.Get([]byte("key3"))
	assert.NoError(t, err3, "")
	assert.Equal(t, "", string(val3))

	keys := []string{}
	itr := db.GetIterator(nil, nil)
	for itr.Next() {
		keys = append(keys, string(itr.Key()))
	}
	assert.Equal(t, []string{"key1", "key2"}, keys)
}

func TestFileLock(t *testing.T) {
	
	fileLockPath := testDBPath + "/fileLock"
	fileLock1 := NewFileLock(fileLockPath)
	assert.Nil(t, fileLock1.db)
	assert.Equal(t, fileLock1.filePath, fileLockPath)

	
	err := fileLock1.Lock()
	assert.NoError(t, err)
	assert.NotNil(t, fileLock1.db)

	
	fileLock2 := NewFileLock(fileLockPath)
	assert.Nil(t, fileLock2.db)
	assert.Equal(t, fileLock2.filePath, fileLockPath)

	
	
	err = fileLock2.Lock()
	expectedErr := fmt.Sprintf("lock is already acquired on file %s", fileLockPath)
	assert.EqualError(t, err, expectedErr)
	assert.Nil(t, fileLock2.db)

	
	fileLock1.Unlock()
	assert.Nil(t, fileLock1.db)

	
	
	err = fileLock2.Lock()
	assert.NoError(t, err)
	assert.NotNil(t, fileLock2.db)

	
	fileLock2.Unlock()
	assert.Nil(t, fileLock1.db)

	
	fileLock2.Unlock()
	assert.Nil(t, fileLock1.db)

	
	assert.NoError(t, os.RemoveAll(fileLockPath))
}

func TestCreateDBInEmptyDir(t *testing.T) {
	assert.NoError(t, os.RemoveAll(testDBPath), "")
	assert.NoError(t, os.MkdirAll(testDBPath, 0775), "")
	db := CreateDB(&Conf{DBPath: testDBPath})
	defer db.Close()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Panic is not expected when opening db in an existing empty dir. %s", r)
		}
	}()
	db.Open()
}

func TestCreateDBInNonEmptyDir(t *testing.T) {
	assert.NoError(t, os.RemoveAll(testDBPath), "")
	assert.NoError(t, os.MkdirAll(testDBPath, 0775), "")
	file, err := os.Create(filepath.Join(testDBPath, "dummyfile.txt"))
	assert.NoError(t, err, "")
	file.Close()
	db := CreateDB(&Conf{DBPath: testDBPath})
	defer db.Close()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("A panic is expected when opening db in an existing non-empty dir. %s", r)
		}
	}()
	db.Open()
}
