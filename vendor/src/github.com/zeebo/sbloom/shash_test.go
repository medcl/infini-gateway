package sbloom

import (
	"hash/fnv"
	"testing"
)

func TestFnv64SeededIndependence(t *testing.T) {
	//two seeds thought to demonstrate collisions
	s1 := []byte{207, 248, 106, 26, 230, 204, 133, 85, 32, 173}
	s2 := []byte{88, 43, 139, 185, 188, 106, 169, 227, 193, 190}

	h1 := sHash{Ha: fnv.New64(), Seed: s1}
	h2 := sHash{Ha: fnv.New64(), Seed: s2}

	const log = 5

	//generate some values and record the collisions
	const iters = 10000
	var coll int
	for i := 0; i < iters; i++ {
		if val := randSeed(); mix(h1.Hash(val), log) == mix(h2.Hash(val), log) {
			coll++
		}
	}

	//clearly a problem
	if coll > iters/2 {
		t.Fatal("Too many collisions:", coll)
	}
}
