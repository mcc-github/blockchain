



package colltab




type Iter struct {
	Weighter Weighter
	Elems    []Elem
	
	
	N int

	bytes []byte
	str   string
	
	
	
	
	pEnd  int 
	pNext int 
}



func (i *Iter) Reset(p int) {
	i.Elems = i.Elems[:0]
	i.N = 0
	i.pEnd = p
	i.pNext = p
}


func (i *Iter) Len() int {
	if i.bytes != nil {
		return len(i.bytes)
	}
	return len(i.str)
}


func (i *Iter) Discard() {
	
	
	i.Elems = i.Elems[:copy(i.Elems, i.Elems[i.N:])]
	i.N = 0
}



func (i *Iter) End() int {
	return i.pEnd
}


func (i *Iter) SetInput(s []byte) {
	i.bytes = s
	i.str = ""
	i.Reset(0)
}


func (i *Iter) SetInputString(s string) {
	i.str = s
	i.bytes = nil
	i.Reset(0)
}

func (i *Iter) done() bool {
	return i.pNext >= len(i.str) && i.pNext >= len(i.bytes)
}

func (i *Iter) appendNext() bool {
	if i.done() {
		return false
	}
	var sz int
	if i.bytes == nil {
		i.Elems, sz = i.Weighter.AppendNextString(i.Elems, i.str[i.pNext:])
	} else {
		i.Elems, sz = i.Weighter.AppendNext(i.Elems, i.bytes[i.pNext:])
	}
	if sz == 0 {
		sz = 1
	}
	i.pNext += sz
	return true
}






func (i *Iter) Next() bool {
	if i.N == len(i.Elems) && !i.appendNext() {
		return false
	}

	
	prevCCC := i.Elems[len(i.Elems)-1].CCC()
	if prevCCC == 0 {
		i.N = len(i.Elems)
		i.pEnd = i.pNext
		return true
	} else if i.Elems[i.N].CCC() == 0 {
		
		
		for i.N++; i.N < len(i.Elems) && i.Elems[i.N].CCC() == 0; i.N++ {
		}
		i.pEnd = i.pNext
		return true
	}

	
	
	for {
		p := len(i.Elems)
		i.pEnd = i.pNext
		if !i.appendNext() {
			break
		}

		if ccc := i.Elems[p].CCC(); ccc == 0 || len(i.Elems)-i.N > maxCombiningCharacters {
			
			
			
			
			
			
			
			i.N = p
			return true
		} else if ccc < prevCCC {
			i.doNorm(p, ccc) 
		} else {
			prevCCC = ccc
		}
	}

	done := len(i.Elems) != i.N
	i.N = len(i.Elems)
	return done
}



func (i *Iter) nextNoNorm() bool {
	
	
	
	if i.done() {
		return false
	}
	i.appendNext()
	i.N = len(i.Elems)
	return true
}

const maxCombiningCharacters = 30






func (i *Iter) doNorm(p int, ccc uint8) {
	n := len(i.Elems)
	k := p
	for p--; p > i.N && ccc < i.Elems[p-1].CCC(); p-- {
	}
	i.Elems = append(i.Elems, i.Elems[p:k]...)
	copy(i.Elems[p:], i.Elems[k:])
	i.Elems = i.Elems[:n]
}
