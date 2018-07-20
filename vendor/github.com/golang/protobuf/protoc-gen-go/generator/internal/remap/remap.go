































package remap

import (
	"fmt"
	"go/scanner"
	"go/token"
)


type Location struct {
	Pos, End int 
}



type Map map[Location]Location




func (m Map) Find(pos, end int) (Location, bool) {
	key := Location{
		Pos: pos,
		End: end,
	}
	if loc, ok := m[key]; ok {
		return loc, true
	}
	return key, false
}

func (m Map) add(opos, oend, npos, nend int) {
	m[Location{Pos: opos, End: oend}] = Location{Pos: npos, End: nend}
}



func Compute(input, output []byte) (Map, error) {
	itok := tokenize(input)
	otok := tokenize(output)
	if len(itok) != len(otok) {
		return nil, fmt.Errorf("wrong number of tokens, %d ≠ %d", len(itok), len(otok))
	}
	m := make(Map)
	for i, ti := range itok {
		to := otok[i]
		if ti.Token != to.Token {
			return nil, fmt.Errorf("token %d type mismatch: %s ≠ %s", i+1, ti, to)
		}
		m.add(ti.pos, ti.end, to.pos, to.end)
	}
	return m, nil
}


type tokinfo struct {
	pos, end int
	token.Token
}

func tokenize(src []byte) []tokinfo {
	fs := token.NewFileSet()
	var s scanner.Scanner
	s.Init(fs.AddFile("src", fs.Base(), len(src)), src, nil, scanner.ScanComments)
	var info []tokinfo
	for {
		pos, next, lit := s.Scan()
		switch next {
		case token.SEMICOLON:
			continue
		}
		info = append(info, tokinfo{
			pos:   int(pos - 1),
			end:   int(pos + token.Pos(len(lit)) - 1),
			Token: next,
		})
		if next == token.EOF {
			break
		}
	}
	return info
}
