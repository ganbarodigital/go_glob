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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobPatternReturnsOriginalPattern(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pattern := "12345"
	g := NewGlob(pattern)
	expectedResult := pattern

	// ----------------------------------------------------------------
	// perform the change

	actualResult := g.Pattern()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestGlobCompileReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	flags := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := g.compile(flags)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Nil(t, actualResult)
}

func TestGlobCompileReturnsErrorWhenFlagsInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pattern := "12345"
	g := NewGlob(pattern)

	flags := 256

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := g.compile(flags)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Nil(t, actualResult)
}

func TestGlobGetCompiledGlobForFlagsReusesCompiledGlobs(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pattern := "12345"
	g := NewGlob(pattern)

	flags := GlobAnchorPrefix
	existingCG, err := g.getCompiledGlobForFlags(flags)
	assert.Nil(t, err)

	// ----------------------------------------------------------------
	// perform the change

	actualCG, err := g.getCompiledGlobForFlags(flags)
	assert.Nil(t, err)

	// ----------------------------------------------------------------
	// test the results

	assert.Same(t, existingCG, actualCG)
}

func TestGlobGetCompiledGlobForFlagsReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	flags := GlobShortestMatch

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := g.getCompiledGlobForFlags(flags)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Nil(t, actualResult)
}

func TestGlobMatchReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	// ----------------------------------------------------------------
	// perform the change

	success, err := g.Match("")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.False(t, success)
}

func TestGlobMatchShortestPrefixReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	// ----------------------------------------------------------------
	// perform the change

	pos, success, err := g.MatchShortestPrefix("")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Equal(t, 0, pos)
	assert.False(t, success)
}

func TestGlobMatchLongestPrefixReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	// ----------------------------------------------------------------
	// perform the change

	pos, success, err := g.MatchLongestPrefix("")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Equal(t, 0, pos)
	assert.False(t, success)
}

func TestGlobMatchShortestSuffixReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	// ----------------------------------------------------------------
	// perform the change

	pos, success, err := g.MatchShortestSuffix("")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Equal(t, 0, pos)
	assert.False(t, success)
}

func TestGlobMatchLongestSuffixReturnsErrorWhenRegexInvalid(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// this pattern is invalid because of the mismatched '['
	pattern := "12345["
	g := NewGlob(pattern)

	// ----------------------------------------------------------------
	// perform the change

	pos, success, err := g.MatchLongestSuffix("")

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
	assert.Equal(t, 0, pos)
	assert.False(t, success)
}
