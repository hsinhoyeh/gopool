package gopool

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestMapPool(t *testing.T) {
	s := NewMapPool(func() map[string]interface{} { return make(map[string]interface{}) })
	b := s.Map()
	b["pi"] = float32(3.14)
	b["golden_ratio"] = float32(1.61803398875)

	exp := make(map[string]interface{})
	exp["pi"] = float32(3.14)
	exp["golden_ratio"] = float32(1.61803398875)

	if !reflect.DeepEqual(b, exp) {
		t.Errorf("exp: %v, but got: %v", exp, b)
	}
	s.Recycle(b)

	b = s.Map()
	exp = make(map[string]interface{})

	if !reflect.DeepEqual(b, exp) {
		t.Errorf("exp: %v, but got: %v", exp, b)
	}
	s.Recycle(b)
}

func TestSlicePool(t *testing.T) {
	s := NewSlicePool(func() []byte { return make([]byte, 8) })
	// use s.Bytes instead of make([]byte, ...)
	b := s.Bytes()
	opByteSlice(b)
	exp := bytes.Repeat([]byte{'a'}, 8)
	if 0 != bytes.Compare(b, exp) {
		t.Errorf("exp: %b, but got: %b", exp, b)
	}
	s.Recycle(b)

	b = s.Bytes()
	exp = make([]byte, 8)
	if len(b) != len(exp) {
		t.Errorf("exp: %b, but got: %b", exp, b)
	}
	s.Recycle(b)
}

var mapPoolTestFactory = func(numRound int, lens int) func() {
	s := NewMapPool(func() map[string]interface{} { return make(map[string]interface{}) })
	return func() {
		for j := 0; j < numRound; j++ {
			m := s.Map()
			opMap(m, lens)
			s.Recycle(m)
		}

	}
}

var nativeMapTestFactory = func(numRound int, lens int) func() {
	return func() {
		for j := 0; j < numRound; j++ {
			m := make(map[string]interface{})
			opMap(m, lens)
			m = nil
		}
	}
}

var slicePoolTestFactory = func(numRound int, sliceLen int) func() {
	s := NewSlicePool(func() []byte { return make([]byte, sliceLen) })
	return func() {
		for j := 0; j < numRound; j++ {
			b := s.Bytes()
			opByteSlice(b)
			s.Recycle(b)
		}
	}
}

var nativeSliceTestFactory = func(numRound int, sliceLen int) func() {
	return func() {
		for j := 0; j < numRound; j++ {
			b := make([]byte, sliceLen)
			opByteSlice(b)
			b = nil
		}
	}
}

func BenchmarkMapPool4K(b *testing.B) {
	doBenchmarkSlice(b, 1, mapPoolTestFactory(128, 4096))
}

func BenchmarkNativeMap4K(b *testing.B) {
	doBenchmarkSlice(b, 1, nativeMapTestFactory(128, 4096))
}

func BenchmarkSlicePool4K(b *testing.B) {
	doBenchmarkSlice(b, 1, slicePoolTestFactory(128, 4096))
}

func BenchmarkNativeSlice4K(b *testing.B) {
	doBenchmarkSlice(b, 1, nativeSliceTestFactory(128, 4096))
}

func BenchmarkSlicePool4K100Routines(b *testing.B) {
	doBenchmarkSlice(b, 100, slicePoolTestFactory(128, 4096))
}

func BenchmarkNativeSlice4K100Routines(b *testing.B) {
	doBenchmarkSlice(b, 100, nativeSliceTestFactory(128, 4096))
}

func doBenchmarkSlice(b *testing.B, numRoutines int, testFactory func()) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}

		// define a tester
		tester := func() {
			testFactory()
			wg.Done()
		}

		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go tester()
		}
		wg.Wait()
	}
}

func opMap(m map[string]interface{}, lens int) {
	for i := 0; i < lens; i++ {
		m[fmt.Sprintf("mapkey_%d", i)] = fmt.Sprintf("mapvalue_%d", i)
	}
}

func opByteSlice(s []byte) {
	for i := range s {
		s[i] = 'a'
	}
}
