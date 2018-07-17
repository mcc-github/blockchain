














package difflib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func calculateRatio(matches, length int) float64 {
	if length > 0 {
		return 2.0 * float64(matches) / float64(length)
	}
	return 1.0
}

type Match struct {
	A    int
	B    int
	Size int
}

type OpCode struct {
	Tag byte
	I1  int
	I2  int
	J1  int
	J2  int
}



























type SequenceMatcher struct {
	a              []string
	b              []string
	b2j            map[string][]int
	IsJunk         func(string) bool
	autoJunk       bool
	bJunk          map[string]struct{}
	matchingBlocks []Match
	fullBCount     map[string]int
	bPopular       map[string]struct{}
	opCodes        []OpCode
}

func NewMatcher(a, b []string) *SequenceMatcher {
	m := SequenceMatcher{autoJunk: true}
	m.SetSeqs(a, b)
	return &m
}

func NewMatcherWithJunk(a, b []string, autoJunk bool,
	isJunk func(string) bool) *SequenceMatcher {

	m := SequenceMatcher{IsJunk: isJunk, autoJunk: autoJunk}
	m.SetSeqs(a, b)
	return &m
}


func (m *SequenceMatcher) SetSeqs(a, b []string) {
	m.SetSeq1(a)
	m.SetSeq2(b)
}










func (m *SequenceMatcher) SetSeq1(a []string) {
	if &a == &m.a {
		return
	}
	m.a = a
	m.matchingBlocks = nil
	m.opCodes = nil
}



func (m *SequenceMatcher) SetSeq2(b []string) {
	if &b == &m.b {
		return
	}
	m.b = b
	m.matchingBlocks = nil
	m.opCodes = nil
	m.fullBCount = nil
	m.chainB()
}

func (m *SequenceMatcher) chainB() {
	
	b2j := map[string][]int{}
	for i, s := range m.b {
		indices := b2j[s]
		indices = append(indices, i)
		b2j[s] = indices
	}

	
	m.bJunk = map[string]struct{}{}
	if m.IsJunk != nil {
		junk := m.bJunk
		for s, _ := range b2j {
			if m.IsJunk(s) {
				junk[s] = struct{}{}
			}
		}
		for s, _ := range junk {
			delete(b2j, s)
		}
	}

	
	popular := map[string]struct{}{}
	n := len(m.b)
	if m.autoJunk && n >= 200 {
		ntest := n/100 + 1
		for s, indices := range b2j {
			if len(indices) > ntest {
				popular[s] = struct{}{}
			}
		}
		for s, _ := range popular {
			delete(b2j, s)
		}
	}
	m.bPopular = popular
	m.b2j = b2j
}

func (m *SequenceMatcher) isBJunk(s string) bool {
	_, ok := m.bJunk[s]
	return ok
}

























func (m *SequenceMatcher) findLongestMatch(alo, ahi, blo, bhi int) Match {
	
	
	
	
	
	
	
	
	
	
	
	besti, bestj, bestsize := alo, blo, 0

	
	
	
	j2len := map[int]int{}
	for i := alo; i != ahi; i++ {
		
		
		newj2len := map[int]int{}
		for _, j := range m.b2j[m.a[i]] {
			
			if j < blo {
				continue
			}
			if j >= bhi {
				break
			}
			k := j2len[j-1] + 1
			newj2len[j] = k
			if k > bestsize {
				besti, bestj, bestsize = i-k+1, j-k+1, k
			}
		}
		j2len = newj2len
	}

	
	
	
	
	for besti > alo && bestj > blo && !m.isBJunk(m.b[bestj-1]) &&
		m.a[besti-1] == m.b[bestj-1] {
		besti, bestj, bestsize = besti-1, bestj-1, bestsize+1
	}
	for besti+bestsize < ahi && bestj+bestsize < bhi &&
		!m.isBJunk(m.b[bestj+bestsize]) &&
		m.a[besti+bestsize] == m.b[bestj+bestsize] {
		bestsize += 1
	}

	
	
	
	
	
	
	
	for besti > alo && bestj > blo && m.isBJunk(m.b[bestj-1]) &&
		m.a[besti-1] == m.b[bestj-1] {
		besti, bestj, bestsize = besti-1, bestj-1, bestsize+1
	}
	for besti+bestsize < ahi && bestj+bestsize < bhi &&
		m.isBJunk(m.b[bestj+bestsize]) &&
		m.a[besti+bestsize] == m.b[bestj+bestsize] {
		bestsize += 1
	}

	return Match{A: besti, B: bestj, Size: bestsize}
}












func (m *SequenceMatcher) GetMatchingBlocks() []Match {
	if m.matchingBlocks != nil {
		return m.matchingBlocks
	}

	var matchBlocks func(alo, ahi, blo, bhi int, matched []Match) []Match
	matchBlocks = func(alo, ahi, blo, bhi int, matched []Match) []Match {
		match := m.findLongestMatch(alo, ahi, blo, bhi)
		i, j, k := match.A, match.B, match.Size
		if match.Size > 0 {
			if alo < i && blo < j {
				matched = matchBlocks(alo, i, blo, j, matched)
			}
			matched = append(matched, match)
			if i+k < ahi && j+k < bhi {
				matched = matchBlocks(i+k, ahi, j+k, bhi, matched)
			}
		}
		return matched
	}
	matched := matchBlocks(0, len(m.a), 0, len(m.b), nil)

	
	
	nonAdjacent := []Match{}
	i1, j1, k1 := 0, 0, 0
	for _, b := range matched {
		
		i2, j2, k2 := b.A, b.B, b.Size
		if i1+k1 == i2 && j1+k1 == j2 {
			
			
			
			k1 += k2
		} else {
			
			
			
			if k1 > 0 {
				nonAdjacent = append(nonAdjacent, Match{i1, j1, k1})
			}
			i1, j1, k1 = i2, j2, k2
		}
	}
	if k1 > 0 {
		nonAdjacent = append(nonAdjacent, Match{i1, j1, k1})
	}

	nonAdjacent = append(nonAdjacent, Match{len(m.a), len(m.b), 0})
	m.matchingBlocks = nonAdjacent
	return m.matchingBlocks
}
















func (m *SequenceMatcher) GetOpCodes() []OpCode {
	if m.opCodes != nil {
		return m.opCodes
	}
	i, j := 0, 0
	matching := m.GetMatchingBlocks()
	opCodes := make([]OpCode, 0, len(matching))
	for _, m := range matching {
		
		
		
		
		
		ai, bj, size := m.A, m.B, m.Size
		tag := byte(0)
		if i < ai && j < bj {
			tag = 'r'
		} else if i < ai {
			tag = 'd'
		} else if j < bj {
			tag = 'i'
		}
		if tag > 0 {
			opCodes = append(opCodes, OpCode{tag, i, ai, j, bj})
		}
		i, j = ai+size, bj+size
		
		
		if size > 0 {
			opCodes = append(opCodes, OpCode{'e', ai, i, bj, j})
		}
	}
	m.opCodes = opCodes
	return m.opCodes
}





func (m *SequenceMatcher) GetGroupedOpCodes(n int) [][]OpCode {
	if n < 0 {
		n = 3
	}
	codes := m.GetOpCodes()
	if len(codes) == 0 {
		codes = []OpCode{OpCode{'e', 0, 1, 0, 1}}
	}
	
	if codes[0].Tag == 'e' {
		c := codes[0]
		i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
		codes[0] = OpCode{c.Tag, max(i1, i2-n), i2, max(j1, j2-n), j2}
	}
	if codes[len(codes)-1].Tag == 'e' {
		c := codes[len(codes)-1]
		i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
		codes[len(codes)-1] = OpCode{c.Tag, i1, min(i2, i1+n), j1, min(j2, j1+n)}
	}
	nn := n + n
	groups := [][]OpCode{}
	group := []OpCode{}
	for _, c := range codes {
		i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
		
		
		if c.Tag == 'e' && i2-i1 > nn {
			group = append(group, OpCode{c.Tag, i1, min(i2, i1+n),
				j1, min(j2, j1+n)})
			groups = append(groups, group)
			group = []OpCode{}
			i1, j1 = max(i1, i2-n), max(j1, j2-n)
		}
		group = append(group, OpCode{c.Tag, i1, i2, j1, j2})
	}
	if len(group) > 0 && !(len(group) == 1 && group[0].Tag == 'e') {
		groups = append(groups, group)
	}
	return groups
}












func (m *SequenceMatcher) Ratio() float64 {
	matches := 0
	for _, m := range m.GetMatchingBlocks() {
		matches += m.Size
	}
	return calculateRatio(matches, len(m.a)+len(m.b))
}





func (m *SequenceMatcher) QuickRatio() float64 {
	
	
	
	if m.fullBCount == nil {
		m.fullBCount = map[string]int{}
		for _, s := range m.b {
			m.fullBCount[s] = m.fullBCount[s] + 1
		}
	}

	
	
	avail := map[string]int{}
	matches := 0
	for _, s := range m.a {
		n, ok := avail[s]
		if !ok {
			n = m.fullBCount[s]
		}
		avail[s] = n - 1
		if n > 0 {
			matches += 1
		}
	}
	return calculateRatio(matches, len(m.a)+len(m.b))
}





func (m *SequenceMatcher) RealQuickRatio() float64 {
	la, lb := len(m.a), len(m.b)
	return calculateRatio(min(la, lb), la+lb)
}


func formatRangeUnified(start, stop int) string {
	
	beginning := start + 1 
	length := stop - start
	if length == 1 {
		return fmt.Sprintf("%d", beginning)
	}
	if length == 0 {
		beginning -= 1 
	}
	return fmt.Sprintf("%d,%d", beginning, length)
}


type UnifiedDiff struct {
	A        []string 
	FromFile string   
	FromDate string   
	B        []string 
	ToFile   string   
	ToDate   string   
	Eol      string   
	Context  int      
}




















func WriteUnifiedDiff(writer io.Writer, diff UnifiedDiff) error {
	buf := bufio.NewWriter(writer)
	defer buf.Flush()
	wf := func(format string, args ...interface{}) error {
		_, err := buf.WriteString(fmt.Sprintf(format, args...))
		return err
	}
	ws := func(s string) error {
		_, err := buf.WriteString(s)
		return err
	}

	if len(diff.Eol) == 0 {
		diff.Eol = "\n"
	}

	started := false
	m := NewMatcher(diff.A, diff.B)
	for _, g := range m.GetGroupedOpCodes(diff.Context) {
		if !started {
			started = true
			fromDate := ""
			if len(diff.FromDate) > 0 {
				fromDate = "\t" + diff.FromDate
			}
			toDate := ""
			if len(diff.ToDate) > 0 {
				toDate = "\t" + diff.ToDate
			}
			if diff.FromFile != "" || diff.ToFile != "" {
				err := wf("--- %s%s%s", diff.FromFile, fromDate, diff.Eol)
				if err != nil {
					return err
				}
				err = wf("+++ %s%s%s", diff.ToFile, toDate, diff.Eol)
				if err != nil {
					return err
				}
			}
		}
		first, last := g[0], g[len(g)-1]
		range1 := formatRangeUnified(first.I1, last.I2)
		range2 := formatRangeUnified(first.J1, last.J2)
		if err := wf("@@ -%s +%s @@%s", range1, range2, diff.Eol); err != nil {
			return err
		}
		for _, c := range g {
			i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
			if c.Tag == 'e' {
				for _, line := range diff.A[i1:i2] {
					if err := ws(" " + line); err != nil {
						return err
					}
				}
				continue
			}
			if c.Tag == 'r' || c.Tag == 'd' {
				for _, line := range diff.A[i1:i2] {
					if err := ws("-" + line); err != nil {
						return err
					}
				}
			}
			if c.Tag == 'r' || c.Tag == 'i' {
				for _, line := range diff.B[j1:j2] {
					if err := ws("+" + line); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}


func GetUnifiedDiffString(diff UnifiedDiff) (string, error) {
	w := &bytes.Buffer{}
	err := WriteUnifiedDiff(w, diff)
	return string(w.Bytes()), err
}


func formatRangeContext(start, stop int) string {
	
	beginning := start + 1 
	length := stop - start
	if length == 0 {
		beginning -= 1 
	}
	if length <= 1 {
		return fmt.Sprintf("%d", beginning)
	}
	return fmt.Sprintf("%d,%d", beginning, beginning+length-1)
}

type ContextDiff UnifiedDiff


















func WriteContextDiff(writer io.Writer, diff ContextDiff) error {
	buf := bufio.NewWriter(writer)
	defer buf.Flush()
	var diffErr error
	wf := func(format string, args ...interface{}) {
		_, err := buf.WriteString(fmt.Sprintf(format, args...))
		if diffErr == nil && err != nil {
			diffErr = err
		}
	}
	ws := func(s string) {
		_, err := buf.WriteString(s)
		if diffErr == nil && err != nil {
			diffErr = err
		}
	}

	if len(diff.Eol) == 0 {
		diff.Eol = "\n"
	}

	prefix := map[byte]string{
		'i': "+ ",
		'd': "- ",
		'r': "! ",
		'e': "  ",
	}

	started := false
	m := NewMatcher(diff.A, diff.B)
	for _, g := range m.GetGroupedOpCodes(diff.Context) {
		if !started {
			started = true
			fromDate := ""
			if len(diff.FromDate) > 0 {
				fromDate = "\t" + diff.FromDate
			}
			toDate := ""
			if len(diff.ToDate) > 0 {
				toDate = "\t" + diff.ToDate
			}
			if diff.FromFile != "" || diff.ToFile != "" {
				wf("*** %s%s%s", diff.FromFile, fromDate, diff.Eol)
				wf("--- %s%s%s", diff.ToFile, toDate, diff.Eol)
			}
		}

		first, last := g[0], g[len(g)-1]
		ws("***************" + diff.Eol)

		range1 := formatRangeContext(first.I1, last.I2)
		wf("*** %s ****%s", range1, diff.Eol)
		for _, c := range g {
			if c.Tag == 'r' || c.Tag == 'd' {
				for _, cc := range g {
					if cc.Tag == 'i' {
						continue
					}
					for _, line := range diff.A[cc.I1:cc.I2] {
						ws(prefix[cc.Tag] + line)
					}
				}
				break
			}
		}

		range2 := formatRangeContext(first.J1, last.J2)
		wf("--- %s ----%s", range2, diff.Eol)
		for _, c := range g {
			if c.Tag == 'r' || c.Tag == 'i' {
				for _, cc := range g {
					if cc.Tag == 'd' {
						continue
					}
					for _, line := range diff.B[cc.J1:cc.J2] {
						ws(prefix[cc.Tag] + line)
					}
				}
				break
			}
		}
	}
	return diffErr
}


func GetContextDiffString(diff ContextDiff) (string, error) {
	w := &bytes.Buffer{}
	err := WriteContextDiff(w, diff)
	return string(w.Bytes()), err
}



func SplitLines(s string) []string {
	lines := strings.SplitAfter(s, "\n")
	lines[len(lines)-1] += "\n"
	return lines
}
