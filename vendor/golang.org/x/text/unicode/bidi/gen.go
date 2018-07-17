





package main

import (
	"flag"
	"log"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/triegen"
	"golang.org/x/text/internal/ucd"
)

var outputFile = flag.String("out", "tables.go", "output file")

func main() {
	gen.Init()
	gen.Repackage("gen_trieval.go", "trieval.go", "bidi")
	gen.Repackage("gen_ranges.go", "ranges_test.go", "bidi")

	genTables()
}



var bidiClass = map[string]Class{
	"AL":  AL,  
	"AN":  AN,  
	"B":   B,   
	"BN":  BN,  
	"CS":  CS,  
	"EN":  EN,  
	"ES":  ES,  
	"ET":  ET,  
	"L":   L,   
	"NSM": NSM, 
	"ON":  ON,  
	"R":   R,   
	"S":   S,   
	"WS":  WS,  

	"FSI": Control,
	"PDF": Control,
	"PDI": Control,
	"LRE": Control,
	"LRI": Control,
	"LRO": Control,
	"RLE": Control,
	"RLI": Control,
	"RLO": Control,
}

func genTables() {
	if numClass > 0x0F {
		log.Fatalf("Too many Class constants (%#x > 0x0F).", numClass)
	}
	w := gen.NewCodeWriter()
	defer w.WriteVersionedGoFile(*outputFile, "bidi")

	gen.WriteUnicodeVersion(w)

	t := triegen.NewTrie("bidi")

	
	
	orMask := map[rune]uint64{}

	xorMap := map[rune]int{}
	xorMasks := []rune{0} 

	ucd.Parse(gen.OpenUCDFile("BidiBrackets.txt"), func(p *ucd.Parser) {
		r1 := p.Rune(0)
		r2 := p.Rune(1)
		xor := r1 ^ r2
		if _, ok := xorMap[xor]; !ok {
			xorMap[xor] = len(xorMasks)
			xorMasks = append(xorMasks, xor)
		}
		entry := uint64(xorMap[xor]) << xorMaskShift
		switch p.String(2) {
		case "o":
			entry |= openMask
		case "c", "n":
		default:
			log.Fatalf("Unknown bracket class %q.", p.String(2))
		}
		orMask[r1] = entry
	})

	w.WriteComment(`
	xorMasks contains masks to be xor-ed with brackets to get the reverse
	version.`)
	w.WriteVar("xorMasks", xorMasks)

	done := map[rune]bool{}

	insert := func(r rune, c Class) {
		if !done[r] {
			t.Insert(r, orMask[r]|uint64(c))
			done[r] = true
		}
	}

	
	ucd.Parse(gen.OpenUCDFile("extracted/DerivedBidiClass.txt"), func(p *ucd.Parser) {
		r := p.Rune(0)
		class, ok := bidiClass[p.String(1)]
		if !ok {
			log.Fatalf("%U: Unknown BiDi class %q", r, p.String(1))
		}
		insert(r, class)
	})
	visitDefaults(insert)

	
	

	sz, err := t.Gen(w)
	if err != nil {
		log.Fatal(err)
	}
	w.Size += sz
}



var (
	xorMasks []rune
)
