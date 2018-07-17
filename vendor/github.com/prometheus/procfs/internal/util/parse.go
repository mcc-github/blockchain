












package util

import "strconv"


func ParseUint32s(ss []string) ([]uint32, error) {
	us := make([]uint32, 0, len(ss))
	for _, s := range ss {
		u, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}

		us = append(us, uint32(u))
	}

	return us, nil
}


func ParseUint64s(ss []string) ([]uint64, error) {
	us := make([]uint64, 0, len(ss))
	for _, s := range ss {
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		us = append(us, u)
	}

	return us, nil
}
