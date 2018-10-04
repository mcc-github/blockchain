










package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {
	
	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH_TARGET")
	if goarch == "" {
		goarch = os.Getenv("GOARCH")
	}
	
	if goos == "linux" && goarch != "sparc64" {
		if os.Getenv("GOLANG_SYS_BUILD") != "docker" {
			os.Stderr.WriteString("In the new build system, mkpost should not be called directly.\n")
			os.Stderr.WriteString("See README.md\n")
			os.Exit(1)
		}
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	
	valRegex := regexp.MustCompile(`type (Fsid|Sigset_t) struct {(\s+)X__val(\s+\S+\s+)}`)
	b = valRegex.ReplaceAll(b, []byte("type $1 struct {${2}Val$3}"))

	
	
	ptraceRexexp := regexp.MustCompile(`type Ptrace((Psw|Fpregs|Per) struct {\s*})`)
	b = ptraceRexexp.ReplaceAll(b, nil)

	
	controlRegsRegex := regexp.MustCompile(`(Control_regs)\s+\[0\]uint64`)
	b = controlRegsRegex.ReplaceAll(b, []byte("_ [0]uint64"))

	
	
	removeFieldsRegex := regexp.MustCompile(`X__glibc\S*`)
	b = removeFieldsRegex.ReplaceAll(b, []byte("_"))

	
	
	convertUtsnameRegex := regexp.MustCompile(`((Sys|Node|Domain)name|Release|Version|Machine)(\s+)\[(\d+)\]u?int8`)
	b = convertUtsnameRegex.ReplaceAll(b, []byte("$1$3[$4]byte"))

	
	spareFieldsRegex := regexp.MustCompile(`X__spare\S*`)
	b = spareFieldsRegex.ReplaceAll(b, []byte("_"))

	
	removePaddingFieldsRegex := regexp.MustCompile(`Pad_cgo_\d+`)
	b = removePaddingFieldsRegex.ReplaceAll(b, []byte("_"))

	
	removeFieldsRegex = regexp.MustCompile(`\b(X_\S+|Padding)`)
	b = removeFieldsRegex.ReplaceAll(b, []byte("_"))

	
	b = b[bytes.IndexByte(b, '\n')+1:]
	
	
	replacement := fmt.Sprintf(`$1 | go run mkpost.go



	cgoCommandRegex := regexp.MustCompile(`(cgo -godefs .*)`)
	b = cgoCommandRegex.ReplaceAll(b, []byte(replacement))

	
	b, err = format.Source(b)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(b)
}
