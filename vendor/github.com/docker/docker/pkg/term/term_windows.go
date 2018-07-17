package term 

import (
	"io"
	"os"
	"os/signal"
	"syscall" 

	"github.com/Azure/go-ansiterm/winterm"
	"github.com/docker/docker/pkg/term/windows"
)


type State struct {
	mode uint32
}


type Winsize struct {
	Height uint16
	Width  uint16
}


var vtInputSupported bool


func StdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	
	
	var emulateStdin, emulateStdout, emulateStderr bool
	fd := os.Stdin.Fd()
	if mode, err := winterm.GetConsoleMode(fd); err == nil {
		
		if err = winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_INPUT); err != nil {
			emulateStdin = true
		} else {
			vtInputSupported = true
		}
		
		
		winterm.SetConsoleMode(fd, mode)
	}

	fd = os.Stdout.Fd()
	if mode, err := winterm.GetConsoleMode(fd); err == nil {
		
		if err = winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING|winterm.DISABLE_NEWLINE_AUTO_RETURN); err != nil {
			emulateStdout = true
		} else {
			winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		}
	}

	fd = os.Stderr.Fd()
	if mode, err := winterm.GetConsoleMode(fd); err == nil {
		
		if err = winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING|winterm.DISABLE_NEWLINE_AUTO_RETURN); err != nil {
			emulateStderr = true
		} else {
			winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		}
	}

	if os.Getenv("ConEmuANSI") == "ON" || os.Getenv("ConsoleZVersion") != "" {
		
		emulateStdin = true
		emulateStdout = false
		emulateStderr = false
	}

	
	
	
	
	if emulateStdin {
		stdIn = windowsconsole.NewAnsiReader(syscall.STD_INPUT_HANDLE)
	} else {
		stdIn = os.Stdin
	}

	if emulateStdout {
		stdOut = windowsconsole.NewAnsiWriter(syscall.STD_OUTPUT_HANDLE)
	} else {
		stdOut = os.Stdout
	}

	if emulateStderr {
		stdErr = windowsconsole.NewAnsiWriter(syscall.STD_ERROR_HANDLE)
	} else {
		stdErr = os.Stderr
	}

	return
}


func GetFdInfo(in interface{}) (uintptr, bool) {
	return windowsconsole.GetHandleInfo(in)
}


func GetWinsize(fd uintptr) (*Winsize, error) {
	info, err := winterm.GetConsoleScreenBufferInfo(fd)
	if err != nil {
		return nil, err
	}

	winsize := &Winsize{
		Width:  uint16(info.Window.Right - info.Window.Left + 1),
		Height: uint16(info.Window.Bottom - info.Window.Top + 1),
	}

	return winsize, nil
}


func IsTerminal(fd uintptr) bool {
	return windowsconsole.IsConsole(fd)
}



func RestoreTerminal(fd uintptr, state *State) error {
	return winterm.SetConsoleMode(fd, state.mode)
}


func SaveState(fd uintptr) (*State, error) {
	mode, e := winterm.GetConsoleMode(fd)
	if e != nil {
		return nil, e
	}

	return &State{mode: mode}, nil
}



func DisableEcho(fd uintptr, state *State) error {
	mode := state.mode
	mode &^= winterm.ENABLE_ECHO_INPUT
	mode |= winterm.ENABLE_PROCESSED_INPUT | winterm.ENABLE_LINE_INPUT
	err := winterm.SetConsoleMode(fd, mode)
	if err != nil {
		return err
	}

	
	restoreAtInterrupt(fd, state)
	return nil
}




func SetRawTerminal(fd uintptr) (*State, error) {
	state, err := MakeRaw(fd)
	if err != nil {
		return nil, err
	}

	
	restoreAtInterrupt(fd, state)
	return state, err
}




func SetRawTerminalOutput(fd uintptr) (*State, error) {
	state, err := SaveState(fd)
	if err != nil {
		return nil, err
	}

	
	
	winterm.SetConsoleMode(fd, state.mode|winterm.DISABLE_NEWLINE_AUTO_RETURN)
	return state, err
}



func MakeRaw(fd uintptr) (*State, error) {
	state, err := SaveState(fd)
	if err != nil {
		return nil, err
	}

	mode := state.mode

	
	
	

	
	mode &^= winterm.ENABLE_ECHO_INPUT
	mode &^= winterm.ENABLE_LINE_INPUT
	mode &^= winterm.ENABLE_MOUSE_INPUT
	mode &^= winterm.ENABLE_WINDOW_INPUT
	mode &^= winterm.ENABLE_PROCESSED_INPUT

	
	mode |= winterm.ENABLE_EXTENDED_FLAGS
	mode |= winterm.ENABLE_INSERT_MODE
	mode |= winterm.ENABLE_QUICK_EDIT_MODE
	if vtInputSupported {
		mode |= winterm.ENABLE_VIRTUAL_TERMINAL_INPUT
	}

	err = winterm.SetConsoleMode(fd, mode)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func restoreAtInterrupt(fd uintptr, state *State) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	go func() {
		_ = <-sigchan
		RestoreTerminal(fd, state)
		os.Exit(0)
	}()
}
