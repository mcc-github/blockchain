





package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb/filter"
)

type iFilter struct {
	filter.Filter
}

func (f iFilter) Contains(filter, key []byte) bool {
	return f.Filter.Contains(filter, internalKey(key).ukey())
}

func (f iFilter) NewGenerator() filter.FilterGenerator {
	return iFilterGenerator{f.Filter.NewGenerator()}
}

type iFilterGenerator struct {
	filter.FilterGenerator
}

func (g iFilterGenerator) Add(key []byte) {
	g.FilterGenerator.Add(internalKey(key).ukey())
}
