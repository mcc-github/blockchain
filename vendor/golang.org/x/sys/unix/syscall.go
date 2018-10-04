























package unix 

import "strings"




func ByteSliceFromString(s string) ([]byte, error) {
	if strings.IndexByte(s, 0) != -1 {
		return nil, EINVAL
	}
	a := make([]byte, len(s)+1)
	copy(a, s)
	return a, nil
}




func BytePtrFromString(s string) (*byte, error) {
	a, err := ByteSliceFromString(s)
	if err != nil {
		return nil, err
	}
	return &a[0], nil
}



var _zero uintptr
