/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"os"

	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/snap"
	"github.com/coreos/etcd/wal"
	"github.com/coreos/etcd/wal/walpb"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)





type MemoryStorage interface {
	raft.Storage
	Append(entries []raftpb.Entry) error
	SetHardState(st raftpb.HardState) error
	CreateSnapshot(i uint64, cs *raftpb.ConfState, data []byte) (raftpb.Snapshot, error)
	Compact(compactIndex uint64) error
	ApplySnapshot(snap raftpb.Snapshot) error
}


type RaftStorage struct {
	SnapshotCatchUpEntries uint64

	lg *flogging.FabricLogger

	ram  MemoryStorage
	wal  *wal.WAL
	snap *snap.Snapshotter
}



func CreateStorage(
	lg *flogging.FabricLogger,
	applied uint64,
	walDir string,
	snapDir string,
	ram MemoryStorage,
) (*RaftStorage, error) {

	sn, err := createSnapshotter(snapDir)
	if err != nil {
		return nil, err
	}

	snapshot, err := sn.Load()
	if err != nil {
		if err == snap.ErrNoSnapshot {
			lg.Debugf("No snapshot found at %s", snapDir)
		} else {
			return nil, errors.Errorf("failed to load snapshot: %s", err)
		}
	} else {
		
		lg.Debugf("Loaded snapshot at Term %d and Index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
	}

	w, err := createWAL(lg, walDir, applied, snapshot)
	if err != nil {
		return nil, err
	}

	_, st, ents, err := w.ReadAll()
	if err != nil {
		return nil, errors.Errorf("failed to read WAL: %s", err)
	}

	if snapshot != nil {
		lg.Debugf("Applying snapshot to raft MemoryStorage")
		if err := ram.ApplySnapshot(*snapshot); err != nil {
			return nil, errors.Errorf("Failed to apply snapshot to memory: %s", err)
		}
	}

	lg.Debugf("Setting HardState to {Term: %d, Commit: %d}", st.Term, st.Commit)
	ram.SetHardState(st) 

	lg.Debugf("Appending %d entries to memory storage", len(ents))
	ram.Append(ents) 

	return &RaftStorage{lg: lg, ram: ram, wal: w, snap: sn}, nil
}

func createSnapshotter(snapDir string) (*snap.Snapshotter, error) {
	if err := os.MkdirAll(snapDir, os.ModePerm); err != nil {
		return nil, errors.Errorf("failed to mkdir '%s' for snapshot: %s", snapDir, err)
	}

	return snap.New(snapDir), nil

}

func createWAL(lg *flogging.FabricLogger, walDir string, applied uint64, snapshot *raftpb.Snapshot) (*wal.WAL, error) {
	hasWAL := wal.Exist(walDir)
	if !hasWAL && applied != 0 {
		return nil, errors.Errorf("applied index is not zero but no WAL data found")
	}

	if !hasWAL {
		lg.Infof("No WAL data found, creating new WAL at path '%s'", walDir)
		
		
		w, err := wal.Create(walDir, nil)
		if err == os.ErrExist {
			lg.Fatalf("programming error, we've just checked that WAL does not exist")
		}

		if err != nil {
			return nil, errors.Errorf("failed to initialize WAL: %s", err)
		}

		if err = w.Close(); err != nil {
			return nil, errors.Errorf("failed to close the WAL just created: %s", err)
		}
	} else {
		lg.Infof("Found WAL data at path '%s', replaying it", walDir)
	}

	walsnap := walpb.Snapshot{}
	if snapshot != nil {
		walsnap.Index, walsnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}

	lg.Debugf("Loading WAL at Term %d and Index %d", walsnap.Term, walsnap.Index)
	w, err := wal.Open(walDir, walsnap)
	if err != nil {
		return nil, errors.Errorf("failed to open existing WAL: %s", err)
	}

	return w, nil
}


func (rs *RaftStorage) Snapshot() raftpb.Snapshot {
	sn, _ := rs.ram.Snapshot() 
	return sn
}


func (rs *RaftStorage) Store(entries []raftpb.Entry, hardstate raftpb.HardState, snapshot raftpb.Snapshot) error {
	if err := rs.wal.Save(hardstate, entries); err != nil {
		return err
	}

	if !raft.IsEmptySnap(snapshot) {
		if err := rs.saveSnap(snapshot); err != nil {
			return err
		}

		if err := rs.ram.ApplySnapshot(snapshot); err != nil {
			if err == raft.ErrSnapOutOfDate {
				rs.lg.Warnf("Attempted to apply out-of-date snapshot at Term %d and Index %d",
					snapshot.Metadata.Term, snapshot.Metadata.Index)
			} else {
				rs.lg.Fatalf("Unexpected programming error: %s", err)
			}
		}
	}

	if err := rs.ram.Append(entries); err != nil {
		return err
	}

	return nil
}

func (rs *RaftStorage) saveSnap(snap raftpb.Snapshot) error {
	
	
	
	walsnap := walpb.Snapshot{
		Index: snap.Metadata.Index,
		Term:  snap.Metadata.Term,
	}

	rs.lg.Debugf("Saving snapshot to WAL")
	if err := rs.wal.SaveSnapshot(walsnap); err != nil {
		return errors.Errorf("failed to save snapshot to WAL: %s", err)
	}

	rs.lg.Debugf("Saving snapshot to disk")
	if err := rs.snap.SaveSnap(snap); err != nil {
		return errors.Errorf("failed to save snapshot to disk: %s", err)
	}

	rs.lg.Debugf("Releasing lock to wal files prior to %d", snap.Metadata.Index)
	if err := rs.wal.ReleaseLockTo(snap.Metadata.Index); err != nil {
		return err
	}

	return nil
}


func (rs *RaftStorage) TakeSnapshot(i uint64, cs *raftpb.ConfState, data []byte) error {
	rs.lg.Debugf("Creating snapshot at index %d from MemoryStorage", i)
	snap, err := rs.ram.CreateSnapshot(i, cs, data)
	if err != nil {
		return errors.Errorf("failed to create snapshot from MemoryStorage: %s", err)
	}

	if err = rs.saveSnap(snap); err != nil {
		return err
	}

	
	if i > rs.SnapshotCatchUpEntries {
		compacti := i - rs.SnapshotCatchUpEntries
		rs.lg.Debugf("Purging in-memory raft entries prior to %d", compacti)
		if err = rs.ram.Compact(compacti); err != nil {
			if err == raft.ErrCompacted {
				rs.lg.Warnf("Raft entries prior to %d are already purged", compacti)
			} else {
				rs.lg.Fatalf("Failed to purg raft entries: %s", err)
			}
		}
	}

	rs.lg.Infof("Snapshot is taken at index %d", i)
	return nil
}


func (rs *RaftStorage) ApplySnapshot(snap raftpb.Snapshot) {
	if err := rs.ram.ApplySnapshot(snap); err != nil {
		if err == raft.ErrSnapOutOfDate {
			rs.lg.Warnf("Attempted to apply out-of-date snapshot at Term %d and Index %d",
				snap.Metadata.Term, snap.Metadata.Index)
		} else {
			rs.lg.Fatalf("Unexpected programming error: %s", err)
		}
	}
}


func (rs *RaftStorage) Close() error {
	if err := rs.wal.Close(); err != nil {
		return err
	}

	return nil
}
