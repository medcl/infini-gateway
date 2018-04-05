package sbloom

import (
	"hash/fnv"
	"testing"
)

func TestFilterFalsePositive(t *testing.T) {
	const (
		prob  = 1
		items = 15
	)
	f := NewFilter(fnv.New64(), prob)

	for i := 0; i < 1<<items; i++ {
		f.Add(randSeed())
	}

	//query the same number of times
	//expect around 1 << 14
	var found int
	for i := 0; i < 1<<items; i++ {
		if f.Lookup(randSeed()) {
			found++
		}
	}

	if exp := 1 << (items - prob); found > exp {
		t.Fatalf("Found too many items. Expected %d <= %d", found, exp)
	}
}

func TestFilterFalseNegative(t *testing.T) {
	f := NewFilter(fnv.New64(), 20) // really hard to get false positive

	const items = 10000
	for i := 0; i < items; i++ {
		p := randSeed()
		f.Add(p)
		if !f.Lookup(p) {
			t.Fatalf("Failed to find %v on iteration %d", p, i)
		}
	}
}

func BenchmarkFnvSFilterAdd5(b *testing.B)  { benchmarkFnvSFilterAdd(b, 5) }
func BenchmarkFnvSFilterAdd10(b *testing.B) { benchmarkFnvSFilterAdd(b, 10) }
func BenchmarkFnvSFilterAdd15(b *testing.B) { benchmarkFnvSFilterAdd(b, 15) }
func BenchmarkFnvSFilterAdd20(b *testing.B) { benchmarkFnvSFilterAdd(b, 20) }

func benchmarkFnvSFilterAdd(b *testing.B, k int) {
	f := NewFilter(fnv.New64(), k)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Add(randSeed())
		b.SetBytes(10)
	}
}
