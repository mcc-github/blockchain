
package gexec

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

const INVALID_EXIT_CODE = 254

type Session struct {
	
	Command *exec.Cmd

	
	Out *gbytes.Buffer

	
	Err *gbytes.Buffer

	
	Exited <-chan struct{}

	lock     *sync.Mutex
	exitCode int
}


func Start(command *exec.Cmd, outWriter io.Writer, errWriter io.Writer) (*Session, error) {
	exited := make(chan struct{})

	session := &Session{
		Command:  command,
		Out:      gbytes.NewBuffer(),
		Err:      gbytes.NewBuffer(),
		Exited:   exited,
		lock:     &sync.Mutex{},
		exitCode: -1,
	}

	var commandOut, commandErr io.Writer

	commandOut, commandErr = session.Out, session.Err

	if outWriter != nil {
		commandOut = io.MultiWriter(commandOut, outWriter)
	}

	if errWriter != nil {
		commandErr = io.MultiWriter(commandErr, errWriter)
	}

	command.Stdout = commandOut
	command.Stderr = commandErr

	err := command.Start()
	if err == nil {
		go session.monitorForExit(exited)
		trackedSessionsMutex.Lock()
		defer trackedSessionsMutex.Unlock()
		trackedSessions = append(trackedSessions, session)
	}

	return session, err
}


func (s *Session) Buffer() *gbytes.Buffer {
	return s.Out
}


func (s *Session) ExitCode() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.exitCode
}


func (s *Session) Wait(timeout ...interface{}) *Session {
	EventuallyWithOffset(1, s, timeout...).Should(Exit())
	return s
}


func (s *Session) Kill() *Session {
	if s.ExitCode() != -1 {
		return s
	}
	s.Command.Process.Kill()
	return s
}


func (s *Session) Interrupt() *Session {
	return s.Signal(syscall.SIGINT)
}


func (s *Session) Terminate() *Session {
	return s.Signal(syscall.SIGTERM)
}


func (s *Session) Signal(signal os.Signal) *Session {
	if s.ExitCode() != -1 {
		return s
	}
	s.Command.Process.Signal(signal)
	return s
}

func (s *Session) monitorForExit(exited chan<- struct{}) {
	err := s.Command.Wait()
	s.lock.Lock()
	s.Out.Close()
	s.Err.Close()
	status := s.Command.ProcessState.Sys().(syscall.WaitStatus)
	if status.Signaled() {
		s.exitCode = 128 + int(status.Signal())
	} else {
		exitStatus := status.ExitStatus()
		if exitStatus == -1 && err != nil {
			s.exitCode = INVALID_EXIT_CODE
		}
		s.exitCode = exitStatus
	}
	s.lock.Unlock()

	close(exited)
}

var trackedSessions = []*Session{}
var trackedSessionsMutex = &sync.Mutex{}


func KillAndWait(timeout ...interface{}) {
	trackedSessionsMutex.Lock()
	defer trackedSessionsMutex.Unlock()
	for _, session := range trackedSessions {
		session.Kill().Wait(timeout...)
	}
	trackedSessions = []*Session{}
}


func TerminateAndWait(timeout ...interface{}) {
	trackedSessionsMutex.Lock()
	defer trackedSessionsMutex.Unlock()
	for _, session := range trackedSessions {
		session.Terminate().Wait(timeout...)
	}
}


func Kill() {
	trackedSessionsMutex.Lock()
	defer trackedSessionsMutex.Unlock()
	for _, session := range trackedSessions {
		session.Kill()
	}
}


func Terminate() {
	trackedSessionsMutex.Lock()
	defer trackedSessionsMutex.Unlock()
	for _, session := range trackedSessions {
		session.Terminate()
	}
}


func Signal(signal os.Signal) {
	trackedSessionsMutex.Lock()
	defer trackedSessionsMutex.Unlock()
	for _, session := range trackedSessions {
		session.Signal(signal)
	}
}


func Interrupt() {
	trackedSessionsMutex.Lock()
	defer trackedSessionsMutex.Unlock()
	for _, session := range trackedSessions {
		session.Interrupt()
	}
}
