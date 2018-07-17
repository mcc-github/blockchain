package system

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)



func GetProcessStartTime(pid int) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(data), " ")
	
	
	
	
	
	
	
	return parts[22-1], nil 
}
