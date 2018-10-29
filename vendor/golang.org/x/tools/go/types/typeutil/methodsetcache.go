





package typeutil

import (
	"go/types"
	"sync"
)




type MethodSetCache struct {
	mu     sync.Mutex
	named  map[*types.Named]struct{ value, pointer *types.MethodSet } 
	others map[types.Type]*types.MethodSet                            
}







func (cache *MethodSetCache) MethodSet(T types.Type) *types.MethodSet {
	if cache == nil {
		return types.NewMethodSet(T)
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	switch T := T.(type) {
	case *types.Named:
		return cache.lookupNamed(T).value

	case *types.Pointer:
		if N, ok := T.Elem().(*types.Named); ok {
			return cache.lookupNamed(N).pointer
		}
	}

	
	
	mset := cache.others[T]
	if mset == nil {
		mset = types.NewMethodSet(T)
		if cache.others == nil {
			cache.others = make(map[types.Type]*types.MethodSet)
		}
		cache.others[T] = mset
	}
	return mset
}

func (cache *MethodSetCache) lookupNamed(named *types.Named) struct{ value, pointer *types.MethodSet } {
	if cache.named == nil {
		cache.named = make(map[*types.Named]struct{ value, pointer *types.MethodSet })
	}
	
	
	msets, ok := cache.named[named]
	if !ok {
		msets.value = types.NewMethodSet(named)
		msets.pointer = types.NewMethodSet(types.NewPointer(named))
		cache.named[named] = msets
	}
	return msets
}
