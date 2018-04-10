package farmhash

// Based on the original C++ basics.cc

const (
	// Some primes between 2^63 and 2^64 for various uses.
	k0 uint64 = 0xc3a5c85c97cb3127
	k1 uint64 = 0xb492b66fbe98f273
	k2 uint64 = 0x9ae16a3b2f90404f

	// Magic numbers for 32-bit hashing.  Copied from Murmur3.
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
)

// fmix is a 32-bit to 32-bit integer hash copied from Murmur3.
func fmix(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

// Mur is a helper from Murmur3 for combining two 32-bit values.
func mur(a, h uint32) uint32 {
	a *= c1
	a = rotate32(a, 17)
	a *= c2
	h ^= a
	h = rotate32(h, 19)
	return h*5 + 0xe6546b64
}

// Couldn't find a version of this in Go
func bswap32(x uint32) uint32 {
	return ((x >> 24) & 0xFF) | ((x >> 8) & 0xFF00) | ((x << 8) & 0xFF0000) | ((x << 24) & 0xFF000000)
}

type Uint128 struct {
	First, Second uint64
}

// hash128to64 is from farmhash.h it is intended to be a reasonably good hash function.
func hash128to64(x Uint128) uint64 {
	// Murmur-inspired hashing.
	const kMul uint64 = 0x9ddfea08eb382d69
	a := (x.First ^ x.Second) * kMul
	a ^= (a >> 47)
	b := (x.Second ^ a) * kMul
	b ^= (b >> 47)
	b *= kMul
	return b
}
