













package pbutil

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)







func WriteDelimited(w io.Writer, m proto.Message) (n int, err error) {
	buffer, err := proto.Marshal(m)
	if err != nil {
		return 0, err
	}

	var buf [binary.MaxVarintLen32]byte
	encodedLength := binary.PutUvarint(buf[:], uint64(len(buffer)))

	sync, err := w.Write(buf[:encodedLength])
	if err != nil {
		return sync, err
	}

	n, err = w.Write(buffer)
	return n + sync, err
}
