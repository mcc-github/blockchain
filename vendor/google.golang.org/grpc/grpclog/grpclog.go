









package grpclog 

import "os"

var logger = newLoggerV2()


func V(l int) bool {
	return logger.V(l)
}


func Info(args ...interface{}) {
	logger.Info(args...)
}


func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}


func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}


func Warning(args ...interface{}) {
	logger.Warning(args...)
}


func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}


func Warningln(args ...interface{}) {
	logger.Warningln(args...)
}


func Error(args ...interface{}) {
	logger.Error(args...)
}


func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}


func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}



func Fatal(args ...interface{}) {
	logger.Fatal(args...)
	
	os.Exit(1)
}



func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
	
	os.Exit(1)
}



func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
	
	os.Exit(1)
}




func Print(args ...interface{}) {
	logger.Info(args...)
}




func Printf(format string, args ...interface{}) {
	logger.Infof(format, args...)
}




func Println(args ...interface{}) {
	logger.Infoln(args...)
}
