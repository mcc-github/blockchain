package system 

import "os"


func IsProcessAlive(pid int) bool {
	_, err := os.FindProcess(pid)

	return err == nil
}


func KillProcess(pid int) {
	p, err := os.FindProcess(pid)
	if err == nil {
		p.Kill()
	}
}
