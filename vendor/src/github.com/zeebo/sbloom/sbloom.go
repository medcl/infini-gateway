package sbloom

import "hash"

//Filter represents a scalable bloom filter.
type Filter struct {
	bh hash.Hash64
	fs []*filter
	hs []sHash
}

//NewFilter returns a scalable bloom filter with a false positive probaility
//less than 1/2**k. It starts with a default size of 1024 bits per k. If you
//know you will use more/less than that, use NewSizedFilter for a better hint.
//It uses the provided hash to operate.
func NewFilter(h hash.Hash64, k int) *Filter {
	return NewSizedFilter(h, k, 10)
}

//NewSizedFilter returns a scalable bloom filter with a false positive
//probability less than 1/2**k. It start with a default size of 2**log bits per
//k. It uses the provided hash to operate.
func NewSizedFilter(h hash.Hash64, k int, log uint) *Filter {
	f := new(Filter)
	f.bh = h
	f.addNewFilter(log, k+1) //need to start at 1/2 to converge up
	return f
}

//Add adds the slice of bytes to the filter so that it will always return true
//when looked up.
func (f *Filter) Add(p []byte) {
	last := f.fs[len(f.fs)-1]
	last.Add(p, f.hs[:last.K])
	if last.Left == 0 {
		f.addNewFilter(last.Log+1, last.K+1)
	}
}

//addNewFilter creates a new filter of the given size and number of bins, making
//sure we have enough hash functions for the bins.
func (f *Filter) addNewFilter(log uint, k int) {
	//make sure we have up to k hashes
	for len(f.hs) < k {
		f.hs = append(f.hs, newsHash(f.bh))
	}

	//add the new bloom filter
	f.fs = append(f.fs, newFilter(log, k))
}

//Lookup looks for the set of bytes in the bloom filter. If it returns false,
//the set of bytes has definitely not been seen before.
func (f *Filter) Lookup(p []byte) bool {
	for _, subfil := range f.fs {
		if subfil.Lookup(p, f.hs[:subfil.K]) {
			return true
		}
	}
	return false
}
