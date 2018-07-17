



package triegen



import "io"





type Compacter interface {
	
	
	Size(v []uint64) (sz int, ok bool)

	
	
	
	Store(v []uint64) uint32

	
	Print(w io.Writer) error

	
	
	
	
	
	
	Handler() string
}



type simpleCompacter builder

func (b *simpleCompacter) Size([]uint64) (sz int, ok bool) {
	return blockSize * b.ValueSize, true
}

func (b *simpleCompacter) Store(v []uint64) uint32 {
	h := uint32(len(b.ValueBlocks) - blockOffset)
	b.ValueBlocks = append(b.ValueBlocks, v)
	return h
}

func (b *simpleCompacter) Print(io.Writer) error {
	
	return nil
}

func (b *simpleCompacter) Handler() string {
	panic("Handler should be special-cased for this Compacter")
}
