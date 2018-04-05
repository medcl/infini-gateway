package farmhash

// Based on the original C++ farmhashmk.cc

// Note: These functions clashed with the versions in farmhashcc
// Some only differ by taking a seed but others are quite different
// To avoid clashes I've added mk to the start of these names

func mkHash32Len13to24(s []byte, seed uint32) uint32 {
	len := uint64(len(s))
	a := fetch32(s[(len>>1)-4:])
	b := fetch32(s[4:])
	c := fetch32(s[len-8:])
	d := fetch32(s[len>>1:])
	e := fetch32(s)
	f := fetch32(s[len-4:])
	h := d*c1 + uint32(len) + seed
	a = rotate32(a, 12) + f
	h = mur(c, h) + a
	a = rotate32(a, 3) + c
	h = mur(e, h) + a
	a = rotate32(a+f, 12) + d
	h = mur(b^seed, h) + a
	return fmix(h)
}

func mkHash32Len0to4(s []byte, seed uint32) uint32 {
	len := uint64(len(s))
	b := seed
	var c uint32 = 9
	var i uint64
	for i = 0; i < len; i++ {
		v := int8(s[i])
		b = b*c1 + uint32(v)
		c ^= b
	}
	return fmix(mur(b, mur(uint32(len), c)))
}

func mkHash32Len5to12(s []byte, seed uint32) uint32 {
	len := uint64(len(s))
	var a, b, c, d uint32
	a = uint32(len)
	b = uint32(len) * 5
	c = 9
	d = b + seed
	a += fetch32(s)
	b += fetch32(s[len-4:])
	c += fetch32(s[(len>>1)&4:])
	return fmix(seed ^ mur(c, mur(b, mur(a, d))))
}

func mkHash32(s []byte) uint32 {
	len := uint64(len(s))
	if len <= 24 {
		if len <= 12 {
			if len <= 4 {
				return mkHash32Len0to4(s, 0)
			}
			return mkHash32Len5to12(s, 0)
		}
		return mkHash32Len13to24(s, 0)
	}

	// len > 24
	var h, g, f uint32
	h = uint32(len)
	g = c1 * uint32(len)
	f = g
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
	f = rotate32(f, 19) + 113
	var iters uint64
	iters = (len - 1) / 20
	for {
		a := fetch32(s)
		b := fetch32(s[4:])
		c := fetch32(s[8:])
		d := fetch32(s[12:])
		e := fetch32(s[16:])
		h += a
		g += b
		f += c
		h = mur(d, h) + e
		g = mur(c, g) + a
		f = mur(b+e*c1, f) + d
		f += g
		g += f
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

func mkHash32WithSeed(s []byte, seed uint32) uint32 {
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
