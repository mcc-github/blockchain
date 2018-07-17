



package bidi

import (
	"container/list"
	"fmt"
	"sort"
)





























type bracketType byte

const (
	bpNone bracketType = iota
	bpOpen
	bpClose
)



type bracketPair struct {
	opener int
	closer int
}

func (b *bracketPair) String() string {
	return fmt.Sprintf("(%v, %v)", b.opener, b.closer)
}


type bracketPairs []bracketPair

func (b bracketPairs) Len() int           { return len(b) }
func (b bracketPairs) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b bracketPairs) Less(i, j int) bool { return b[i].opener < b[j].opener }









func resolvePairedBrackets(s *isolatingRunSequence) {
	p := bracketPairer{
		sos:              s.sos,
		openers:          list.New(),
		codesIsolatedRun: s.types,
		indexes:          s.indexes,
	}
	dirEmbed := L
	if s.level&1 != 0 {
		dirEmbed = R
	}
	p.locateBrackets(s.p.pairTypes, s.p.pairValues)
	p.resolveBrackets(dirEmbed, s.p.initialTypes)
}

type bracketPairer struct {
	sos Class 

	
	
	
	
	
	
	
	
	
	
	

	
	
	
	
	
	

	openers *list.List 

	
	pairPositions bracketPairs

	codesIsolatedRun []Class 
	indexes          []int   

}



func (p *bracketPairer) matchOpener(pairValues []rune, opener, closer int) bool {
	return pairValues[p.indexes[opener]] == pairValues[p.indexes[closer]]
}

const maxPairingDepth = 63






func (p *bracketPairer) locateBrackets(pairTypes []bracketType, pairValues []rune) {
	
	
	for i, index := range p.indexes {

		
		if pairTypes[index] == bpNone || p.codesIsolatedRun[i] != ON {
			
			continue
		}
		switch pairTypes[index] {
		case bpOpen:
			
			if p.openers.Len() == maxPairingDepth {
				p.openers.Init()
				return
			}
			
			p.openers.PushFront(i)

		case bpClose:
			
			count := 0
			for elem := p.openers.Front(); elem != nil; elem = elem.Next() {
				count++
				opener := elem.Value.(int)
				if p.matchOpener(pairValues, opener, i) {
					
					p.pairPositions = append(p.pairPositions, bracketPair{opener, i})
					
					for ; count > 0; count-- {
						p.openers.Remove(p.openers.Front())
					}
					break
				}
			}
			sort.Sort(p.pairPositions)
			
			
		}
	}
}
























































func (p *bracketPairer) getStrongTypeN0(index int) Class {
	switch p.codesIsolatedRun[index] {
	
	case EN, AN, AL, R:
		return R
	case L:
		return L
	default:
		return ON
	}
}








func (p *bracketPairer) classifyPairContent(loc bracketPair, dirEmbed Class) Class {
	dirOpposite := ON
	for i := loc.opener + 1; i < loc.closer; i++ {
		dir := p.getStrongTypeN0(i)
		if dir == ON {
			continue
		}
		if dir == dirEmbed {
			return dir 
		}
		dirOpposite = dir
	}
	
	return dirOpposite
}



func (p *bracketPairer) classBeforePair(loc bracketPair) Class {
	for i := loc.opener - 1; i >= 0; i-- {
		if dir := p.getStrongTypeN0(i); dir != ON {
			return dir
		}
	}
	
	return p.sos
}


func (p *bracketPairer) assignBracketType(loc bracketPair, dirEmbed Class, initialTypes []Class) {
	
	dirPair := p.classifyPairContent(loc, dirEmbed)

	

	
	
	if dirPair == ON {
		return 
	}

	if dirPair != dirEmbed {
		
		dirPair = p.classBeforePair(loc)
		if dirPair == dirEmbed || dirPair == ON {
			
			dirPair = dirEmbed
		}
	}
	
	
	

	
	p.setBracketsToType(loc, dirPair, initialTypes)
}

func (p *bracketPairer) setBracketsToType(loc bracketPair, dirPair Class, initialTypes []Class) {
	p.codesIsolatedRun[loc.opener] = dirPair
	p.codesIsolatedRun[loc.closer] = dirPair

	for i := loc.opener + 1; i < loc.closer; i++ {
		index := p.indexes[i]
		if initialTypes[index] != NSM {
			break
		}
		p.codesIsolatedRun[i] = dirPair
	}

	for i := loc.closer + 1; i < len(p.indexes); i++ {
		index := p.indexes[i]
		if initialTypes[index] != NSM {
			break
		}
		p.codesIsolatedRun[i] = dirPair
	}
}


func (p *bracketPairer) resolveBrackets(dirEmbed Class, initialTypes []Class) {
	for _, loc := range p.pairPositions {
		p.assignBracketType(loc, dirEmbed, initialTypes)
	}
}
