// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colornames

import (
	"image/color"
	"testing"
)

func TestColornames(t *testing.T) {
	if len(Map) != len(Names) {
		t.Fatalf("Map and Names have different length: %d vs %d", len(Map), len(Names))
	}

	for name, want := range testCases {
		got, ok := Map[name]
		if !ok {
			t.Errorf("Did not find %s", name)
			continue
		}
		if got != want {
			t.Errorf("%s:\ngot  %v\nwant %v", name, got, want)
		}
	}
}

var testCases = map[string]color.RGBA{
	"aliceblue":      {240, 248, 255, 255},
	"crimson":        {220, 20, 60, 255},
	"darkorange":     {255, 140, 0, 255},
	"deepskyblue":    {0, 191, 255, 255},
	"greenyellow":    {173, 255, 47, 255},
	"lightgrey":      {211, 211, 211, 255},
	"lightpink":      {255, 182, 193, 255},
	"mediumseagreen": {60, 179, 113, 255},
	"olivedrab":      {107, 142, 35, 255},
	"purple":         {128, 0, 128, 255},
	"slategrey":      {112, 128, 144, 255},
	"yellowgreen":    {154, 205, 50, 255},
}
