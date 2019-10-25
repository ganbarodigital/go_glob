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
	"strings"
)

// buildRegex turns a parsed pattern into the equivalent Golang regex
// expression
func buildRegex(pattern []parsedPattern, flags int) string {
	rawRegex := strings.Builder{}

	if flags&GlobAnchorPrefix != 0 {
		rawRegex.WriteRune('^')
	}

	for pos, part := range pattern {
		switch part.patternType {
		case patternTypeSingleMatch:
			rawRegex.WriteString(".")
		case patternTypeMultiMatch:
			// special case - a * at the end of the pattern must always
			// be longest match
			if pos == len(pattern)-1 {
				rawRegex.WriteString(".*")
			} else if flags&GlobLongestMatch != 0 {
				rawRegex.WriteString(".*")
			} else {
				rawRegex.WriteString(".*?")
			}
		case patternTypeStatic:
			// TODO - we need to escape characters in the pattern
			rawRegex.WriteString(part.pattern)
		}
	}

	if flags&GlobAnchorSuffix != 0 {
		rawRegex.WriteRune('$')
	}

	return rawRegex.String()
}

// Compile attempts to convert your glob pattern into a Golang regex
// that's ready to use
func Compile(pattern string, flags int) (*Glob, error) {
	retval := Glob{
		pattern: pattern,
		flags:   flags,
	}

	parsedPattern := parsePattern(pattern)
	rawRegex := buildRegex(parsedPattern, flags)

	var err error
	retval.regex, err = regexp.Compile(rawRegex)
	if err != nil {
		return nil, err
	}

	err = assignMatcherToGlob(&retval, flags)
	if err != nil {
		return nil, err
	}

	// all done
	return &retval, nil
}

func assignMatcherToGlob(g *Glob, flags int) error {
	switch flags {
	case GlobAnchorPrefix + GlobShortestMatch:
		g.matcher = g.matchGlobShortestPrefix
	case GlobAnchorPrefix + GlobLongestMatch:
		g.matcher = g.matchGlobLongestPrefix
	case GlobAnchorSuffix + GlobShortestMatch:
		g.matcher = g.matchGlobShortestSuffix
	case GlobAnchorSuffix + GlobLongestMatch:
		g.matcher = g.matchGlobLongestSuffix
	case GlobAnchorPrefix + GlobAnchorSuffix:
		g.matcher = g.matchGlobWholeString
	default:
		return fmt.Errorf("unsupported flags combination %d", flags)
	}

	return nil
}

// MustCompile is like Compile, but it panics if the underlying regex
// cannot be compiled
func MustCompile(pattern string, flags int) *Glob {
	retval := Glob{
		pattern: pattern,
		flags:   flags,
	}

	parsedPattern := parsePattern(pattern)
	rawRegex := buildRegex(parsedPattern, flags)

	retval.regex = regexp.MustCompile(rawRegex)

	err := assignMatcherToGlob(&retval, flags)
	if err != nil {
		panic(err)
	}

	// all done
	return &retval
}
