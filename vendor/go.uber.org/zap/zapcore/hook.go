



















package zapcore

import "go.uber.org/multierr"

type hooked struct {
	Core
	funcs []func(Entry) error
}






func RegisterHooks(core Core, hooks ...func(Entry) error) Core {
	funcs := append([]func(Entry) error{}, hooks...)
	return &hooked{
		Core:  core,
		funcs: funcs,
	}
}

func (h *hooked) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	
	
	
	if downstream := h.Core.Check(ent, ce); downstream != nil {
		return downstream.AddCore(ent, h)
	}
	return ce
}

func (h *hooked) With(fields []Field) Core {
	return &hooked{
		Core:  h.Core.With(fields),
		funcs: h.funcs,
	}
}

func (h *hooked) Write(ent Entry, _ []Field) error {
	
	
	var err error
	for i := range h.funcs {
		err = multierr.Append(err, h.funcs[i](ent))
	}
	return err
}
