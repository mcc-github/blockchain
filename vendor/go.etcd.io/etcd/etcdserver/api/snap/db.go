













package snap

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/etcd/pkg/fileutil"

	humanize "github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

var ErrNoDBSnapshot = errors.New("snap: snapshot file doesn't exist")



func (s *Snapshotter) SaveDBFrom(r io.Reader, id uint64) (int64, error) {
	start := time.Now()

	f, err := ioutil.TempFile(s.dir, "tmp")
	if err != nil {
		return 0, err
	}
	var n int64
	n, err = io.Copy(f, r)
	if err == nil {
		fsyncStart := time.Now()
		err = fileutil.Fsync(f)
		snapDBFsyncSec.Observe(time.Since(fsyncStart).Seconds())
	}
	f.Close()
	if err != nil {
		os.Remove(f.Name())
		return n, err
	}
	fn := s.dbFilePath(id)
	if fileutil.Exist(fn) {
		os.Remove(f.Name())
		return n, nil
	}
	err = os.Rename(f.Name(), fn)
	if err != nil {
		os.Remove(f.Name())
		return n, err
	}

	if s.lg != nil {
		s.lg.Info(
			"saved database snapshot to disk",
			zap.String("path", fn),
			zap.Int64("bytes", n),
			zap.String("size", humanize.Bytes(uint64(n))),
		)
	} else {
		plog.Infof("saved database snapshot to disk [total bytes: %d]", n)
	}

	snapDBSaveSec.Observe(time.Since(start).Seconds())
	return n, nil
}



func (s *Snapshotter) DBFilePath(id uint64) (string, error) {
	if _, err := fileutil.ReadDir(s.dir); err != nil {
		return "", err
	}
	fn := s.dbFilePath(id)
	if fileutil.Exist(fn) {
		return fn, nil
	}
	if s.lg != nil {
		s.lg.Warn(
			"failed to find [SNAPSHOT-INDEX].snap.db",
			zap.Uint64("snapshot-index", id),
			zap.String("snapshot-file-path", fn),
			zap.Error(ErrNoDBSnapshot),
		)
	}
	return "", ErrNoDBSnapshot
}

func (s *Snapshotter) dbFilePath(id uint64) string {
	return filepath.Join(s.dir, fmt.Sprintf("%016x.snap.db", id))
}
