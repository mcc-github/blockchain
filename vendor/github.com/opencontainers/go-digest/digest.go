













package digest

import (
	"fmt"
	"hash"
	"io"
	"regexp"
	"strings"
)












type Digest string


func NewDigest(alg Algorithm, h hash.Hash) Digest {
	return NewDigestFromBytes(alg, h.Sum(nil))
}





func NewDigestFromBytes(alg Algorithm, p []byte) Digest {
	return NewDigestFromEncoded(alg, alg.Encode(p))
}


func NewDigestFromHex(alg, hex string) Digest {
	return NewDigestFromEncoded(Algorithm(alg), hex)
}


func NewDigestFromEncoded(alg Algorithm, encoded string) Digest {
	return Digest(fmt.Sprintf("%s:%s", alg, encoded))
}


var DigestRegexp = regexp.MustCompile(`[a-z0-9]+(?:[.+_-][a-z0-9]+)*:[a-zA-Z0-9=_-]+`)


var DigestRegexpAnchored = regexp.MustCompile(`^` + DigestRegexp.String() + `$`)

var (
	
	ErrDigestInvalidFormat = fmt.Errorf("invalid checksum digest format")

	
	ErrDigestInvalidLength = fmt.Errorf("invalid checksum digest length")

	
	ErrDigestUnsupported = fmt.Errorf("unsupported digest algorithm")
)



func Parse(s string) (Digest, error) {
	d := Digest(s)
	return d, d.Validate()
}


func FromReader(rd io.Reader) (Digest, error) {
	return Canonical.FromReader(rd)
}


func FromBytes(p []byte) Digest {
	return Canonical.FromBytes(p)
}


func FromString(s string) Digest {
	return Canonical.FromString(s)
}



func (d Digest) Validate() error {
	s := string(d)
	i := strings.Index(s, ":")
	if i <= 0 || i+1 == len(s) {
		return ErrDigestInvalidFormat
	}
	algorithm, encoded := Algorithm(s[:i]), s[i+1:]
	if !algorithm.Available() {
		if !DigestRegexpAnchored.MatchString(s) {
			return ErrDigestInvalidFormat
		}
		return ErrDigestUnsupported
	}
	return algorithm.Validate(encoded)
}



func (d Digest) Algorithm() Algorithm {
	return Algorithm(d[:d.sepIndex()])
}



func (d Digest) Verifier() Verifier {
	return hashVerifier{
		hash:   d.Algorithm().Hash(),
		digest: d,
	}
}



func (d Digest) Encoded() string {
	return string(d[d.sepIndex()+1:])
}


func (d Digest) Hex() string {
	return d.Encoded()
}

func (d Digest) String() string {
	return string(d)
}

func (d Digest) sepIndex() int {
	i := strings.Index(string(d), ":")

	if i < 0 {
		panic(fmt.Sprintf("no ':' separator in digest %q", d))
	}

	return i
}
