





package encoding

import (
	"io"
	"strings"
)



const Identity = "identity"



type Compressor interface {
	
	
	
	Compress(w io.Writer) (io.WriteCloser, error)
	
	
	
	Decompress(r io.Reader) (io.Reader, error)
	
	
	
	Name() string
}

var registeredCompressor = make(map[string]Compressor)










func RegisterCompressor(c Compressor) {
	registeredCompressor[c.Name()] = c
}


func GetCompressor(name string) Compressor {
	return registeredCompressor[name]
}




type Codec interface {
	
	Marshal(v interface{}) ([]byte, error)
	
	Unmarshal(data []byte, v interface{}) error
	
	
	
	Name() string
}

var registeredCodecs = make(map[string]Codec)















func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	contentSubtype := strings.ToLower(codec.Name())
	if contentSubtype == "" {
		panic("cannot register Codec with empty string result for String()")
	}
	registeredCodecs[contentSubtype] = codec
}





func GetCodec(contentSubtype string) Codec {
	return registeredCodecs[contentSubtype]
}
