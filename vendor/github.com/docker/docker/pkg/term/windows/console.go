

package windowsconsole 

import (
	"os"

	"github.com/Azure/go-ansiterm/winterm"
)


func GetHandleInfo(in interface{}) (uintptr, bool) {
	switch t := in.(type) {
	case *ansiReader:
		return t.Fd(), true
	case *ansiWriter:
		return t.Fd(), true
	}

	var inFd uintptr
	var isTerminal bool

	if file, ok := in.(*os.File); ok {
		inFd = file.Fd()
		isTerminal = IsConsole(inFd)
	}
	return inFd, isTerminal
}



func IsConsole(fd uintptr) bool {
	_, e := winterm.GetConsoleMode(fd)
	return e == nil
}
