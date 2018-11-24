













package snap

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreos/etcd/pkg/fileutil"
)

var ErrNoDBSnapshot = errors.New("snap: snapshot file doesn't exist")



func (s *Snapshotter) SaveDBFrom(r io.Reader, id uint64) (int64, error) {
	f, err := ioutil.TempFile(s.dir, "tmp")
	if err != nil {
		return 0, err
	}
	var n int64
	n, err = io.Copy(f, r)
	if err == nil {
		err = fileutil.Fsync(f)
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

	plog.Infof("saved database snapshot to disk [total bytes: %d]", n)

	return n, nil
}



func (s *Snapshotter) DBFilePath(id uint64) (string, error) {
	if _, err := fileutil.ReadDir(s.dir); err != nil {
		return "", err
	}
	if fn := s.dbFilePath(id); fileutil.Exist(fn) {
		return fn, nil
	}
	return "", ErrNoDBSnapshot
}

func (s *Snapshotter) dbFilePath(id uint64) string {
	return filepath.Join(s.dir, fmt.Sprintf("%016x.snap.db", id))
}
