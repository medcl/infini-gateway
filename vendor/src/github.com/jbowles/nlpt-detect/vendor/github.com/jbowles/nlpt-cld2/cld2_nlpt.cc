#include "cld2_nlpt.h"

#include <cstddef>
#include <string.h>
#include <stdio.h>
#include <string>

#include "compact_lang_det.h"

const char* CLD2_Static_ExtDetectLanguageSummary(char *data) {

    bool is_plain_text = true;
    CLD2::CLDHints cldhints = {NULL, NULL, 0, CLD2::UNKNOWN_LANGUAGE};
    bool allow_extended_lang = true;
    int flags = 0;
    CLD2::Language language3[3];
    int percent3[3];
    double normalized_score3[3];
    CLD2::ResultChunkVector resultchunkvector;
    int text_bytes;
    bool is_reliable;

    int length = strlen(data);

    CLD2::Language summary_lang = CLD2::UNKNOWN_LANGUAGE;

    summary_lang = CLD2::ExtDetectLanguageSummary(data, 
            length,
            is_plain_text,
            &cldhints,
            flags,
            language3,
            percent3,
            normalized_score3,
            &resultchunkvector,
            &text_bytes,
            &is_reliable);

    return CLD2::LanguageCode(summary_lang);
}

/*
ResultChunks *CLD2_ResultChunkVector_new() {
    return static_cast<ResultChunks*>(new CLD2::ResultChunkVector());
}

const ResultChunk *CLD2_ResultChunkVector_data(const ResultChunks *chunks) {
    const CLD2::ResultChunkVector *vec =
        static_cast<const CLD2::ResultChunkVector *>(chunks);
    return static_cast<const ResultChunk*>(static_cast<const void*>(&(*vec)[0]));
}

size_t CLD2_ResultChunkVector_size(const ResultChunks *chunks) {
    const CLD2::ResultChunkVector *vec =
        static_cast<const CLD2::ResultChunkVector *>(chunks);
    return vec->size();
}

void CLD2_ResultChunkVector_delete(ResultChunks *chunks) {
    CLD2::ResultChunkVector *vec =
        static_cast<CLD2::ResultChunkVector *>(chunks);
    delete vec;
}
*/

const char* CLD2_LanguageName(Language lang) {
    return CLD2::LanguageName(lang);
}

const char* CLD2_LanguageCode(Language lang) {
    return CLD2::LanguageCode(lang);
}

const char* CLD2_LanguageDeclaredName(Language lang) {
    return CLD2::LanguageDeclaredName(lang);
}

const char* CLD2_DetectLanguageVersion() {
    return CLD2::DetectLanguageVersion();
}

Language CLD2_GetLanguageFromName(const char* src) {
    return CLD2::GetLanguageFromName(src);
}

Language CLD2_DetectLanguage(const char* buffer,int buffer_length) {
  bool is_plain_text = true;
  bool is_reliable = true;
  if (buffer_length <= 0) {
      buffer_length = strlen(buffer);
  }
  return CLD2::DetectLanguage(buffer, buffer_length, is_plain_text, &is_reliable);
}

// TODO: add support for passing language hints. will require mapping table for the c++ table of supported languages.
Language CLD2_DetectExtendLanguageSummary(const char *buffer, int buffer_length, int rank, int percent, int normal_score) {
    bool is_plain_text = true;
    bool allow_extended_lang = true;
    bool is_reliable = true;
    CLD2::CLDHints cldhints = {NULL, NULL, 0, CLD2::UNKNOWN_LANGUAGE};
    int flags = 0;
    CLD2::Language language3[rank]; //3
    int percent3[percent]; //3
    double normalized_score3[normal_score]; //3
    CLD2::ResultChunkVector resultchunkvector;
    int text_bytes;
    if (buffer_length <= 0) {
      buffer_length = strlen(buffer);
    }

    return CLD2::ExtDetectLanguageSummary(buffer, 
            buffer_length,
            is_plain_text,
            &cldhints,
            flags,
            language3,
            percent3,
            normalized_score3,
            &resultchunkvector,
            &text_bytes,
            &is_reliable);
}
/*
Language CLD2_DetectLanguageSummary(const char* buffer, int buffer_length, Language* language3, int* percent3, int* text_bytes) {
  bool is_plain_text = true;
  bool is_reliable = true;
  if (buffer_length <= 0) {
      buffer_length = strlen(buffer);
  }
  return CLD2::DetectLanguageSummary(buffer,buffer_length,is_plain_text,language3,percent3,text_bytes,&is_reliable);
}

Language CLD2_DetectLanguageSummary2(const char* buffer,
                                     int buffer_length,
                                     bool is_plain_text,
                                     const char* tld_hint,
                                     int encoding_hint,
                                     Language language_hint,
                                     Language* language3,
                                     int* percent3,
                                     int* text_bytes,
                                     bool* is_reliable)
{
    return CLD2::DetectLanguageSummary(buffer, buffer_length, is_plain_text,
                                       tld_hint, encoding_hint, language_hint,
                                       language3, percent3, text_bytes,
                                       is_reliable);
}

Language CLD2_ExtDetectLanguageSummary(const char* buffer,
                                       int buffer_length,
                                       bool is_plain_text,
                                       Language* language3,
                                       int* percent3,
                                       int* text_bytes,
                                       bool* is_reliable)
{
    return CLD2::ExtDetectLanguageSummary(buffer, buffer_length, is_plain_text,
                                          language3, percent3, text_bytes,
                                          is_reliable);
}

Language CLD2_ExtDetectLanguageSummary2(const char* buffer,
                                        int buffer_length,
                                        bool is_plain_text,
                                        const char* tld_hint,
                                        int encoding_hint,
                                        Language language_hint,
                                        Language* language3,
                                        int* percent3,
                                        int* text_bytes,
                                        bool* is_reliable)
{
    return CLD2::ExtDetectLanguageSummary(buffer, buffer_length, is_plain_text,
                                          tld_hint, encoding_hint,
                                          language_hint, language3, percent3,
                                          text_bytes, is_reliable);
}

Language CLD2_ExtDetectLanguageSummary3(const char* buffer,
                                        int buffer_length,
                                        bool is_plain_text,
                                        const char* tld_hint,
                                        int encoding_hint,
                                        Language language_hint,
                                        Language* language3,
                                        int* percent3,
                                        double* normalized_score3,
                                        int* text_bytes,
                                        bool* is_reliable)
{
    return CLD2::ExtDetectLanguageSummary(buffer, buffer_length, is_plain_text,
                                          tld_hint, encoding_hint, language_hint,
                                          language3, percent3, normalized_score3,
                                          text_bytes, is_reliable);
}

Language CLD2_ExtDetectLanguageSummary4(const char* buffer,
                                        int buffer_length,
                                        bool is_plain_text,
                                        const CLDHints* cld_hints,
                                        int flags,
                                        Language* language3,
                                        int* percent3,
                                        double* normalized_score3,
                                        ResultChunks *resultchunkvector,
                                        int* text_bytes,
                                        bool* is_reliable)
{
    const CLD2::CLDHints *hints =
        static_cast<const CLD2::CLDHints *>(static_cast<const void*>(cld_hints));
    CLD2::ResultChunkVector *vec =
        static_cast<CLD2::ResultChunkVector *>(resultchunkvector);
    return CLD2::ExtDetectLanguageSummary(buffer, buffer_length, is_plain_text,
                                          hints, flags, language3, percent3,
                                          normalized_score3, vec, text_bytes,
                                          is_reliable);
}
*/
