





package unix

import (
	errorspkg "errors"
	"fmt"
)



const (
	
	capRightsGoVersion = CAP_RIGHTS_VERSION_00
	capArSizeMin       = CAP_RIGHTS_VERSION_00 + 2
	capArSizeMax       = capRightsGoVersion + 2
)

var (
	bit2idx = []int{
		-1, 0, 1, -1, 2, -1, -1, -1, 3, -1, -1, -1, -1, -1, -1, -1,
		4, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}
)

func capidxbit(right uint64) int {
	return int((right >> 57) & 0x1f)
}

func rightToIndex(right uint64) (int, error) {
	idx := capidxbit(right)
	if idx < 0 || idx >= len(bit2idx) {
		return -2, fmt.Errorf("index for right 0x%x out of range", right)
	}
	return bit2idx[idx], nil
}

func caprver(right uint64) int {
	return int(right >> 62)
}

func capver(rights *CapRights) int {
	return caprver(rights.Rights[0])
}

func caparsize(rights *CapRights) int {
	return capver(rights) + 2
}


func CapRightsSet(rights *CapRights, setrights []uint64) error {
	
	if capver(rights) != CAP_RIGHTS_VERSION_00 {
		return fmt.Errorf("bad rights version %d", capver(rights))
	}

	n := caparsize(rights)
	if n < capArSizeMin || n > capArSizeMax {
		return errorspkg.New("bad rights size")
	}

	for _, right := range setrights {
		if caprver(right) != CAP_RIGHTS_VERSION_00 {
			return errorspkg.New("bad right version")
		}
		i, err := rightToIndex(right)
		if err != nil {
			return err
		}
		if i >= n {
			return errorspkg.New("index overflow")
		}
		if capidxbit(rights.Rights[i]) != capidxbit(right) {
			return errorspkg.New("index mismatch")
		}
		rights.Rights[i] |= right
		if capidxbit(rights.Rights[i]) != capidxbit(right) {
			return errorspkg.New("index mismatch (after assign)")
		}
	}

	return nil
}


func CapRightsClear(rights *CapRights, clearrights []uint64) error {
	
	if capver(rights) != CAP_RIGHTS_VERSION_00 {
		return fmt.Errorf("bad rights version %d", capver(rights))
	}

	n := caparsize(rights)
	if n < capArSizeMin || n > capArSizeMax {
		return errorspkg.New("bad rights size")
	}

	for _, right := range clearrights {
		if caprver(right) != CAP_RIGHTS_VERSION_00 {
			return errorspkg.New("bad right version")
		}
		i, err := rightToIndex(right)
		if err != nil {
			return err
		}
		if i >= n {
			return errorspkg.New("index overflow")
		}
		if capidxbit(rights.Rights[i]) != capidxbit(right) {
			return errorspkg.New("index mismatch")
		}
		rights.Rights[i] &= ^(right & 0x01FFFFFFFFFFFFFF)
		if capidxbit(rights.Rights[i]) != capidxbit(right) {
			return errorspkg.New("index mismatch (after assign)")
		}
	}

	return nil
}


func CapRightsIsSet(rights *CapRights, setrights []uint64) (bool, error) {
	
	if capver(rights) != CAP_RIGHTS_VERSION_00 {
		return false, fmt.Errorf("bad rights version %d", capver(rights))
	}

	n := caparsize(rights)
	if n < capArSizeMin || n > capArSizeMax {
		return false, errorspkg.New("bad rights size")
	}

	for _, right := range setrights {
		if caprver(right) != CAP_RIGHTS_VERSION_00 {
			return false, errorspkg.New("bad right version")
		}
		i, err := rightToIndex(right)
		if err != nil {
			return false, err
		}
		if i >= n {
			return false, errorspkg.New("index overflow")
		}
		if capidxbit(rights.Rights[i]) != capidxbit(right) {
			return false, errorspkg.New("index mismatch")
		}
		if (rights.Rights[i] & right) != right {
			return false, nil
		}
	}

	return true, nil
}

func capright(idx uint64, bit uint64) uint64 {
	return ((1 << (57 + idx)) | bit)
}



func CapRightsInit(rights []uint64) (*CapRights, error) {
	var r CapRights
	r.Rights[0] = (capRightsGoVersion << 62) | capright(0, 0)
	r.Rights[1] = capright(1, 0)

	err := CapRightsSet(&r, rights)
	if err != nil {
		return nil, err
	}
	return &r, nil
}




func CapRightsLimit(fd uintptr, rights *CapRights) error {
	return capRightsLimit(int(fd), rights)
}



func CapRightsGet(fd uintptr) (*CapRights, error) {
	r, err := CapRightsInit(nil)
	if err != nil {
		return nil, err
	}
	err = capRightsGet(capRightsGoVersion, int(fd), r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
