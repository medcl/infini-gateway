package sbloom

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"testing"
)

func newTestFilter(k int, log uint) (f *filter, hashes []sHash) {
	bh := fnv.New64()
	for i := 0; i < k; i++ {
		hashes = append(hashes, newsHash(bh))
	}

	f = newFilter(log, k)
	return
}

func BenchmarkFilterAddK10L5(b *testing.B)  { benchmarkFilterAdd(b, 10, 5) }
func BenchmarkFilterAddK10L10(b *testing.B) { benchmarkFilterAdd(b, 10, 10) }
func BenchmarkFilterAddK10L15(b *testing.B) { benchmarkFilterAdd(b, 10, 15) }

func BenchmarkFilterAddK20L5(b *testing.B)  { benchmarkFilterAdd(b, 20, 5) }
func BenchmarkFilterAddK20L10(b *testing.B) { benchmarkFilterAdd(b, 20, 10) }
func BenchmarkFilterAddK20L15(b *testing.B) { benchmarkFilterAdd(b, 20, 15) }

func BenchmarkFilterAddK30L5(b *testing.B)  { benchmarkFilterAdd(b, 30, 5) }
func BenchmarkFilterAddK30L10(b *testing.B) { benchmarkFilterAdd(b, 30, 10) }
func BenchmarkFilterAddK30L15(b *testing.B) { benchmarkFilterAdd(b, 30, 15) }

func benchmarkFilterAdd(b *testing.B, k int, log uint) {
	f, hashes := newTestFilter(k, log)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f.Add(randSeed(), hashes)
		b.SetBytes(10)
	}
}

func TestNoFalseNegatives(t *testing.T) {
	f, hashes := newTestFilter(8, 10)

	if p := randSeed(); f.Lookup(p, hashes) {
		t.Fatalf("false positive on zero: %v", p)
	}

	//while we're under half full
	for f.Left > 0 {
		item := randSeed()
		f.Add(item, hashes)
		if !f.Lookup(item, hashes) {
			t.Fatalf("false negative: %v", item)
		}
	}
}

func TestOneBitSet(t *testing.T) {
	f, hashes := newTestFilter(8, 5)
	f.Add(randSeed(), hashes)

	//check each bin
	for i, bin := range f.Bins {

		var ones int
		for _, b := range bin {
			for i := uint(0); i < 1<<elemSize; i++ {
				if (1<<i)&b > 0 {
					ones++
				}
			}
		}

		if ones > 1 {
			for i, bin := range f.Bins {
				t.Logf("%d: %s", i, dumpBin(bin))
			}
			t.Fatalf("got %d bits set in bin %d", ones, i)
		}

	}
}

func TestFalsePositiveRate(t *testing.T) {
	for k := 1; k <= 16; k++ {
		testFalsePostive(t, k)
	}
}

func testFalsePostive(t *testing.T, k int) {
	f, hashes := newTestFilter(k, 5)

	//fill it halfway up
	var items [][]byte
	for f.Left > 0 {
		item := randSeed()
		f.Add(item, hashes)
		items = append(items, item)
	}

	in := func(p []byte) bool {
		for _, c := range items {
			if bytes.Equal(p, c) {
				return true
			}
		}
		return false
	}

	iters := 1 << uint(k+1)
	var found int
	for i := 0; i < iters; i++ {
		var item []byte
		for item == nil || in(item) {
			item = randSeed()
		}
		if f.Lookup(item, hashes) {
			found++
		}
	}

	if found > 5 {
		t.Errorf("[%d] %d out of %d false positives", k, found, iters)
		for i, bin := range f.Bins {
			t.Logf("[%d] %d: %s", k, i, dumpBin(bin))
		}
	}
}

func TestFilterLookupSmall(t *testing.T) {
	f, hashes := newTestFilter(8, 3)
	p := randSeed()
	f.Add(p, hashes)

	if !f.Lookup(p, hashes) {
		t.Fatal("False negative")
	}
}

func dumpBin(buf []uint8) string {
	fmtr := fmt.Sprintf("%%0%db", 1<<elemSize)
	sbuf := new(bytes.Buffer)

	var ones int
	for _, b := range buf {
		fmt.Fprintf(sbuf, fmtr, b)
		for i := uint(0); i < 1<<elemSize; i++ {
			if (1<<i)&b > 0 {
				ones++
			}
		}
	}

	ratio := float64(ones) / float64(len(buf)*1<<elemSize)
	fmt.Fprintf(sbuf, " [%d %.2f]", ones, ratio)
	return sbuf.String()
}
