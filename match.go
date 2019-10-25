// glob brings UNIX shell-like pattern matching support to Golang
//
// Copyright 2019-present Ganbaro Digital Ltd
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
//   * Neither the names of the copyright holders nor the names of his
//     contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package glob

// Match determines if the whole input string matches the given glob
// pattern.
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
func Match(input, pattern string) (bool, error) {
	g := NewGlob(pattern)
	return g.Match(input)
}

// MatchPrefix returns the prefix of input that matches the glob pattern
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// flags can be:
// - GlobShortestMatch (default)
// - GlobLongestMatch
//
// Returns
// - length of prefix that matches, or zero otherwise
// - `true` if the input has prefix that matched the pattern
func MatchPrefix(input, pattern string, flags int) (int, bool, error) {
	if flags&GlobLongestMatch != 0 {
		return MatchLongestPrefix(input, pattern)
	}

	return MatchShortestPrefix(input, pattern)
}

// MatchShortestPrefix returns the prefix of input that matches the glob
// pattern. It treats '*' as matching minimum number of characters.
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// Returns
// - length of prefix that matches, or zero otherwise
// - `true` if the input has prefix that matched the pattern
func MatchShortestPrefix(input, pattern string) (int, bool, error) {
	g := NewGlob(pattern)
	return g.MatchShortestPrefix(input)
}

// MatchLongestPrefix returns the prefix of input that matches the glob
// pattern. It treats '*' as matching maximum number of characters.
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// Returns
// - length of prefix that matches, or zero otherwise
// - `true` if the input has prefix tath matched the pattern
func MatchLongestPrefix(input, pattern string) (int, bool, error) {
	g := NewGlob(pattern)
	return g.MatchLongestPrefix(input)
}

// MatchSuffix returns the start of input that matches the glob pattern.
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// flags can be:
// - GlobShortestMatch (default)
// - GlobLongestMatch
//
// Returns
// - start of suffix that matches (can be len(input)), or zero otherwise
// - `true` if the input has suffix that matched the pattern
func MatchSuffix(input, pattern string, flags int) (int, bool, error) {
	if flags&GlobLongestMatch != 0 {
		return MatchLongestSuffix(input, pattern)
	}

	return MatchShortestSuffix(input, pattern)
}

// MatchShortestSuffix returns the suffix of input that matches the glob
// pattern. It treats '*' as matching minimum number of characters.
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// It is computationally more expensive than the other MatchXXX() functions,
// due to Golang's leftmost-match mechanics (which we have to compensate for).
//
// Returns
// - start of suffix that matches (can be len(input)), or zero otherwise
// - `true` if the input has suffix that matched the pattern
func MatchShortestSuffix(input, pattern string) (int, bool, error) {
	g := NewGlob(pattern)
	return g.MatchShortestSuffix(input)
}

// MatchLongestSuffix returns the suffix of input that matches the glob
// pattern. It treats '*' as matching maximum number of characters.
//
// Pattern can be built from:
//
//  *     Matches zero or more characters
//  ?     Matches exactly one character
//  [...] Matches any one character within the brackets
//
// any other character matches itself
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// Returns
// - start of suffix that matches (can be len(input)), or zero otherwise
// - `true` if the input has suffix that matched the pattern
func MatchLongestSuffix(input, pattern string) (int, bool, error) {
	g := NewGlob(pattern)
	return g.MatchLongestSuffix(input)
}
