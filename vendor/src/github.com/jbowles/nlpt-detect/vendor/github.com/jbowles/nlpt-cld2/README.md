# cld2 nlpt
This is a separate fork of [rainycape's wrapper of cld2](https://github.com/rainycape/cld2).

From the README of rainycape wrapper of cld2:

```sh
Go wrapper for the cld2 language detection library by Google Chrome.

Package cld2 implements language detection using the Compact Language Detector.

This package includes the relevant sources from the cld2 project, so it doesn't
require any external dependencies. For more information about CLD2, see
https://code.google.com/p/cld2/.
```

The `nlpt` part is a side project of mine for a Natural Language Processing Toolkit in go.

## External Sources
This wrapper owes its existence to 3 projects:

* cld2 project -- original code [cld2](https://code.google.com/p/cld2/)
* cld2 go wrapper -- [cld2](https://github.com/rainycape/cld2)
* rust-cld2 wrapper [rust-cld2](https://github.com/emk/rust-cld2)

I'm not very good at C/C++ so I leaved heavily on wrapper projects in Go and Rust mentioned.


## Get it

```sh
go get github.com/jbowles/cld2_nlpt
```

## Using it
See tests for full usage. This package consists of 5 public functions only.

The function `Detect` is the preferred way of using this package, but it requires many arguments that depend on user familiarity with the cld2 project. It uses the full set of options for cld2, including extended language detection. It will eventually support passing a struct of language hints.

If user has no familiarty with cld2 or doesn't want to be bothered with complex usage then `StaticDetect` is for you: it uses the extended language feature and pre-defines all the options to cld2; it requires 1 argument: text to be identified. 

The 3 remaining public functions can be used to return language code, name, or displayed name; they do not use extended language features and use the most basic options in cld2.

In terms of accuracy and reliability `Detect` and `StaticDetect` are most reliable.

```go
package main

import "github.com/jbowles/cld2_nlpt"

func main() {

  string := "This is an english sentence"
  cld2_NLPT.GetLanguageName(s)
   // => ENGLISH
```

## Complex usage
In a project called `smallgear` a language detection API defines the endpoint: `r.HandleFunc("/lang/detect/{text}", LanguageDetect).Methods("GET")`.

LanguageDetect returns a JSON object with timestamp along with the language name and code for the highest ranking detected languages:

```go
func LanguageDetect(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	text := params["text"]
	detected_lang := nlpt_detect.Detect(text, "name", len(text), 3, 3, 3)
	second_rank_lang := nlpt_detect.Detect(text, "name", len(text), 2, 2, 2)
	third_rank_lang := nlpt_detect.Detect(text, "name", len(text), 1, 1, 1)
	four_rank_lang := nlpt_detect.Detect(text, "name", len(text), 0, 0, 0)

	detected_code := nlpt_detect.Detect(text, "code", len(text), 3, 3, 3)
	second_rank_code := nlpt_detect.Detect(text, "code", len(text), 2, 2, 2)
	third_rank_code := nlpt_detect.Detect(text, "code", len(text), 1, 1, 1)
	four_rank_code := nlpt_detect.Detect(text, "code", len(text), 0, 0, 0)

	//w.Write([]byte(m))

	langres := &LangDetectResponse{
		Timestamp: time.Now(),
		Detected1: detected_lang,
		Code1:     detected_code,
		Detected2: second_rank_lang,
		Code2:     second_rank_code,
		Detected3: third_rank_lang,
		Code3:     third_rank_code,
		Detected4: four_rank_lang,
		Code4:     four_rank_code,
		Input:     text,
	}

	response, err := json.Marshal(langres)
	if err != nil {
		log.Printf("Error", err)
	}
	w.Write(response)
}
```


## Misc.
The first version used the original go wrapper

```go
// Package cld2 implements language detection using the
// Compact Language Detector.
//
// This package includes the relevant sources from the cld2
// project, so it doesn't require any external dependencies.
// For more information about CLD2, see https://code.google.com/p/cld2/.

package cld2_nlpt

// #include <stdlib.h>
// #include "cld2_min/cld2.h"
import "C"
import "unsafe"

// Detect returns the language code for detected language
// in the given text.
func Detect(text string) (lang string) {
	cs := C.CString(text)
	res := C.DetectLang(cs, -1)
	defer C.free(unsafe.Pointer(cs))
	if res != nil {
		lang = C.GoString(res)
	}
	return
}
```
