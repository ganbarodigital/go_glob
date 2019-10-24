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
// pattern
func Match(input, pattern string) bool {
	flags := GlobMatchWholeString + GlobShortestMatch
	regex := MustCompile(pattern, flags)

	loc := regex.FindStringIndex(input)
	return !(loc == nil)
}

// MatchPrefix returns the prefix of input that matches the glob pattern
//
// flags can be:
// - GlobShortestMatch (default)
// - GlobLongestMatch
func MatchPrefix(input, pattern string, flags int) (int, bool) {
	if flags&GlobLongestMatch != 0 {
		return MatchLongestPrefix(input, pattern)
	}

	return MatchShortestPrefix(input, pattern)
}

// MatchShortestPrefix treats '*' as matching minimum number of
// characters
func MatchShortestPrefix(input, pattern string) (int, bool) {
	flags := GlobAnchorPrefix + GlobShortestMatch
	regex := MustCompile(pattern, flags)

	loc := regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return loc[1], true
}

// MatchLongestPrefix treats '*' as matching maximum number of
// characters
func MatchLongestPrefix(input, pattern string) (int, bool) {
	flags := GlobAnchorPrefix + GlobLongestMatch
	regex := MustCompile(pattern, flags)

	loc := regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return loc[1], true
}

// MatchSuffix returns the start of input that matches the glob pattern
//
// flags can be:
// - GlobShortestMatch (default)
// - GlobLongestMatch
func MatchSuffix(input, pattern string, flags int) (int, bool) {
	if flags&GlobLongestMatch != 0 {
		return MatchLongestSuffix(input, pattern)
	}

	return MatchShortestSuffix(input, pattern)
}

// MatchShortestSuffix treats '*' as matching minimum number of
// characters
func MatchShortestSuffix(input, pattern string) (int, bool) {
	flags := GlobAnchorSuffix + GlobShortestMatch
	regex := MustCompile(pattern, flags)

	loc := regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	// Golang's regexes return the left-most result ... which may not
	// be the shortest result when we're anchoring to a suffix
	//
	// once we have found a match, we need to see if there is a shorter
	// string that will also match
	//
	// I'm sure this can be optimised in the future. PRs most welcome!
	lastLoc := loc
	i := lastLoc[0] + 1
	for i < len(input)+1 {
		var subLoc []int
		if i == len(input) {
			subLoc = regex.FindStringIndex("")
		} else {
			subLoc = regex.FindStringIndex(input[i:])
		}
		if subLoc == nil {
			return lastLoc[0], true
		}
		copy(lastLoc, subLoc)
		lastLoc[0] += i
		lastLoc[1] += i
		i += subLoc[0] + 1
	}
	return lastLoc[0], true
}

// MatchLongestSuffix treats '*' as matching maximum number of
// characters
func MatchLongestSuffix(input, pattern string) (int, bool) {
	flags := GlobAnchorSuffix + GlobLongestMatch
	regex := MustCompile(pattern, flags)

	loc := regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return loc[0], true
}
