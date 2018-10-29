





package gcimporter

import "go/types"

func newInterface(methods []*types.Func, embeddeds []types.Type) *types.Interface {
	named := make([]*types.Named, len(embeddeds))
	for i, e := range embeddeds {
		var ok bool
		named[i], ok = e.(*types.Named)
		if !ok {
			panic("embedding of non-defined interfaces in interfaces is not supported before Go 1.11")
		}
	}
	return types.NewInterface(methods, named)
}
