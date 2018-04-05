// Package cld2 implements language detection using the
// Compact Language Detector.
//
// The `nlpt` part is a side project of mine for a Natural Language Processing Toolkit in go.
//
// This package includes the relevant sources from the cld2
// project, so it doesn't require any external dependencies.
// It also uses the cld2_nlpt.h and cld2_nlpt.cc files to create a specific CLD2_* namespace
// to distinguish usage here with the original c++ project.
//
// For more information about CLD2, see https://code.google.com/p/cld2/.
//
// This package leaned heavily on two existing projects:
// cld2 go wrapper for extracting relevant code: https://github.com/rainycape/cld2
// rust-cld2 wrapper for creating custom header and c++ files for the CLD2_* namespace https://github.com/emk/rust-cld2
//

package cld2_nlpt

// #include <stdlib.h>
// #include "cld2_nlpt.h"
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

/*
type LanguageDialect string

type LanguageInfo struct {
	Language Language         // language code; "en"
	Dialect  LanguageDialect  // languace code + dialect; "en-uk"
	Scores   [3]LanguageScore // the 3 most likely languages
	Reliable bool             // is the result reliable?
}

type LanguageScore struct {
	Dialect LanguageDialect
	Percent int // probability/"confidence"
}

func DetectLanguageInfo() LanguageInfo {
	...
}

type Cld2Hints struct {
	ContentLanguageHint string
	TldHint             string
	EncodingHint        int
	LanguageHint        int
}
		c_hints := Cld2Hints{
			ContentLanguageHint: C.CString(""),
			TldHint:             C.CString(""),
			EncodingHint:        C.int(0),
			LanguageHint:        C.int(0),
		}
*/

type Language string
type Cld2NlptError struct {
	When time.Time
	Msg  string
}

func (e Cld2NlptError) Error() string {
	return fmt.Sprint("%v, %v", e.When, e.Msg)
}

// SimpleDetect returns the language code for detected language in the given text.
// It uses nlpt_wrapper.h and returns the language code, eg. 'en'.
// C++ function sets the buffer length to the that of the text.
//
// By default it defines plain text as true, reliable as true. This means it will not strip out HTML tags, returns only detection rankings that are know to be reliable.
//
// It also defaults return value to UNKOWN_LANGUAGE before running detection. If a language cannot be reliably detected then UNKOWN will be returned.
//	  bool is_plain_text = true;
//	  bool is_reliable = true;
//    ...
//	  CLD2::Language summary_lang = CLD2::UNKNOWN_LANGUAGE;
//
func SimpleDetect(text string) (lang Language, err error) {
	cs := C.CString(text)
	res := C.CLD2_Static_ExtDetectLanguageSummary(cs)
	defer C.free(unsafe.Pointer(cs))
	if res != nil {
		lang = Language(C.GoString(res))
		return lang, err
	} else {
		err = Cld2NlptError{
			time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
			"result returned nil: C.CLD2_Static_ExtDetectLanguageSummary(cs)",
		}
		return lang, err
	}
	return
}

// DetectLanguage uses nlpt_wrapper.h and returns a format of the output.
// cld2 defualts languages to ENGLISH, and so any unreliability returns default; which can yeild wrong results, expecially for small data sets.
// By default it defines plain text as true and reliable as true. This means it will not strip out HTML tags and returns only detection rankings that are know to be reliable.
//	  bool is_plain_text = true;
//	  bool is_reliable = true;
//
// Format Choices
//	  'name' returns 'ENGLISH'
//	  'code' returns 'en'
//	  'declname' returns 'ENGLISH'
//
// If the buffer_length is less than or equal to zero the C++ code will set the length to the that of the text.
// See cld2_nlpt_test for usage.
func DetectLanguage(buffer_length int, text, format string) (lang Language, err error) {
	c_buffer := C.int(buffer_length)
	c_string := C.CString(text)

	var c_char = C.CString("")
	defer C.free(unsafe.Pointer(c_char))
	defer C.free(unsafe.Pointer(c_string))

	var lang_result C.Language = C.CLD2_DetectLanguage(c_string, c_buffer)

	switch {
	case format == "name":
		c_char = C.CLD2_LanguageName(lang_result)
	case format == "code":
		c_char = C.CLD2_LanguageCode(lang_result)
	case format == "declname":
		c_char = C.CLD2_LanguageDeclaredName(lang_result)
	default:
		c_char = C.CLD2_LanguageCode(lang_result)
	}

	if c_char != nil {
		lang = Language(C.GoString(c_char))
		return lang, err
	} else {
		err = Cld2NlptError{
			time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
			"result returned nil: C.CLD2_LanguageName(C.CLD2_DetectLanguage(cs, b_length))",
		}
		return lang, err
	}
	return
}

// DetectExtendedLanguage uses nlpt_wrapper.h and returns a format of the output.
// By default it defines plain text as true, reliable as true, and language hints as unknown. This means it will not strip out HTML tags, returns only detection rankings that are know to be reliable, and in the case of no reliability returns UNKNOWN.
//	  bool is_plain_text = true;
//	  bool is_reliable = true;
//    ...
//	  CLD2::CLDHints cldhints = {NULL, NULL, 0, CLD2::UNKNOWN_LANGUAGE};
//
// It provides the choice to select out of index of ranked languages, the percent, and normal_score. See cld2 for more info on what these choices mean. Default to 3 for most accuracy.
// Format Choices
//	  'name' returns 'ENGLISH'
//	  'code' returns 'en'
//	  'declname' returns 'ENGLISH'
//
// If the buffer_length is less than or equal to zero the C++ code will set the length to the that of the text.
// See cld2_nlpt_test for usage.
//
// TODO: add support for passing language hints. will require mapping table for the c++ table of supported languages.
func DetectExtendedLanguage(text string, format string, buffer_length, rank, percent, normal_score int) (lang Language, err error) {
	c_buffer := C.int(buffer_length)
	c_string := C.CString(text)

	var c_char = C.CString("")
	defer C.free(unsafe.Pointer(c_char))
	defer C.free(unsafe.Pointer(c_string))

	c_rank := C.int(rank)
	c_percent := C.int(percent)
	c_normal_score := C.int(normal_score)

	var lang_result C.Language = C.CLD2_DetectExtendLanguageSummary(
		c_string,
		c_buffer,
		c_rank,
		c_percent,
		c_normal_score)

	switch {
	case format == "name":
		c_char = C.CLD2_LanguageName(lang_result)
	case format == "code":
		c_char = C.CLD2_LanguageCode(lang_result)
	case format == "declname":
		c_char = C.CLD2_LanguageDeclaredName(lang_result)
	default:
		c_char = C.CLD2_LanguageCode(lang_result)
	}

	if c_char != nil {
		lang = Language(C.GoString(c_char))
		return lang, err
	} else {
		err = Cld2NlptError{
			time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
			"result returned nil: C.CLD2_DetectExtendLanguageSummary",
		}
		return lang, err
	}
	return
}
