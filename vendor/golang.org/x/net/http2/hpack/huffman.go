



package hpack

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}




func HuffmanDecode(w io.Writer, v []byte) (int, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if err := huffmanDecode(buf, 0, v); err != nil {
		return 0, err
	}
	return w.Write(buf.Bytes())
}


func HuffmanDecodeToString(v []byte) (string, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if err := huffmanDecode(buf, 0, v); err != nil {
		return "", err
	}
	return buf.String(), nil
}



var ErrInvalidHuffman = errors.New("hpack: invalid Huffman-encoded data")




func huffmanDecode(buf *bytes.Buffer, maxLen int, v []byte) error {
	n := rootHuffmanNode
	
	
	
	cur, cbits, sbits := uint(0), uint8(0), uint8(0)
	for _, b := range v {
		cur = cur<<8 | uint(b)
		cbits += 8
		sbits += 8
		for cbits >= 8 {
			idx := byte(cur >> (cbits - 8))
			n = n.children[idx]
			if n == nil {
				return ErrInvalidHuffman
			}
			if n.children == nil {
				if maxLen != 0 && buf.Len() == maxLen {
					return ErrStringLength
				}
				buf.WriteByte(n.sym)
				cbits -= n.codeLen
				n = rootHuffmanNode
				sbits = cbits
			} else {
				cbits -= 8
			}
		}
	}
	for cbits > 0 {
		n = n.children[byte(cur<<(8-cbits))]
		if n == nil {
			return ErrInvalidHuffman
		}
		if n.children != nil || n.codeLen > cbits {
			break
		}
		if maxLen != 0 && buf.Len() == maxLen {
			return ErrStringLength
		}
		buf.WriteByte(n.sym)
		cbits -= n.codeLen
		n = rootHuffmanNode
		sbits = cbits
	}
	if sbits > 7 {
		
		
		return ErrInvalidHuffman
	}
	if mask := uint(1<<cbits - 1); cur&mask != mask {
		
		return ErrInvalidHuffman
	}

	return nil
}

type node struct {
	
	children []*node

	
	codeLen uint8 
	sym     byte  
}

func newInternalNode() *node {
	return &node{children: make([]*node, 256)}
}

var rootHuffmanNode = newInternalNode()

func init() {
	if len(huffmanCodes) != 256 {
		panic("unexpected size")
	}
	for i, code := range huffmanCodes {
		addDecoderNode(byte(i), code, huffmanCodeLen[i])
	}
}

func addDecoderNode(sym byte, code uint32, codeLen uint8) {
	cur := rootHuffmanNode
	for codeLen > 8 {
		codeLen -= 8
		i := uint8(code >> codeLen)
		if cur.children[i] == nil {
			cur.children[i] = newInternalNode()
		}
		cur = cur.children[i]
	}
	shift := 8 - codeLen
	start, end := int(uint8(code<<shift)), int(1<<shift)
	for i := start; i < start+end; i++ {
		cur.children[i] = &node{sym: sym, codeLen: codeLen}
	}
}



func AppendHuffmanString(dst []byte, s string) []byte {
	rembits := uint8(8)

	for i := 0; i < len(s); i++ {
		if rembits == 8 {
			dst = append(dst, 0)
		}
		dst, rembits = appendByteToHuffmanCode(dst, rembits, s[i])
	}

	if rembits < 8 {
		
		code := uint32(0x3fffffff)
		nbits := uint8(30)

		t := uint8(code >> (nbits - rembits))
		dst[len(dst)-1] |= t
	}

	return dst
}



func HuffmanEncodeLength(s string) uint64 {
	n := uint64(0)
	for i := 0; i < len(s); i++ {
		n += uint64(huffmanCodeLen[s[i]])
	}
	return (n + 7) / 8
}





func appendByteToHuffmanCode(dst []byte, rembits uint8, c byte) ([]byte, uint8) {
	code := huffmanCodes[c]
	nbits := huffmanCodeLen[c]

	for {
		if rembits > nbits {
			t := uint8(code << (rembits - nbits))
			dst[len(dst)-1] |= t
			rembits -= nbits
			break
		}

		t := uint8(code >> (nbits - rembits))
		dst[len(dst)-1] |= t

		nbits -= rembits
		rembits = 8

		if nbits == 0 {
			break
		}

		dst = append(dst, 0)
	}

	return dst, rembits
}
