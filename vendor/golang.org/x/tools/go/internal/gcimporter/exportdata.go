







package gcimporter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readGopackHeader(r *bufio.Reader) (name string, size int, err error) {
	
	hdr := make([]byte, 16+12+6+6+8+10+2)
	_, err = io.ReadFull(r, hdr)
	if err != nil {
		return
	}
	
	if false {
		fmt.Printf("header: %s", hdr)
	}
	s := strings.TrimSpace(string(hdr[16+12+6+6+8:][:10]))
	size, err = strconv.Atoi(s)
	if err != nil || hdr[len(hdr)-2] != '`' || hdr[len(hdr)-1] != '\n' {
		err = fmt.Errorf("invalid archive header")
		return
	}
	name = strings.TrimSpace(string(hdr[:16]))
	return
}







func FindExportData(r *bufio.Reader) (hdr string, err error) {
	
	line, err := r.ReadSlice('\n')
	if err != nil {
		err = fmt.Errorf("can't find export data (%v)", err)
		return
	}

	if string(line) == "!<arch>\n" {
		
		var name string
		if name, _, err = readGopackHeader(r); err != nil {
			return
		}

		
		if name != "__.PKGDEF" {
			err = fmt.Errorf("go archive is missing __.PKGDEF")
			return
		}

		
		
		if line, err = r.ReadSlice('\n'); err != nil {
			err = fmt.Errorf("can't find export data (%v)", err)
			return
		}
	}

	
	
	if !strings.HasPrefix(string(line), "go object ") {
		err = fmt.Errorf("not a Go object file")
		return
	}

	
	
	for line[0] != '$' {
		if line, err = r.ReadSlice('\n'); err != nil {
			err = fmt.Errorf("can't find export data (%v)", err)
			return
		}
	}
	hdr = string(line)

	return
}
