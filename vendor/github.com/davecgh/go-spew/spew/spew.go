

package spew

import (
	"fmt"
	"io"
)









func Errorf(format string, a ...interface{}) (err error) {
	return fmt.Errorf(format, convertArgs(a)...)
}









func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, convertArgs(a)...)
}









func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format, convertArgs(a)...)
}








func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, convertArgs(a)...)
}









func Print(a ...interface{}) (n int, err error) {
	return fmt.Print(convertArgs(a)...)
}









func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, convertArgs(a)...)
}









func Println(a ...interface{}) (n int, err error) {
	return fmt.Println(convertArgs(a)...)
}








func Sprint(a ...interface{}) string {
	return fmt.Sprint(convertArgs(a)...)
}








func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, convertArgs(a)...)
}








func Sprintln(a ...interface{}) string {
	return fmt.Sprintln(convertArgs(a)...)
}



func convertArgs(args []interface{}) (formatters []interface{}) {
	formatters = make([]interface{}, len(args))
	for index, arg := range args {
		formatters[index] = NewFormatter(arg)
	}
	return formatters
}
