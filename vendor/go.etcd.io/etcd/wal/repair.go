













package wal

import (
	"io"
	"os"
	"path/filepath"

	"go.etcd.io/etcd/pkg/fileutil"
	"go.etcd.io/etcd/wal/walpb"

	"go.uber.org/zap"
)



func Repair(lg *zap.Logger, dirpath string) bool {
	f, err := openLast(lg, dirpath)
	if err != nil {
		return false
	}
	defer f.Close()

	if lg != nil {
		lg.Info("repairing", zap.String("path", f.Name()))
	} else {
		plog.Noticef("repairing %v", f.Name())
	}

	rec := &walpb.Record{}
	decoder := newDecoder(f)
	for {
		lastOffset := decoder.lastOffset()
		err := decoder.decode(rec)
		switch err {
		case nil:
			
			switch rec.Type {
			case crcType:
				crc := decoder.crc.Sum32()
				
				
				if crc != 0 && rec.Validate(crc) != nil {
					return false
				}
				decoder.updateCRC(rec.Crc)
			}
			continue

		case io.EOF:
			if lg != nil {
				lg.Info("repaired", zap.String("path", f.Name()), zap.Error(io.EOF))
			}
			return true

		case io.ErrUnexpectedEOF:
			bf, bferr := os.Create(f.Name() + ".broken")
			if bferr != nil {
				if lg != nil {
					lg.Warn("failed to create backup file", zap.String("path", f.Name()+".broken"), zap.Error(bferr))
				} else {
					plog.Errorf("could not repair %v, failed to create backup file", f.Name())
				}
				return false
			}
			defer bf.Close()

			if _, err = f.Seek(0, io.SeekStart); err != nil {
				if lg != nil {
					lg.Warn("failed to read file", zap.String("path", f.Name()), zap.Error(err))
				} else {
					plog.Errorf("could not repair %v, failed to read file", f.Name())
				}
				return false
			}

			if _, err = io.Copy(bf, f); err != nil {
				if lg != nil {
					lg.Warn("failed to copy", zap.String("from", f.Name()+".broken"), zap.String("to", f.Name()), zap.Error(err))
				} else {
					plog.Errorf("could not repair %v, failed to copy file", f.Name())
				}
				return false
			}

			if err = f.Truncate(lastOffset); err != nil {
				if lg != nil {
					lg.Warn("failed to truncate", zap.String("path", f.Name()), zap.Error(err))
				} else {
					plog.Errorf("could not repair %v, failed to truncate file", f.Name())
				}
				return false
			}

			if err = fileutil.Fsync(f.File); err != nil {
				if lg != nil {
					lg.Warn("failed to fsync", zap.String("path", f.Name()), zap.Error(err))
				} else {
					plog.Errorf("could not repair %v, failed to sync file", f.Name())
				}
				return false
			}

			if lg != nil {
				lg.Info("repaired", zap.String("path", f.Name()), zap.Error(io.ErrUnexpectedEOF))
			}
			return true

		default:
			if lg != nil {
				lg.Warn("failed to repair", zap.String("path", f.Name()), zap.Error(err))
			} else {
				plog.Errorf("could not repair error (%v)", err)
			}
			return false
		}
	}
}


func openLast(lg *zap.Logger, dirpath string) (*fileutil.LockedFile, error) {
	names, err := readWALNames(lg, dirpath)
	if err != nil {
		return nil, err
	}
	last := filepath.Join(dirpath, names[len(names)-1])
	return fileutil.LockFile(last, os.O_RDWR, fileutil.PrivateFileMode)
}
