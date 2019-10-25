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

// Glob is a compiled Glob expression, which can safely be reused.
//
// Call `NewGlob()` to create your Glob structure
type Glob struct {
	pattern       string
	patternParts  []parsedPattern
	compiledGlobs map[int]*compiledGlob
}

// NewGlob turns your pattern into a reusable Glob
func NewGlob(pattern string, options ...func(*Glob)) *Glob {
	// create the Glob we're going to send back
	retval := Glob{
		pattern:       pattern,
		patternParts:  parsePattern(pattern),
		compiledGlobs: make(map[int]*compiledGlob, 5),
	}

	// apply any options we've been given
	for _, option := range options {
		option(&retval)
	}

	// all done
	return &retval
}

// Pattern returns a copy of the original glob pattern that was compiled
// into the given Glob
func (g Glob) Pattern() string {
	return g.pattern
}

// compile creates a new regex from the previously parsed pattern, that will
// satisfy the given flags.
func (g *Glob) compile(flags int) (*compiledGlob, error) {
	retval := compiledGlob{}
	rawRegex := buildRegex(g.patternParts, flags)

	var err error
	retval.regex, err = regexp.Compile(rawRegex)
	if err != nil {
		return nil, fmt.Errorf("bad or unsupported glob pattern '%s': %s", g.pattern, err.Error())
	}

	err = retval.assignMatcher(flags)
	if err != nil {
		return nil, err
	}

	// all done
	return &retval, nil

}

// getCompiledGlobForFlags will return a compiled glob that will satisfy
// the given flags.
//
// It will build a new compiled glob if a suitable one does not already
// exist.
func (g *Glob) getCompiledGlobForFlags(flags int) (*compiledGlob, error) {
	// do we already have a compiled glob?
	existingGlob, ok := g.compiledGlobs[flags]
	if ok {
		return existingGlob, nil
	}

	// no, we need to make one
	compiledGlob, err := g.compile(flags)
	if err != nil {
		return nil, err
	}

	g.compiledGlobs[flags] = compiledGlob
	return compiledGlob, nil
}

// Match determines if the whole input string matches the given glob
// pattern.
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
func (g *Glob) Match(input string) (bool, error) {
	compiledGlob, err := g.getCompiledGlobForFlags(GlobMatchWholeString)
	if err != nil {
		return false, err
	}

	_, success, err := compiledGlob.matcher(input)
	return success, err
}

// MatchShortestPrefix returns the prefix of input that matches the glob
// pattern. It treats '*' as matching minimum number of characters.
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// Returns
// - length of prefix that matches, or zero otherwise
// - `true` if the input has prefix that matched the pattern
func (g *Glob) MatchShortestPrefix(input string) (int, bool, error) {
	compiledGlob, err := g.getCompiledGlobForFlags(GlobAnchorPrefix + GlobShortestMatch)
	if err != nil {
		return 0, false, err
	}

	return compiledGlob.matcher(input)
}

// MatchLongestPrefix returns the prefix of input that matches the glob
// pattern. It treats '*' as matching maximum number of characters.
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// Returns
// - length of prefix that matches, or zero otherwise
// - `true` if the input has prefix tath matched the pattern
func (g *Glob) MatchLongestPrefix(input string) (int, bool, error) {
	compiledGlob, err := g.getCompiledGlobForFlags(GlobAnchorPrefix + GlobLongestMatch)
	if err != nil {
		return 0, false, err
	}

	return compiledGlob.matcher(input)
}

// MatchShortestSuffix returns the suffix of input that matches the glob
// pattern. It treats '*' as matching minimum number of characters.
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
func (g *Glob) MatchShortestSuffix(input string) (int, bool, error) {
	compiledGlob, err := g.getCompiledGlobForFlags(GlobAnchorSuffix + GlobShortestMatch)
	if err != nil {
		return 0, false, err
	}

	return compiledGlob.matcher(input)
}

// MatchLongestSuffix returns the suffix of input that matches the glob
// pattern. It treats '*' as matching maximum number of characters.
//
// Intent is to be 100% compatible with UNIX shell globbing. Please open
// a GitHub issue if you find any test cases that show up compatibility
// problems.
//
// Returns
// - start of suffix that matches (can be len(input)), or zero otherwise
// - `true` if the input has suffix that matched the pattern
func (g *Glob) MatchLongestSuffix(input string) (int, bool, error) {
	compiledGlob, err := g.getCompiledGlobForFlags(GlobAnchorSuffix + GlobLongestMatch)
	if err != nil {
		return 0, false, err
	}

	return compiledGlob.matcher(input)
}
