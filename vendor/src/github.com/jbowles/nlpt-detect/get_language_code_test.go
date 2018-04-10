package nlpt_detect

import (
	"testing"
)

func TestCodeEnglishDetection(t *testing.T) {
	english := "This is an english sentence"
	expect := "en"
	got := GetLanguageCode(english)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestCodeGreekDetection(t *testing.T) {
	greek := "Αυτή είναι μια ελληνική πρόταση"
	expect := "el"
	got := GetLanguageCode(greek)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestCodeMandarinDetection(t *testing.T) {
	mandarin := "这是中国一句"
	expect := "zh"
	got := GetLanguageCode(mandarin)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}
