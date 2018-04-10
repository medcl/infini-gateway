package sbloom

import (
	"hash/fnv"
	"reflect"
	"testing"
)

func TestGobDeepEqual(t *testing.T) {
	f := NewFilter(fnv.New64(), 8)

	for i := 0; i < 100; i++ {
		f.Add(randSeed())
	}

	m, err := f.GobEncode()
	if err != nil {
		t.Fatal(err)
	}

	g := new(Filter)
	if err := g.GobDecode(m); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(f, g) {
		t.Fatal("Not equal after gob")
	}
}

func TestGobFunctionalEqual(t *testing.T) {
	f := NewFilter(fnv.New64(), 2) //only 1/4 chance of postitives

	for i := 0; i < 1000; i++ {
		f.Add(randSeed())
	}

	m, err := f.GobEncode()
	if err != nil {
		t.Fatal(err)
	}

	g := new(Filter)
	if err := g.GobDecode(m); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		if p := randSeed(); g.Lookup(p) != f.Lookup(p) {
			t.Errorf("differed on %v", p)
		}
	}
}
