



package bidi

import "log"













































type level int8

const implicitLevel level = -1


func (c Class) in(set ...Class) bool {
	for _, s := range set {
		if c == s {
			return true
		}
	}
	return false
}


type paragraph struct {
	initialTypes []Class

	
	pairTypes  []bracketType 
	pairValues []rune        

	embeddingLevel level 

	
	resultTypes  []Class
	resultLevels []level

	
	
	
	
	matchingPDI []int

	
	
	
	matchingIsolateInitiator []int
}










func newParagraph(types []Class, pairTypes []bracketType, pairValues []rune, levels level) *paragraph {
	validateTypes(types)
	validatePbTypes(pairTypes)
	validatePbValues(pairValues, pairTypes)
	validateParagraphEmbeddingLevel(levels)

	p := &paragraph{
		initialTypes:   append([]Class(nil), types...),
		embeddingLevel: levels,

		pairTypes:  pairTypes,
		pairValues: pairValues,

		resultTypes: append([]Class(nil), types...),
	}
	p.run()
	return p
}

func (p *paragraph) Len() int { return len(p.initialTypes) }



func (p *paragraph) run() {
	p.determineMatchingIsolates()

	
	
	
	
	if p.embeddingLevel == implicitLevel {
		p.embeddingLevel = p.determineParagraphEmbeddingLevel(0, p.Len())
	}

	
	p.resultLevels = make([]level, p.Len())
	setLevels(p.resultLevels, p.embeddingLevel)

	
	
	p.determineExplicitEmbeddingLevels()

	
	
	
	
	

	
	
	for _, seq := range p.determineIsolatingRunSequences() {
		
		
		seq.resolveWeakTypes()

		
		
		resolvePairedBrackets(seq)

		
		
		seq.resolveNeutralTypes()

		
		
		seq.resolveImplicitLevels()

		
		seq.applyLevelsAndTypes()
	}

	
	
	
	p.assignLevelsToCharactersRemovedByX9()
}
















func (p *paragraph) determineMatchingIsolates() {
	p.matchingPDI = make([]int, p.Len())
	p.matchingIsolateInitiator = make([]int, p.Len())

	for i := range p.matchingIsolateInitiator {
		p.matchingIsolateInitiator[i] = -1
	}

	for i := range p.matchingPDI {
		p.matchingPDI[i] = -1

		if t := p.resultTypes[i]; t.in(LRI, RLI, FSI) {
			depthCounter := 1
			for j := i + 1; j < p.Len(); j++ {
				if u := p.resultTypes[j]; u.in(LRI, RLI, FSI) {
					depthCounter++
				} else if u == PDI {
					if depthCounter--; depthCounter == 0 {
						p.matchingPDI[i] = j
						p.matchingIsolateInitiator[j] = i
						break
					}
				}
			}
			if p.matchingPDI[i] == -1 {
				p.matchingPDI[i] = p.Len()
			}
		}
	}
}






func (p *paragraph) determineParagraphEmbeddingLevel(start, end int) level {
	var strongType Class = unknownClass

	
	for i := start; i < end; i++ {
		if t := p.resultTypes[i]; t.in(L, AL, R) {
			strongType = t
			break
		} else if t.in(FSI, LRI, RLI) {
			i = p.matchingPDI[i] 
			if i > end {
				log.Panic("assert (i <= end)")
			}
		}
	}
	
	switch strongType {
	case unknownClass: 
		
		return 0
	case L:
		return 0
	default: 
		return 1
	}
}

const maxDepth = 125



type directionalStatusStack struct {
	stackCounter        int
	embeddingLevelStack [maxDepth + 1]level
	overrideStatusStack [maxDepth + 1]Class
	isolateStatusStack  [maxDepth + 1]bool
}

func (s *directionalStatusStack) empty()     { s.stackCounter = 0 }
func (s *directionalStatusStack) pop()       { s.stackCounter-- }
func (s *directionalStatusStack) depth() int { return s.stackCounter }

func (s *directionalStatusStack) push(level level, overrideStatus Class, isolateStatus bool) {
	s.embeddingLevelStack[s.stackCounter] = level
	s.overrideStatusStack[s.stackCounter] = overrideStatus
	s.isolateStatusStack[s.stackCounter] = isolateStatus
	s.stackCounter++
}

func (s *directionalStatusStack) lastEmbeddingLevel() level {
	return s.embeddingLevelStack[s.stackCounter-1]
}

func (s *directionalStatusStack) lastDirectionalOverrideStatus() Class {
	return s.overrideStatusStack[s.stackCounter-1]
}

func (s *directionalStatusStack) lastDirectionalIsolateStatus() bool {
	return s.isolateStatusStack[s.stackCounter-1]
}


func (p *paragraph) determineExplicitEmbeddingLevels() {
	var stack directionalStatusStack
	var overflowIsolateCount, overflowEmbeddingCount, validIsolateCount int

	
	stack.push(p.embeddingLevel, ON, false)

	for i, t := range p.resultTypes {
		
		switch t {
		case RLE, LRE, RLO, LRO, RLI, LRI, FSI:
			isIsolate := t.in(RLI, LRI, FSI)
			isRTL := t.in(RLE, RLO, RLI)

			
			if t == FSI {
				isRTL = (p.determineParagraphEmbeddingLevel(i+1, p.matchingPDI[i]) == 1)
			}
			if isIsolate {
				p.resultLevels[i] = stack.lastEmbeddingLevel()
				if stack.lastDirectionalOverrideStatus() != ON {
					p.resultTypes[i] = stack.lastDirectionalOverrideStatus()
				}
			}

			var newLevel level
			if isRTL {
				
				newLevel = (stack.lastEmbeddingLevel() + 1) | 1
			} else {
				
				newLevel = (stack.lastEmbeddingLevel() + 2) &^ 1
			}

			if newLevel <= maxDepth && overflowIsolateCount == 0 && overflowEmbeddingCount == 0 {
				if isIsolate {
					validIsolateCount++
				}
				
				
				
				
				switch t {
				case LRO:
					stack.push(newLevel, L, isIsolate)
				case RLO:
					stack.push(newLevel, R, isIsolate)
				default:
					stack.push(newLevel, ON, isIsolate)
				}
				
				if !isIsolate {
					p.resultLevels[i] = newLevel
				}
			} else {
				
				
				if isIsolate {
					overflowIsolateCount++
				} else { 
					if overflowIsolateCount == 0 {
						overflowEmbeddingCount++
					}
				}
			}

		
		case PDI:
			if overflowIsolateCount > 0 {
				overflowIsolateCount--
			} else if validIsolateCount == 0 {
				
			} else {
				overflowEmbeddingCount = 0
				for !stack.lastDirectionalIsolateStatus() {
					stack.pop()
				}
				stack.pop()
				validIsolateCount--
			}
			p.resultLevels[i] = stack.lastEmbeddingLevel()

		
		case PDF:
			
			p.resultLevels[i] = stack.lastEmbeddingLevel()

			if overflowIsolateCount > 0 {
				
			} else if overflowEmbeddingCount > 0 {
				overflowEmbeddingCount--
			} else if !stack.lastDirectionalIsolateStatus() && stack.depth() >= 2 {
				stack.pop()
			}

		case B: 
			

			
			
			stack.empty()
			overflowIsolateCount = 0
			overflowEmbeddingCount = 0
			validIsolateCount = 0
			p.resultLevels[i] = p.embeddingLevel

		default:
			p.resultLevels[i] = stack.lastEmbeddingLevel()
			if stack.lastDirectionalOverrideStatus() != ON {
				p.resultTypes[i] = stack.lastDirectionalOverrideStatus()
			}
		}
	}
}

type isolatingRunSequence struct {
	p *paragraph

	indexes []int 

	types          []Class 
	resolvedLevels []level 
	level          level
	sos, eos       Class
}

func (i *isolatingRunSequence) Len() int { return len(i.indexes) }

func maxLevel(a, b level) level {
	if a > b {
		return a
	}
	return b
}



func (p *paragraph) isolatingRunSequence(indexes []int) *isolatingRunSequence {
	length := len(indexes)
	types := make([]Class, length)
	for i, x := range indexes {
		types[i] = p.resultTypes[x]
	}

	
	prevChar := indexes[0] - 1
	for prevChar >= 0 && isRemovedByX9(p.initialTypes[prevChar]) {
		prevChar--
	}
	prevLevel := p.embeddingLevel
	if prevChar >= 0 {
		prevLevel = p.resultLevels[prevChar]
	}

	var succLevel level
	lastType := types[length-1]
	if lastType.in(LRI, RLI, FSI) {
		succLevel = p.embeddingLevel
	} else {
		
		limit := indexes[length-1] + 1
		for ; limit < p.Len() && isRemovedByX9(p.initialTypes[limit]); limit++ {

		}
		succLevel = p.embeddingLevel
		if limit < p.Len() {
			succLevel = p.resultLevels[limit]
		}
	}
	level := p.resultLevels[indexes[0]]
	return &isolatingRunSequence{
		p:       p,
		indexes: indexes,
		types:   types,
		level:   level,
		sos:     typeForLevel(maxLevel(prevLevel, level)),
		eos:     typeForLevel(maxLevel(succLevel, level)),
	}
}





func (s *isolatingRunSequence) resolveWeakTypes() {

	
	s.assertOnly(L, R, AL, EN, ES, ET, AN, CS, B, S, WS, ON, NSM, LRI, RLI, FSI, PDI)

	
	
	preceedingCharacterType := s.sos
	for i, t := range s.types {
		if t == NSM {
			s.types[i] = preceedingCharacterType
		} else {
			if t.in(LRI, RLI, FSI, PDI) {
				preceedingCharacterType = ON
			}
			preceedingCharacterType = t
		}
	}

	
	
	for i, t := range s.types {
		if t == EN {
			for j := i - 1; j >= 0; j-- {
				if t := s.types[j]; t.in(L, R, AL) {
					if t == AL {
						s.types[i] = AN
					}
					break
				}
			}
		}
	}

	
	for i, t := range s.types {
		if t == AL {
			s.types[i] = R
		}
	}

	
	
	
	
	
	
	
	
	
	
	
	

	for i := 1; i < s.Len()-1; i++ {
		t := s.types[i]
		if t == ES || t == CS {
			prevSepType := s.types[i-1]
			succSepType := s.types[i+1]
			if prevSepType == EN && succSepType == EN {
				s.types[i] = EN
			} else if s.types[i] == CS && prevSepType == AN && succSepType == AN {
				s.types[i] = AN
			}
		}
	}

	
	for i, t := range s.types {
		if t == ET {
			
			runStart := i
			runEnd := s.findRunLimit(runStart, ET)

			
			t := s.sos
			if runStart > 0 {
				t = s.types[runStart-1]
			}
			if t != EN {
				t = s.eos
				if runEnd < len(s.types) {
					t = s.types[runEnd]
				}
			}
			if t == EN {
				setTypes(s.types[runStart:runEnd], EN)
			}
			
			i = runEnd
		}
	}

	
	for i, t := range s.types {
		if t.in(ES, ET, CS) {
			s.types[i] = ON
		}
	}

	
	for i, t := range s.types {
		if t == EN {
			
			prevStrongType := s.sos
			for j := i - 1; j >= 0; j-- {
				t = s.types[j]
				if t == L || t == R { 
					prevStrongType = t
					break
				}
			}
			if prevStrongType == L {
				s.types[i] = L
			}
		}
	}
}


func (s *isolatingRunSequence) resolveNeutralTypes() {

	
	s.assertOnly(L, R, EN, AN, B, S, WS, ON, RLI, LRI, FSI, PDI)

	for i, t := range s.types {
		switch t {
		case WS, ON, B, S, RLI, LRI, FSI, PDI:
			
			runStart := i
			runEnd := s.findRunLimit(runStart, B, S, WS, ON, RLI, LRI, FSI, PDI)

			
			var leadType, trailType Class

			
			
			if runStart == 0 {
				leadType = s.sos
			} else {
				leadType = s.types[runStart-1]
				if leadType.in(AN, EN) {
					leadType = R
				}
			}
			if runEnd == len(s.types) {
				trailType = s.eos
			} else {
				trailType = s.types[runEnd]
				if trailType.in(AN, EN) {
					trailType = R
				}
			}

			var resolvedType Class
			if leadType == trailType {
				
				resolvedType = leadType
			} else {
				
				
				
				resolvedType = typeForLevel(s.level)
			}

			setTypes(s.types[runStart:runEnd], resolvedType)

			
			i = runEnd
		}
	}
}

func setLevels(levels []level, newLevel level) {
	for i := range levels {
		levels[i] = newLevel
	}
}

func setTypes(types []Class, newType Class) {
	for i := range types {
		types[i] = newType
	}
}


func (s *isolatingRunSequence) resolveImplicitLevels() {

	
	s.assertOnly(L, R, EN, AN)

	s.resolvedLevels = make([]level, len(s.types))
	setLevels(s.resolvedLevels, s.level)

	if (s.level & 1) == 0 { 
		for i, t := range s.types {
			
			if t == L {
				
			} else if t == R {
				s.resolvedLevels[i] += 1
			} else { 
				s.resolvedLevels[i] += 2
			}
		}
	} else { 
		for i, t := range s.types {
			
			if t == R {
				
			} else { 
				s.resolvedLevels[i] += 1
			}
		}
	}
}



func (s *isolatingRunSequence) applyLevelsAndTypes() {
	for i, x := range s.indexes {
		s.p.resultTypes[x] = s.types[i]
		s.p.resultLevels[x] = s.resolvedLevels[i]
	}
}




func (s *isolatingRunSequence) findRunLimit(index int, validSet ...Class) int {
loop:
	for ; index < len(s.types); index++ {
		t := s.types[index]
		for _, valid := range validSet {
			if t == valid {
				continue loop
			}
		}
		return index 
	}
	return len(s.types)
}



func (s *isolatingRunSequence) assertOnly(codes ...Class) {
loop:
	for i, t := range s.types {
		for _, c := range codes {
			if t == c {
				continue loop
			}
		}
		log.Panicf("invalid bidi code %v present in assertOnly at position %d", t, s.indexes[i])
	}
}







func (p *paragraph) determineLevelRuns() [][]int {
	run := []int{}
	allRuns := [][]int{}
	currentLevel := implicitLevel

	for i := range p.initialTypes {
		if !isRemovedByX9(p.initialTypes[i]) {
			if p.resultLevels[i] != currentLevel {
				
				if currentLevel >= 0 { 
					allRuns = append(allRuns, run)
					run = nil
				}
				
				currentLevel = p.resultLevels[i]
			}
			run = append(run, i)
		}
	}
	
	if len(run) > 0 {
		allRuns = append(allRuns, run)
	}
	return allRuns
}


func (p *paragraph) determineIsolatingRunSequences() []*isolatingRunSequence {
	levelRuns := p.determineLevelRuns()

	
	runForCharacter := make([]int, p.Len())
	for i, run := range levelRuns {
		for _, index := range run {
			runForCharacter[index] = i
		}
	}

	sequences := []*isolatingRunSequence{}

	var currentRunSequence []int

	for _, run := range levelRuns {
		first := run[0]
		if p.initialTypes[first] != PDI || p.matchingIsolateInitiator[first] == -1 {
			currentRunSequence = nil
			
			for {
				
				currentRunSequence = append(currentRunSequence, run...)

				last := currentRunSequence[len(currentRunSequence)-1]
				lastT := p.initialTypes[last]
				if lastT.in(LRI, RLI, FSI) && p.matchingPDI[last] != p.Len() {
					run = levelRuns[runForCharacter[p.matchingPDI[last]]]
				} else {
					break
				}
			}
			sequences = append(sequences, p.isolatingRunSequence(currentRunSequence))
		}
	}
	return sequences
}





func (p *paragraph) assignLevelsToCharactersRemovedByX9() {
	for i, t := range p.initialTypes {
		if t.in(LRE, RLE, LRO, RLO, PDF, BN) {
			p.resultTypes[i] = t
			p.resultLevels[i] = -1
		}
	}
	
	
	

	if p.resultLevels[0] == -1 {
		p.resultLevels[0] = p.embeddingLevel
	}
	for i := 1; i < len(p.initialTypes); i++ {
		if p.resultLevels[i] == -1 {
			p.resultLevels[i] = p.resultLevels[i-1]
		}
	}
	
	
}











func (p *paragraph) getLevels(linebreaks []int) []level {
	
	
	
	
	
	
	
	
	
	

	validateLineBreaks(linebreaks, p.Len())

	result := append([]level(nil), p.resultLevels...)

	
	
	
	for i, t := range p.initialTypes {
		if t.in(B, S) {
			
			result[i] = p.embeddingLevel

			
			for j := i - 1; j >= 0; j-- {
				if isWhitespace(p.initialTypes[j]) { 
					result[j] = p.embeddingLevel
				} else {
					break
				}
			}
		}
	}

	
	start := 0
	for _, limit := range linebreaks {
		for j := limit - 1; j >= start; j-- {
			if isWhitespace(p.initialTypes[j]) { 
				result[j] = p.embeddingLevel
			} else {
				break
			}
		}
		start = limit
	}

	return result
}
















func (p *paragraph) getReordering(linebreaks []int) []int {
	validateLineBreaks(linebreaks, p.Len())

	return computeMultilineReordering(p.getLevels(linebreaks), linebreaks)
}



func computeMultilineReordering(levels []level, linebreaks []int) []int {
	result := make([]int, len(levels))

	start := 0
	for _, limit := range linebreaks {
		tempLevels := make([]level, limit-start)
		copy(tempLevels, levels[start:])

		for j, order := range computeReordering(tempLevels) {
			result[start+j] = order + start
		}
		start = limit
	}
	return result
}




func computeReordering(levels []level) []int {
	result := make([]int, len(levels))
	
	for i := range result {
		result[i] = i
	}

	
	
	
	highestLevel := level(0)
	lowestOddLevel := level(maxDepth + 2)
	for _, level := range levels {
		if level > highestLevel {
			highestLevel = level
		}
		if level&1 != 0 && level < lowestOddLevel {
			lowestOddLevel = level
		}
	}

	for level := highestLevel; level >= lowestOddLevel; level-- {
		for i := 0; i < len(levels); i++ {
			if levels[i] >= level {
				
				start := i
				limit := i + 1
				for limit < len(levels) && levels[limit] >= level {
					limit++
				}

				for j, k := start, limit-1; j < k; j, k = j+1, k-1 {
					result[j], result[k] = result[k], result[j]
				}
				
				i = limit
			}
		}
	}

	return result
}



func isWhitespace(c Class) bool {
	switch c {
	case LRE, RLE, LRO, RLO, PDF, LRI, RLI, FSI, PDI, BN, WS:
		return true
	}
	return false
}


func isRemovedByX9(c Class) bool {
	switch c {
	case LRE, RLE, LRO, RLO, PDF, BN:
		return true
	}
	return false
}


func typeForLevel(level level) Class {
	if (level & 0x1) == 0 {
		return L
	}
	return R
}



func validateTypes(types []Class) {
	if len(types) == 0 {
		log.Panic("types is null")
	}
	for i, t := range types[:len(types)-1] {
		if t == B {
			log.Panicf("B type before end of paragraph at index: %d", i)
		}
	}
}

func validateParagraphEmbeddingLevel(embeddingLevel level) {
	if embeddingLevel != implicitLevel &&
		embeddingLevel != 0 &&
		embeddingLevel != 1 {
		log.Panicf("illegal paragraph embedding level: %d", embeddingLevel)
	}
}

func validateLineBreaks(linebreaks []int, textLength int) {
	prev := 0
	for i, next := range linebreaks {
		if next <= prev {
			log.Panicf("bad linebreak: %d at index: %d", next, i)
		}
		prev = next
	}
	if prev != textLength {
		log.Panicf("last linebreak was %d, want %d", prev, textLength)
	}
}

func validatePbTypes(pairTypes []bracketType) {
	if len(pairTypes) == 0 {
		log.Panic("pairTypes is null")
	}
	for i, pt := range pairTypes {
		switch pt {
		case bpNone, bpOpen, bpClose:
		default:
			log.Panicf("illegal pairType value at %d: %v", i, pairTypes[i])
		}
	}
}

func validatePbValues(pairValues []rune, pairTypes []bracketType) {
	if pairValues == nil {
		log.Panic("pairValues is null")
	}
	if len(pairTypes) != len(pairValues) {
		log.Panic("pairTypes is different length from pairValues")
	}
}
