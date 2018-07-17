













package pbutil

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/golang/protobuf/proto"
)

var errInvalidVarint = errors.New("invalid varint32 encountered")












func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	
	
	var headerBuf [binary.MaxVarintLen32]byte
	var bytesRead, varIntBytes int
	var messageLength uint64
	for varIntBytes == 0 { 
		if bytesRead >= len(headerBuf) {
			return bytesRead, errInvalidVarint
		}
		
		
		
		newBytesRead, err := r.Read(headerBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if err != nil {
				return bytesRead, err
			}
			
			
			
			continue
		}
		bytesRead += newBytesRead
		
		
		messageLength, varIntBytes = proto.DecodeVarint(headerBuf[:bytesRead])
	}

	messageBuf := make([]byte, messageLength)
	newBytesRead, err := io.ReadFull(r, messageBuf)
	bytesRead += newBytesRead
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, proto.Unmarshal(messageBuf, m)
}
