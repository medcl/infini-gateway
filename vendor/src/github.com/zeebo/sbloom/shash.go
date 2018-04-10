package sbloom

import (
	"hash"
	"math/rand"
)

//sHash is a hash with an initial seed.
type sHash struct {
	Ha   hash.Hash64
	Seed []byte
}

func newsHash(ha hash.Hash64) (s sHash) {
	s.Ha = ha
	s.Seed = randSeed()
	return
}

func randSeed() (p []byte) {
	for i := 0; i < 10; i++ {
		p = append(p, byte(rand.Intn(256)))
	}
	return
}

func (s sHash) Hash(p []byte) uint64 {
	s.Ha.Reset()
	s.Ha.Write(s.Seed)
	s.Ha.Write(p)
	return s.Ha.Sum64()
}
