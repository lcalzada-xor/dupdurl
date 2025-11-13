package pool

import (
	"strings"
	"sync"
)

// StringBuilderPool is a pool of strings.Builder for reducing allocations
var StringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// GetBuilder gets a builder from the pool
func GetBuilder() *strings.Builder {
	sb := StringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// PutBuilder returns a builder to the pool
func PutBuilder(sb *strings.Builder) {
	StringBuilderPool.Put(sb)
}

// ByteSlicePool is a pool of byte slices for reducing allocations
var ByteSlicePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 256) // Default capacity
		return &b
	},
}

// GetByteSlice gets a byte slice from the pool
func GetByteSlice() *[]byte {
	return ByteSlicePool.Get().(*[]byte)
}

// PutByteSlice returns a byte slice to the pool
func PutByteSlice(b *[]byte) {
	*b = (*b)[:0] // Reset length but keep capacity
	ByteSlicePool.Put(b)
}

// MapPool provides pre-sized maps for common operations
type MapPool struct {
	pool sync.Pool
	size int
}

// NewMapPool creates a new map pool with initial size hint
func NewMapPool(initialSize int) *MapPool {
	return &MapPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make(map[string]struct{}, initialSize)
			},
		},
		size: initialSize,
	}
}

// Get gets a map from the pool
func (mp *MapPool) Get() map[string]struct{} {
	return mp.pool.Get().(map[string]struct{})
}

// Put returns a map to the pool (clears it first)
func (mp *MapPool) Put(m map[string]struct{}) {
	// Clear the map
	for k := range m {
		delete(m, k)
	}
	mp.pool.Put(m)
}
