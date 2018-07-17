package mount 

import (
	"golang.org/x/sys/unix"
)

const (
	
	ptypes = unix.MS_SHARED | unix.MS_PRIVATE | unix.MS_SLAVE | unix.MS_UNBINDABLE

	
	pflags = ptypes | unix.MS_REC | unix.MS_SILENT

	
	broflags = unix.MS_BIND | unix.MS_RDONLY
)


func isremount(device string, flags uintptr) bool {
	switch {
	
	
	case flags&unix.MS_REMOUNT != 0, device == "", device == "none":
		return true
	default:
		return false
	}
}

func mount(device, target, mType string, flags uintptr, data string) error {
	oflags := flags &^ ptypes
	if !isremount(device, flags) || data != "" {
		
		
		if err := unix.Mount(device, target, mType, oflags, data); err != nil {
			return err
		}
	}

	if flags&ptypes != 0 {
		
		if err := unix.Mount("", target, "", flags&pflags, ""); err != nil {
			return err
		}
	}

	if oflags&broflags == broflags {
		
		return unix.Mount("", target, "", oflags|unix.MS_REMOUNT, "")
	}

	return nil
}

func unmount(target string, flag int) error {
	return unix.Unmount(target, flag)
}
