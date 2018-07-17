

package winterm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/Azure/go-ansiterm"
)



const (
	VK_PRIOR    = 0x21 
	VK_NEXT     = 0x22 
	VK_END      = 0x23 
	VK_HOME     = 0x24 
	VK_LEFT     = 0x25 
	VK_UP       = 0x26 
	VK_RIGHT    = 0x27 
	VK_DOWN     = 0x28 
	VK_SELECT   = 0x29 
	VK_PRINT    = 0x2A 
	VK_EXECUTE  = 0x2B 
	VK_SNAPSHOT = 0x2C 
	VK_INSERT   = 0x2D 
	VK_DELETE   = 0x2E 
	VK_HELP     = 0x2F 
	VK_F1       = 0x70 
	VK_F2       = 0x71 
	VK_F3       = 0x72 
	VK_F4       = 0x73 
	VK_F5       = 0x74 
	VK_F6       = 0x75 
	VK_F7       = 0x76 
	VK_F8       = 0x77 
	VK_F9       = 0x78 
	VK_F10      = 0x79 
	VK_F11      = 0x7A 
	VK_F12      = 0x7B 

	RIGHT_ALT_PRESSED  = 0x0001
	LEFT_ALT_PRESSED   = 0x0002
	RIGHT_CTRL_PRESSED = 0x0004
	LEFT_CTRL_PRESSED  = 0x0008
	SHIFT_PRESSED      = 0x0010
	NUMLOCK_ON         = 0x0020
	SCROLLLOCK_ON      = 0x0040
	CAPSLOCK_ON        = 0x0080
	ENHANCED_KEY       = 0x0100
)

type ansiCommand struct {
	CommandBytes []byte
	Command      string
	Parameters   []string
	IsSpecial    bool
}

func newAnsiCommand(command []byte) *ansiCommand {

	if isCharacterSelectionCmdChar(command[1]) {
		
		return &ansiCommand{
			CommandBytes: command,
			Command:      string(command),
			IsSpecial:    true,
		}
	}

	
	lastCharIndex := len(command) - 1

	ac := &ansiCommand{
		CommandBytes: command,
		Command:      string(command[lastCharIndex]),
		IsSpecial:    false,
	}

	
	if lastCharIndex != 0 {
		start := 1
		
		if command[0] == ansiterm.ANSI_ESCAPE_PRIMARY && command[1] == ansiterm.ANSI_ESCAPE_SECONDARY {
			start++
		}
		
		ac.Parameters = strings.Split(string(command[start:lastCharIndex]), ansiterm.ANSI_PARAMETER_SEP)
	}

	return ac
}

func (ac *ansiCommand) paramAsSHORT(index int, defaultValue int16) int16 {
	if index < 0 || index >= len(ac.Parameters) {
		return defaultValue
	}

	param, err := strconv.ParseInt(ac.Parameters[index], 10, 16)
	if err != nil {
		return defaultValue
	}

	return int16(param)
}

func (ac *ansiCommand) String() string {
	return fmt.Sprintf("0x%v \"%v\" (\"%v\")",
		bytesToHex(ac.CommandBytes),
		ac.Command,
		strings.Join(ac.Parameters, "\",\""))
}



func isAnsiCommandChar(b byte) bool {
	switch {
	case ansiterm.ANSI_COMMAND_FIRST <= b && b <= ansiterm.ANSI_COMMAND_LAST && b != ansiterm.ANSI_ESCAPE_SECONDARY:
		return true
	case b == ansiterm.ANSI_CMD_G1 || b == ansiterm.ANSI_CMD_OSC || b == ansiterm.ANSI_CMD_DECPAM || b == ansiterm.ANSI_CMD_DECPNM:
		
		return true
	case b == ansiterm.ANSI_CMD_STR_TERM || b == ansiterm.ANSI_BEL:
		
		return true
	}
	return false
}

func isXtermOscSequence(command []byte, current byte) bool {
	return (len(command) >= 2 && command[0] == ansiterm.ANSI_ESCAPE_PRIMARY && command[1] == ansiterm.ANSI_CMD_OSC && current != ansiterm.ANSI_BEL)
}

func isCharacterSelectionCmdChar(b byte) bool {
	return (b == ansiterm.ANSI_CMD_G0 || b == ansiterm.ANSI_CMD_G1 || b == ansiterm.ANSI_CMD_G2 || b == ansiterm.ANSI_CMD_G3)
}


func bytesToHex(b []byte) string {
	hex := make([]string, len(b))
	for i, ch := range b {
		hex[i] = fmt.Sprintf("%X", ch)
	}
	return strings.Join(hex, "")
}



func ensureInRange(n int16, min int16, max int16) int16 {
	if n < min {
		return min
	} else if n > max {
		return max
	} else {
		return n
	}
}

func GetStdFile(nFile int) (*os.File, uintptr) {
	var file *os.File
	switch nFile {
	case syscall.STD_INPUT_HANDLE:
		file = os.Stdin
	case syscall.STD_OUTPUT_HANDLE:
		file = os.Stdout
	case syscall.STD_ERROR_HANDLE:
		file = os.Stderr
	default:
		panic(fmt.Errorf("Invalid standard handle identifier: %v", nFile))
	}

	fd, err := syscall.GetStdHandle(nFile)
	if err != nil {
		panic(fmt.Errorf("Invalid standard handle identifier: %v -- %v", nFile, err))
	}

	return file, uintptr(fd)
}
