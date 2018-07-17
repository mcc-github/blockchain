













package digest

import "hash"




type Digester interface {
	Hash() hash.Hash 
	Digest() Digest
}


type digester struct {
	alg  Algorithm
	hash hash.Hash
}

func (d *digester) Hash() hash.Hash {
	return d.hash
}

func (d *digester) Digest() Digest {
	return NewDigest(d.alg, d.hash)
}
