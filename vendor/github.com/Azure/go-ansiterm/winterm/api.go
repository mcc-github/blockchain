

package winterm

import (
	"fmt"
	"syscall"
	"unsafe"
)





















var (
	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")

	getConsoleCursorInfoProc       = kernel32DLL.NewProc("GetConsoleCursorInfo")
	setConsoleCursorInfoProc       = kernel32DLL.NewProc("SetConsoleCursorInfo")
	setConsoleCursorPositionProc   = kernel32DLL.NewProc("SetConsoleCursorPosition")
	setConsoleModeProc             = kernel32DLL.NewProc("SetConsoleMode")
	getConsoleScreenBufferInfoProc = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")
	setConsoleScreenBufferSizeProc = kernel32DLL.NewProc("SetConsoleScreenBufferSize")
	scrollConsoleScreenBufferProc  = kernel32DLL.NewProc("ScrollConsoleScreenBufferA")
	setConsoleTextAttributeProc    = kernel32DLL.NewProc("SetConsoleTextAttribute")
	setConsoleWindowInfoProc       = kernel32DLL.NewProc("SetConsoleWindowInfo")
	writeConsoleOutputProc         = kernel32DLL.NewProc("WriteConsoleOutputW")
	readConsoleInputProc           = kernel32DLL.NewProc("ReadConsoleInputW")
	waitForSingleObjectProc        = kernel32DLL.NewProc("WaitForSingleObject")
)


const (
	
	
	ENABLE_PROCESSED_INPUT        = 0x0001
	ENABLE_LINE_INPUT             = 0x0002
	ENABLE_ECHO_INPUT             = 0x0004
	ENABLE_WINDOW_INPUT           = 0x0008
	ENABLE_MOUSE_INPUT            = 0x0010
	ENABLE_INSERT_MODE            = 0x0020
	ENABLE_QUICK_EDIT_MODE        = 0x0040
	ENABLE_EXTENDED_FLAGS         = 0x0080
	ENABLE_AUTO_POSITION          = 0x0100
	ENABLE_VIRTUAL_TERMINAL_INPUT = 0x0200

	ENABLE_PROCESSED_OUTPUT            = 0x0001
	ENABLE_WRAP_AT_EOL_OUTPUT          = 0x0002
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	DISABLE_NEWLINE_AUTO_RETURN        = 0x0008
	ENABLE_LVB_GRID_WORLDWIDE          = 0x0010

	
	
	
	
	
	FOREGROUND_BLUE      uint16 = 0x0001
	FOREGROUND_GREEN     uint16 = 0x0002
	FOREGROUND_RED       uint16 = 0x0004
	FOREGROUND_INTENSITY uint16 = 0x0008
	FOREGROUND_MASK      uint16 = 0x000F

	BACKGROUND_BLUE      uint16 = 0x0010
	BACKGROUND_GREEN     uint16 = 0x0020
	BACKGROUND_RED       uint16 = 0x0040
	BACKGROUND_INTENSITY uint16 = 0x0080
	BACKGROUND_MASK      uint16 = 0x00F0

	COMMON_LVB_MASK          uint16 = 0xFF00
	COMMON_LVB_REVERSE_VIDEO uint16 = 0x4000
	COMMON_LVB_UNDERSCORE    uint16 = 0x8000

	
	
	KEY_EVENT                = 0x0001
	MOUSE_EVENT              = 0x0002
	WINDOW_BUFFER_SIZE_EVENT = 0x0004
	MENU_EVENT               = 0x0008
	FOCUS_EVENT              = 0x0010

	
	WAIT_ABANDONED = 0x00000080
	WAIT_FAILED    = 0xFFFFFFFF
	WAIT_SIGNALED  = 0x0000000
	WAIT_TIMEOUT   = 0x00000102

	
	WAIT_INFINITE       = 0xFFFFFFFF
	WAIT_ONE_SECOND     = 1000
	WAIT_HALF_SECOND    = 500
	WAIT_QUARTER_SECOND = 250
)




type (
	CHAR_INFO struct {
		UnicodeChar uint16
		Attributes  uint16
	}

	CONSOLE_CURSOR_INFO struct {
		Size    uint32
		Visible int32
	}

	CONSOLE_SCREEN_BUFFER_INFO struct {
		Size              COORD
		CursorPosition    COORD
		Attributes        uint16
		Window            SMALL_RECT
		MaximumWindowSize COORD
	}

	COORD struct {
		X int16
		Y int16
	}

	SMALL_RECT struct {
		Left   int16
		Top    int16
		Right  int16
		Bottom int16
	}

	
	
	INPUT_RECORD struct {
		EventType uint16
		KeyEvent  KEY_EVENT_RECORD
	}

	KEY_EVENT_RECORD struct {
		KeyDown         int32
		RepeatCount     uint16
		VirtualKeyCode  uint16
		VirtualScanCode uint16
		UnicodeChar     uint16
		ControlKeyState uint32
	}

	WINDOW_BUFFER_SIZE struct {
		Size COORD
	}
)


func boolToBOOL(f bool) int32 {
	if f {
		return int32(1)
	} else {
		return int32(0)
	}
}



func GetConsoleCursorInfo(handle uintptr, cursorInfo *CONSOLE_CURSOR_INFO) error {
	r1, r2, err := getConsoleCursorInfoProc.Call(handle, uintptr(unsafe.Pointer(cursorInfo)), 0)
	return checkError(r1, r2, err)
}



func SetConsoleCursorInfo(handle uintptr, cursorInfo *CONSOLE_CURSOR_INFO) error {
	r1, r2, err := setConsoleCursorInfoProc.Call(handle, uintptr(unsafe.Pointer(cursorInfo)), 0)
	return checkError(r1, r2, err)
}



func SetConsoleCursorPosition(handle uintptr, coord COORD) error {
	r1, r2, err := setConsoleCursorPositionProc.Call(handle, coordToPointer(coord))
	use(coord)
	return checkError(r1, r2, err)
}



func GetConsoleMode(handle uintptr) (mode uint32, err error) {
	err = syscall.GetConsoleMode(syscall.Handle(handle), &mode)
	return mode, err
}



func SetConsoleMode(handle uintptr, mode uint32) error {
	r1, r2, err := setConsoleModeProc.Call(handle, uintptr(mode), 0)
	use(mode)
	return checkError(r1, r2, err)
}



func GetConsoleScreenBufferInfo(handle uintptr) (*CONSOLE_SCREEN_BUFFER_INFO, error) {
	info := CONSOLE_SCREEN_BUFFER_INFO{}
	err := checkError(getConsoleScreenBufferInfoProc.Call(handle, uintptr(unsafe.Pointer(&info)), 0))
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func ScrollConsoleScreenBuffer(handle uintptr, scrollRect SMALL_RECT, clipRect SMALL_RECT, destOrigin COORD, char CHAR_INFO) error {
	r1, r2, err := scrollConsoleScreenBufferProc.Call(handle, uintptr(unsafe.Pointer(&scrollRect)), uintptr(unsafe.Pointer(&clipRect)), coordToPointer(destOrigin), uintptr(unsafe.Pointer(&char)))
	use(scrollRect)
	use(clipRect)
	use(destOrigin)
	use(char)
	return checkError(r1, r2, err)
}



func SetConsoleScreenBufferSize(handle uintptr, coord COORD) error {
	r1, r2, err := setConsoleScreenBufferSizeProc.Call(handle, coordToPointer(coord))
	use(coord)
	return checkError(r1, r2, err)
}




func SetConsoleTextAttribute(handle uintptr, attribute uint16) error {
	r1, r2, err := setConsoleTextAttributeProc.Call(handle, uintptr(attribute), 0)
	use(attribute)
	return checkError(r1, r2, err)
}




func SetConsoleWindowInfo(handle uintptr, isAbsolute bool, rect SMALL_RECT) error {
	r1, r2, err := setConsoleWindowInfoProc.Call(handle, uintptr(boolToBOOL(isAbsolute)), uintptr(unsafe.Pointer(&rect)))
	use(isAbsolute)
	use(rect)
	return checkError(r1, r2, err)
}



func WriteConsoleOutput(handle uintptr, buffer []CHAR_INFO, bufferSize COORD, bufferCoord COORD, writeRegion *SMALL_RECT) error {
	r1, r2, err := writeConsoleOutputProc.Call(handle, uintptr(unsafe.Pointer(&buffer[0])), coordToPointer(bufferSize), coordToPointer(bufferCoord), uintptr(unsafe.Pointer(writeRegion)))
	use(buffer)
	use(bufferSize)
	use(bufferCoord)
	return checkError(r1, r2, err)
}



func ReadConsoleInput(handle uintptr, buffer []INPUT_RECORD, count *uint32) error {
	r1, r2, err := readConsoleInputProc.Call(handle, uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)), uintptr(unsafe.Pointer(count)))
	use(buffer)
	return checkError(r1, r2, err)
}




func WaitForSingleObject(handle uintptr, msWait uint32) (bool, error) {
	r1, _, err := waitForSingleObjectProc.Call(handle, uintptr(uint32(msWait)))
	switch r1 {
	case WAIT_ABANDONED, WAIT_TIMEOUT:
		return false, nil
	case WAIT_SIGNALED:
		return true, nil
	}
	use(msWait)
	return false, err
}


func (info CONSOLE_SCREEN_BUFFER_INFO) String() string {
	return fmt.Sprintf("Size(%v) Cursor(%v) Window(%v) Max(%v)", info.Size, info.CursorPosition, info.Window, info.MaximumWindowSize)
}

func (coord COORD) String() string {
	return fmt.Sprintf("%v,%v", coord.X, coord.Y)
}

func (rect SMALL_RECT) String() string {
	return fmt.Sprintf("(%v,%v),(%v,%v)", rect.Left, rect.Top, rect.Right, rect.Bottom)
}


func checkError(r1, r2 uintptr, err error) error {
	
	if r1 != 0 {
		return nil
	}

	
	if err != nil {
		return err
	}
	return syscall.EINVAL
}


func coordToPointer(c COORD) uintptr {
	
	return uintptr(*((*uint32)(unsafe.Pointer(&c))))
}



func use(p interface{}) {}
