




package unix

import "unsafe"


type PtraceRegsArm struct {
	Uregs [18]uint32
}


func PtraceGetRegsArm(pid int, regsout *PtraceRegsArm) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}


func PtraceSetRegsArm(pid int, regs *PtraceRegsArm) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}


type PtraceRegsArm64 struct {
	Regs   [31]uint64
	Sp     uint64
	Pc     uint64
	Pstate uint64
}


func PtraceGetRegsArm64(pid int, regsout *PtraceRegsArm64) error {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}


func PtraceSetRegsArm64(pid int, regs *PtraceRegsArm64) error {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}
