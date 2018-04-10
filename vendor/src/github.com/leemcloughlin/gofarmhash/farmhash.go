package farmhash

// farmhash.go is the Public interface to farmhash.
// There are some additional functions in farmhousecc.go

const (
	Version = "1.0.0"
)

// These functions based on original C++ namespace util

func Hash32(s []byte) uint32 {
	return mkHash32(s)
}

func Hash32WithSeed(s []byte, seed uint32) uint32 {
	return mkHash32WithSeed(s, seed)
}

func Hash64(s []byte) uint64 {
	return naHash64(s)
}

func Hash64WithSeed(s []byte, seed uint64) uint64 {
	return naHash64WithSeed(s, seed)
}

func Hash64WithSeeds(s []byte, seed0, seed1 uint64) uint64 {
	return naHash64WithSeeds(s, seed0, seed1)
}

func Hash128(s []byte) Uint128 {
	return Fingerprint128(s)
}

func Hash128WithSeed(s []byte, seed Uint128) Uint128 {
	return CityHash128WithSeed(s, seed)
}

func FingerPrint32(s []byte) uint32 {
	return mkHash32(s)
}

func FingerPrint64(s []byte) uint64 {
	return naHash64(s)
}
