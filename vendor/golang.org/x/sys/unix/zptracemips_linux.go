




package unix

import "unsafe"


type PtraceRegsMips struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}


func PtraceGetRegsMips(pid int, regsout *PtraceRegsMips) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}


func PtraceSetRegsMips(pid int, regs *PtraceRegsMips) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}


type PtraceRegsMips64 struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}


func PtraceGetRegsMips64(pid int, regsout *PtraceRegsMips64) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}


func PtraceSetRegsMips64(pid int, regs *PtraceRegsMips64) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}
