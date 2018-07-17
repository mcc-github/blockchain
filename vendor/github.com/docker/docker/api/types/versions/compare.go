package versions 

import (
	"strconv"
	"strings"
)



func compare(v1, v2 string) int {
	var (
		currTab  = strings.Split(v1, ".")
		otherTab = strings.Split(v2, ".")
	)

	max := len(currTab)
	if len(otherTab) > max {
		max = len(otherTab)
	}
	for i := 0; i < max; i++ {
		var currInt, otherInt int

		if len(currTab) > i {
			currInt, _ = strconv.Atoi(currTab[i])
		}
		if len(otherTab) > i {
			otherInt, _ = strconv.Atoi(otherTab[i])
		}
		if currInt > otherInt {
			return 1
		}
		if otherInt > currInt {
			return -1
		}
	}
	return 0
}


func LessThan(v, other string) bool {
	return compare(v, other) == -1
}


func LessThanOrEqualTo(v, other string) bool {
	return compare(v, other) <= 0
}


func GreaterThan(v, other string) bool {
	return compare(v, other) == 1
}


func GreaterThanOrEqualTo(v, other string) bool {
	return compare(v, other) >= 0
}


func Equal(v, other string) bool {
	return compare(v, other) == 0
}
