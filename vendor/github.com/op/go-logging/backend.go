



package logging


var defaultBackend LeveledBackend



type Backend interface {
	Log(Level, int, *Record) error
}



func SetBackend(backends ...Backend) LeveledBackend {
	var backend Backend
	if len(backends) == 1 {
		backend = backends[0]
	} else {
		backend = MultiLogger(backends...)
	}

	defaultBackend = AddModuleLevel(backend)
	return defaultBackend
}



func SetLevel(level Level, module string) {
	defaultBackend.SetLevel(level, module)
}


func GetLevel(module string) Level {
	return defaultBackend.GetLevel(module)
}
