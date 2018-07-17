

package system


import "C"

func GetClockTicks() int {
	return int(C.sysconf(C._SC_CLK_TCK))
}
