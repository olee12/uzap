package bufferpool

import "github.com/olee12/zap/buffer"

var (
	_pool = buffer.NewPool()
	Get   = _pool.Get
)
