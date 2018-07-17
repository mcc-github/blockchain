
package sarama

import (
	"io/ioutil"
	"log"
)




var Logger StdLogger = log.New(ioutil.Discard, "[Sarama] ", log.LstdFlags)


type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}



var PanicHandler func(interface{})





var MaxRequestSize int32 = 100 * 1024 * 1024






var MaxResponseSize int32 = 100 * 1024 * 1024
