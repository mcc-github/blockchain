/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/pkg/fileutil"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/wal"
	"go.uber.org/zap"
)

var (
	err    error
	logger *flogging.FabricLogger

	dataDir, walDir, snapDir string

	ram   *raft.MemoryStorage
	store *RaftStorage
)

func setup(t *testing.T) {
	logger = flogging.NewFabricLogger(zap.NewExample())
	ram = raft.NewMemoryStorage()
	dataDir, err = ioutil.TempDir("", "etcdraft-")
	assert.NoError(t, err)
	walDir, snapDir = path.Join(dataDir, "wal"), path.Join(dataDir, "snapshot")
	store, err = CreateStorage(logger, walDir, snapDir, ram)
	assert.NoError(t, err)
}

func clean(t *testing.T) {
	err = store.Close()
	assert.NoError(t, err)
	err = os.RemoveAll(dataDir)
	assert.NoError(t, err)
}

func fileCount(files []string, suffix string) (c int) {
	for _, f := range files {
		if strings.HasSuffix(f, suffix) {
			c++
		}
	}
	return
}

func assertFileCount(t *testing.T, wal, snap int) {
	files, err := fileutil.ReadDir(walDir)
	assert.NoError(t, err)
	assert.Equal(t, wal, fileCount(files, ".wal"), "WAL file count mismatch")

	files, err = fileutil.ReadDir(snapDir)
	assert.NoError(t, err)
	assert.Equal(t, snap, fileCount(files, ".snap"), "Snap file count mismatch")
}

func TestOpenWAL(t *testing.T) {
	t.Run("Last WAL file is broken", func(t *testing.T) {
		setup(t)
		defer clean(t)

		
		for i := 0; i < 10; i++ {
			store.Store(
				[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 10)}},
				raftpb.HardState{},
				raftpb.Snapshot{},
			)
		}
		assertFileCount(t, 1, 0)
		lasti, _ := store.ram.LastIndex() 

		
		err = store.Close()
		assert.NoError(t, err)

		
		w := func() string {
			files, err := fileutil.ReadDir(walDir)
			assert.NoError(t, err)
			for _, f := range files {
				if strings.HasSuffix(f, ".wal") {
					return path.Join(walDir, f)
				}
			}
			t.FailNow()
			return ""
		}()
		err = os.Truncate(w, 200)
		assert.NoError(t, err)

		
		ram = raft.NewMemoryStorage()
		store, err = CreateStorage(logger, walDir, snapDir, ram)
		require.NoError(t, err)
		lastI, _ := store.ram.LastIndex()
		assert.True(t, lastI > 0)     
		assert.True(t, lasti > lastI) 
	})
}

func TestTakeSnapshot(t *testing.T) {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	t.Run("Good", func(t *testing.T) {
		t.Run("MaxSnapshotFiles==1", func(t *testing.T) {
			backup := MaxSnapshotFiles
			MaxSnapshotFiles = 1
			defer func() { MaxSnapshotFiles = backup }()

			setup(t)
			defer clean(t)

			
			
			
			oldSegmentSizeBytes := wal.SegmentSizeBytes
			wal.SegmentSizeBytes = 10
			defer func() {
				wal.SegmentSizeBytes = oldSegmentSizeBytes
			}()

			
			for i := 0; i < 10; i++ {
				store.Store(
					[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
					raftpb.HardState{},
					raftpb.Snapshot{},
				)
			}

			assertFileCount(t, 11, 0)

			err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			
			
			assertFileCount(t, 9, 1)

			err = store.TakeSnapshot(uint64(5), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			
			
			assertFileCount(t, 7, 1)

			t.Logf("Close the storage and create a new one based on existing files")

			err = store.Close()
			assert.NoError(t, err)
			ram := raft.NewMemoryStorage()
			store, err = CreateStorage(logger, walDir, snapDir, ram)
			assert.NoError(t, err)

			err = store.TakeSnapshot(uint64(7), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			
			
			assertFileCount(t, 5, 1)

			err = store.TakeSnapshot(uint64(9), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			
			
			assertFileCount(t, 3, 1)
		})

		t.Run("MaxSnapshotFiles==2", func(t *testing.T) {
			backup := MaxSnapshotFiles
			MaxSnapshotFiles = 2
			defer func() { MaxSnapshotFiles = backup }()

			setup(t)
			defer clean(t)

			
			
			
			oldSegmentSizeBytes := wal.SegmentSizeBytes
			wal.SegmentSizeBytes = 10
			defer func() {
				wal.SegmentSizeBytes = oldSegmentSizeBytes
			}()

			
			for i := 0; i < 10; i++ {
				store.Store(
					[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
					raftpb.HardState{},
					raftpb.Snapshot{},
				)
			}

			assertFileCount(t, 11, 0)

			
			err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 11, 1)

			
			
			err = store.TakeSnapshot(uint64(5), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			t.Logf("Close the storage and create a new one based on existing files")

			err = store.Close()
			assert.NoError(t, err)
			ram := raft.NewMemoryStorage()
			store, err = CreateStorage(logger, walDir, snapDir, ram)
			assert.NoError(t, err)

			
			
			err = store.TakeSnapshot(uint64(7), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 7, 2)

			
			
			err = store.TakeSnapshot(uint64(9), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 5, 2)
		})
	})

	t.Run("Bad", func(t *testing.T) {

		t.Run("MaxSnapshotFiles==2", func(t *testing.T) {
			
			

			backup := MaxSnapshotFiles
			MaxSnapshotFiles = 2
			defer func() { MaxSnapshotFiles = backup }()

			setup(t)
			defer clean(t)

			
			
			
			oldSegmentSizeBytes := wal.SegmentSizeBytes
			wal.SegmentSizeBytes = 10
			defer func() {
				wal.SegmentSizeBytes = oldSegmentSizeBytes
			}()

			
			for i := 0; i < 10; i++ {
				store.Store(
					[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
					raftpb.HardState{},
					raftpb.Snapshot{},
				)
			}

			assertFileCount(t, 11, 0)

			
			err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 11, 1)

			
			
			err = store.TakeSnapshot(uint64(5), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			d, err := os.Open(snapDir)
			assert.NoError(t, err)
			defer d.Close()
			names, err := d.Readdirnames(-1)
			assert.NoError(t, err)
			sort.Sort(sort.Reverse(sort.StringSlice(names)))

			corrupted := filepath.Join(snapDir, names[0])
			t.Logf("Corrupt latest snapshot file: %s", corrupted)
			f, err := os.OpenFile(corrupted, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			assert.NoError(t, err)
			_, err = f.WriteString("Corrupted Snapshot")
			assert.NoError(t, err)
			f.Close()

			
			_ = ListSnapshots(logger, snapDir)
			assertFileCount(t, 9, 1)

			
			broken := corrupted + ".broken"
			err = os.Rename(broken, corrupted)
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			t.Logf("Close the storage and create a new one based on existing files")

			err = store.Close()
			assert.NoError(t, err)
			ram := raft.NewMemoryStorage()
			store, err = CreateStorage(logger, walDir, snapDir, ram)
			assert.NoError(t, err)

			
			assertFileCount(t, 9, 1)

			files, err := fileutil.ReadDir(snapDir)
			assert.NoError(t, err)
			assert.Equal(t, 1, fileCount(files, ".broken"))

			err = store.TakeSnapshot(uint64(7), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 9, 2)

			err = store.TakeSnapshot(uint64(9), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
			assert.NoError(t, err)
			assertFileCount(t, 5, 2)
		})
	})
}

func TestApplyOutOfDateSnapshot(t *testing.T) {
	t.Run("Apply out of date snapshot", func(t *testing.T) {
		setup(t)
		defer clean(t)

		
		
		
		oldSegmentSizeBytes := wal.SegmentSizeBytes
		wal.SegmentSizeBytes = 10
		defer func() {
			wal.SegmentSizeBytes = oldSegmentSizeBytes
		}()

		
		for i := 0; i < 10; i++ {
			store.Store(
				[]raftpb.Entry{{Index: uint64(i), Data: make([]byte, 100)}},
				raftpb.HardState{},
				raftpb.Snapshot{},
			)
		}
		assertFileCount(t, 11, 0)

		err = store.TakeSnapshot(uint64(3), raftpb.ConfState{Nodes: []uint64{1}}, make([]byte, 10))
		assert.NoError(t, err)
		assertFileCount(t, 11, 1)

		snapshot := store.Snapshot()
		assert.NotNil(t, snapshot)

		
		store.ApplySnapshot(snapshot)

		
		err := store.Store(
			[]raftpb.Entry{{Index: uint64(10), Data: make([]byte, 100)}},
			raftpb.HardState{},
			snapshot,
		)
		assert.NoError(t, err)
		assertFileCount(t, 12, 1)
	})
}
