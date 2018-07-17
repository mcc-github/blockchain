












package procfs

import (
	"fmt"
	"io/ioutil"
	"os"
)


type ProcIO struct {
	
	RChar uint64
	
	WChar uint64
	
	SyscR uint64
	
	SyscW uint64
	
	ReadBytes uint64
	
	WriteBytes uint64
	
	
	
	CancelledWriteBytes int64
}


func (p Proc) NewIO() (ProcIO, error) {
	pio := ProcIO{}

	f, err := os.Open(p.path("io"))
	if err != nil {
		return pio, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return pio, err
	}

	ioFormat := "rchar: %d\nwchar: %d\nsyscr: %d\nsyscw: %d\n" +
		"read_bytes: %d\nwrite_bytes: %d\n" +
		"cancelled_write_bytes: %d\n"

	_, err = fmt.Sscanf(string(data), ioFormat, &pio.RChar, &pio.WChar, &pio.SyscR,
		&pio.SyscW, &pio.ReadBytes, &pio.WriteBytes, &pio.CancelledWriteBytes)

	return pio, err
}
