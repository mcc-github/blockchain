package mount 


import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)



func parseMountTable(filter FilterFunc) ([]*Info, error) {
	var rawEntries *C.struct_statfs

	count := int(C.getmntinfo(&rawEntries, C.MNT_WAIT))
	if count == 0 {
		return nil, fmt.Errorf("Failed to call getmntinfo")
	}

	var entries []C.struct_statfs
	header := (*reflect.SliceHeader)(unsafe.Pointer(&entries))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(rawEntries))

	var out []*Info
	for _, entry := range entries {
		var mountinfo Info
		var skip, stop bool
		mountinfo.Mountpoint = C.GoString(&entry.f_mntonname[0])

		if filter != nil {
			
			skip, stop = filter(p)
			if skip {
				continue
			}
		}

		mountinfo.Source = C.GoString(&entry.f_mntfromname[0])
		mountinfo.Fstype = C.GoString(&entry.f_fstypename[0])

		out = append(out, &mountinfo)
		if stop {
			break
		}
	}
	return out, nil
}
