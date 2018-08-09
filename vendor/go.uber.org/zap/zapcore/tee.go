



















package zapcore

import "go.uber.org/multierr"

type multiCore []Core






func NewTee(cores ...Core) Core {
	switch len(cores) {
	case 0:
		return NewNopCore()
	case 1:
		return cores[0]
	default:
		return multiCore(cores)
	}
}

func (mc multiCore) With(fields []Field) Core {
	clone := make(multiCore, len(mc))
	for i := range mc {
		clone[i] = mc[i].With(fields)
	}
	return clone
}

func (mc multiCore) Enabled(lvl Level) bool {
	for i := range mc {
		if mc[i].Enabled(lvl) {
			return true
		}
	}
	return false
}

func (mc multiCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	for i := range mc {
		ce = mc[i].Check(ent, ce)
	}
	return ce
}

func (mc multiCore) Write(ent Entry, fields []Field) error {
	var err error
	for i := range mc {
		err = multierr.Append(err, mc[i].Write(ent, fields))
	}
	return err
}

func (mc multiCore) Sync() error {
	var err error
	for i := range mc {
		err = multierr.Append(err, mc[i].Sync())
	}
	return err
}
