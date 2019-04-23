














package ioutil

import "io"



func NewLimitedBufferReader(r io.Reader, n int) io.Reader {
	return &limitedBufferReader{
		r: r,
		n: n,
	}
}

type limitedBufferReader struct {
	r io.Reader
	n int
}

func (r *limitedBufferReader) Read(p []byte) (n int, err error) {
	np := p
	if len(np) > r.n {
		np = np[:r.n]
	}
	return r.r.Read(np)
}
