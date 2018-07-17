



package properties

import (
	"fmt"
	"math"
)


var is32Bit = ^uint(0) == math.MaxUint32



func intRangeCheck(key string, v int64) int {
	if is32Bit && (v < math.MinInt32 || v > math.MaxInt32) {
		panic(fmt.Sprintf("Value %d for key %s out of range", v, key))
	}
	return int(v)
}



func uintRangeCheck(key string, v uint64) uint {
	if is32Bit && v > math.MaxUint32 {
		panic(fmt.Sprintf("Value %d for key %s out of range", v, key))
	}
	return uint(v)
}
