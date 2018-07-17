

package grpc

import (
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto" 
)




type baseCodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var _ baseCodec = Codec(nil)
var _ baseCodec = encoding.Codec(nil)






type Codec interface {
	
	Marshal(v interface{}) ([]byte, error)
	
	Unmarshal(data []byte, v interface{}) error
	
	
	String() string
}
