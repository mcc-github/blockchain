package archive 

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)


type ChangeType int

const (
	
	ChangeModify = iota
	
	ChangeAdd
	
	ChangeDelete
)

func (c ChangeType) String() string {
	switch c {
	case ChangeModify:
		return "C"
	case ChangeAdd:
		return "A"
	case ChangeDelete:
		return "D"
	}
	return ""
}





type Change struct {
	Path string
	Kind ChangeType
}

func (change *Change) String() string {
	return fmt.Sprintf("%s %s", change.Kind, change.Path)
}


type changesByPath []Change

func (c changesByPath) Less(i, j int) bool { return c[i].Path < c[j].Path }
func (c changesByPath) Len() int           { return len(c) }
func (c changesByPath) Swap(i, j int)      { c[j], c[i] = c[i], c[j] }





func sameFsTime(a, b time.Time) bool {
	return a == b ||
		(a.Unix() == b.Unix() &&
			(a.Nanosecond() == 0 || b.Nanosecond() == 0))
}

func sameFsTimeSpec(a, b syscall.Timespec) bool {
	return a.Sec == b.Sec &&
		(a.Nsec == b.Nsec || a.Nsec == 0 || b.Nsec == 0)
}



func Changes(layers []string, rw string) ([]Change, error) {
	return changes(layers, rw, aufsDeletedFile, aufsMetadataSkip)
}

func aufsMetadataSkip(path string) (skip bool, err error) {
	skip, err = filepath.Match(string(os.PathSeparator)+WhiteoutMetaPrefix+"*", path)
	if err != nil {
		skip = true
	}
	return
}

func aufsDeletedFile(root, path string, fi os.FileInfo) (string, error) {
	f := filepath.Base(path)

	
	if strings.HasPrefix(f, WhiteoutPrefix) {
		originalFile := f[len(WhiteoutPrefix):]
		return filepath.Join(filepath.Dir(path), originalFile), nil
	}

	return "", nil
}

type skipChange func(string) (bool, error)
type deleteChange func(string, string, os.FileInfo) (string, error)

func changes(layers []string, rw string, dc deleteChange, sc skipChange) ([]Change, error) {
	var (
		changes     []Change
		changedDirs = make(map[string]struct{})
	)

	err := filepath.Walk(rw, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		
		path, err = filepath.Rel(rw, path)
		if err != nil {
			return err
		}

		
		path = filepath.Join(string(os.PathSeparator), path)

		
		if path == string(os.PathSeparator) {
			return nil
		}

		if sc != nil {
			if skip, err := sc(path); skip {
				return err
			}
		}

		change := Change{
			Path: path,
		}

		deletedFile, err := dc(rw, path, f)
		if err != nil {
			return err
		}

		
		if deletedFile != "" {
			change.Path = deletedFile
			change.Kind = ChangeDelete
		} else {
			
			change.Kind = ChangeAdd

			
			for _, layer := range layers {
				stat, err := os.Stat(filepath.Join(layer, path))
				if err != nil && !os.IsNotExist(err) {
					return err
				}
				if err == nil {
					

					
					
					if stat.IsDir() && f.IsDir() {
						if f.Size() == stat.Size() && f.Mode() == stat.Mode() && sameFsTime(f.ModTime(), stat.ModTime()) {
							
							return nil
						}
					}
					change.Kind = ChangeModify
					break
				}
			}
		}

		
		
		
		
		if f.IsDir() {
			changedDirs[path] = struct{}{}
		}
		if change.Kind == ChangeAdd || change.Kind == ChangeDelete {
			parent := filepath.Dir(path)
			if _, ok := changedDirs[parent]; !ok && parent != "/" {
				changes = append(changes, Change{Path: parent, Kind: ChangeModify})
				changedDirs[parent] = struct{}{}
			}
		}

		
		changes = append(changes, change)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return changes, nil
}


type FileInfo struct {
	parent     *FileInfo
	name       string
	stat       *system.StatT
	children   map[string]*FileInfo
	capability []byte
	added      bool
}


func (info *FileInfo) LookUp(path string) *FileInfo {
	
	parent := info
	if path == string(os.PathSeparator) {
		return info
	}

	pathElements := strings.Split(path, string(os.PathSeparator))
	for _, elem := range pathElements {
		if elem != "" {
			child := parent.children[elem]
			if child == nil {
				return nil
			}
			parent = child
		}
	}
	return parent
}

func (info *FileInfo) path() string {
	if info.parent == nil {
		
		return string(os.PathSeparator)
	}
	return filepath.Join(info.parent.path(), info.name)
}

func (info *FileInfo) addChanges(oldInfo *FileInfo, changes *[]Change) {

	sizeAtEntry := len(*changes)

	if oldInfo == nil {
		
		change := Change{
			Path: info.path(),
			Kind: ChangeAdd,
		}
		*changes = append(*changes, change)
		info.added = true
	}

	
	
	
	oldChildren := make(map[string]*FileInfo)
	if oldInfo != nil && info.isDir() {
		for k, v := range oldInfo.children {
			oldChildren[k] = v
		}
	}

	for name, newChild := range info.children {
		oldChild := oldChildren[name]
		if oldChild != nil {
			
			oldStat := oldChild.stat
			newStat := newChild.stat
			
			
			
			
			
			
			if statDifferent(oldStat, newStat) ||
				!bytes.Equal(oldChild.capability, newChild.capability) {
				change := Change{
					Path: newChild.path(),
					Kind: ChangeModify,
				}
				*changes = append(*changes, change)
				newChild.added = true
			}

			
			delete(oldChildren, name)
		}

		newChild.addChanges(oldChild, changes)
	}
	for _, oldChild := range oldChildren {
		
		change := Change{
			Path: oldChild.path(),
			Kind: ChangeDelete,
		}
		*changes = append(*changes, change)
	}

	
	
	
	if len(*changes) > sizeAtEntry && info.isDir() && !info.added && info.path() != string(os.PathSeparator) {
		change := Change{
			Path: info.path(),
			Kind: ChangeModify,
		}
		
		*changes = append(*changes, change) 
		copy((*changes)[sizeAtEntry+1:], (*changes)[sizeAtEntry:])
		(*changes)[sizeAtEntry] = change
	}

}


func (info *FileInfo) Changes(oldInfo *FileInfo) []Change {
	var changes []Change

	info.addChanges(oldInfo, &changes)

	return changes
}

func newRootFileInfo() *FileInfo {
	
	root := &FileInfo{
		name:     string(os.PathSeparator),
		children: make(map[string]*FileInfo),
	}
	return root
}



func ChangesDirs(newDir, oldDir string) ([]Change, error) {
	var (
		oldRoot, newRoot *FileInfo
	)
	if oldDir == "" {
		emptyDir, err := ioutil.TempDir("", "empty")
		if err != nil {
			return nil, err
		}
		defer os.Remove(emptyDir)
		oldDir = emptyDir
	}
	oldRoot, newRoot, err := collectFileInfoForChanges(oldDir, newDir)
	if err != nil {
		return nil, err
	}

	return newRoot.Changes(oldRoot), nil
}


func ChangesSize(newDir string, changes []Change) int64 {
	var (
		size int64
		sf   = make(map[uint64]struct{})
	)
	for _, change := range changes {
		if change.Kind == ChangeModify || change.Kind == ChangeAdd {
			file := filepath.Join(newDir, change.Path)
			fileInfo, err := os.Lstat(file)
			if err != nil {
				logrus.Errorf("Can not stat %q: %s", file, err)
				continue
			}

			if fileInfo != nil && !fileInfo.IsDir() {
				if hasHardlinks(fileInfo) {
					inode := getIno(fileInfo)
					if _, ok := sf[inode]; !ok {
						size += fileInfo.Size()
						sf[inode] = struct{}{}
					}
				} else {
					size += fileInfo.Size()
				}
			}
		}
	}
	return size
}


func ExportChanges(dir string, changes []Change, uidMaps, gidMaps []idtools.IDMap) (io.ReadCloser, error) {
	reader, writer := io.Pipe()
	go func() {
		ta := newTarAppender(idtools.NewIDMappingsFromMaps(uidMaps, gidMaps), writer, nil)

		
		defer pools.BufioWriter32KPool.Put(ta.Buffer)

		sort.Sort(changesByPath(changes))

		
		
		
		
		for _, change := range changes {
			if change.Kind == ChangeDelete {
				whiteOutDir := filepath.Dir(change.Path)
				whiteOutBase := filepath.Base(change.Path)
				whiteOut := filepath.Join(whiteOutDir, WhiteoutPrefix+whiteOutBase)
				timestamp := time.Now()
				hdr := &tar.Header{
					Name:       whiteOut[1:],
					Size:       0,
					ModTime:    timestamp,
					AccessTime: timestamp,
					ChangeTime: timestamp,
				}
				if err := ta.TarWriter.WriteHeader(hdr); err != nil {
					logrus.Debugf("Can't write whiteout header: %s", err)
				}
			} else {
				path := filepath.Join(dir, change.Path)
				if err := ta.addTarFile(path, change.Path[1:]); err != nil {
					logrus.Debugf("Can't add file %s to tar: %s", path, err)
				}
			}
		}

		
		if err := ta.TarWriter.Close(); err != nil {
			logrus.Debugf("Can't close layer: %s", err)
		}
		if err := writer.Close(); err != nil {
			logrus.Debugf("failed close Changes writer: %s", err)
		}
	}()
	return reader, nil
}
