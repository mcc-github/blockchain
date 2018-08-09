





















package exit

import "os"

var real = func() { os.Exit(1) }



func Exit() {
	real()
}


type StubbedExit struct {
	Exited bool
	prev   func()
}


func Stub() *StubbedExit {
	s := &StubbedExit{prev: real}
	real = s.exit
	return s
}



func WithStub(f func()) *StubbedExit {
	s := Stub()
	defer s.Unstub()
	f()
	return s
}


func (se *StubbedExit) Unstub() {
	real = se.prev
}

func (se *StubbedExit) exit() {
	se.Exited = true
}
