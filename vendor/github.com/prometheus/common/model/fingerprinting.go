












package model

import (
	"fmt"
	"strconv"
)



type Fingerprint uint64


func FingerprintFromString(s string) (Fingerprint, error) {
	num, err := strconv.ParseUint(s, 16, 64)
	return Fingerprint(num), err
}


func ParseFingerprint(s string) (Fingerprint, error) {
	num, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return Fingerprint(num), nil
}

func (f Fingerprint) String() string {
	return fmt.Sprintf("%016x", uint64(f))
}



type Fingerprints []Fingerprint


func (f Fingerprints) Len() int {
	return len(f)
}


func (f Fingerprints) Less(i, j int) bool {
	return f[i] < f[j]
}


func (f Fingerprints) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}


type FingerprintSet map[Fingerprint]struct{}


func (s FingerprintSet) Equal(o FingerprintSet) bool {
	if len(s) != len(o) {
		return false
	}

	for k := range s {
		if _, ok := o[k]; !ok {
			return false
		}
	}

	return true
}


func (s FingerprintSet) Intersection(o FingerprintSet) FingerprintSet {
	myLength, otherLength := len(s), len(o)
	if myLength == 0 || otherLength == 0 {
		return FingerprintSet{}
	}

	subSet := s
	superSet := o

	if otherLength < myLength {
		subSet = o
		superSet = s
	}

	out := FingerprintSet{}

	for k := range subSet {
		if _, ok := superSet[k]; ok {
			out[k] = struct{}{}
		}
	}

	return out
}
