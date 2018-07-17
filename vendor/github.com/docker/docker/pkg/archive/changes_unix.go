

package archive 

import (
	"os"
	"syscall"

	"github.com/docker/docker/pkg/system"
	"golang.org/x/sys/unix"
)

func statDifferent(oldStat *system.StatT, newStat *system.StatT) bool {
	
	if oldStat.Mode() != newStat.Mode() ||
		oldStat.UID() != newStat.UID() ||
		oldStat.GID() != newStat.GID() ||
		oldStat.Rdev() != newStat.Rdev() ||
		
		(oldStat.Mode()&unix.S_IFDIR != unix.S_IFDIR &&
			(!sameFsTimeSpec(oldStat.Mtim(), newStat.Mtim()) || (oldStat.Size() != newStat.Size()))) {
		return true
	}
	return false
}

func (info *FileInfo) isDir() bool {
	return info.parent == nil || info.stat.Mode()&unix.S_IFDIR != 0
}

func getIno(fi os.FileInfo) uint64 {
	return fi.Sys().(*syscall.Stat_t).Ino
}

func hasHardlinks(fi os.FileInfo) bool {
	return fi.Sys().(*syscall.Stat_t).Nlink > 1
}
