













package wal

import (
	"errors"
	"fmt"
	"strings"

	"go.etcd.io/etcd/pkg/fileutil"

	"go.uber.org/zap"
)

var errBadWALName = errors.New("bad wal name")


func Exist(dir string) bool {
	names, err := fileutil.ReadDir(dir, fileutil.WithExt(".wal"))
	if err != nil {
		return false
	}
	return len(names) != 0
}




func searchIndex(lg *zap.Logger, names []string, index uint64) (int, bool) {
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		_, curIndex, err := parseWALName(name)
		if err != nil {
			if lg != nil {
				lg.Panic("failed to parse WAL file name", zap.String("path", name), zap.Error(err))
			} else {
				plog.Panicf("parse correct name should never fail: %v", err)
			}
		}
		if index >= curIndex {
			return i, true
		}
	}
	return -1, false
}



func isValidSeq(lg *zap.Logger, names []string) bool {
	var lastSeq uint64
	for _, name := range names {
		curSeq, _, err := parseWALName(name)
		if err != nil {
			if lg != nil {
				lg.Panic("failed to parse WAL file name", zap.String("path", name), zap.Error(err))
			} else {
				plog.Panicf("parse correct name should never fail: %v", err)
			}
		}
		if lastSeq != 0 && lastSeq != curSeq-1 {
			return false
		}
		lastSeq = curSeq
	}
	return true
}

func readWALNames(lg *zap.Logger, dirpath string) ([]string, error) {
	names, err := fileutil.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}
	wnames := checkWalNames(lg, names)
	if len(wnames) == 0 {
		return nil, ErrFileNotFound
	}
	return wnames, nil
}

func checkWalNames(lg *zap.Logger, names []string) []string {
	wnames := make([]string, 0)
	for _, name := range names {
		if _, _, err := parseWALName(name); err != nil {
			
			if !strings.HasSuffix(name, ".tmp") {
				if lg != nil {
					lg.Warn(
						"ignored file in WAL directory",
						zap.String("path", name),
					)
				} else {
					plog.Warningf("ignored file %v in wal", name)
				}
			}
			continue
		}
		wnames = append(wnames, name)
	}
	return wnames
}

func parseWALName(str string) (seq, index uint64, err error) {
	if !strings.HasSuffix(str, ".wal") {
		return 0, 0, errBadWALName
	}
	_, err = fmt.Sscanf(str, "%016x-%016x.wal", &seq, &index)
	return seq, index, err
}

func walName(seq, index uint64) string {
	return fmt.Sprintf("%016x-%016x.wal", seq, index)
}
