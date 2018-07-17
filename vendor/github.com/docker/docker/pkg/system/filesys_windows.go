package system 

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	winio "github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

const (
	
	SddlAdministratorsLocalSystem = "D:P(A;OICI;GA;;;BA)(A;OICI;GA;;;SY)"
	
	SddlNtvmAdministratorsLocalSystem = "D:P(A;OICI;GA;;;S-1-5-83-0)(A;OICI;GA;;;BA)(A;OICI;GA;;;SY)"
)



func MkdirAllWithACL(path string, perm os.FileMode, sddl string) error {
	return mkdirall(path, true, sddl)
}


func MkdirAll(path string, _ os.FileMode, sddl string) error {
	return mkdirall(path, false, sddl)
}




func mkdirall(path string, applyACL bool, sddl string) error {
	if re := regexp.MustCompile(`^\\\\\?\\Volume{[a-z0-9-]+}$`); re.MatchString(path) {
		return nil
	}

	
	

	
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{
			Op:   "mkdir",
			Path: path,
			Err:  syscall.ENOTDIR,
		}
	}

	
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { 
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { 
		j--
	}

	if j > 1 {
		
		err = mkdirall(path[0:j-1], false, sddl)
		if err != nil {
			return err
		}
	}

	
	if applyACL {
		err = mkdirWithACL(path, sddl)
	} else {
		err = os.Mkdir(path, 0)
	}

	if err != nil {
		
		
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}








func mkdirWithACL(name string, sddl string) error {
	sa := windows.SecurityAttributes{Length: 0}
	sd, err := winio.SddlToSecurityDescriptor(sddl)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	sa.SecurityDescriptor = uintptr(unsafe.Pointer(&sd[0]))

	namep, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: err}
	}

	e := windows.CreateDirectory(namep, &sa)
	if e != nil {
		return &os.PathError{Op: "mkdir", Path: name, Err: e}
	}
	return nil
}








func IsAbs(path string) bool {
	if !filepath.IsAbs(path) {
		if !strings.HasPrefix(path, string(os.PathSeparator)) {
			return false
		}
	}
	return true
}














func CreateSequential(name string) (*os.File, error) {
	return OpenFileSequential(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
}





func OpenSequential(name string) (*os.File, error) {
	return OpenFileSequential(name, os.O_RDONLY, 0)
}




func OpenFileSequential(name string, flag int, _ os.FileMode) (*os.File, error) {
	if name == "" {
		return nil, &os.PathError{Op: "open", Path: name, Err: syscall.ENOENT}
	}
	r, errf := windowsOpenFileSequential(name, flag, 0)
	if errf == nil {
		return r, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: errf}
}

func windowsOpenFileSequential(name string, flag int, _ os.FileMode) (file *os.File, err error) {
	r, e := windowsOpenSequential(name, flag|windows.O_CLOEXEC, 0)
	if e != nil {
		return nil, e
	}
	return os.NewFile(uintptr(r), name), nil
}

func makeInheritSa() *windows.SecurityAttributes {
	var sa windows.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	return &sa
}

func windowsOpenSequential(path string, mode int, _ uint32) (fd windows.Handle, err error) {
	if len(path) == 0 {
		return windows.InvalidHandle, windows.ERROR_FILE_NOT_FOUND
	}
	pathp, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	var access uint32
	switch mode & (windows.O_RDONLY | windows.O_WRONLY | windows.O_RDWR) {
	case windows.O_RDONLY:
		access = windows.GENERIC_READ
	case windows.O_WRONLY:
		access = windows.GENERIC_WRITE
	case windows.O_RDWR:
		access = windows.GENERIC_READ | windows.GENERIC_WRITE
	}
	if mode&windows.O_CREAT != 0 {
		access |= windows.GENERIC_WRITE
	}
	if mode&windows.O_APPEND != 0 {
		access &^= windows.GENERIC_WRITE
		access |= windows.FILE_APPEND_DATA
	}
	sharemode := uint32(windows.FILE_SHARE_READ | windows.FILE_SHARE_WRITE)
	var sa *windows.SecurityAttributes
	if mode&windows.O_CLOEXEC == 0 {
		sa = makeInheritSa()
	}
	var createmode uint32
	switch {
	case mode&(windows.O_CREAT|windows.O_EXCL) == (windows.O_CREAT | windows.O_EXCL):
		createmode = windows.CREATE_NEW
	case mode&(windows.O_CREAT|windows.O_TRUNC) == (windows.O_CREAT | windows.O_TRUNC):
		createmode = windows.CREATE_ALWAYS
	case mode&windows.O_CREAT == windows.O_CREAT:
		createmode = windows.OPEN_ALWAYS
	case mode&windows.O_TRUNC == windows.O_TRUNC:
		createmode = windows.TRUNCATE_EXISTING
	default:
		createmode = windows.OPEN_EXISTING
	}
	
	
	const fileFlagSequentialScan = 0x08000000 
	h, e := windows.CreateFile(pathp, access, sharemode, sa, createmode, fileFlagSequentialScan, 0)
	return h, e
}


var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}
func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}












func TempFileSequential(dir, prefix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix())
		f, err = OpenFileSequential(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}
