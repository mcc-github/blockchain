package version

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)


var versionRegexp *regexp.Regexp



const VersionRegexpRaw string = `v?([0-9]+(\.[0-9]+)*?)` +
	`(-?([0-9A-Za-z\-~]+(\.[0-9A-Za-z\-~]+)*))?` +
	`(\+([0-9A-Za-z\-~]+(\.[0-9A-Za-z\-~]+)*))?` +
	`?`


type Version struct {
	metadata string
	pre      string
	segments []int64
	si       int
}

func init() {
	versionRegexp = regexp.MustCompile("^" + VersionRegexpRaw + "$")
}



func NewVersion(v string) (*Version, error) {
	matches := versionRegexp.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("Malformed version: %s", v)
	}
	segmentsStr := strings.Split(matches[1], ".")
	segments := make([]int64, len(segmentsStr))
	si := 0
	for i, str := range segmentsStr {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"Error parsing version: %s", err)
		}

		segments[i] = int64(val)
		si++
	}

	
	
	
	for i := len(segments); i < 3; i++ {
		segments = append(segments, 0)
	}

	return &Version{
		metadata: matches[7],
		pre:      matches[4],
		segments: segments,
		si:       si,
	}, nil
}



func Must(v *Version, err error) *Version {
	if err != nil {
		panic(err)
	}

	return v
}







func (v *Version) Compare(other *Version) int {
	
	if v.String() == other.String() {
		return 0
	}

	segmentsSelf := v.Segments64()
	segmentsOther := other.Segments64()

	
	if reflect.DeepEqual(segmentsSelf, segmentsOther) {
		preSelf := v.Prerelease()
		preOther := other.Prerelease()
		if preSelf == "" && preOther == "" {
			return 0
		}
		if preSelf == "" {
			return 1
		}
		if preOther == "" {
			return -1
		}

		return comparePrereleases(preSelf, preOther)
	}

	
	lenSelf := len(segmentsSelf)
	lenOther := len(segmentsOther)
	hS := lenSelf
	if lenSelf < lenOther {
		hS = lenOther
	}
	
	
	
	for i := 0; i < hS; i++ {
		if i > lenSelf-1 {
			
			
			if !allZero(segmentsOther[i:]) {
				
				return -1
			}
			break
		} else if i > lenOther-1 {
			
			
			if !allZero(segmentsSelf[i:]) {
				
				return 1
			}
			break
		}
		lhs := segmentsSelf[i]
		rhs := segmentsOther[i]
		if lhs == rhs {
			continue
		} else if lhs < rhs {
			return -1
		}
		
		return 1
	}

	
	return 0
}

func allZero(segs []int64) bool {
	for _, s := range segs {
		if s != 0 {
			return false
		}
	}
	return true
}

func comparePart(preSelf string, preOther string) int {
	if preSelf == preOther {
		return 0
	}

	var selfInt int64
	selfNumeric := true
	selfInt, err := strconv.ParseInt(preSelf, 10, 64)
	if err != nil {
		selfNumeric = false
	}

	var otherInt int64
	otherNumeric := true
	otherInt, err = strconv.ParseInt(preOther, 10, 64)
	if err != nil {
		otherNumeric = false
	}

	
	if preSelf == "" {
		if otherNumeric {
			return -1
		}
		return 1
	}

	if preOther == "" {
		if selfNumeric {
			return 1
		}
		return -1
	}

	if selfNumeric && !otherNumeric {
		return -1
	} else if !selfNumeric && otherNumeric {
		return 1
	} else if !selfNumeric && !otherNumeric && preSelf > preOther {
		return 1
	} else if selfInt > otherInt {
		return 1
	}

	return -1
}

func comparePrereleases(v string, other string) int {
	
	if v == other {
		return 0
	}

	
	selfPreReleaseMeta := strings.Split(v, ".")
	otherPreReleaseMeta := strings.Split(other, ".")

	selfPreReleaseLen := len(selfPreReleaseMeta)
	otherPreReleaseLen := len(otherPreReleaseMeta)

	biggestLen := otherPreReleaseLen
	if selfPreReleaseLen > otherPreReleaseLen {
		biggestLen = selfPreReleaseLen
	}

	
	for i := 0; i < biggestLen; i = i + 1 {
		partSelfPre := ""
		if i < selfPreReleaseLen {
			partSelfPre = selfPreReleaseMeta[i]
		}

		partOtherPre := ""
		if i < otherPreReleaseLen {
			partOtherPre = otherPreReleaseMeta[i]
		}

		compare := comparePart(partSelfPre, partOtherPre)
		
		if compare != 0 {
			return compare
		}
	}

	return 0
}


func (v *Version) Equal(o *Version) bool {
	return v.Compare(o) == 0
}


func (v *Version) GreaterThan(o *Version) bool {
	return v.Compare(o) > 0
}


func (v *Version) LessThan(o *Version) bool {
	return v.Compare(o) < 0
}






func (v *Version) Metadata() string {
	return v.metadata
}







func (v *Version) Prerelease() string {
	return v.pre
}






func (v *Version) Segments() []int {
	segmentSlice := make([]int, len(v.segments))
	for i, v := range v.segments {
		segmentSlice[i] = int(v)
	}
	return segmentSlice
}






func (v *Version) Segments64() []int64 {
	return v.segments
}



func (v *Version) String() string {
	var buf bytes.Buffer
	fmtParts := make([]string, len(v.segments))
	for i, s := range v.segments {
		
		str := strconv.FormatInt(s, 10)
		fmtParts[i] = str
	}
	fmt.Fprintf(&buf, strings.Join(fmtParts, "."))
	if v.pre != "" {
		fmt.Fprintf(&buf, "-%s", v.pre)
	}
	if v.metadata != "" {
		fmt.Fprintf(&buf, "+%s", v.metadata)
	}

	return buf.String()
}
