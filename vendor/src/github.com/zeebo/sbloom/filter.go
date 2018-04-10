package sbloom

import "fmt"

type filter struct {
	Log  uint // size == 1 << log
	K    int
	Bins [][]uint8 // [k][1<<size]uint8 bins
	Left uint64    // number of additions left until new filter
}

//newFilter returns a new bloom filter of the specified bitsize and given mask
//and number of bins.
func newFilter(log uint, k int) (f *filter) {
	if log < elemSize {
		panic(fmt.Sprintf("log must be at least %d", elemSize))
	}

	size := uint64(1 << log)
	bins := make([][]uint8, k)  //allocate one more bin
	binSize := size >> elemSize //size of each bin to have size bits
	for i := range bins {
		bins[i] = make([]uint8, binSize)
	}

	return &filter{
		Log:  log,
		Bins: bins,
		K:    k,
		Left: uint64(k) * (1 << (log - 1)),
	}
}

const (
	elemSize = 3               // 1 << 3 bits in uint8
	idxShift = elemSize        // divide by uint8
	idxMask  = 1<<elemSize - 1 // modulus uint8
)

//set turns the nth bit in the slice of bytes to one.
func set(m []uint8, n uint64) {
	idx, mask := n>>idxShift, n&idxMask
	m[idx] |= 1 << mask
}

//get returns true if the nth bit in the slice of bytes is one.
func get(m []uint8, n uint64) bool {
	idx, mask := n>>idxShift, n&idxMask
	return m[idx]&(1<<mask) != 0
}

func (f *filter) Add(p []byte, hs []sHash) {
	for i, h := range hs {
		val := mix(h.Hash(p), f.Log)
		if f.Left > 0 && !get(f.Bins[i], val) {
			f.Left--
		}
		set(f.Bins[i], val)
	}
}

func (f *filter) Lookup(p []byte, hs []sHash) bool {
	for i, h := range hs {
		val := mix(h.Hash(p), f.Log)
		if !get(f.Bins[i], val) {
			return false
		}
	}
	return true
}

func mix(val uint64, log uint) (tmp uint64) {
	var mask uint64 = 1<<log - 1
	for val > 0 {
		tmp ^= val & mask
		val >>= log
	}
	return
}
