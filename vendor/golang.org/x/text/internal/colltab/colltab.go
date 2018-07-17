





package colltab 

import (
	"sort"

	"golang.org/x/text/language"
)











func MatchLang(t language.Tag, tags []language.Tag) int {
	
	t, _ = language.All.Canonicalize(t)

	base, conf := t.Base()
	
	if conf < language.High {
		
		
		return 0
	}

	
	if _, s, r := t.Raw(); (r != language.Region{}) {
		p, _ := language.Raw.Compose(base, s, r)
		
		p = p.Parent()
		
		t, _ = language.Raw.Compose(p, r, t.Extensions())
	} else {
		
		t, _ = language.Raw.Compose(base, s, t.Extensions())
	}

	
	start := 1 + sort.Search(len(tags)-1, func(i int) bool {
		b, _, _ := tags[i+1].Raw()
		return base.String() <= b.String()
	})
	if start < len(tags) {
		if b, _, _ := tags[start].Raw(); b != base {
			return 0
		}
	}

	
	
	
	
	tdef, _ := language.Raw.Compose(t.Raw())
	tdef, _ = tdef.SetTypeForKey("va", t.TypeForKey("va"))

	
	try := []language.Tag{tdef}
	if co := t.TypeForKey("co"); co != "" {
		tco, _ := tdef.SetTypeForKey("co", co)
		try = []language.Tag{tco, tdef}
	}

	for _, tx := range try {
		for ; tx != language.Und; tx = parent(tx) {
			for i, t := range tags[start:] {
				if b, _, _ := t.Raw(); b != base {
					break
				}
				if tx == t {
					return start + i
				}
			}
		}
	}
	return 0
}



func parent(t language.Tag) language.Tag {
	if t.TypeForKey("va") != "" {
		t, _ = t.SetTypeForKey("va", "")
		return t
	}
	result := language.Und
	if b, s, r := t.Raw(); (r != language.Region{}) {
		result, _ = language.Raw.Compose(b, s, t.Extensions())
	} else if (s != language.Script{}) {
		result, _ = language.Raw.Compose(b, t.Extensions())
	} else if (b != language.Base{}) {
		result, _ = language.Raw.Compose(t.Extensions())
	}
	return result
}
