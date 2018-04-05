package sbloom

import (
	"bytes"
	"encoding/gob"
	"errors"
	"hash"
	"hash/fnv"
)

func init() {
	//register the fnv type as it is commonly used
	gob.Register(fnv.New64())
}

//gobFilter is an internal type that gob will use to represent a Filter.
type gobFilter struct {
	Hash    hash.Hash64
	Seeds   [][]byte
	Filters []*filter
}

//GobEncode returns the gob marshalled value of the filter.
func (f *Filter) GobEncode() (p []byte, err error) {
	if f.bh == nil {
		err = errors.New("no hash function specified")
	}

	gf := gobFilter{
		Hash:    f.bh,
		Filters: f.fs,
		Seeds:   make([][]byte, 0, len(f.hs)),
	}
	for _, sh := range f.hs {
		gf.Seeds = append(gf.Seeds, sh.Seed)
	}
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(gf)
	if err == nil {
		p = buf.Bytes()
	}
	return
}

//GobDecode sets the filters state to the gob marshalled value in the buffer.
func (f *Filter) GobDecode(p []byte) (err error) {
	var gf gobFilter
	buf := bytes.NewReader(p)
	err = gob.NewDecoder(buf).Decode(&gf)
	if err != nil {
		return
	}

	if gf.Hash == nil {
		err = errors.New("no hash function specified")
		return
	}

	//create the hashes from the seeds
	f.hs = make([]sHash, 0, len(gf.Seeds))
	for _, s := range gf.Seeds {
		f.hs = append(f.hs, sHash{
			Ha:   gf.Hash,
			Seed: s,
		})
	}

	//set the other fields
	f.bh = gf.Hash
	f.fs = gf.Filters
	return
}
