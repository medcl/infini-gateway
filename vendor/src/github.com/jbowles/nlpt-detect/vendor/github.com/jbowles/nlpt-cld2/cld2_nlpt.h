// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//
// Author: dsites@google.com (Dick Sites)
// Converted to C by Eric Kidd <git@randomhacks.net>
//
// This is a stripped down version of cld2/public/compact_lang_det.h with
// an extern wrapper and "CLD_" preprended to all the symbols we want to
// link.
//

#ifdef __cplusplus

#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include "compact_lang_det.h"
typedef CLD2::Language Language;
typedef CLD2::uint16 uint16;
typedef CLD2::int32 int32;

#else

//// We're being run through bindgen or a similar tool, so just stub all
//// this.
////typedef long long size_t;
////////////////////// WOAH HACK BECAUSE ... APPLE!! ///////////////////////////
////
//// github.com/jbowles/cld2_nlpt
//// In file included from ./cld2_nlpt.go:10:
//// ./cld2_nlpt.h:49:19: error: typedef redefinition with different types ('long long' vs '__darwin_size_t' (aka 'unsigned long'))
//// typedef long long size_t;
//                  ^
//// /usr/include/sys/_types/_size_t.h:30:32: note: previous definition is here
//// typedef __darwin_size_t        size_t;/
////
////http://stackoverflow.com/questions/11603818/why-is-there-ambiguity-between-uint32-t-and-uint64-t-when-using-size-t-on-mac-os
////////////////////// DARWIN: /usr/include/sys/_types ln: 30 ///////////////////////////
////typedef __darwin_size_t size_t;
////
////
//
//typedef unsigned long size_t;
typedef int int32;
typedef unsigned short uint16;
//typedef char bool;
// Corresponds to the Language enumeration in the C++ headers.
typedef int Language;
typedef Language language_hint;

#endif

#ifdef __cplusplus
extern "C" {
#endif

  /*
  typedef struct {
    const char* content_language_hint;      // "mi,en" boosts Maori and English
    const char* tld_hint;                   // "id" boosts Indonesian
    int encoding_hint;                      // SJS boosts Japanese
    Language language_hint;                 // ITALIAN boosts it
  } CLDHints;

  static const int kMaxResultChunkBytes = 65535;

  typedef struct {
    int offset;                 // Starting byte offset in original buffer
    int32 bytes;                // Number of bytes in chunk
    uint16 lang1;               // Top lang, as full Language. Apply
                                //  static_cast<Language>() to this short value.
    uint16 pad;                 // Make multiple of 4 bytes
  } ResultChunk;

  typedef std::vector<ResultChunk> ResultChunkVector;
  // A wrapper around ResultChunkVector, which is a C++ type.
  typedef void ResultChunks;
  ResultChunks *CLD2_ResultChunkVector_new();
  const ResultChunk *CLD2_ResultChunkVector_data(const ResultChunks *chunks);
  size_t CLD2_ResultChunkVector_size(const ResultChunks *chunks);
  void CLD2_ResultChunkVector_delete(ResultChunks *chunks);
  */

  // These APIs are in a private header included by a public header, but
  // they're really useful, so let's assume they're public.
  const char* CLD2_LanguageName(Language lang);
  const char* CLD2_LanguageCode(Language lang);
  const char* CLD2_LanguageDeclaredName(Language lang);
  const char* CLD2_Static_ExtDetectLanguageSummary(char *data);

  // Return version text string, String is "code_version - data_build_date"
  const char* CLD2_DetectLanguageVersion();
  Language CLD2_GetLanguageFromName(const char* src);

  // Scan interchange-valid UTF-8 bytes and detect most likely language
  Language CLD2_DetectLanguage(const char* buffer,int buffer_length);

  // TODO: add support for passing language hints. will require mapping table for the c++ table of supported languages.
  Language CLD2_DetectExtendLanguageSummary(
      const char *buffer,
      int buffer_length,
      int rank,
      int percent,
      int normal_score);

  /*
  // Scan interchange-valid UTF-8 bytes and detect list of top 3 languages.
  // language3[0] is usually also the return value
  Language CLD2_DetectLanguageSummary(const char* buffer,
                          int buffer_length,
                          Language* language3, // int return n-number top langs
                          int* percent3,
                          int* text_bytes);

  // Same as above, with hints supplied
  // Scan interchange-valid UTF-8 bytes and detect list of top 3 languages.
  // language3[0] is usually also the return value
  Language CLD2_DetectLanguageSummary2(const char* buffer,
                          int buffer_length,
                          bool is_plain_text,
                          const char* tld_hint,       // "id" boosts Indonesian
                          int encoding_hint,          // SJS boosts Japanese
                          Language language_hint,     // ITALIAN boosts it
                          Language* language3,
                          int* percent3,
                          int* text_bytes,
                          bool* is_reliable);

  // Scan interchange-valid UTF-8 bytes and detect list of top 3 extended
  // languages.
  //
  // Extended languages are additional interface languages and Unicode
  // single-language scripts, from lang_script.h
  //
  // language3[0] is usually also the return value
  Language CLD2_ExtDetectLanguageSummary(
                          const char* buffer,
                          int buffer_length,
                          bool is_plain_text,
                          Language* language3,
                          int* percent3,
                          int* text_bytes,
                          bool* is_reliable);

  // Same as above, with hints supplied
  // Scan interchange-valid UTF-8 bytes and detect list of top 3 extended
  // languages.
  //
  // Extended languages are additional Google interface languages and Unicode
  // single-language scripts, from lang_script.h
  //
  // language3[0] is usually also the return value
  Language CLD2_ExtDetectLanguageSummary2(
                          const char* buffer,
                          int buffer_length,
                          bool is_plain_text,
                          const char* tld_hint,       // "id" boosts Indonesian
                          int encoding_hint,          // SJS boosts Japanese
                          Language language_hint,     // ITALIAN boosts it
                          Language* language3,
                          int* percent3,
                          int* text_bytes,
                          bool* is_reliable);

  // Same as above, and also returns 3 internal language scores as a ratio to
  // normal score for real text in that language. Scores close to 1.0 indicate
  // normal text, while scores far away from 1.0 indicate badly-skewed text or
  // gibberish
  //
  Language CLD2_ExtDetectLanguageSummary3(
                          const char* buffer,
                          int buffer_length,
                          bool is_plain_text,
                          const char* tld_hint,       // "id" boosts Indonesian
                          int encoding_hint,          // SJS boosts Japanese
                          Language language_hint,     // ITALIAN boosts it
                          Language* language3,
                          int* percent3,
                          double* normalized_score3,
                          int* text_bytes,
                          bool* is_reliable);

  // Use this one.
  // Hints are collected into a struct.
  // Flags are passed in (normally zero).
  //
  // Also returns 3 internal language scores as a ratio to
  // normal score for real text in that language. Scores close to 1.0 indicate
  // normal text, while scores far away from 1.0 indicate badly-skewed text or
  // gibberish
  //
  // Returns a vector of chunks in different languages, so that caller may
  // spell-check, translate, or otherwaise process different parts of the input
  // buffer in language-dependant ways.
  //
  Language CLD2_ExtDetectLanguageSummary4(
                          const char* buffer,
                          int buffer_length,
                          bool is_plain_text,
                          const CLDHints* cld_hints,
                          int flags,
                          Language* language3,
                          int* percent3,
                          double* normalized_score3,
                          ResultChunks *resultchunkvector,
                          int* text_bytes,
                          bool* is_reliable);
  // Public use flags, debug output controls
  static const int kCLDFlagScoreAsQuads = 0x0100;  // Force Greek, etc. => quads
  static const int kCLDFlagHtml =         0x0200;  // Debug HTML => stderr
  static const int kCLDFlagCr =           0x0400;  // <cr> per chunk if HTML
  static const int kCLDFlagVerbose =      0x0800;  // More debug HTML => stderr
  static const int kCLDFlagQuiet =        0x1000;  // Less debug HTML => stderr
  static const int kCLDFlagEcho =         0x2000;  // Echo input => stderr
  */

#ifdef __cplusplus
};
#endif
