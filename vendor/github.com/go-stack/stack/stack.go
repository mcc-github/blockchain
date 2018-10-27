










package stack

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)


type Call struct {
	frame runtime.Frame
}




func Caller(skip int) Call {
	
	
	
	
	
	
	
	
	
	var pcs [3]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	frame, _ = frames.Next()

	return Call{
		frame: frame,
	}
}


func (c Call) String() string {
	return fmt.Sprint(c)
}



func (c Call) MarshalText() ([]byte, error) {
	if c.frame == (runtime.Frame{}) {
		return nil, ErrNoFunc
	}

	buf := bytes.Buffer{}
	fmt.Fprint(&buf, c)
	return buf.Bytes(), nil
}



var ErrNoFunc = errors.New("no call stack information")



















func (c Call) Format(s fmt.State, verb rune) {
	if c.frame == (runtime.Frame{}) {
		fmt.Fprintf(s, "%%!%c(NOFUNC)", verb)
		return
	}

	switch verb {
	case 's', 'v':
		file := c.frame.File
		switch {
		case s.Flag('#'):
			
		case s.Flag('+'):
			file = pkgFilePath(&c.frame)
		default:
			const sep = "/"
			if i := strings.LastIndex(file, sep); i != -1 {
				file = file[i+len(sep):]
			}
		}
		io.WriteString(s, file)
		if verb == 'v' {
			buf := [7]byte{':'}
			s.Write(strconv.AppendInt(buf[:1], int64(c.frame.Line), 10))
		}

	case 'd':
		buf := [6]byte{}
		s.Write(strconv.AppendInt(buf[:0], int64(c.frame.Line), 10))

	case 'k':
		name := c.frame.Function
		const pathSep = "/"
		start, end := 0, len(name)
		if i := strings.LastIndex(name, pathSep); i != -1 {
			start = i + len(pathSep)
		}
		const pkgSep = "."
		if i := strings.Index(name[start:], pkgSep); i != -1 {
			end = start + i
		}
		if s.Flag('+') {
			start = 0
		}
		io.WriteString(s, name[start:end])

	case 'n':
		name := c.frame.Function
		if !s.Flag('+') {
			const pathSep = "/"
			if i := strings.LastIndex(name, pathSep); i != -1 {
				name = name[i+len(pathSep):]
			}
			const pkgSep = "."
			if i := strings.Index(name, pkgSep); i != -1 {
				name = name[i+len(pkgSep):]
			}
		}
		io.WriteString(s, name)
	}
}


func (c Call) Frame() runtime.Frame {
	return c.frame
}





func (c Call) PC() uintptr {
	return c.frame.PC
}



type CallStack []Call


func (cs CallStack) String() string {
	return fmt.Sprint(cs)
}

var (
	openBracketBytes  = []byte("[")
	closeBracketBytes = []byte("]")
	spaceBytes        = []byte(" ")
)



func (cs CallStack) MarshalText() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write(openBracketBytes)
	for i, pc := range cs {
		if i > 0 {
			buf.Write(spaceBytes)
		}
		fmt.Fprint(&buf, pc)
	}
	buf.Write(closeBracketBytes)
	return buf.Bytes(), nil
}




func (cs CallStack) Format(s fmt.State, verb rune) {
	s.Write(openBracketBytes)
	for i, pc := range cs {
		if i > 0 {
			s.Write(spaceBytes)
		}
		pc.Format(s, verb)
	}
	s.Write(closeBracketBytes)
}



func Trace() CallStack {
	var pcs [512]uintptr
	n := runtime.Callers(1, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	cs := make(CallStack, 0, n)

	
	
	frame, more := frames.Next()

	for more {
		frame, more = frames.Next()
		cs = append(cs, Call{frame: frame})
	}

	return cs
}



func (cs CallStack) TrimBelow(c Call) CallStack {
	for len(cs) > 0 && cs[0] != c {
		cs = cs[1:]
	}
	return cs
}



func (cs CallStack) TrimAbove(c Call) CallStack {
	for len(cs) > 0 && cs[len(cs)-1] != c {
		cs = cs[:len(cs)-1]
	}
	return cs
}





func pkgIndex(file, funcName string) int {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	const sep = "/"
	i := len(file)
	for n := strings.Count(funcName, sep) + 2; n > 0; n-- {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			i = -len(sep)
			break
		}
	}
	
	return i + len(sep)
}













































func pkgFilePath(frame *runtime.Frame) string {
	pre := pkgPrefix(frame.Function)
	post := pathSuffix(frame.File)
	if pre == "" {
		return post
	}
	return pre + "/" + post
}



func pkgPrefix(funcName string) string {
	const pathSep = "/"
	end := strings.LastIndex(funcName, pathSep)
	if end == -1 {
		return ""
	}
	return funcName[:end]
}


func pathSuffix(path string) string {
	const pathSep = "/"
	lastSep := strings.LastIndex(path, pathSep)
	if lastSep == -1 {
		return path
	}
	return path[strings.LastIndex(path[:lastSep], pathSep)+1:]
}

var runtimePath string

func init() {
	var pcs [3]uintptr
	runtime.Callers(0, pcs[:])
	frames := runtime.CallersFrames(pcs[:])
	frame, _ := frames.Next()
	file := frame.File

	idx := pkgIndex(frame.File, frame.Function)

	runtimePath = file[:idx]
	if runtime.GOOS == "windows" {
		runtimePath = strings.ToLower(runtimePath)
	}
}

func inGoroot(c Call) bool {
	file := c.frame.File
	if len(file) == 0 || file[0] == '?' {
		return true
	}
	if runtime.GOOS == "windows" {
		file = strings.ToLower(file)
	}
	return strings.HasPrefix(file, runtimePath) || strings.HasSuffix(file, "/_testmain.go")
}




func (cs CallStack) TrimRuntime() CallStack {
	for len(cs) > 0 && inGoroot(cs[len(cs)-1]) {
		cs = cs[:len(cs)-1]
	}
	return cs
}
