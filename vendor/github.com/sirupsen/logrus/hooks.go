package logrus






type Hook interface {
	Levels() []Level
	Fire(*Entry) error
}


type LevelHooks map[Level][]Hook



func (hooks LevelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}



func (hooks LevelHooks) Fire(level Level, entry *Entry) error {
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}
