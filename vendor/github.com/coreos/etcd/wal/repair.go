













package wal

import (
	"io"
	"os"
	"path/filepath"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/coreos/etcd/wal/walpb"
)



func Repair(dirpath string) bool {
	f, err := openLast(dirpath)
	if err != nil {
		return false
	}
	defer f.Close()

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
			return true
		case io.ErrUnexpectedEOF:
			plog.Noticef("repairing %v", f.Name())
			bf, bferr := os.Create(f.Name() + ".broken")
			if bferr != nil {
				plog.Errorf("could not repair %v, failed to create backup file", f.Name())
				return false
			}
			defer bf.Close()

			if _, err = f.Seek(0, io.SeekStart); err != nil {
				plog.Errorf("could not repair %v, failed to read file", f.Name())
				return false
			}

			if _, err = io.Copy(bf, f); err != nil {
				plog.Errorf("could not repair %v, failed to copy file", f.Name())
				return false
			}

			if err = f.Truncate(int64(lastOffset)); err != nil {
				plog.Errorf("could not repair %v, failed to truncate file", f.Name())
				return false
			}
			if err = fileutil.Fsync(f.File); err != nil {
				plog.Errorf("could not repair %v, failed to sync file", f.Name())
				return false
			}
			return true
		default:
			plog.Errorf("could not repair error (%v)", err)
			return false
		}
	}
}


func openLast(dirpath string) (*fileutil.LockedFile, error) {
	names, err := readWalNames(dirpath)
	if err != nil {
		return nil, err
	}
	last := filepath.Join(dirpath, names[len(names)-1])
	return fileutil.LockFile(last, os.O_RDWR, fileutil.PrivateFileMode)
}
