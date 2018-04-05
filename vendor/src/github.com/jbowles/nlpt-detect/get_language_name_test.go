package nlpt_detect

import (
	"testing"
)

func TestBasicEnglishDetection(t *testing.T) {
	english := "This is an english sentence"
	expect := "ENGLISH"
	got := GetLanguageName(english)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestBasicGreekDetection(t *testing.T) {
	greek := "Αυτή είναι μια ελληνική πρόταση"
	expect := "GREEK"
	got := GetLanguageName(greek)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestBasicMandarinDetection(t *testing.T) {
	mandarin := "这是中国一句"
	expect := "Chinese"
	got := GetLanguageName(mandarin)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestBasicRussianDetection(t *testing.T) {
	russian := "Эторусская предложение"
	expect := "RUSSIAN"
	got := GetLanguageName(russian)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestBasicUkrainianDetection(t *testing.T) {
	ukrainian := "Це український пропозицію"
	expect := "UKRAINIAN"
	got := GetLanguageName(ukrainian)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}
