## Overview

gopool implements a resource pool of golang, which utilizes sync.Pool introduced in golang 1.3.

In current implementation, we support byte slice and map types in the pool. Here is an example showing how to use it

```
import (
	github.com/hsinhoyeh/gopool
)

...
mapFactory := func() map[string]interface{} {
	return make(map[string]interface{})
}

s := NewMapPool(mapFactory)

// call s.Map to obtain a new map instance
m := s.Map()

// call s.Recycle to recycle the map instance
s.Recycle(m)

```

Similarity, you can define a byte slice factory to operate the slicePool.

```
import (
	github.com/hsinhoyeh/gopool
)

...
sliceFactory := func() ]byte {
	return make([]byte, 128)
}

s := NewSlicePool(sliceFactory)

// call s.Bytes to obtain a new byte slice instance
b := s.Bytes()

// call s.Recycle to recycle the slice
s.Recycle(b)

```

######map benchmark
here we also conduct the map benchmark result on mac air for the reference:

name       | cost
-----------|-------------
BenchmarkMapPool4K | 443726525 ns/op
BenchmarkNativeMap4K | 532768829 ns/op

######slice benchmark

name       | cost
-----------|-------------
BenchmarkSlicePool4K | 380693 ns/op
BenchmarkNativeSlice4K | 834490 ns/op

######benchmark on goroutines
name       | cost
-----------|-------------
BenchmarkSlicePool4K100Routines | 39230265 ns/op
BenchmarkNativeSlice4K100Routines| 86492977 ns/op

