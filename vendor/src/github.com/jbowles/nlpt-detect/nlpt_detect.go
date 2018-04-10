// Package nlpt_detect uses a go wrapper around the cld2 (Compact Language Detection) project for language detection (see cld2 info here https://code.google.com/p/cld2).
package nlpt_detect

import (
	"github.com/jbowles/nlpt-cld2"
	"log"
)

// Detect is the default function for detecting a language using the cld2_nlpt wrapper.
//
// It requires the text, the format, size of the buffer, ranked choice index, reliability percent index, and normal score index.
//
//  Format options return
//	  name 'ENGLISH'
//	  code 'en'
//	  declname 'ENGLISH'
//
// NOTE: cld2 defines indexes its own way. That is, if you want accuracy you should query index 3.
//
// From the description of CLD2:
//   language3 is an array of the top 3 languages or UNKNOWN_LANGUAGE
//   percent3 is an array of the text percentages 0..100 of the top 3 languages
// CLD2 returns the 3 highest ranked languages, the 3 best percentages, and the 3 best normalized scores... all of which are returned from CLD2 as arrays. The integer value here is the index of the array to return.
func Detect(s string, format string, buffer_length, rank, percent, normal_score int) string {
	lang, err := cld2_nlpt.DetectExtendedLanguage(s, format, buffer_length, rank, percent, normal_score)
	if err != nil {
		log.Fatal(err)
	}
	return string(cld2_nlpt.Language(lang))
}

func StaticDetect(s string) string {
	lang, err := cld2_nlpt.SimpleDetect(s)
	if err != nil {
		log.Fatal(err)
	}
	return string(cld2_nlpt.Language(lang))
}

// GetLanguageName returns the the name (for example 'ENGLISH') of detected text.
// If it cannot determine the text then 'ENGLISH' is returned by default.
// It does guarantee the greatest amount of accuracy and will return ENGLISH if it probable identification is not reliable.
func GetLanguageName(s string) string {
	lang, err := cld2_nlpt.DetectLanguage(len(s), s, "name")
	if err != nil {
		log.Fatal(err)
	}
	return string(cld2_nlpt.Language(lang))
}

// GetLanguageCode returns the the code('en') of detected text.
// It should be used for testing or demos or simple text.
// It does guarantee the greatest amount of accuracy and will return 'en' if it probable identification is not reliable.
func GetLanguageCode(s string) string {
	lang, err := cld2_nlpt.DetectLanguage(len(s), s, "code")
	if err != nil {
		log.Fatal(err)
	}
	return string(cld2_nlpt.Language(lang))
}

// GetLanguageDeclaredName returns the the name('ENGLISH') of detected text.
// It should be used for testing or demos or simple text.
// It does guarantee the greatest amount of accuracy.
func GetLanguageDeclaredName(s string) string {
	lang, err := cld2_nlpt.DetectLanguage(len(s), s, "declname")
	if err != nil {
		log.Fatal(err)
	}
	return string(cld2_nlpt.Language(lang))
}
