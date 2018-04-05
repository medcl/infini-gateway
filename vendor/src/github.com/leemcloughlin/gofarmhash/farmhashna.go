package farmhash

// Based on the original C++ farmhashna.cc

func shiftMix(val uint64) uint64 {
	return val ^ (val >> 47)
}

func hashLen16(u, v uint64) uint64 {
	return hash128to64(Uint128{u, v})
}

// Note: the C++ original was overloaded hashLen16()
func hashLen16Mul(u, v, mul uint64) uint64 {
	// Murmur-inspired hashing.
	a := (u ^ v) * mul
	a ^= (a >> 47)
	b := (v ^ a) * mul
	b ^= (b >> 47)
	b *= mul
	return b
}

func hashLen0to16(s []byte) uint64 {
	len := uint64(len(s))
	if len >= 8 {
		mul := k2 + len*2
		a := fetch64(s) + k2
		b := fetch64(s[len-8:])
		c := rotate64(b, 37)*mul + a
		d := (rotate64(a, 25) + b) * mul
		return hashLen16Mul(c, d, mul)
	}
	if len >= 4 {
		mul := k2 + len*2
		a := fetch32(s)
		return hashLen16Mul(len+(uint64(a)<<3), uint64(fetch32(s[len-4:])), mul)
	}
	if len > 0 {
		a := s[0]
		b := s[len>>1]
		c := s[len-1]
		y := uint32(a) + (uint32(b) << 8)
		var z uint64 = len + (uint64(c) << 2)
		return shiftMix(uint64(y)*k2^z*k0) * k2
	}
	return k2
}

// This probably works well for 16-byte strings as well, but it may be overkill
// in that case.
func hashLen17to32(s []byte) uint64 {
	len := uint64(len(s))
	mul := k2 + len*2
	a := fetch64(s) * k1
	b := fetch64(s[8:])
	c := fetch64(s[len-8:]) * mul
	d := fetch64(s[len-16:]) * k2
	return hashLen16Mul(rotate64(a+b, 43)+rotate64(c, 30)+d,
		a+rotate64(b+k2, 18)+c, mul)
}

/*
Use the identical versions of these two from farmhashcc.go

// Return a 16-byte hash for 48 bytes.  Quick and dirty.
// Callers do best to use "random-looking" values for a and b.
// Note: C++ returned pair<uint64_t, uint64_t>
func weakHashLen32WithSeeds(w, x, y, z, a, b uint64) uint128 {
	a += w
	b = rotate64(b+a+z, 21)
	c := a
	a += x
	a += y
	b += rotate64(a, 44)
	return uint128{a + z, b + c}
}

// Return a 16-byte hash for s[0] ... s[31], a, and b.  Quick and dirty.
// Note: original C++ returned pair<uint64_t, uint64_t>
func weakHashLen32WithSeedsBytes(s []byte, a, b uint64) uint128 {
	return WeakHashLen32WithSeeds(fetch64(s),
		fetch64(s[8:]),
		fetch64(s[16:]),
		fetch64(s[24:]),
		a,
		b)
}
*/

// Return an 8-byte hash for 33 to 64 bytes.
func hashLen33to64(s []byte) uint64 {
	len := uint64(len(s))
	mul := k2 + len*2
	a := fetch64(s) * k2
	b := fetch64(s[8:])
	c := fetch64(s[len-8:]) * mul
	d := fetch64(s[len-16:]) * k2
	y := rotate64(a+b, 43) + rotate64(c, 30) + d
	z := hashLen16Mul(y, a+rotate64(b+k2, 18)+c, mul)
	e := fetch64(s[16:]) * mul
	f := fetch64(s[24:])
	g := (y + fetch64(s[len-32:])) * mul
	h := (z + fetch64(s[len-24:])) * mul
	return hashLen16Mul(rotate64(e+f, 43)+rotate64(g, 30)+h,
		e+rotate64(f+a, 18)+g, mul)
}

// renamed from hash64 to make it clearer elsewhere which is being called
func naHash64(s []byte) uint64 {
	len := uint64(len(s))
	const seed uint64 = 81
	if len <= 32 {
		if len <= 16 {
			return hashLen0to16(s)
		} else {
			return hashLen17to32(s)
		}
	} else if len <= 64 {
		return hashLen33to64(s)
	}

	// For strings over 64 bytes we loop.  Internal state consists of
	// 56 bytes: v, w, x, y, and z.
	x := seed
	// This overflows a qint64 which causes a Go compiler error so split it up
	// y := (seed*k1) + 113
	y := seed
	y *= k1
	y += 113
	z := shiftMix(y*k2+113) * k2
	// v and w were pair<uint64_t, uint64_t>
	var v, w Uint128
	x = x*k2 + fetch64(s)

	// Set end so that after the loop we have 1 to 64 bytes left to process.
	var end uint64 = ((len - 1) / 64) * 64
	var last64 uint64 = end + ((len - 1) & 63) - 63
	var i uint64
	origS := s[:]
	for {
		x = rotate64(x+y+v.First+fetch64(s[8:]), 37) * k1
		y = rotate64(y+v.Second+fetch64(s[48:]), 42) * k1
		x ^= w.Second
		y += v.First + fetch64(s[40:])
		z = rotate64(z+w.First, 33) * k1
		v = weakHashLen32WithSeedsBytes(s, v.Second*k1, x+w.First)
		w = weakHashLen32WithSeedsBytes(s[32:], z+w.Second, y+fetch64(s[16:]))
		z, x = x, z
		s = s[64:]
		i += 64
		if i == end {
			break
		}
	}
	mul := k1 + ((z & 0xff) << 1)
	// Make s point to the last 64 bytes of input.
	s = origS[last64:]
	w.First += ((len - 1) & 63)
	v.First += w.First
	w.First += v.First
	x = rotate64(x+y+v.First+fetch64(s[8:]), 37) * mul
	y = rotate64(y+v.Second+fetch64(s[48:]), 42) * mul
	x ^= w.Second * 9
	y += v.First*9 + fetch64(s[40:])
	z = rotate64(z+w.First, 33) * mul
	v = weakHashLen32WithSeedsBytes(s, v.Second*mul, x+w.First)
	w = weakHashLen32WithSeedsBytes(s[32:], z+w.Second, y+fetch64(s[16:]))
	z, x = x, z
	return hashLen16Mul(hashLen16Mul(v.First, w.First, mul)+shiftMix(y)*k0+z,
		hashLen16Mul(v.Second, w.Second, mul)+x,
		mul)
}

func naHash64WithSeed(s []byte, seed uint64) uint64 {
	return naHash64WithSeeds(s, k2, seed)
}

func naHash64WithSeeds(s []byte, seed0, seed1 uint64) uint64 {
	return hashLen16(naHash64(s)-seed0, seed1)
}
