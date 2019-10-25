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

import (
	"regexp"
)

// Glob is a compiled Glob expression, which can safely be reused.
//
// Call either 'Compile' or 'MustCompile' to create a Glob.
type Glob struct {
	pattern string
	regex   *regexp.Regexp
	matcher func(string) (int, bool)
	// the flags passed into 'Compile' or 'MustCompile'
	flags int
}

// Pattern returns a copy of the original glob pattern that was compiled
// into the given Glob
func (g Glob) Pattern() string {
	return g.pattern
}

// Flags returns a copy of the original flags that were passed into the
// Glob compiler
func (g Glob) Flags() int {
	return g.flags
}

// Match returns true if the input string satisfies the pre-compiled
// Glob pattern and flags.
func (g *Glob) Match(input string) bool {
	_, success := g.matcher(input)
	return success
}

// MatchWithPosition returns the position that matches the pre-compiled
// Glob pattern and flags.
func (g *Glob) MatchWithPosition(input string) (int, bool) {
	return g.matcher(input)
}

func (g *Glob) matchGlobWholeString(input string) (int, bool) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return len(input), true
}

func (g *Glob) matchGlobShortestPrefix(input string) (int, bool) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return loc[1], true
}

func (g *Glob) matchGlobLongestPrefix(input string) (int, bool) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return loc[1], true
}

func (g *Glob) matchGlobShortestSuffix(input string) (int, bool) {
	loc := g.regex.FindStringIndex(input)
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
			subLoc = g.regex.FindStringIndex("")
		} else {
			subLoc = g.regex.FindStringIndex(input[i:])
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

func (g *Glob) matchGlobLongestSuffix(input string) (int, bool) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false
	}

	return loc[0], true
}
