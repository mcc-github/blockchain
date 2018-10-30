





















package journal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)


type Priority int

const (
	PriEmerg Priority = iota
	PriAlert
	PriCrit
	PriErr
	PriWarning
	PriNotice
	PriInfo
	PriDebug
)

var conn net.Conn

func init() {
	var err error
	conn, err = net.Dial("unixgram", "/run/systemd/journal/socket")
	if err != nil {
		conn = nil
	}
}


func Enabled() bool {
	return conn != nil
}








func Send(message string, priority Priority, vars map[string]string) error {
	if conn == nil {
		return journalError("could not connect to journald socket")
	}

	data := new(bytes.Buffer)
	appendVariable(data, "PRIORITY", strconv.Itoa(int(priority)))
	appendVariable(data, "MESSAGE", message)
	for k, v := range vars {
		appendVariable(data, k, v)
	}

	_, err := io.Copy(conn, data)
	if err != nil && isSocketSpaceError(err) {
		file, err := tempFd()
		if err != nil {
			return journalError(err.Error())
		}
		defer file.Close()
		_, err = io.Copy(file, data)
		if err != nil {
			return journalError(err.Error())
		}

		rights := syscall.UnixRights(int(file.Fd()))

		
		unixConn, ok := conn.(*net.UnixConn)
		if !ok {
			return journalError("can't send file through non-Unix connection")
		}
		_, _, err = unixConn.WriteMsgUnix([]byte{}, rights, nil)
		if err != nil {
			return journalError(err.Error())
		}
	} else if err != nil {
		return journalError(err.Error())
	}
	return nil
}


func Print(priority Priority, format string, a ...interface{}) error {
	return Send(fmt.Sprintf(format, a...), priority, nil)
}

func appendVariable(w io.Writer, name, value string) {
	if !validVarName(name) {
		journalError("variable name contains invalid character, ignoring")
	}
	if strings.ContainsRune(value, '\n') {
		
		fmt.Fprintln(w, name)
		binary.Write(w, binary.LittleEndian, uint64(len(value)))
		fmt.Fprintln(w, value)
	} else {
		
		fmt.Fprintf(w, "%s=%s\n", name, value)
	}
}

func validVarName(name string) bool {
	

	valid := name[0] != '_'
	for _, c := range name {
		valid = valid && ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '_'
	}
	return valid
}

func isSocketSpaceError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	sysErr, ok := opErr.Err.(syscall.Errno)
	if !ok {
		return false
	}

	return sysErr == syscall.EMSGSIZE || sysErr == syscall.ENOBUFS
}

func tempFd() (*os.File, error) {
	file, err := ioutil.TempFile("/dev/shm/", "journal.XXXXX")
	if err != nil {
		return nil, err
	}
	err = syscall.Unlink(file.Name())
	if err != nil {
		return nil, err
	}
	return file, nil
}

func journalError(s string) error {
	s = "journal error: " + s
	fmt.Fprintln(os.Stderr, s)
	return errors.New(s)
}
