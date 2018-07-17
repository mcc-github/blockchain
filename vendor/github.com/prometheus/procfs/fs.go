












package procfs

import (
	"fmt"
	"os"
	"path"

	"github.com/prometheus/procfs/nfs"
	"github.com/prometheus/procfs/xfs"
)



type FS string


const DefaultMountPoint = "/proc"



func NewFS(mountPoint string) (FS, error) {
	info, err := os.Stat(mountPoint)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %s", mountPoint, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("mount point %s is not a directory", mountPoint)
	}

	return FS(mountPoint), nil
}


func (fs FS) Path(p ...string) string {
	return path.Join(append([]string{string(fs)}, p...)...)
}


func (fs FS) XFSStats() (*xfs.Stats, error) {
	f, err := os.Open(fs.Path("fs/xfs/stat"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return xfs.ParseStats(f)
}


func (fs FS) NFSClientRPCStats() (*nfs.ClientRPCStats, error) {
	f, err := os.Open(fs.Path("net/rpc/nfs"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nfs.ParseClientRPCStats(f)
}


func (fs FS) NFSdServerRPCStats() (*nfs.ServerRPCStats, error) {
	f, err := os.Open(fs.Path("net/rpc/nfsd"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nfs.ParseServerRPCStats(f)
}
