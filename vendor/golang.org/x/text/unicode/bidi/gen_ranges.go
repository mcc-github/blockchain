





package main

import (
	"unicode"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/ucd"
	"golang.org/x/text/unicode/rangetable"
)



func visitDefaults(fn func(r rune, c Class)) {
	
	visitRunes(fn, AL, []rune{
		0x0600, 0x07BF, 
		0x08A0, 0x08FF, 
		0xFB50, 0xFDCF, 
		0xFDF0, 0xFDFF,
		0xFE70, 0xFEFF,
		0x0001EE00, 0x0001EEFF, 
	})
	visitRunes(fn, R, []rune{
		0x0590, 0x05FF, 
		0x07C0, 0x089F, 
		0xFB1D, 0xFB4F,
		0x00010800, 0x00010FFF, 
		0x0001E800, 0x0001EDFF,
		0x0001EF00, 0x0001EFFF,
	})
	visitRunes(fn, ET, []rune{ 
		0x20A0, 0x20Cf, 
	})
	rangetable.Visit(unicode.Noncharacter_Code_Point, func(r rune) {
		fn(r, BN) 
	})
	ucd.Parse(gen.OpenUCDFile("DerivedCoreProperties.txt"), func(p *ucd.Parser) {
		if p.String(1) == "Default_Ignorable_Code_Point" {
			fn(p.Rune(0), BN) 
		}
	})
}

func visitRunes(fn func(r rune, c Class), c Class, runes []rune) {
	for i := 0; i < len(runes); i += 2 {
		lo, hi := runes[i], runes[i+1]
		for j := lo; j <= hi; j++ {
			fn(j, c)
		}
	}
}
