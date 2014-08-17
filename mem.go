package gopool

import "sync"

// ByteSliceFactory defines a factory function used to allocate a slice of bytes
type ByteSliceFactory func() []byte

// MapFactory defines a factory function used to allocates a map
type MapFactory func() map[string]interface{}

type genericPool struct {
	pool *sync.Pool
	// Newed is s stats used to record the number of times of allocation
	Newed int64
}

func newGenericPool(newfactory func() interface{}) *genericPool {
	g := &genericPool{
		Newed: 0,
	}
	g.pool = &sync.Pool{
		New: func() interface{} {
			g.Newed = g.Newed + 1
			return newfactory()
		},
	}
	return g
}

type slicePool struct {
	gp *genericPool
}

// NewSlicePool accepts a byteslice factory b and returns a reference of slicePool
// we use the byteslice factory to allocates new slice if we donot have one.
func NewSlicePool(b ByteSliceFactory) *slicePool {
	return &slicePool{
		gp: newGenericPool(func() interface{} { return b() }),
	}
}

// Bytes returns a slice of byte from the resource pool.
// TODO: initialize the values in the bytes to be zero
func (s *slicePool) Bytes() []byte {
	b := s.gp.pool.Get().([]byte)
	return b[:cap(b)]
}

// Recycle recycles the byte slices into the resource pool.
func (s *slicePool) Recycle(b []byte) {
	s.gp.pool.Put(b)
}

type mapPool struct {
	gp *genericPool
}

// NewMapPool accepts a mapfactory m and returns a reference of mapPool
// we use the mapfactory to allocates new map if we donot have one.
func NewMapPool(m MapFactory) *mapPool {
	return &mapPool{
		gp: newGenericPool(func() interface{} { return m() }),
	}
}

// Map returns a map reference from the resource pool.
func (s *mapPool) Map() map[string]interface{} {
	m := s.gp.pool.Get().(map[string]interface{})
	// cleanup the previous keys
	for k := range m {
		delete(m, k)
	}
	return m
}

// Recycle recycles the byte slices into the resource pool.
func (s *mapPool) Recycle(m map[string]interface{}) {
	s.gp.pool.Put(m)
}
