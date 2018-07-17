




package gotty



import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
)




func OpenTermInfo(termName string) (*TermInfo, error) {
	var term *TermInfo
	var err error
	
	termloc := os.Getenv("TERMINFO")
	if len(termloc) == 0 {
		
		locations := []string{os.Getenv("HOME") + "/.terminfo/", "/etc/terminfo/",
			"/lib/terminfo/", "/usr/share/terminfo/"}
		var path string
		for _, str := range locations {
			
			path = str + string(termName[0]) + "/" + termName
			
			file, _ := os.Open(path)
			if file != nil {
				
				file.Close()
				break
			}
		}
		if len(path) > 0 {
			term, err = readTermInfo(path)
		} else {
			err = errors.New(fmt.Sprintf("No terminfo file(-location) found"))
		}
	}
	return term, err
}




func OpenTermInfoEnv() (*TermInfo, error) {
	termenv := os.Getenv("TERM")
	return OpenTermInfo(termenv)
}



func (term *TermInfo) GetAttribute(attr string) (stacker, error) {
	
	var value stacker
	
	var block sync.WaitGroup
	
	written := false
	
	f := func(ats interface{}) {
		var ok bool
		var v stacker
		
		switch reflect.TypeOf(ats).Elem().Kind() {
		case reflect.Bool:
			v, ok = ats.(map[string]bool)[attr]
		case reflect.Int16:
			v, ok = ats.(map[string]int16)[attr]
		case reflect.String:
			v, ok = ats.(map[string]string)[attr]
		}
		
		if ok {
			value = v
			written = true
		}
		
		block.Done()
	}
	block.Add(3)
	
	go f(term.boolAttributes)
	go f(term.numAttributes)
	go f(term.strAttributes)
	
	block.Wait()
	
	if written {
		return value, nil
	}
	
	return nil, fmt.Errorf("Erorr finding attribute")
}



func (term *TermInfo) GetAttributeName(name string) (stacker, error) {
	tc := GetTermcapName(name)
	return term.GetAttribute(tc)
}



func GetTermcapName(name string) string {
	
	var tc string
	
	var wait sync.WaitGroup
	
	f := func(attrs []string) {
		
		for i, s := range attrs {
			if s == name {
				tc = attrs[i+1]
			}
		}
		
		wait.Done()
	}
	wait.Add(3)
	
	go f(BoolAttr[:])
	go f(NumAttr[:])
	go f(StrAttr[:])
	
	wait.Wait()
	
	return tc
}



func readTermInfo(path string) (*TermInfo, error) {
	
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	
	
	
	
	var header [6]int16
	
	var byteArray []byte
	
	var shArray []int16
	
	var term TermInfo

	
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}
	
	if header[0] != 0432 {
		return nil, errors.New(fmt.Sprintf("Wrong filetype"))
	}

	
	byteArray = make([]byte, header[1])
	err = binary.Read(file, binary.LittleEndian, &byteArray)
	if err != nil {
		return nil, err
	}
	term.Names = strings.Split(string(byteArray), "|")

	
	byteArray = make([]byte, header[2])
	err = binary.Read(file, binary.LittleEndian, &byteArray)
	if err != nil {
		return nil, err
	}
	term.boolAttributes = make(map[string]bool)
	for i, b := range byteArray {
		if b == 1 {
			term.boolAttributes[BoolAttr[i*2+1]] = true
		}
	}
	
	if len(byteArray)%2 != 0 {
		err = binary.Read(file, binary.LittleEndian, make([]byte, 1))
		if err != nil {
			return nil, err
		}
	}

	
	shArray = make([]int16, header[3])
	err = binary.Read(file, binary.LittleEndian, &shArray)
	if err != nil {
		return nil, err
	}
	term.numAttributes = make(map[string]int16)
	for i, n := range shArray {
		if n != 0377 && n > -1 {
			term.numAttributes[NumAttr[i*2+1]] = n
		}
	}

	
	shArray = make([]int16, header[4])
	err = binary.Read(file, binary.LittleEndian, &shArray)
	if err != nil {
		return nil, err
	}
	
	byteArray = make([]byte, header[5])
	err = binary.Read(file, binary.LittleEndian, &byteArray)
	if err != nil {
		return nil, err
	}
	term.strAttributes = make(map[string]string)
	
	for i, offset := range shArray {
		if offset > -1 {
			r := offset
			for ; byteArray[r] != 0; r++ {
			}
			term.strAttributes[StrAttr[i*2+1]] = string(byteArray[offset:r])
		}
	}
	return &term, nil
}
