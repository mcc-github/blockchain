

package spew

import (
	"bytes"
	"fmt"
	"io"
	"os"
)












type ConfigState struct {
	
	
	
	
	Indent string

	
	
	
	
	
	
	MaxDepth int

	
	
	DisableMethods bool

	
	
	
	
	
	
	
	
	
	
	
	
	DisablePointerMethods bool

	
	
	DisablePointerAddresses bool

	
	
	
	DisableCapacities bool

	
	
	
	
	
	
	
	
	ContinueOnMethod bool

	
	
	
	
	
	
	SortKeys bool

	
	
	
	SpewKeys bool
}



var Config = ConfigState{Indent: " "}









func (c *ConfigState) Errorf(format string, a ...interface{}) (err error) {
	return fmt.Errorf(format, c.convertArgs(a)...)
}









func (c *ConfigState) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, c.convertArgs(a)...)
}









func (c *ConfigState) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, c.convertArgs(a)...)
}








func (c *ConfigState) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, c.convertArgs(a)...)
}









func (c *ConfigState) Print(a ...interface{}) (n int, err error) {
	return fmt.Print(c.convertArgs(a)...)
}









func (c *ConfigState) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, c.convertArgs(a)...)
}









func (c *ConfigState) Println(a ...interface{}) (n int, err error) {
	return fmt.Println(c.convertArgs(a)...)
}








func (c *ConfigState) Sprint(a ...interface{}) string {
	return fmt.Sprint(c.convertArgs(a)...)
}








func (c *ConfigState) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, c.convertArgs(a)...)
}








func (c *ConfigState) Sprintln(a ...interface{}) string {
	return fmt.Sprintln(c.convertArgs(a)...)
}


func (c *ConfigState) NewFormatter(v interface{}) fmt.Formatter {
	return newFormatter(c, v)
}



func (c *ConfigState) Fdump(w io.Writer, a ...interface{}) {
	fdump(c, w, a...)
}


func (c *ConfigState) Dump(a ...interface{}) {
	fdump(c, os.Stdout, a...)
}



func (c *ConfigState) Sdump(a ...interface{}) string {
	var buf bytes.Buffer
	fdump(c, &buf, a...)
	return buf.String()
}




func (c *ConfigState) convertArgs(args []interface{}) (formatters []interface{}) {
	formatters = make([]interface{}, len(args))
	for index, arg := range args {
		formatters[index] = newFormatter(c, arg)
	}
	return formatters
}









func NewDefaultConfig() *ConfigState {
	return &ConfigState{Indent: " "}
}
