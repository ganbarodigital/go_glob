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
	"unicode/utf8"
)

// MatchPrefix returns the prefix of input that matches the glob pattern
//
// flags can be:
// - GlobShortestMatch (default)
// - GlobLongestMatch
func MatchPrefix(input, pattern string, flags int) (int, bool) {
	// shorthand
	pLen := len(pattern)
	iLen := len(input)

	// where are we in the pattern?
	p := 0

	// where are we in the input?
	i := 0

	// if we need to restart, where do we restart from?
	nextP := 0
	nextI := -1

	// we need to know if we are currently matching a variable-length
	// wildcard
	amMatching := false

	// `i` will track our position in the input
	for i < iLen {
		c := pattern[p]

		if c == '*' {
			// special case - variable length wildcard match
			// we will 'fall forwards', and restart from the next
			// character in our input, until we find the match
			nextP = p
			nextI = i + 1
			amMatching = true
			p++

			// we want to try and match the current input char
			// against our *next* pattern
			if p < pLen {
				c = pattern[p]
			} else {
				c = '?'
			}
		}

		switch c {
		case '?':
			p++
			i++
		default:
			// do we match?
			if input[i] == c {
				p++
				i++
				amMatching = false
			} else if nextI >= 0 {
				i = nextI
				p = nextP
			} else {
				// no, we do not
				return 0, false
			}
		}

		// have we reached the end of the pattern?
		if p >= pLen {
			if amMatching {
				p = pLen - 1
			} else {
				return i, true
			}
		}
	}

	// did we fall off the end, matching a variable-length wildcard?
	if amMatching {
		return len(input), true
	}

	// all done
	return 0, false
}

// MatchGreedyPrefix will treat '*' as match-most
func MatchGreedyPrefix(input, pattern string, flags int) (int, bool) {
	// special case - empty input, empty pattern
	if len(input) == 0 && len(pattern) == 0 {
		return 0, true
	}

	// we need to break down our pattern
	pR := []rune(pattern)
	iR := []rune(input)

	// how many runes are there in the pattern?
	pLen := len(pR)

	// and how many in the input?
	iLen := len(iR)

	// where are we in the pattern?
	p := 0

	// where are we in the input?
	i := 0
	iPos := 0

	// if we need to restart, where do we restart from?
	nextP := 0
	nextI := -1
	nextIPos := -1

	// we need to know if we are currently matching a variable-length
	// wildcard
	amMatching := false

	// `i` will track our position in the input
	for i < iLen || p < pLen {
		if p < pLen {
			c := pR[p]

			if c == '*' {
				if i < iLen {
					// special case - variable length wildcard match
					// we will 'fall backwards', and restart from the next
					// character in our input, until we find the match
					nextP = p
					nextI = i + 1
					nextIPos = iPos + utf8.RuneLen(iR[i])
					amMatching = true
					p++

					// we want to try and match the current input char
					// against our *next* pattern
					if p < pLen {
						c = pR[p]
					} else {
						c = '?'
					}
				} else if p == pLen-1 {
					// we've reached the end of the pattern
					return iPos, true
				}
			}

			switch c {
			case '?':
				if i < iLen {
					p++
					iPos += utf8.RuneLen(iR[i])
					i++
				}
			default:
				// do we match?
				if i < iLen && iR[i] == c {
					p++
					iPos += utf8.RuneLen(iR[i])
					i++
					amMatching = false
				} else if nextI >= 0 {
					i = nextI
					p = nextP
					iPos = nextIPos
				} else {
					// no, we do not
					return 0, false
				}
			}
		}

		// have we reached the end of the pattern?
		if p >= pLen {
			if amMatching {
				p = pLen - 1
			} else {
				return iPos, true
			}
		}
	}

	// did we fall off the end, matching a variable-length wildcard?
	if amMatching {
		return iLen, true
	}

	// all done
	return 0, false
}
