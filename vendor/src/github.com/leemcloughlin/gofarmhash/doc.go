// farmhash Google hash algorithm

/*
Farmhash is a successor to Cityhash (both from Google)

Original Copyright

Copyright (c) 2014 Google, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

FarmHash, by Geoff Pike

Conversion Notes

Converted by Lee McLoughlin, LMMRTech

Converted from the original C++ source code by building it (on a Ubuntu
based system with an Intel CPU). Then copying the build command and editing it to
generate the output of the C pre-processor stage. This showed a version of what code was
required. I then copied code from the original files to convert to Go
in order to preserve original comments.

Note: If you want to compare results between this Go library and the original then
when building the C++ its important to build with -DFARMHASH_DEBUG=0
(or edit src/farmhash.cc and add a #define) otherwise the results are byte swapped
for reasons I don't understand. Of course a byte swapped hash is still a hash.

To test I wrote a small program in C++ to generate both hashes and results from
internal routines to add to the test routines here in the Go version. This ensures
these func's work the same as the C++ versions.

To obey Go export rules some functions had their first character case changed.

TODO: Sort out all Public vs private names & rationalise my use of
prefixes (cc, mk, na) that I use to avoid clashes.

TODO: Figure out how to hash incrementally to use with the Go standard hash package.

TODO: More testing!

Note: An earlier version was a more literal conversion and lots of functions
passed a len parameter after every slice passed.

Note: I'm sure others have already converted farmhash to Go but I'm improving
my Go skills and wanted the experience.
*/
package farmhash
