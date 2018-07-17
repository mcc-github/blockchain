

package idtools 

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/docker/pkg/system"
	"github.com/opencontainers/runc/libcontainer/user"
)

var (
	entOnce   sync.Once
	getentCmd string
)

func mkdirAs(path string, mode os.FileMode, ownerUID, ownerGID int, mkAll, chownExisting bool) error {
	
	
	
	
	var paths []string

	stat, err := system.Stat(path)
	if err == nil {
		if !stat.IsDir() {
			return &os.PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
		}
		if !chownExisting {
			return nil
		}

		
		return lazyChown(path, ownerUID, ownerGID, stat)
	}

	if os.IsNotExist(err) {
		paths = []string{path}
	}

	if mkAll {
		
		
		dirPath := path
		for {
			dirPath = filepath.Dir(dirPath)
			if dirPath == "/" {
				break
			}
			if _, err := os.Stat(dirPath); err != nil && os.IsNotExist(err) {
				paths = append(paths, dirPath)
			}
		}
		if err := system.MkdirAll(path, mode, ""); err != nil {
			return err
		}
	} else {
		if err := os.Mkdir(path, mode); err != nil && !os.IsExist(err) {
			return err
		}
	}
	
	
	for _, pathComponent := range paths {
		if err := lazyChown(pathComponent, ownerUID, ownerGID, nil); err != nil {
			return err
		}
	}
	return nil
}



func CanAccess(path string, pair IDPair) bool {
	statInfo, err := system.Stat(path)
	if err != nil {
		return false
	}
	fileMode := os.FileMode(statInfo.Mode())
	permBits := fileMode.Perm()
	return accessible(statInfo.UID() == uint32(pair.UID),
		statInfo.GID() == uint32(pair.GID), permBits)
}

func accessible(isOwner, isGroup bool, perms os.FileMode) bool {
	if isOwner && (perms&0100 == 0100) {
		return true
	}
	if isGroup && (perms&0010 == 0010) {
		return true
	}
	if perms&0001 == 0001 {
		return true
	}
	return false
}



func LookupUser(username string) (user.User, error) {
	
	usr, err := user.LookupUser(username)
	if err == nil {
		return usr, nil
	}
	
	usr, err = getentUser(fmt.Sprintf("%s %s", "passwd", username))
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}



func LookupUID(uid int) (user.User, error) {
	
	usr, err := user.LookupUid(uid)
	if err == nil {
		return usr, nil
	}
	
	return getentUser(fmt.Sprintf("%s %d", "passwd", uid))
}

func getentUser(args string) (user.User, error) {
	reader, err := callGetent(args)
	if err != nil {
		return user.User{}, err
	}
	users, err := user.ParsePasswd(reader)
	if err != nil {
		return user.User{}, err
	}
	if len(users) == 0 {
		return user.User{}, fmt.Errorf("getent failed to find passwd entry for %q", strings.Split(args, " ")[1])
	}
	return users[0], nil
}



func LookupGroup(groupname string) (user.Group, error) {
	
	group, err := user.LookupGroup(groupname)
	if err == nil {
		return group, nil
	}
	
	return getentGroup(fmt.Sprintf("%s %s", "group", groupname))
}



func LookupGID(gid int) (user.Group, error) {
	
	group, err := user.LookupGid(gid)
	if err == nil {
		return group, nil
	}
	
	return getentGroup(fmt.Sprintf("%s %d", "group", gid))
}

func getentGroup(args string) (user.Group, error) {
	reader, err := callGetent(args)
	if err != nil {
		return user.Group{}, err
	}
	groups, err := user.ParseGroup(reader)
	if err != nil {
		return user.Group{}, err
	}
	if len(groups) == 0 {
		return user.Group{}, fmt.Errorf("getent failed to find groups entry for %q", strings.Split(args, " ")[1])
	}
	return groups[0], nil
}

func callGetent(args string) (io.Reader, error) {
	entOnce.Do(func() { getentCmd, _ = resolveBinary("getent") })
	
	if getentCmd == "" {
		return nil, fmt.Errorf("")
	}
	out, err := execCmd(getentCmd, args)
	if err != nil {
		exitCode, errC := system.GetExitCode(err)
		if errC != nil {
			return nil, err
		}
		switch exitCode {
		case 1:
			return nil, fmt.Errorf("getent reported invalid parameters/database unknown")
		case 2:
			terms := strings.Split(args, " ")
			return nil, fmt.Errorf("getent unable to find entry %q in %s database", terms[1], terms[0])
		case 3:
			return nil, fmt.Errorf("getent database doesn't support enumeration")
		default:
			return nil, err
		}

	}
	return bytes.NewReader(out), nil
}




func lazyChown(p string, uid, gid int, stat *system.StatT) error {
	if stat == nil {
		var err error
		stat, err = system.Stat(p)
		if err != nil {
			return err
		}
	}
	if stat.UID() == uint32(uid) && stat.GID() == uint32(gid) {
		return nil
	}
	return os.Chown(p, uid, gid)
}
