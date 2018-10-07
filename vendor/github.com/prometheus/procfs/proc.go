












package procfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)


type Proc struct {
	
	PID int

	fs FS
}


type Procs []Proc

func (p Procs) Len() int           { return len(p) }
func (p Procs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Procs) Less(i, j int) bool { return p[i].PID < p[j].PID }


func Self() (Proc, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Proc{}, err
	}
	return fs.Self()
}


func NewProc(pid int) (Proc, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Proc{}, err
	}
	return fs.NewProc(pid)
}


func AllProcs() (Procs, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Procs{}, err
	}
	return fs.AllProcs()
}


func (fs FS) Self() (Proc, error) {
	p, err := os.Readlink(fs.Path("self"))
	if err != nil {
		return Proc{}, err
	}
	pid, err := strconv.Atoi(strings.Replace(p, string(fs), "", -1))
	if err != nil {
		return Proc{}, err
	}
	return fs.NewProc(pid)
}


func (fs FS) NewProc(pid int) (Proc, error) {
	if _, err := os.Stat(fs.Path(strconv.Itoa(pid))); err != nil {
		return Proc{}, err
	}
	return Proc{PID: pid, fs: fs}, nil
}


func (fs FS) AllProcs() (Procs, error) {
	d, err := os.Open(fs.Path())
	if err != nil {
		return Procs{}, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return Procs{}, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	p := Procs{}
	for _, n := range names {
		pid, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			continue
		}
		p = append(p, Proc{PID: int(pid), fs: fs})
	}

	return p, nil
}


func (p Proc) CmdLine() ([]string, error) {
	f, err := os.Open(p.path("cmdline"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return []string{}, nil
	}

	return strings.Split(string(bytes.TrimRight(data, string("\x00"))), string(byte(0))), nil
}


func (p Proc) Comm() (string, error) {
	f, err := os.Open(p.path("comm"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}


func (p Proc) Executable() (string, error) {
	exe, err := os.Readlink(p.path("exe"))
	if os.IsNotExist(err) {
		return "", nil
	}

	return exe, err
}


func (p Proc) Cwd() (string, error) {
	wd, err := os.Readlink(p.path("cwd"))
	if os.IsNotExist(err) {
		return "", nil
	}

	return wd, err
}


func (p Proc) RootDir() (string, error) {
	rdir, err := os.Readlink(p.path("root"))
	if os.IsNotExist(err) {
		return "", nil
	}

	return rdir, err
}


func (p Proc) FileDescriptors() ([]uintptr, error) {
	names, err := p.fileDescriptors()
	if err != nil {
		return nil, err
	}

	fds := make([]uintptr, len(names))
	for i, n := range names {
		fd, err := strconv.ParseInt(n, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("could not parse fd %s: %s", n, err)
		}
		fds[i] = uintptr(fd)
	}

	return fds, nil
}



func (p Proc) FileDescriptorTargets() ([]string, error) {
	names, err := p.fileDescriptors()
	if err != nil {
		return nil, err
	}

	targets := make([]string, len(names))

	for i, name := range names {
		target, err := os.Readlink(p.path("fd", name))
		if err == nil {
			targets[i] = target
		}
	}

	return targets, nil
}



func (p Proc) FileDescriptorsLen() (int, error) {
	fds, err := p.fileDescriptors()
	if err != nil {
		return 0, err
	}

	return len(fds), nil
}



func (p Proc) MountStats() ([]*Mount, error) {
	f, err := os.Open(p.path("mountstats"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseMountStats(f)
}

func (p Proc) fileDescriptors() ([]string, error) {
	d, err := os.Open(p.path("fd"))
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	return names, nil
}

func (p Proc) path(pa ...string) string {
	return p.fs.Path(append([]string{strconv.Itoa(p.PID)}, pa...)...)
}
