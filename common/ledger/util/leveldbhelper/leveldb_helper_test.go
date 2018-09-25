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

package leveldbhelper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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

	val3, err3 := db.Get([]byte("key3"))
	assert.Error(t, err3)

	db.Open()
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

func TestCreateDBInEmptyDir(t *testing.T) {
	assert.NoError(t, os.RemoveAll(testDBPath), "")
	assert.NoError(t, os.MkdirAll(testDBPath, 0775), "")
	db := CreateDB(&Conf{testDBPath})
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
	db := CreateDB(&Conf{testDBPath})
	defer db.Close()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("A panic is expected when opening db in an existing non-empty dir. %s", r)
		}
	}()
	db.Open()
}
