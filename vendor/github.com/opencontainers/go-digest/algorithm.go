













package digest

import (
	"crypto"
	"fmt"
	"hash"
	"io"
	"regexp"
)




type Algorithm string


const (
	SHA256 Algorithm = "sha256" 
	SHA384 Algorithm = "sha384" 
	SHA512 Algorithm = "sha512" 

	
	
	
	Canonical = SHA256
)

var (
	
	
	

	
	
	algorithms = map[Algorithm]crypto.Hash{
		SHA256: crypto.SHA256,
		SHA384: crypto.SHA384,
		SHA512: crypto.SHA512,
	}

	
	
	anchoredEncodedRegexps = map[Algorithm]*regexp.Regexp{
		SHA256: regexp.MustCompile(`^[a-f0-9]{64}$`),
		SHA384: regexp.MustCompile(`^[a-f0-9]{96}$`),
		SHA512: regexp.MustCompile(`^[a-f0-9]{128}$`),
	}
)



func (a Algorithm) Available() bool {
	h, ok := algorithms[a]
	if !ok {
		return false
	}

	
	return h.Available()
}

func (a Algorithm) String() string {
	return string(a)
}


func (a Algorithm) Size() int {
	h, ok := algorithms[a]
	if !ok {
		return 0
	}
	return h.Size()
}


func (a *Algorithm) Set(value string) error {
	if value == "" {
		*a = Canonical
	} else {
		
		*a = Algorithm(value)
	}

	if !a.Available() {
		return ErrDigestUnsupported
	}

	return nil
}




func (a Algorithm) Digester() Digester {
	return &digester{
		alg:  a,
		hash: a.Hash(),
	}
}



func (a Algorithm) Hash() hash.Hash {
	if !a.Available() {
		
		if a == "" {
			panic(fmt.Sprintf("empty digest algorithm, validate before calling Algorithm.Hash()"))
		}

		
		
		
		
		
		
		
		panic(fmt.Sprintf("%v not available (make sure it is imported)", a))
	}

	return algorithms[a].New()
}



func (a Algorithm) Encode(d []byte) string {
	
	
	return fmt.Sprintf("%x", d)
}


func (a Algorithm) FromReader(rd io.Reader) (Digest, error) {
	digester := a.Digester()

	if _, err := io.Copy(digester.Hash(), rd); err != nil {
		return "", err
	}

	return digester.Digest(), nil
}


func (a Algorithm) FromBytes(p []byte) Digest {
	digester := a.Digester()

	if _, err := digester.Hash().Write(p); err != nil {
		
		
		
		
		
		panic("write to hash function returned error: " + err.Error())
	}

	return digester.Digest()
}


func (a Algorithm) FromString(s string) Digest {
	return a.FromBytes([]byte(s))
}


func (a Algorithm) Validate(encoded string) error {
	r, ok := anchoredEncodedRegexps[a]
	if !ok {
		return ErrDigestUnsupported
	}
	
	
	if a.Size()*2 != len(encoded) {
		return ErrDigestInvalidLength
	}
	if r.MatchString(encoded) {
		return nil
	}
	return ErrDigestInvalidFormat
}
