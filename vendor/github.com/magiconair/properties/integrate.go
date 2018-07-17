



package properties

import "flag"







func (p *Properties) MustFlag(dst *flag.FlagSet) {
	m := make(map[string]*flag.Flag)
	dst.VisitAll(func(f *flag.Flag) {
		m[f.Name] = f
	})
	dst.Visit(func(f *flag.Flag) {
		delete(m, f.Name) 
	})

	for name, f := range m {
		v, ok := p.Get(name)
		if !ok {
			continue
		}

		if err := f.Value.Set(v); err != nil {
			ErrorHandler(err)
		}
	}
}
