package sbloom

import (
	"hash/crc64"
	"hash/fnv"
	"math/rand"
	"testing"
)

func randIndex(top uint64) uint64 {
	hi, low := rand.Uint32(), rand.Uint32()
	return (uint64(hi)<<32 | uint64(low)) % top
}

const maxSize = 10

var sizes []uint64

func init() {
	for i := uint(0); i < maxSize; i++ {
		sizes = append(sizes, 1<<i)
	}
}

func BenchmarkSet(b *testing.B) {
	const size = 1 << maxSize
	x := make([]uint8, size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := uint64(0); i < size*(1<<elemSize); i++ {
			set(x, i)
		}
		b.SetBytes(size * (1 << (elemSize - 3)))
	}
}

func BenchmarkGet(b *testing.B) {
	const size = 1 << maxSize
	x := make([]uint8, size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := uint64(0); i < size*(1<<elemSize); i++ {
			get(x, i)
		}
		b.SetBytes(size * (1 << (elemSize - 3)))
	}
}

func BenchmarkFnvHash(b *testing.B) {
	const bytes = 1024
	s := newsHash(fnv.New64())
	dat := make([]byte, bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Hash(dat)
		b.SetBytes(bytes)
	}
}

func BenchmarkCrc64Hash(b *testing.B) {
	const bytes = 1024
	s := newsHash(crc64.New(crc64.MakeTable(crc64.ISO)))
	dat := make([]byte, bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Hash(dat)
		b.SetBytes(bytes)
	}
}

func TestGetAndSetOnEveryBit(t *testing.T) {
	for _, size := range sizes {
		x := make([]uint8, size)
		for j := uint64(0); j < size*(1<<elemSize); j++ {
			if get(x, j) {
				t.Errorf("0[%d] %08b", j, x[0])
			}
		}
		for i := uint64(0); i < size*(1<<elemSize); i++ {
			set(x, i)

			//make sure every bit is right
			for j := uint64(0); j < size*(1<<elemSize); j++ {
				if get(x, j) != (j <= i) { //get is true when j <= i
					t.Errorf("%d[%d] %08b", i, j, x[0])
				}
			}
		}
	}
}

func TestPanicGetLarge(t *testing.T) {
	recov := func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}

	for _, size := range sizes {
		func() {
			defer recov()
			x := make([]uint8, size)
			get(x, size*(1<<elemSize))
		}()

		func() {
			defer recov()
			x := make([]uint8, size)
			set(x, size*(1<<elemSize))
		}()
	}
}
