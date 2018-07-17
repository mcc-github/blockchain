



package term 

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

var (
	
	ErrInvalidState = errors.New("Invalid terminal state")
)


type State struct {
	termios Termios
}


type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16
	y      uint16
}


func StdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) {
	return os.Stdin, os.Stdout, os.Stderr
}


func GetFdInfo(in interface{}) (uintptr, bool) {
	var inFd uintptr
	var isTerminalIn bool
	if file, ok := in.(*os.File); ok {
		inFd = file.Fd()
		isTerminalIn = IsTerminal(inFd)
	}
	return inFd, isTerminalIn
}


func IsTerminal(fd uintptr) bool {
	var termios Termios
	return tcget(fd, &termios) == 0
}



func RestoreTerminal(fd uintptr, state *State) error {
	if state == nil {
		return ErrInvalidState
	}
	if err := tcset(fd, &state.termios); err != 0 {
		return err
	}
	return nil
}


func SaveState(fd uintptr) (*State, error) {
	var oldState State
	if err := tcget(fd, &oldState.termios); err != 0 {
		return nil, err
	}

	return &oldState, nil
}



func DisableEcho(fd uintptr, state *State) error {
	newState := state.termios
	newState.Lflag &^= unix.ECHO

	if err := tcset(fd, &newState); err != 0 {
		return err
	}
	handleInterrupt(fd, state)
	return nil
}




func SetRawTerminal(fd uintptr) (*State, error) {
	oldState, err := MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	handleInterrupt(fd, oldState)
	return oldState, err
}




func SetRawTerminalOutput(fd uintptr) (*State, error) {
	return nil, nil
}

func handleInterrupt(fd uintptr, state *State) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		for range sigchan {
			
			fmt.Println()
			signal.Stop(sigchan)
			close(sigchan)
			RestoreTerminal(fd, state)
			os.Exit(1)
		}
	}()
}
