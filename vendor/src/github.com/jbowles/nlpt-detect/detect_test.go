package nlpt_detect

import (
	"testing"
)

func TestDetectEnglishDetection(t *testing.T) {
	english := "This is an english sentence"
	expect := "ENGLISH"
	got := Detect(english, "name", len(english), 3, 3, 3)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestDetectGreekDetection(t *testing.T) {
	greek := "Αυτή είναι μια ελληνική πρόταση"
	expect := "GREEK"
	got := Detect(greek, "name", len(greek), 3, 3, 3)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestDetectMandarinDetection(t *testing.T) {
	mandarin := "这是中国一句"
	expect := "Chinese"
	got := Detect(mandarin, "name", len(mandarin), 3, 3, 3)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestDetectRussianDetection(t *testing.T) {
	russian := "Эторусская предложение"
	expect := "RUSSIAN"
	got := Detect(russian, "name", len(russian), 3, 3, 3)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}

func TestDetectUkrainianDetection(t *testing.T) {
	ukrainian := "Це український пропозицію"
	expect := "UKRAINIAN"
	got := Detect(ukrainian, "name", len(ukrainian), 3, 3, 3)
	if expect != got {
		t.Log("expected: ", expect, "got... ", got)
		t.Fail()
	}
}
