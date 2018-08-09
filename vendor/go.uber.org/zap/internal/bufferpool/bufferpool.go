





















package bufferpool

import "go.uber.org/zap/buffer"

var (
	_pool = buffer.NewPool()
	
	Get = _pool.Get
)
