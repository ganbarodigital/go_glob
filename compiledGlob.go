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
	"fmt"
	"regexp"
)

type compiledGlob struct {
	regex   *regexp.Regexp
	matcher func(string) (int, bool, error)
	flags   int
}

func (g *compiledGlob) assignMatcher(flags int) error {
	switch flags {
	case GlobAnchorPrefix + GlobShortestMatch:
		g.matcher = g.matchShortestPrefix
	case GlobAnchorPrefix + GlobLongestMatch:
		g.matcher = g.matchLongestPrefix
	case GlobAnchorSuffix + GlobShortestMatch:
		g.matcher = g.matchShortestSuffix
	case GlobAnchorSuffix + GlobLongestMatch:
		g.matcher = g.matchLongestSuffix
	case GlobAnchorPrefix + GlobAnchorSuffix, GlobAnchorPrefix + GlobAnchorSuffix + GlobLongestMatch:
		g.matcher = g.matchWholeString
	default:
		return fmt.Errorf("glob compiler: unsupported flags combination %d", flags)
	}

	return nil
}

func (g *compiledGlob) matchWholeString(input string) (int, bool, error) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false, nil
	}

	return len(input), true, nil
}

func (g *compiledGlob) matchShortestPrefix(input string) (int, bool, error) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false, nil
	}

	return loc[1], true, nil
}

func (g *compiledGlob) matchLongestPrefix(input string) (int, bool, error) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false, nil
	}

	return loc[1], true, nil
}

func (g *compiledGlob) matchShortestSuffix(input string) (int, bool, error) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false, nil
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
			return lastLoc[0], true, nil
		}
		copy(lastLoc, subLoc)
		lastLoc[0] += i
		lastLoc[1] += i
		i += subLoc[0] + 1
	}
	return lastLoc[0], true, nil
}

func (g *compiledGlob) matchLongestSuffix(input string) (int, bool, error) {
	loc := g.regex.FindStringIndex(input)
	if loc == nil {
		return 0, false, nil
	}

	return loc[0], true, nil
}
