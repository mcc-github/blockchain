





package unix

import "syscall"



func seek(fd int, offset int64, whence int) (newoffset int64, err syscall.Errno)

func socketcall(call int, a0, a1, a2, a3, a4, a5 uintptr) (n int, err syscall.Errno)
func rawsocketcall(call int, a0, a1, a2, a3, a4, a5 uintptr) (n int, err syscall.Errno)
