



















package tally

import (
	"bytes"
)

var (
	
	
	DefaultReplacementCharacter = '_'

	
	AlphanumericRange = []SanitizeRange{
		{rune('a'), rune('z')},
		{rune('A'), rune('Z')},
		{rune('0'), rune('9')}}

	
	
	UnderscoreDashCharacters = []rune{
		'-',
		'_'}

	
	
	UnderscoreDashDotCharacters = []rune{
		'.',
		'-',
		'_'}
)


type SanitizeFn func(string) string


type SanitizeRange [2]rune


type ValidCharacters struct {
	Ranges     []SanitizeRange
	Characters []rune
}


type SanitizeOptions struct {
	NameCharacters       ValidCharacters
	KeyCharacters        ValidCharacters
	ValueCharacters      ValidCharacters
	ReplacementCharacter rune
}


type Sanitizer interface {
	
	Name(n string) string

	
	Key(k string) string

	
	Value(v string) string
}


func NewSanitizer(opts SanitizeOptions) Sanitizer {
	return sanitizer{
		nameFn:  opts.NameCharacters.sanitizeFn(opts.ReplacementCharacter),
		keyFn:   opts.KeyCharacters.sanitizeFn(opts.ReplacementCharacter),
		valueFn: opts.ValueCharacters.sanitizeFn(opts.ReplacementCharacter),
	}
}


func NoOpSanitizeFn(v string) string { return v }


func NewNoOpSanitizer() Sanitizer {
	return sanitizer{
		nameFn:  NoOpSanitizeFn,
		keyFn:   NoOpSanitizeFn,
		valueFn: NoOpSanitizeFn,
	}
}

type sanitizer struct {
	nameFn  SanitizeFn
	keyFn   SanitizeFn
	valueFn SanitizeFn
}

func (s sanitizer) Name(n string) string {
	return s.nameFn(n)
}

func (s sanitizer) Key(k string) string {
	return s.keyFn(k)
}

func (s sanitizer) Value(v string) string {
	return s.valueFn(v)
}

func (c *ValidCharacters) sanitizeFn(repChar rune) SanitizeFn {
	return func(value string) string {
		var buf *bytes.Buffer
		for idx, ch := range value {
			
			validCurr := false
			for i := 0; !validCurr && i < len(c.Ranges); i++ {
				if ch >= c.Ranges[i][0] && ch <= c.Ranges[i][1] {
					validCurr = true
					break
				}
			}
			for i := 0; !validCurr && i < len(c.Characters); i++ {
				if c.Characters[i] == ch {
					validCurr = true
					break
				}
			}

			
			if validCurr {
				if buf == nil {
					continue 
				}
				buf.WriteRune(ch) 
				continue
			}

			
			
			if buf == nil {
				buf = bytes.NewBuffer(make([]byte, 0, len(value)))
				if idx > 0 {
					buf.WriteString(value[:idx])
				}
			}

			
			buf.WriteRune(repChar)
		}

		
		if buf == nil {
			return value
		}

		
		return buf.String()
	}
}
