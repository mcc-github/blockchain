






package pretty

import (
	"fmt"
	"io"
	"log"
	"reflect"
)





func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, wrap(a, false)...)
}





func Fprintf(w io.Writer, format string, a ...interface{}) (n int, error error) {
	return fmt.Fprintf(w, format, wrap(a, false)...)
}






func Log(a ...interface{}) {
	log.Print(wrap(a, true)...)
}





func Logf(format string, a ...interface{}) {
	log.Printf(format, wrap(a, false)...)
}






func Logln(a ...interface{}) {
	log.Println(wrap(a, true)...)
}






func Print(a ...interface{}) (n int, errno error) {
	return fmt.Print(wrap(a, true)...)
}





func Printf(format string, a ...interface{}) (n int, errno error) {
	return fmt.Printf(format, wrap(a, false)...)
}






func Println(a ...interface{}) (n int, errno error) {
	return fmt.Println(wrap(a, true)...)
}






func Sprint(a ...interface{}) string {
	return fmt.Sprint(wrap(a, true)...)
}





func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, wrap(a, false)...)
}

func wrap(a []interface{}, force bool) []interface{} {
	w := make([]interface{}, len(a))
	for i, x := range a {
		w[i] = formatter{v: reflect.ValueOf(x), force: force}
	}
	return w
}
