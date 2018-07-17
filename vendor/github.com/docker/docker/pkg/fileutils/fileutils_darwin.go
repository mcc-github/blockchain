package fileutils 

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)



func GetTotalUsedFds() int {
	pid := os.Getpid()

	cmd := exec.Command("lsof", "-p", strconv.Itoa(pid))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1
	}

	outputStr := strings.TrimSpace(string(output))

	fds := strings.Split(outputStr, "\n")

	return len(fds) - 1
}
