











package bidi 
















type Direction int

const (
	
	
	
	LeftToRight Direction = iota

	
	
	
	RightToLeft

	
	
	Mixed

	
	
	Neutral
)

type options struct{}


type Option func(*options)












func DefaultDirection(d Direction) Option {
	panic("unimplemented")
}


type Paragraph struct {
	
}






func (p *Paragraph) SetBytes(b []byte, opts ...Option) (n int, err error) {
	panic("unimplemented")
}






func (p *Paragraph) SetString(s string, opts ...Option) (n int, err error) {
	panic("unimplemented")
}




func (p *Paragraph) IsLeftToRight() bool {
	panic("unimplemented")
}




func (p *Paragraph) Direction() Direction {
	panic("unimplemented")
}




func (p *Paragraph) RunAt(pos int) Run {
	panic("unimplemented")
}


func (p *Paragraph) Order() (Ordering, error) {
	panic("unimplemented")
}



func (p *Paragraph) Line(start, end int) (Ordering, error) {
	panic("unimplemented")
}




type Ordering struct{}




func (o *Ordering) Direction() Direction {
	panic("unimplemented")
}


func (o *Ordering) NumRuns() int {
	panic("unimplemented")
}


func (o *Ordering) Run(i int) Run {
	panic("unimplemented")
}









type Run struct {
}


func (r *Run) String() string {
	panic("unimplemented")
}


func (r *Run) Bytes() []byte {
	panic("unimplemented")
}







func (r *Run) Direction() Direction {
	panic("unimplemented")
}



func (r *Run) Pos() (start, end int) {
	panic("unimplemented")
}




func AppendReverse(out, in []byte) []byte {
	panic("unimplemented")
}




func ReverseString(s string) string {
	panic("unimplemented")
}
