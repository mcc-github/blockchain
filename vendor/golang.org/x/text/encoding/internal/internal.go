




package internal

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/transform"
)



type Encoding struct {
	encoding.Encoding
	Name string
	MIB  identifier.MIB
}


var _ identifier.Interface = (*Encoding)(nil)

func (e *Encoding) String() string {
	return e.Name
}

func (e *Encoding) ID() (mib identifier.MIB, other string) {
	return e.MIB, ""
}


type SimpleEncoding struct {
	Decoder transform.Transformer
	Encoder transform.Transformer
}

func (e *SimpleEncoding) NewDecoder() *encoding.Decoder {
	return &encoding.Decoder{Transformer: e.Decoder}
}

func (e *SimpleEncoding) NewEncoder() *encoding.Encoder {
	return &encoding.Encoder{Transformer: e.Encoder}
}



type FuncEncoding struct {
	Decoder func() transform.Transformer
	Encoder func() transform.Transformer
}

func (e FuncEncoding) NewDecoder() *encoding.Decoder {
	return &encoding.Decoder{Transformer: e.Decoder()}
}

func (e FuncEncoding) NewEncoder() *encoding.Encoder {
	return &encoding.Encoder{Transformer: e.Encoder()}
}




type RepertoireError byte


func (r RepertoireError) Error() string {
	return "encoding: rune not supported by encoding."
}


func (r RepertoireError) Replacement() byte { return byte(r) }

var ErrASCIIReplacement = RepertoireError(encoding.ASCIISub)
