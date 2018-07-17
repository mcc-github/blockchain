







package htmlindex














import (
	"errors"
	"strings"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/language"
)

var (
	errInvalidName = errors.New("htmlindex: invalid encoding name")
	errUnknown     = errors.New("htmlindex: unknown Encoding")
	errUnsupported = errors.New("htmlindex: this encoding is not supported")
)

var (
	matcherOnce sync.Once
	matcher     language.Matcher
)



func LanguageDefault(tag language.Tag) string {
	matcherOnce.Do(func() {
		tags := []language.Tag{}
		for _, t := range strings.Split(locales, " ") {
			tags = append(tags, language.MustParse(t))
		}
		matcher = language.NewMatcher(tags, language.PreferSameScript(true))
	})
	_, i, _ := matcher.Match(tag)
	return canonical[localeMap[i]] 
}




func Get(name string) (encoding.Encoding, error) {
	x, ok := nameMap[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return nil, errInvalidName
	}
	return encodings[x], nil
}



func Name(e encoding.Encoding) (string, error) {
	id, ok := e.(identifier.Interface)
	if !ok {
		return "", errUnknown
	}
	mib, _ := id.ID()
	if mib == 0 {
		return "", errUnknown
	}
	v, ok := mibMap[mib]
	if !ok {
		return "", errUnsupported
	}
	return canonical[v], nil
}
