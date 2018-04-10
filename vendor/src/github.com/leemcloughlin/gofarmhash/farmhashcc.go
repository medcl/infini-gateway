package farmhash

// Based on the original C++ farmhashcc.cc

// This file provides a 32-bit hash equivalent to CityHash32 (v1.1.1)
// and a 128-bit hash equivalent to CityHash128 (v1.1.1).  It also provides
// a seeded 32-bit hash function similar to CityHash32.

func hash32Len13to24(s []byte) uint32 {
	len := uint64(len(s))
	a := fetch32(s[-4+int64(len>>1):])
	b := fetch32(s[4:])
	c := fetch32(s[len-8:])
	d := fetch32(s[len>>1:])
	e := fetch32(s[0:])
	f := fetch32(s[len-4:])
	h := uint32(len)
	return fmix(mur(f, mur(e, mur(d, mur(c, mur(b, mur(a, h)))))))
	return 0
}

func hash32Len0to4(s []byte) uint32 {
	len := uint64(len(s))
	var b uint32 = 0
	var c uint32 = 9
	var i uint64
	for i = 0; i < len; i++ {
		v := int8(s[i])
		b = b*c1 + uint32(v)
		c ^= b
	}
	return fmix(mur(b, mur(uint32(len), c)))
}

func hash32Len5to12(s []byte) uint32 {
	len := uint64(len(s))
	a := uint32(len)
	b := uint32(len) * 5
	var c uint32 = 9
	var d uint32 = b
	a += fetch32(s[0:])
	b += fetch32(s[len-4:])
	c += fetch32(s[((len >> 1) & 4):])
	return fmix(mur(c, mur(b, mur(a, d))))
}

func hash32(s []byte) uint32 {
	len := uint64(len(s))
	if len <= 24 {
		if len <= 12 {
			if len <= 4 {
				return hash32Len0to4(s)
			}
			return hash32Len5to12(s)
		}
		return hash32Len13to24(s)
	}

	// len > 24
	h := uint32(len)
	g := c1 * uint32(len)
	f := g
	a0 := rotate32(fetch32(s[len-4:])*c1, 17) * c2
	a1 := rotate32(fetch32(s[len-8:])*c1, 17) * c2
	a2 := rotate32(fetch32(s[len-16:])*c1, 17) * c2
	a3 := rotate32(fetch32(s[len-12:])*c1, 17) * c2
	a4 := rotate32(fetch32(s[len-20:])*c1, 17) * c2
	h ^= a0
	h = rotate32(h, 19)
	h = h*5 + 0xe6546b64
	h ^= a2
	h = rotate32(h, 19)
	h = h*5 + 0xe6546b64
	g ^= a1
	g = rotate32(g, 19)
	g = g*5 + 0xe6546b64
	g ^= a3
	g = rotate32(g, 19)
	g = g*5 + 0xe6546b64
	f += a4
	f = rotate32(f, 19)
	f = f*5 + 0xe6546b64

	var iters uint64 = (len - 1) / 20
	for {
		a0 := rotate32(fetch32(s[:])*c1, 17) * c2
		a1 := fetch32(s[4:])
		a2 := rotate32(fetch32(s[8:])*c1, 17) * c2
		a3 := rotate32(fetch32(s[12:])*c1, 17) * c2
		a4 := fetch32(s[16:])
		h ^= a0
		h = rotate32(h, 18)
		h = h*5 + 0xe6546b64
		f += a1
		f = rotate32(f, 19)
		f = f * c1
		g += a2
		g = rotate32(g, 18)
		g = g*5 + 0xe6546b64
		h ^= a3 + a1
		h = rotate32(h, 19)
		h = h*5 + 0xe6546b64
		g ^= a4
		g = bswap32(g) * 5
		h += a4 * 5
		h = bswap32(h)
		f += a0
		//PERMUTE3(f, h, g) - swap(a,b);swap(b,c)
		f, h = h, f
		f, g = g, f
		s = s[20:]
		if iters--; iters == 0 {
			break
		}
	}
	g = rotate32(g, 11) * c1
	g = rotate32(g, 17) * c1
	f = rotate32(f, 11) * c1
	f = rotate32(f, 17) * c1
	h = rotate32(h+g, 19)
	h = h*5 + 0xe6546b64
	h = rotate32(h, 17) * c1
	h = rotate32(h+f, 19)
	h = h*5 + 0xe6546b64
	h = rotate32(h, 17) * c1
	return h
}

func hash32WithSeed(s []byte, seed uint32) uint32 {
	len := uint64(len(s))
	if len <= 24 {
		if len >= 13 {
			return mkHash32Len13to24(s, seed*c1)
		} else if len >= 5 {
			return mkHash32Len5to12(s, seed)
		} else {
			return mkHash32Len0to4(s, seed)
		}
	}
	h := mkHash32Len13to24(s[0:24], seed^uint32(len))
	return mur(hash32(s[24:])+seed, h)
}

// Note: use identical ShiftMix from farmhashna.go

func hashLen16NoMul(u, v uint64) uint64 {
	return hash128to64(Uint128{u, v})
}

// Note: avoid clashing with the same names in farmhashna by adding cc to start

func ccHashLen16(u, v, mul uint64) uint64 {
	// Murmur-inspired hashing.
	a := (u ^ v) * mul
	a ^= (a >> 47)
	b := (v ^ a) * mul
	b ^= (b >> 47)
	b *= mul
	return b
}

func ccHashLen0to16(s []byte) uint64 {
	len := uint64(len(s))
	if len >= 8 {
		mul := k2 + len*2
		a := fetch64(s) + k2
		b := fetch64(s[len-8:])
		c := rotate64(b, 37)*mul + a
		d := (rotate64(a, 25) + b) * mul
		return ccHashLen16(c, d, mul)
	}
	if len >= 4 {
		mul := k2 + len*2
		a := uint64(fetch32(s))
		return ccHashLen16(len+(a<<3), uint64(fetch32(s[len-4:])), mul)
	}
	if len > 0 {
		a := s[0]
		b := s[len>>1]
		c := s[len-1]
		y := uint32(a) + (uint32(b) << 8)
		z := uint32(len) + (uint32(c) << 2)
		return shiftMix(uint64(y)*k2^uint64(z)*k0) * uint64(k2)
	}
	return k2
}

// Return a 16-byte hash for 48 bytes.  Quick and dirty.
// Callers do best to use "random-looking" values for a and b.
// Note: original C++ returned a pair<uint64_t, uint64_t>
func weakHashLen32WithSeeds(w, x, y, z, a, b uint64) Uint128 {
	a += w
	b = rotate64(b+a+z, 21)
	c := a
	a += x
	a += y
	b += rotate64(a, 44)
	return Uint128{a + z, b + c}
}

// Return a 16-byte hash for s[0] ... s[31], a, and b.  Quick and dirty.
// Note: original C++ returned a pair<uint64_t, uint64_t>
func weakHashLen32WithSeedsBytes(s []byte, a, b uint64) Uint128 {
	return weakHashLen32WithSeeds(fetch64(s),
		fetch64(s[8:]),
		fetch64(s[16:]),
		fetch64(s[24:]),
		a,
		b)
}

// A subroutine for CityHash128().  Returns a decent 128-bit hash for strings
// of any length representable in signed long.  Based on City and Murmur.
func CityMurmur(s []byte, seed Uint128) Uint128 {
	len := uint64(len(s))
	a := seed.First
	b := seed.Second
	var c uint64 = 0
	var d uint64 = 0
	l := int64(len - 16)
	if l <= 0 { // len <= 16
		a = shiftMix(a*k1) * k1
		c = b*k1 + ccHashLen0to16(s)
		if len >= 8 {
			d = shiftMix(a + fetch64(s))
		} else {
			d = shiftMix(a + c)
		}
	} else { // len > 16
		c = hashLen16NoMul(fetch64(s[(len-8):])+k1, a)
		d = hashLen16NoMul(b+len, c+fetch64(s[len-16:]))
		a += d
		for {
			a ^= shiftMix(fetch64(s)*k1) * k1
			a *= k1
			b ^= a
			c ^= shiftMix(fetch64(s[8:])*k1) * k1
			c *= k1
			d ^= c
			s = s[16:]
			l -= 16
			if l <= 0 {
				break
			}
		}
	}
	a = hashLen16NoMul(a, c)
	b = hashLen16NoMul(d, b)
	return Uint128{a ^ b, hashLen16NoMul(b, a)}
}

func CityHash128WithSeed(s []byte, seed Uint128) Uint128 {
	len := uint64(len(s))
	if len < 128 {
		return CityMurmur(s, seed)
	}

	// We expect len >= 128 to be the common case.  Keep 56 bytes of state:
	// v, w, x, y, and z.
	// Note: v and w were pair<uint64_t, uint64_t> in the C++
	var v, w Uint128
	x := seed.First
	y := seed.Second
	z := len * k1
	v.First = rotate64(y^k1, 49)*k1 + fetch64(s)
	v.Second = rotate64(v.First, 42)*k1 + fetch64(s[8:])
	w.First = rotate64(y+z, 35)*k1 + x
	w.Second = rotate64(x+fetch64(s[88:]), 53) * k1

	// This is the same inner loop as CityHash64(), manually unrolled.
	for {
		x = rotate64(x+y+v.First+fetch64(s[8:]), 37) * k1
		y = rotate64(y+v.Second+fetch64(s[48:]), 42) * k1
		x ^= w.Second
		y += v.First + fetch64(s[40:])
		z = rotate64(z+w.First, 33) * k1
		v = weakHashLen32WithSeedsBytes(s, v.Second*k1, x+w.First)
		w = weakHashLen32WithSeedsBytes(s[32:], z+w.Second, y+fetch64(s[16:]))
		z, x = x, y
		s = s[64:]
		x = rotate64(x+y+v.First+fetch64(s[8:]), 37) * k1
		y = rotate64(y+v.Second+fetch64(s[48:]), 42) * k1
		x ^= w.Second
		y += v.First + fetch64(s[40:])
		z = rotate64(z+w.First, 33) * k1
		v = weakHashLen32WithSeedsBytes(s, v.Second*k1, x+w.First)
		w = weakHashLen32WithSeedsBytes(s[32:], z+w.Second, y+fetch64(s[16:]))
		z, x = x, y
		s = s[64:]
		len -= 128
		if len < 128 {
			break
		}
	}
	x += rotate64(v.First+z, 49) * k0
	y = y*k0 + rotate64(w.Second, 37)
	z = z*k0 + rotate64(w.First, 27)
	w.First *= 9
	v.First *= k0
	// If 0 < len < 128, hash up to 4 chunks of 32 bytes each from the end of s.
	var tail_done uint64
	for tail_done = 0; tail_done < len; {
		tail_done += 32
		y = rotate64(x+y, 42)*k0 + v.Second
		w.First += fetch64(s[len-tail_done+16:])
		x = x*k0 + w.First
		z += w.Second + fetch64(s[len-tail_done:])
		w.Second += v.First
		v = weakHashLen32WithSeedsBytes(s[len-tail_done:], v.First+z, v.Second)
		v.First *= k0
	}
	// At this point our 56 bytes of state should contain more than
	// enough information for a strong 128-bit hash.  We use two
	// different 56-byte-to-8-byte hashes to get a 16-byte final result.
	x = hashLen16NoMul(x, v.First)
	y = hashLen16NoMul(y+z, w.First)
	return Uint128{hashLen16NoMul(x+v.Second, w.Second) + y,
		hashLen16NoMul(x+w.Second, y+v.Second)}
}

func CityHash128(s []byte) Uint128 {
	len := uint64(len(s))
	if len >= 16 {
		return CityHash128WithSeed(s[16:],
			Uint128{fetch64(s), fetch64(s[8:]) + k0})

	}
	return CityHash128WithSeed(s, Uint128{k0, k1})
}

func Fingerprint128(s []byte) Uint128 {
	return CityHash128(s)
}
