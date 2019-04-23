













package fileutil

import (
	"os"
	"path/filepath"
	"sort"
)


type ReadDirOp struct {
	ext string
}


type ReadDirOption func(*ReadDirOp)



func WithExt(ext string) ReadDirOption {
	return func(op *ReadDirOp) { op.ext = ext }
}

func (op *ReadDirOp) applyOpts(opts []ReadDirOption) {
	for _, opt := range opts {
		opt(op)
	}
}


func ReadDir(d string, opts ...ReadDirOption) ([]string, error) {
	op := &ReadDirOp{}
	op.applyOpts(opts)

	dir, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)

	if op.ext != "" {
		tss := make([]string, 0)
		for _, v := range names {
			if filepath.Ext(v) == op.ext {
				tss = append(tss, v)
			}
		}
		names = tss
	}
	return names, nil
}
