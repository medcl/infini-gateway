package nlpt_detect

import (
	"testing"
)

func TestHindiDetection(t *testing.T) {
	hindi := "मनमोहन सिंह मंत्रिमंडल में मंत्रालयों का बँटवारा पूरा हो गया है. कपिल सिब्बल को मानव संसाधन विकास और आनंद शर्मा को वाणिज्य एवं उद्योग विभाग दिया"
	get := GetLanguageName(hindi)
	static := StaticDetect(hindi)
	det := Detect(hindi, "name", len(hindi), 3, 3, 3)
	if det != "HINDI" || static != "hi" || get != "HINDI" {
		t.Log("expected: hindi names and code as 'HINDI', 'hi' and 'HINDI' but got... ", get, static, det)
		t.Fail()
	}
}

func TestJordanianArabicSnippetDetection(t *testing.T) {
	jordan_arabic := "(ف) (ع)"
	get := GetLanguageName(jordan_arabic)                             // ENGLISH
	static := StaticDetect(jordan_arabic)                             // un
	det := Detect(jordan_arabic, "name", len(jordan_arabic), 3, 3, 3) // Unknown
	if get != "ENGLISH" || static != "un" || det != "Unknown" {
		t.Log("expected: Jordanian Arabic names and code as get == 'ENGLISH', static == 'un' and det == 'Unknown' but got... ", get, static, det)
		t.Fail()
	}
}
