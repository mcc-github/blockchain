package logrus























import (
	"fmt"
	"os"
)

var handlers = []func(){}

func runHandler(handler func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Logrus exit handler error:", err)
		}
	}()

	handler()
}

func runHandlers() {
	for _, handler := range handlers {
		runHandler(handler)
	}
}


func Exit(code int) {
	runHandlers()
	os.Exit(code)
}









func RegisterExitHandler(handler func()) {
	handlers = append(handlers, handler)
}
