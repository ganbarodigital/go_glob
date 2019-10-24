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

import "strings"

const (
	patternTypeNone = iota
	patternTypeStatic
	patternTypeSingleMatch
	patternTypeMultiMatch
)

const (
	patternTokenNone = iota
	patternTokenEscape
	patternTokenStatic
	patternTokenSingleMatch
	patternTokenMultiMatch
)

type parsedPattern struct {
	pattern     string
	patternType int
}

func parsePattern(pattern string) []parsedPattern {
	// what we'll be sending back
	var retval []parsedPattern

	// what are we currently looking at?
	lastTokenType := patternTokenNone
	currentTokenType := patternTokenNone

	patternBuf := strings.Builder{}

	// iterate over the runes
	for _, p := range pattern {
		// special case - have we just seen the start of an escape
		// sequence?
		if lastTokenType == patternTokenEscape {
			patternBuf.WriteRune(p)
			lastTokenType = patternTokenStatic
			continue
		}

		// classify the pattern
		switch p {
		case '\\':
			currentTokenType = patternTokenEscape
			patternBuf.WriteRune(p)
		case '?':
			currentTokenType = patternTokenSingleMatch
			if lastTokenType == patternTokenStatic {
				retval = append(
					retval,
					parsedPattern{
						pattern:     patternBuf.String(),
						patternType: patternTypeStatic,
					},
				)
				patternBuf.Reset()
			}

			retval = append(
				retval,
				parsedPattern{
					pattern:     "?",
					patternType: patternTypeSingleMatch,
				},
			)
		case '*':
			currentTokenType = patternTokenMultiMatch

			if lastTokenType == patternTokenStatic {
				retval = append(
					retval,
					parsedPattern{
						pattern:     patternBuf.String(),
						patternType: patternTypeStatic,
					},
				)
				patternBuf.Reset()
			}

			retval = append(
				retval,
				parsedPattern{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
			)
		case '.':
			// this character needs escaping
			currentTokenType = patternTokenStatic
			patternBuf.WriteString("\\.")
		case '+':
			// this character needs escaping
			currentTokenType = patternTokenStatic
			patternBuf.WriteString("\\+")
		case '(':
			// this character needs escaping
			currentTokenType = patternTokenStatic
			patternBuf.WriteString("\\(")
		case '{':
			// this character needs escaping
			currentTokenType = patternTokenStatic
			patternBuf.WriteString("\\{")
		default:
			currentTokenType = patternTokenStatic
			patternBuf.WriteRune(p)
		}

		lastTokenType = currentTokenType
	}

	// deal with last char in the pattern
	if lastTokenType == patternTokenStatic {
		retval = append(
			retval,
			parsedPattern{
				pattern:     patternBuf.String(),
				patternType: patternTypeStatic,
			},
		)
		patternBuf.Reset()
	}

	// all done
	return retval
}
