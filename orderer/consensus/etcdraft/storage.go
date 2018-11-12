/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package etcdraft

import (
	"os"

	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal"
	"github.com/coreos/etcd/wal/walpb"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/pkg/errors"
)





type MemoryStorage interface {
	raft.Storage
	Append(entries []raftpb.Entry) error
	SetHardState(st raftpb.HardState) error
}


type RaftStorage struct {
	ram MemoryStorage
	wal *wal.WAL
}



func Restore(lg *flogging.FabricLogger, applied uint64, walDir string, ram MemoryStorage) (*RaftStorage, bool, error) {
	hasWAL := wal.Exist(walDir)
	if !hasWAL && applied != 0 {
		return nil, hasWAL, errors.Errorf("applied index is not zero but no WAL data found")
	}

	if !hasWAL {
		lg.Infof("No WAL data found, creating new WAL at path '%s'", walDir)
		
		
		w, err := wal.Create(walDir, nil)
		if err == os.ErrExist {
			lg.Fatalf("programming error, we've just checked that WAL does not exist")
		}

		if err != nil {
			return nil, hasWAL, errors.Errorf("failed to initialize WAL: %s", err)
		}

		if err = w.Close(); err != nil {
			return nil, hasWAL, errors.Errorf("failed to close the WAL just created: %s", err)
		}
	} else {
		lg.Infof("Found WAL data at path '%s', replaying it", walDir)
	}

	w, err := wal.Open(walDir, walpb.Snapshot{})
	if err != nil {
		return nil, hasWAL, errors.Errorf("failed to open existing WAL: %s", err)
	}

	_, st, ents, err := w.ReadAll() 
	if err != nil {
		return nil, hasWAL, errors.Errorf("failed to read WAL: %s", err)
	}

	lg.Debugf("Setting HardState to {Term: %d, Commit: %d}", st.Term, st.Commit)
	ram.SetHardState(st) 

	lg.Debugf("Appending %d entries to memory storage", len(ents))
	ram.Append(ents) 

	return &RaftStorage{ram: ram, wal: w}, hasWAL, nil
}


func (rs *RaftStorage) Store(entries []raftpb.Entry, hardstate raftpb.HardState) error {
	if err := rs.ram.Append(entries); err != nil {
		return err
	}

	if err := rs.wal.Save(hardstate, entries); err != nil {
		return err
	}

	return nil
}


func (rs *RaftStorage) Close() error {
	if err := rs.wal.Close(); err != nil {
		return err
	}

	return nil
}
