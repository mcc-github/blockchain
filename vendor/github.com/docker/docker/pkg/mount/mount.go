package mount 

import (
	"sort"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)







type FilterFunc func(*Info) (skip, stop bool)



func PrefixFilter(prefix string) FilterFunc {
	return func(m *Info) (bool, bool) {
		skip := !strings.HasPrefix(m.Mountpoint, prefix)
		return skip, false
	}
}


func SingleEntryFilter(mp string) FilterFunc {
	return func(m *Info) (bool, bool) {
		if m.Mountpoint == mp {
			return false, true 
		}
		return true, false 
	}
}





func ParentsFilter(path string) FilterFunc {
	return func(m *Info) (bool, bool) {
		skip := !strings.HasPrefix(path, m.Mountpoint)
		return skip, false
	}
}



func GetMounts(f FilterFunc) ([]*Info, error) {
	return parseMountTable(f)
}



func Mounted(mountpoint string) (bool, error) {
	entries, err := GetMounts(SingleEntryFilter(mountpoint))
	if err != nil {
		return false, err
	}

	return len(entries) > 0, nil
}





func Mount(device, target, mType, options string) error {
	flag, _ := parseOptions(options)
	if flag&REMOUNT != REMOUNT {
		if mounted, err := Mounted(target); err != nil || mounted {
			return err
		}
	}
	return ForceMount(device, target, mType, options)
}





func ForceMount(device, target, mType, options string) error {
	flag, data := parseOptions(options)
	return mount(device, target, mType, uintptr(flag), data)
}



func Unmount(target string) error {
	err := unmount(target, mntDetach)
	if err == syscall.EINVAL {
		
		err = nil
	}
	return err
}



func RecursiveUnmount(target string) error {
	mounts, err := parseMountTable(PrefixFilter(target))
	if err != nil {
		return err
	}

	
	sort.Slice(mounts, func(i, j int) bool {
		return len(mounts[i].Mountpoint) > len(mounts[j].Mountpoint)
	})

	for i, m := range mounts {
		logrus.Debugf("Trying to unmount %s", m.Mountpoint)
		err = unmount(m.Mountpoint, mntDetach)
		if err != nil {
			
			
			
			
			
			
			
			if err == syscall.EINVAL {
				continue
			}
			if i == len(mounts)-1 {
				if mounted, e := Mounted(m.Mountpoint); e != nil || mounted {
					return err
				}
				continue
			}
			
			logrus.WithError(err).Warnf("Failed to unmount submount %s", m.Mountpoint)
			continue
		}

		logrus.Debugf("Unmounted %s", m.Mountpoint)
	}
	return nil
}
