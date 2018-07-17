





package unix


func SyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)



func RawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr)
