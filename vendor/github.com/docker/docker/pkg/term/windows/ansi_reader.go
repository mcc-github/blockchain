

package windowsconsole 

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	ansiterm "github.com/Azure/go-ansiterm"
	"github.com/Azure/go-ansiterm/winterm"
)

const (
	escapeSequence = ansiterm.KEY_ESC_CSI
)


type ansiReader struct {
	file     *os.File
	fd       uintptr
	buffer   []byte
	cbBuffer int
	command  []byte
}



func NewAnsiReader(nFile int) io.ReadCloser {
	initLogger()
	file, fd := winterm.GetStdFile(nFile)
	return &ansiReader{
		file:    file,
		fd:      fd,
		command: make([]byte, 0, ansiterm.ANSI_MAX_CMD_LENGTH),
		buffer:  make([]byte, 0),
	}
}


func (ar *ansiReader) Close() (err error) {
	return ar.file.Close()
}


func (ar *ansiReader) Fd() uintptr {
	return ar.fd
}


func (ar *ansiReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	
	if len(ar.buffer) > 0 {
		logger.Debugf("Reading previously cached bytes")

		originalLength := len(ar.buffer)
		copiedLength := copy(p, ar.buffer)

		if copiedLength == originalLength {
			ar.buffer = make([]byte, 0, len(p))
		} else {
			ar.buffer = ar.buffer[copiedLength:]
		}

		logger.Debugf("Read from cache p[%d]: % x", copiedLength, p)
		return copiedLength, nil
	}

	
	events, err := readInputEvents(ar.fd, len(p))
	if err != nil {
		return 0, err
	} else if len(events) == 0 {
		logger.Debug("No input events detected")
		return 0, nil
	}

	keyBytes := translateKeyEvents(events, []byte(escapeSequence))

	
	if len(keyBytes) > len(p) {
		logger.Debugf("Received %d keyBytes, only room for %d bytes", len(keyBytes), len(p))
		ar.buffer = keyBytes[len(p):]
		keyBytes = keyBytes[:len(p)]
	} else if len(keyBytes) == 0 {
		logger.Debug("No key bytes returned from the translator")
		return 0, nil
	}

	copiedLength := copy(p, keyBytes)
	if copiedLength != len(keyBytes) {
		return 0, errors.New("unexpected copy length encountered")
	}

	logger.Debugf("Read        p[%d]: % x", copiedLength, p)
	logger.Debugf("Read keyBytes[%d]: % x", copiedLength, keyBytes)
	return copiedLength, nil
}


func readInputEvents(fd uintptr, maxBytes int) ([]winterm.INPUT_RECORD, error) {
	
	
	
	
	recordSize := int(unsafe.Sizeof(*((*winterm.INPUT_RECORD)(unsafe.Pointer(&maxBytes)))))
	countRecords := maxBytes / recordSize
	if countRecords > ansiterm.MAX_INPUT_EVENTS {
		countRecords = ansiterm.MAX_INPUT_EVENTS
	} else if countRecords == 0 {
		countRecords = 1
	}
	logger.Debugf("[windows] readInputEvents: Reading %v records (buffer size %v, record size %v)", countRecords, maxBytes, recordSize)

	
	events := make([]winterm.INPUT_RECORD, countRecords)
	nEvents := uint32(0)
	eventsExist, err := winterm.WaitForSingleObject(fd, winterm.WAIT_INFINITE)
	if err != nil {
		return nil, err
	}

	if eventsExist {
		err = winterm.ReadConsoleInput(fd, events, &nEvents)
		if err != nil {
			return nil, err
		}
	}

	
	logger.Debugf("[windows] readInputEvents: Read %v events", nEvents)
	return events[:nEvents], nil
}



var arrowKeyMapPrefix = map[uint16]string{
	winterm.VK_UP:    "%s%sA",
	winterm.VK_DOWN:  "%s%sB",
	winterm.VK_RIGHT: "%s%sC",
	winterm.VK_LEFT:  "%s%sD",
}

var keyMapPrefix = map[uint16]string{
	winterm.VK_UP:     "\x1B[%sA",
	winterm.VK_DOWN:   "\x1B[%sB",
	winterm.VK_RIGHT:  "\x1B[%sC",
	winterm.VK_LEFT:   "\x1B[%sD",
	winterm.VK_HOME:   "\x1B[1%s~", 
	winterm.VK_END:    "\x1B[4%s~", 
	winterm.VK_INSERT: "\x1B[2%s~",
	winterm.VK_DELETE: "\x1B[3%s~",
	winterm.VK_PRIOR:  "\x1B[5%s~",
	winterm.VK_NEXT:   "\x1B[6%s~",
	winterm.VK_F1:     "",
	winterm.VK_F2:     "",
	winterm.VK_F3:     "\x1B[13%s~",
	winterm.VK_F4:     "\x1B[14%s~",
	winterm.VK_F5:     "\x1B[15%s~",
	winterm.VK_F6:     "\x1B[17%s~",
	winterm.VK_F7:     "\x1B[18%s~",
	winterm.VK_F8:     "\x1B[19%s~",
	winterm.VK_F9:     "\x1B[20%s~",
	winterm.VK_F10:    "\x1B[21%s~",
	winterm.VK_F11:    "\x1B[23%s~",
	winterm.VK_F12:    "\x1B[24%s~",
}


func translateKeyEvents(events []winterm.INPUT_RECORD, escapeSequence []byte) []byte {
	var buffer bytes.Buffer
	for _, event := range events {
		if event.EventType == winterm.KEY_EVENT && event.KeyEvent.KeyDown != 0 {
			buffer.WriteString(keyToString(&event.KeyEvent, escapeSequence))
		}
	}

	return buffer.Bytes()
}


func keyToString(keyEvent *winterm.KEY_EVENT_RECORD, escapeSequence []byte) string {
	if keyEvent.UnicodeChar == 0 {
		return formatVirtualKey(keyEvent.VirtualKeyCode, keyEvent.ControlKeyState, escapeSequence)
	}

	_, alt, control := getControlKeys(keyEvent.ControlKeyState)
	if control {
		
		
		
		
		
		
		

	}

	
	if !control && alt {
		return ansiterm.KEY_ESC_N + strings.ToLower(string(keyEvent.UnicodeChar))
	}

	return string(keyEvent.UnicodeChar)
}


func formatVirtualKey(key uint16, controlState uint32, escapeSequence []byte) string {
	shift, alt, control := getControlKeys(controlState)
	modifier := getControlKeysModifier(shift, alt, control)

	if format, ok := arrowKeyMapPrefix[key]; ok {
		return fmt.Sprintf(format, escapeSequence, modifier)
	}

	if format, ok := keyMapPrefix[key]; ok {
		return fmt.Sprintf(format, modifier)
	}

	return ""
}


func getControlKeys(controlState uint32) (shift, alt, control bool) {
	shift = 0 != (controlState & winterm.SHIFT_PRESSED)
	alt = 0 != (controlState & (winterm.LEFT_ALT_PRESSED | winterm.RIGHT_ALT_PRESSED))
	control = 0 != (controlState & (winterm.LEFT_CTRL_PRESSED | winterm.RIGHT_CTRL_PRESSED))
	return shift, alt, control
}


func getControlKeysModifier(shift, alt, control bool) string {
	if shift && alt && control {
		return ansiterm.KEY_CONTROL_PARAM_8
	}
	if alt && control {
		return ansiterm.KEY_CONTROL_PARAM_7
	}
	if shift && control {
		return ansiterm.KEY_CONTROL_PARAM_6
	}
	if control {
		return ansiterm.KEY_CONTROL_PARAM_5
	}
	if shift && alt {
		return ansiterm.KEY_CONTROL_PARAM_4
	}
	if alt {
		return ansiterm.KEY_CONTROL_PARAM_3
	}
	if shift {
		return ansiterm.KEY_CONTROL_PARAM_2
	}
	return ""
}
