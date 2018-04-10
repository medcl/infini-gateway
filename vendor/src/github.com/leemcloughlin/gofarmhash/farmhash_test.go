package farmhash

// go test
// Test both the Public and internal func's

import (
	"crypto/sha1"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
	"testing"
)

type StrToHash32 struct {
	str      string
	expected uint32
}

type StrToHash64 struct {
	str      string
	expected uint64
}

type NumInOut struct {
	in  uint32
	out uint32
}

type NumIn2Out struct {
	in  uint32
	in2 uint32
	out uint32
}

func TestRotate32(t *testing.T) {
	rots := []NumIn2Out{
		{0, 0, 0},
		// Rotate32
		{100, 0, 100},
		{100, 1, 50},
		{100, 2, 25},
		{100, 3, 2147483660},
		{100, 4, 1073741830},
		{100, 5, 536870915},
		{100, 6, 2415919105},
		{100, 7, 3355443200},
		{100, 8, 1677721600},
		{100, 9, 838860800},
		{100, 10, 419430400},
		{100, 11, 209715200},
		{100, 12, 104857600},
		{100, 13, 52428800},
		{100, 14, 26214400},
		{100, 15, 13107200},
		{100, 16, 6553600},
		{100, 17, 3276800},
		{100, 18, 1638400},
		{100, 19, 819200},
		{100, 20, 409600},
		{100, 21, 204800},
		{100, 22, 102400},
		{100, 23, 51200},
		{100, 24, 25600},
		{100, 25, 12800},
		{100, 26, 6400},
		{100, 27, 3200},
		{100, 28, 1600},
		{100, 29, 800},
		{100, 30, 400},
		{100, 31, 200},
		{1000, 0, 1000},
		{1000, 1, 500},
		{1000, 2, 250},
		{1000, 3, 125},
		{1000, 4, 2147483710},
		{1000, 5, 1073741855},
		{1000, 6, 2684354575},
		{1000, 7, 3489660935},
		{1000, 8, 3892314115},
		{1000, 9, 4093640705},
		{1000, 10, 4194304000},
		{1000, 11, 2097152000},
		{1000, 12, 1048576000},
		{1000, 13, 524288000},
		{1000, 14, 262144000},
		{1000, 15, 131072000},
		{1000, 16, 65536000},
		{1000, 17, 32768000},
		{1000, 18, 16384000},
		{1000, 19, 8192000},
		{1000, 20, 4096000},
		{1000, 21, 2048000},
		{1000, 22, 1024000},
		{1000, 23, 512000},
		{1000, 24, 256000},
		{1000, 25, 128000},
		{1000, 26, 64000},
		{1000, 27, 32000},
		{1000, 28, 16000},
		{1000, 29, 8000},
		{1000, 30, 4000},
		{1000, 31, 2000},
		{100000, 0, 100000},
		{100000, 1, 50000},
		{100000, 2, 25000},
		{100000, 3, 12500},
		{100000, 4, 6250},
		{100000, 5, 3125},
		{100000, 6, 2147485210},
		{100000, 7, 1073742605},
		{100000, 8, 2684354950},
		{100000, 9, 1342177475},
		{100000, 10, 2818572385},
		{100000, 11, 3556769840},
		{100000, 12, 1778384920},
		{100000, 13, 889192460},
		{100000, 14, 444596230},
		{100000, 15, 222298115},
		{100000, 16, 2258632705},
		{100000, 17, 3276800000},
		{100000, 18, 1638400000},
		{100000, 19, 819200000},
		{100000, 20, 409600000},
		{100000, 21, 204800000},
		{100000, 22, 102400000},
		{100000, 23, 51200000},
		{100000, 24, 25600000},
		{100000, 25, 12800000},
		{100000, 26, 6400000},
		{100000, 27, 3200000},
		{100000, 28, 1600000},
		{100000, 29, 800000},
		{100000, 30, 400000},
		{100000, 31, 200000},
		{1000000, 0, 1000000},
		{1000000, 1, 500000},
		{1000000, 2, 250000},
		{1000000, 3, 125000},
		{1000000, 4, 62500},
		{1000000, 5, 31250},
		{1000000, 6, 15625},
		{1000000, 7, 2147491460},
		{1000000, 8, 1073745730},
		{1000000, 9, 536872865},
		{1000000, 10, 2415920080},
		{1000000, 11, 1207960040},
		{1000000, 12, 603980020},
		{1000000, 13, 301990010},
		{1000000, 14, 150995005},
		{1000000, 15, 2222981150},
		{1000000, 16, 1111490575},
		{1000000, 17, 2703228935},
		{1000000, 18, 3499098115},
		{1000000, 19, 3897032705},
		{1000000, 20, 4096000000},
		{1000000, 21, 2048000000},
		{1000000, 22, 1024000000},
		{1000000, 23, 512000000},
		{1000000, 24, 256000000},
		{1000000, 25, 128000000},
		{1000000, 26, 64000000},
		{1000000, 27, 32000000},
		{1000000, 28, 16000000},
		{1000000, 29, 8000000},
		{1000000, 30, 4000000},
		{1000000, 31, 2000000},
		{10000000, 0, 10000000},
		{10000000, 1, 5000000},
		{10000000, 2, 2500000},
		{10000000, 3, 1250000},
		{10000000, 4, 625000},
		{10000000, 5, 312500},
		{10000000, 6, 156250},
		{10000000, 7, 78125},
		{10000000, 8, 2147522710},
		{10000000, 9, 1073761355},
		{10000000, 10, 2684364325},
		{10000000, 11, 3489665810},
		{10000000, 12, 1744832905},
		{10000000, 13, 3019900100},
		{10000000, 14, 1509950050},
		{10000000, 15, 754975025},
		{10000000, 16, 2524971160},
		{10000000, 17, 1262485580},
		{10000000, 18, 631242790},
		{10000000, 19, 315621395},
		{10000000, 20, 2305294345},
		{10000000, 21, 3300130820},
		{10000000, 22, 1650065410},
		{10000000, 23, 825032705},
		{10000000, 24, 2560000000},
		{10000000, 25, 1280000000},
		{10000000, 26, 640000000},
		{10000000, 27, 320000000},
		{10000000, 28, 160000000},
		{10000000, 29, 80000000},
		{10000000, 30, 40000000},
		{10000000, 31, 20000000},
		{100000000, 0, 100000000},
		{100000000, 1, 50000000},
		{100000000, 2, 25000000},
		{100000000, 3, 12500000},
		{100000000, 4, 6250000},
		{100000000, 5, 3125000},
		{100000000, 6, 1562500},
		{100000000, 7, 781250},
		{100000000, 8, 390625},
		{100000000, 9, 2147678960},
		{100000000, 10, 1073839480},
		{100000000, 11, 536919740},
		{100000000, 12, 268459870},
		{100000000, 13, 134229935},
		{100000000, 14, 2214598615},
		{100000000, 15, 3254782955},
		{100000000, 16, 3774875125},
		{100000000, 17, 4034921210},
		{100000000, 18, 2017460605},
		{100000000, 19, 3156213950},
		{100000000, 20, 1578106975},
		{100000000, 21, 2936537135},
		{100000000, 22, 3615752215},
		{100000000, 23, 3955359755},
		{100000000, 24, 4125163525},
		{100000000, 25, 4210065410},
		{100000000, 26, 2105032705},
		{100000000, 27, 3200000000},
		{100000000, 28, 1600000000},
		{100000000, 29, 800000000},
		{100000000, 30, 400000000},
		{100000000, 31, 200000000},
		{1000000000, 0, 1000000000},
		{1000000000, 1, 500000000},
		{1000000000, 2, 250000000},
		{1000000000, 3, 125000000},
		{1000000000, 4, 62500000},
		{1000000000, 5, 31250000},
		{1000000000, 6, 15625000},
		{1000000000, 7, 7812500},
		{1000000000, 8, 3906250},
		{1000000000, 9, 1953125},
		{1000000000, 10, 2148460210},
		{1000000000, 11, 1074230105},
		{1000000000, 12, 2684598700},
		{1000000000, 13, 1342299350},
		{1000000000, 14, 671149675},
		{1000000000, 15, 2483058485},
		{1000000000, 16, 3389012890},
		{1000000000, 17, 1694506445},
		{1000000000, 18, 2994736870},
		{1000000000, 19, 1497368435},
		{1000000000, 20, 2896167865},
		{1000000000, 21, 3595567580},
		{1000000000, 22, 1797783790},
		{1000000000, 23, 898891895},
		{1000000000, 24, 2596929595},
		{1000000000, 25, 3445948445},
		{1000000000, 26, 3870457870},
		{1000000000, 27, 1935228935},
		{1000000000, 28, 3115098115},
		{1000000000, 29, 3705032705},
		{1000000000, 30, 4000000000},
		{1000000000, 31, 2000000000},
	}
	for _, f := range rots {
		u := rotate32(f.in, uint(f.in2))
		if u != f.out {
			t.Errorf("expected %d got %d", f.out, u)
		} else {
			t.Logf("OK expected %d got %d", f.out, u)
		}
	}
}

func TestMux(t *testing.T) {
	murs := []NumIn2Out{
		// Mur
		{100, 100, 296331858},
		{100, 1000, 322382418},
		{100, 100000, 3309578834},
		{100, 1000000, 2332043863},
		{100, 10000000, 1532504573},
		{100, 100000000, 109586476},
		{100, 1000000000, 4129826525},
		{1000, 100, 904209336},
		{1000, 1000, 909616056},
		{1000, 100000, 2222629816},
		{1000, 1000000, 4010714045},
		{1000, 10000000, 3999704087},
		{1000, 100000000, 3377899334},
		{1000, 1000000000, 2218174583},
		{100000, 100, 193992030},
		{100000, 1000, 157783390},
		{100000, 100000, 4239037790},
		{100000, 1000000, 2456196451},
		{100000, 10000000, 1388221705},
		{100000, 100000000, 2402195232},
		{100000, 1000000000, 3180234409},
		{1000000, 100, 3909538526},
		{1000000, 1000, 3946730206},
		{1000000, 100000, 3622916830},
		{1000000, 1000000, 1918718691},
		{1000000, 10000000, 1107645065},
		{1000000, 100000000, 2509590432},
		{1000000, 1000000000, 2805288489},
		{10000000, 100, 1497840043},
		{10000000, 1000, 1471134123},
		{10000000, 100000, 1794947499},
		{10000000, 1000000, 3491281318},
		{10000000, 10000000, 353417548},
		{10000000, 100000000, 2900408157},
		{10000000, 1000000000, 2604713900},
		{100000000, 100, 3858265169},
		{100000000, 1000, 3843028049},
		{100000000, 100000, 2534601809},
		{100000000, 1000000, 2896622668},
		{100000000, 10000000, 754906278},
		{100000000, 100000000, 716108471},
		{100000000, 1000000000, 567217414},
		{1000000000, 100, 4215452687},
		{1000000000, 1000, 4179244047},
		{1000000000, 100000, 1192047631},
		{1000000000, 1000000, 32584714},
		{1000000000, 10000000, 2969121872},
		{1000000000, 100000000, 107557113},
		{1000000000, 1000000000, 2529766768},
	}
	for _, f := range murs {
		u := mur(f.in, f.in2)
		if u != f.out {
			t.Errorf("expected %d got %d", f.out, u)
		} else {
			t.Logf("OK expected %d got %d", f.out, u)
		}
	}
}

func TestFmix(t *testing.T) {
	fmixInOut := []NumInOut{
		{0, 0},
		{100, 4258159850},
		{1000, 1718167128},
		{100000, 1391155934},
		{1000000, 37787785},
		{10000000, 3568206535},
		{100000000, 701900797},
		{1000000000, 4234498180},
	}
	for _, f := range fmixInOut {
		u := fmix(f.in)
		if u != f.out {
			t.Errorf("expected %d got %d", f.out, u)
		} else {
			t.Logf("OK expected %d got %d", f.out, u)
		}
	}
}

func TestHash32Len0to4(t *testing.T) {
	strToHash32 := []StrToHash32{
		// Hash32
		{"hi", 4063302914},
		{"hello world", 4181326496},
	}
	for _, s := range strToHash32 {
		if len(s.str) > 4 {
			continue
		}
		hash := Hash32([]byte(s.str))
		if hash != s.expected {
			t.Errorf("failed expected %d got %d", s.expected, hash)
		} else {
			t.Logf("OK expected %d got %d", s.expected, hash)
		}
	}
}

func TestHash32(t *testing.T) {
	strToHash32 := []StrToHash32{
		// Hash32
		// Hash32
		{"", 0xdc56d17a},
		{"a", 0x3c973d4d},
		{"hi", 0xf2311502},
		{"hello world", 0x19a7581a},
		{"lee@lmmrtech.com", 0xaf0a30fe},
		{"docklandsman@gmail.com", 0x5d8cdbf4},
		{"fred@example.com", 0x7acdc357},
		{"Go is a tool for managing Go source code.Usage:	go command [arguments]The commands are:    build       compile packages and dependencies    clean       remove object files    env         print Go environment information    fix         run go tool fix on packages    fmt         run gofmt on package sources    generate    generate Go files by processing source    get         download and install packages and dependencies    install     compile and install packages and dependencies    list        list packages    run         compile and run Go program    test        test packages    tool        run specified go tool    version     print Go version    vet         run go tool vet on packagesUse go help [command] for more information about a command.Additional help topics:    c           calling between Go and C    filetype    file types    gopath      GOPATH environment variable    importpath  import path syntax    packages    description of package lists    testflag    description of testing flags    testfunc    description of testing functionsUse go help [topic] for more information about that topic.", 0x9c8f96f3},
	}
	for _, s := range strToHash32 {
		hash := Hash32([]byte(s.str))
		if hash != s.expected {
			t.Errorf("failed expected %x got %x", s.expected, hash)
		} else {
			t.Logf("OK expected %x got %x", s.expected, hash)
		}
	}
}

func TestHash64(t *testing.T) {
	strToHash64 := []StrToHash64{
		// Hash64
		{"", 0x9ae16a3b2f90404f},
		{"a", 0xb3454265b6df75e3},
		{"hi", 0x6a5d2fba44f012f8},
		{"hello world", 0x588fb7478bd6b01b},
		{"lee@lmmrtech.com", 0x61bec68db00fa2ff},
		{"docklandsman@gmail.com", 0xb678cf3842309f40},
		{"fred@example.com", 0x7fbbcd6191d8dce0},
		{"Go is a tool for managing Go source code.Usage:	go command [arguments]The commands are:    build       compile packages and dependencies    clean       remove object files    env         print Go environment information    fix         run go tool fix on packages    fmt         run gofmt on package sources    generate    generate Go files by processing source    get         download and install packages and dependencies    install     compile and install packages and dependencies    list        list packages    run         compile and run Go program    test        test packages    tool        run specified go tool    version     print Go version    vet         run go tool vet on packagesUse go help [command] for more information about a command.Additional help topics:    c           calling between Go and C    filetype    file types    gopath      GOPATH environment variable    importpath  import path syntax    packages    description of package lists    testflag    description of testing flags    testfunc    description of testing functionsUse go help [topic] for more information about that topic.", 0xafe256550e4567c9},
	}
	for _, s := range strToHash64 {
		hash := Hash64([]byte(s.str))
		if hash != s.expected {
			t.Errorf("failed expected %x got %x", s.expected, hash)
		} else {
			t.Logf("OK expected %x got %x", s.expected, hash)
		}
	}
}

// Testing my doc example code!
func TestExample(t *testing.T) {
	str := "hello world"
	bytes := []byte(str)
	hash := Hash32(bytes)
	fmt.Printf("Hash32(%s) is %x\n", str, hash)
}

func BenchmarkHash32(b *testing.B) {
	str := "docklandsman@gmail.com"
	bytes := []byte(str)
	for i := 0; i < b.N; i++ {
		_ = Hash32(bytes)
	}
}

// FNV1 is run just to compare against farmhash
func BenchmarkFNV1(b *testing.B) {
	str := "docklandsman@gmail.com"
	bytes := []byte(str)
	h := fnv.New32()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(bytes)
		_ = h.Sum([]byte{})
	}
}

// FNV1a is run just to compare against farmhash
func BenchmarkFNV1a(b *testing.B) {
	str := "docklandsman@gmail.com"
	bytes := []byte(str)
	h := fnv.New32a()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(bytes)
		_ = h.Sum([]byte{})
	}
}

// Adler32 is run just to compare against farmhash
func BenchmarkAdler32(b *testing.B) {
	str := "docklandsman@gmail.com"
	bytes := []byte(str)
	for i := 0; i < b.N; i++ {
		_ = adler32.Checksum(bytes)
	}
}

// Crc32 (IEEE) is run just to compare against farmhash
func BenchmarkCrc32(b *testing.B) {
	str := "docklandsman@gmail.com"
	bytes := []byte(str)
	for i := 0; i < b.N; i++ {
		_ = crc32.ChecksumIEEE(bytes)
	}
}

// SHA1 is run just to compare against farmhash
func BenchmarkSHA1(b *testing.B) {
	str := "docklandsman@gmail.com"
	bytes := []byte(str)
	h := sha1.New()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(bytes)
		_ = h.Sum([]byte{})
	}
}
