



package logging





type multiLogger struct {
	backends []LeveledBackend
}


func MultiLogger(backends ...Backend) LeveledBackend {
	var leveledBackends []LeveledBackend
	for _, backend := range backends {
		leveledBackends = append(leveledBackends, AddModuleLevel(backend))
	}
	return &multiLogger{leveledBackends}
}


func (b *multiLogger) Log(level Level, calldepth int, rec *Record) (err error) {
	for _, backend := range b.backends {
		if backend.IsEnabledFor(level, rec.Module) {
			
			
			r2 := *rec
			if e := backend.Log(level, calldepth+1, &r2); e != nil {
				err = e
			}
		}
	}
	return
}


func (b *multiLogger) GetLevel(module string) Level {
	var level Level
	for _, backend := range b.backends {
		if backendLevel := backend.GetLevel(module); backendLevel > level {
			level = backendLevel
		}
	}
	return level
}


func (b *multiLogger) SetLevel(level Level, module string) {
	for _, backend := range b.backends {
		backend.SetLevel(level, module)
	}
}


func (b *multiLogger) IsEnabledFor(level Level, module string) bool {
	for _, backend := range b.backends {
		if backend.IsEnabledFor(level, module) {
			return true
		}
	}
	return false
}
