




package unix

import "unsafe"


type PtraceRegsMipsle struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}


func PtraceGetRegsMipsle(pid int, regsout *PtraceRegsMipsle) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}


func PtraceSetRegsMipsle(pid int, regs *PtraceRegsMipsle) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}


type PtraceRegsMips64le struct {
	Regs     [32]uint64
	Lo       uint64
	Hi       uint64
	Epc      uint64
	Badvaddr uint64
	Status   uint64
	Cause    uint64
}


func PtraceGetRegsMips64le(pid int, regsout *PtraceRegsMips64le) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}


func PtraceSetRegsMips64le(pid int, regs *PtraceRegsMips64le) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}
