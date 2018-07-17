package ioutils 

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)




func NewAtomicFileWriter(filename string, perm os.FileMode) (io.WriteCloser, error) {
	f, err := ioutil.TempFile(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil {
		return nil, err
	}

	abspath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	return &atomicFileWriter{
		f:    f,
		fn:   abspath,
		perm: perm,
	}, nil
}


func AtomicWriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := NewAtomicFileWriter(filename, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
		f.(*atomicFileWriter).writeErr = err
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

type atomicFileWriter struct {
	f        *os.File
	fn       string
	writeErr error
	perm     os.FileMode
}

func (w *atomicFileWriter) Write(dt []byte) (int, error) {
	n, err := w.f.Write(dt)
	if err != nil {
		w.writeErr = err
	}
	return n, err
}

func (w *atomicFileWriter) Close() (retErr error) {
	defer func() {
		if retErr != nil || w.writeErr != nil {
			os.Remove(w.f.Name())
		}
	}()
	if err := w.f.Sync(); err != nil {
		w.f.Close()
		return err
	}
	if err := w.f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(w.f.Name(), w.perm); err != nil {
		return err
	}
	if w.writeErr == nil {
		return os.Rename(w.f.Name(), w.fn)
	}
	return nil
}




type AtomicWriteSet struct {
	root string
}






func NewAtomicWriteSet(tmpDir string) (*AtomicWriteSet, error) {
	td, err := ioutil.TempDir(tmpDir, "write-set-")
	if err != nil {
		return nil, err
	}

	return &AtomicWriteSet{
		root: td,
	}, nil
}



func (ws *AtomicWriteSet) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := ws.FileWriter(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

type syncFileCloser struct {
	*os.File
}

func (w syncFileCloser) Close() error {
	err := w.File.Sync()
	if err1 := w.File.Close(); err == nil {
		err = err1
	}
	return err
}



func (ws *AtomicWriteSet) FileWriter(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	f, err := os.OpenFile(filepath.Join(ws.root, name), flag, perm)
	if err != nil {
		return nil, err
	}
	return syncFileCloser{f}, nil
}



func (ws *AtomicWriteSet) Cancel() error {
	return os.RemoveAll(ws.root)
}




func (ws *AtomicWriteSet) Commit(target string) error {
	return os.Rename(ws.root, target)
}


func (ws *AtomicWriteSet) String() string {
	return ws.root
}
