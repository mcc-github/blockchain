

package grpclog



type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}




func SetLogger(l Logger) {
	logger = &loggerWrapper{Logger: l}
}


type loggerWrapper struct {
	Logger
}

func (g *loggerWrapper) Info(args ...interface{}) {
	g.Logger.Print(args...)
}

func (g *loggerWrapper) Infoln(args ...interface{}) {
	g.Logger.Println(args...)
}

func (g *loggerWrapper) Infof(format string, args ...interface{}) {
	g.Logger.Printf(format, args...)
}

func (g *loggerWrapper) Warning(args ...interface{}) {
	g.Logger.Print(args...)
}

func (g *loggerWrapper) Warningln(args ...interface{}) {
	g.Logger.Println(args...)
}

func (g *loggerWrapper) Warningf(format string, args ...interface{}) {
	g.Logger.Printf(format, args...)
}

func (g *loggerWrapper) Error(args ...interface{}) {
	g.Logger.Print(args...)
}

func (g *loggerWrapper) Errorln(args ...interface{}) {
	g.Logger.Println(args...)
}

func (g *loggerWrapper) Errorf(format string, args ...interface{}) {
	g.Logger.Printf(format, args...)
}

func (g *loggerWrapper) V(l int) bool {
	
	return true
}
