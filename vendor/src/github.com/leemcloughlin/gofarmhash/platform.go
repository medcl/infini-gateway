package farmhash

// Based on the original C++ platform.cc

// A small optimisation I switched from binary.LittleEndian.Uint32/64 to
// hand coded. The benchmark, on my system, went from
// around 64ns/op to around 51ns/op
/*
import (
	"encoding/binary"
)
*/

// Note: I used to call binary.LittleEndian.Uint32 and Uint64 inline but it
// made comparing the code to the original much trickier

func fetch32(p []byte) uint32 {
	//return binary.LittleEndian.Uint32(p)
	return uint32(p[0]) | uint32(p[1])<<8 | uint32(p[2])<<16 | uint32(p[3])<<24
}

func fetch64(p []byte) uint64 {
	//return binary.LittleEndian.Uint64(p)
	return uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 |
		uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56
}

// rotate32 is a bitwise rotate
func rotate32(val uint32, shift uint) uint32 {
	if shift == 0 {
		return val
	}
	return val>>shift | val<<(32-shift)
}

// rotate64 is a bitwise rotate
func rotate64(val uint64, shift uint) uint64 {
	if shift == 0 {
		return val
	}
	return val>>shift | val<<(64-shift)
}
