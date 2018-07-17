package archive 

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"unsafe"

	"github.com/docker/docker/pkg/system"
	"golang.org/x/sys/unix"
)









type walker struct {
	dir1  string
	dir2  string
	root1 *FileInfo
	root2 *FileInfo
}







func collectFileInfoForChanges(dir1, dir2 string) (*FileInfo, *FileInfo, error) {
	w := &walker{
		dir1:  dir1,
		dir2:  dir2,
		root1: newRootFileInfo(),
		root2: newRootFileInfo(),
	}

	i1, err := os.Lstat(w.dir1)
	if err != nil {
		return nil, nil, err
	}
	i2, err := os.Lstat(w.dir2)
	if err != nil {
		return nil, nil, err
	}

	if err := w.walk("/", i1, i2); err != nil {
		return nil, nil, err
	}

	return w.root1, w.root2, nil
}



func walkchunk(path string, fi os.FileInfo, dir string, root *FileInfo) error {
	if fi == nil {
		return nil
	}
	parent := root.LookUp(filepath.Dir(path))
	if parent == nil {
		return fmt.Errorf("walkchunk: Unexpectedly no parent for %s", path)
	}
	info := &FileInfo{
		name:     filepath.Base(path),
		children: make(map[string]*FileInfo),
		parent:   parent,
	}
	cpath := filepath.Join(dir, path)
	stat, err := system.FromStatT(fi.Sys().(*syscall.Stat_t))
	if err != nil {
		return err
	}
	info.stat = stat
	info.capability, _ = system.Lgetxattr(cpath, "security.capability") 
	parent.children[info.name] = info
	return nil
}



func (w *walker) walk(path string, i1, i2 os.FileInfo) (err error) {
	
	
	if path != "/" {
		if err := walkchunk(path, i1, w.dir1, w.root1); err != nil {
			return err
		}
		if err := walkchunk(path, i2, w.dir2, w.root2); err != nil {
			return err
		}
	}

	is1Dir := i1 != nil && i1.IsDir()
	is2Dir := i2 != nil && i2.IsDir()

	sameDevice := false
	if i1 != nil && i2 != nil {
		si1 := i1.Sys().(*syscall.Stat_t)
		si2 := i2.Sys().(*syscall.Stat_t)
		if si1.Dev == si2.Dev {
			sameDevice = true
		}
	}

	
	if !is1Dir && !is2Dir {
		return nil
	}

	
	var names1, names2 []nameIno
	if is1Dir {
		names1, err = readdirnames(filepath.Join(w.dir1, path)) 
		if err != nil {
			return err
		}
	}
	if is2Dir {
		names2, err = readdirnames(filepath.Join(w.dir2, path)) 
		if err != nil {
			return err
		}
	}

	
	
	
	var names []string
	ix1 := 0
	ix2 := 0

	for {
		if ix1 >= len(names1) {
			break
		}
		if ix2 >= len(names2) {
			break
		}

		ni1 := names1[ix1]
		ni2 := names2[ix2]

		switch bytes.Compare([]byte(ni1.name), []byte(ni2.name)) {
		case -1: 
			
			names = append(names, ni1.name)
			ix1++
		case 0: 
			if ni1.ino != ni2.ino || !sameDevice {
				names = append(names, ni1.name)
			}
			ix1++
			ix2++
		case 1: 
			
			names = append(names, ni2.name)
			ix2++
		}
	}
	for ix1 < len(names1) {
		names = append(names, names1[ix1].name)
		ix1++
	}
	for ix2 < len(names2) {
		names = append(names, names2[ix2].name)
		ix2++
	}

	
	
	for _, name := range names {
		fname := filepath.Join(path, name)
		var cInfo1, cInfo2 os.FileInfo
		if is1Dir {
			cInfo1, err = os.Lstat(filepath.Join(w.dir1, fname)) 
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		if is2Dir {
			cInfo2, err = os.Lstat(filepath.Join(w.dir2, fname)) 
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		if err = w.walk(fname, cInfo1, cInfo2); err != nil {
			return err
		}
	}
	return nil
}


type nameIno struct {
	name string
	ino  uint64
}

type nameInoSlice []nameIno

func (s nameInoSlice) Len() int           { return len(s) }
func (s nameInoSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s nameInoSlice) Less(i, j int) bool { return s[i].name < s[j].name }





func readdirnames(dirname string) (names []nameIno, err error) {
	var (
		size = 100
		buf  = make([]byte, 4096)
		nbuf int
		bufp int
		nb   int
	)

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names = make([]nameIno, 0, size) 
	for {
		
		if bufp >= nbuf {
			bufp = 0
			nbuf, err = unix.ReadDirent(int(f.Fd()), buf) 
			if nbuf < 0 {
				nbuf = 0
			}
			if err != nil {
				return nil, os.NewSyscallError("readdirent", err)
			}
			if nbuf <= 0 {
				break 
			}
		}

		
		nb, names = parseDirent(buf[bufp:nbuf], names)
		bufp += nb
	}

	sl := nameInoSlice(names)
	sort.Sort(sl)
	return sl, nil
}



func parseDirent(buf []byte, names []nameIno) (consumed int, newnames []nameIno) {
	origlen := len(buf)
	for len(buf) > 0 {
		dirent := (*unix.Dirent)(unsafe.Pointer(&buf[0]))
		buf = buf[dirent.Reclen:]
		if dirent.Ino == 0 { 
			continue
		}
		bytes := (*[10000]byte)(unsafe.Pointer(&dirent.Name[0]))
		var name = string(bytes[0:clen(bytes[:])])
		if name == "." || name == ".." { 
			continue
		}
		names = append(names, nameIno{name, dirent.Ino})
	}
	return origlen - len(buf), names
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}



func OverlayChanges(layers []string, rw string) ([]Change, error) {
	return changes(layers, rw, overlayDeletedFile, nil)
}

func overlayDeletedFile(root, path string, fi os.FileInfo) (string, error) {
	if fi.Mode()&os.ModeCharDevice != 0 {
		s := fi.Sys().(*syscall.Stat_t)
		if unix.Major(uint64(s.Rdev)) == 0 && unix.Minor(uint64(s.Rdev)) == 0 { 
			return path, nil
		}
	}
	if fi.Mode()&os.ModeDir != 0 {
		opaque, err := system.Lgetxattr(filepath.Join(root, path), "trusted.overlay.opaque")
		if err != nil {
			return "", err
		}
		if len(opaque) == 1 && opaque[0] == 'y' {
			return path, nil
		}
	}

	return "", nil

}
