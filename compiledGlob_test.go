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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompiledGlobAssignMatcherSupportsMatchWholeString(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	g := compiledGlob{}
	flags := GlobMatchWholeString
	expectedResult := fmt.Sprintf("%v", reflect.ValueOf(g.matchWholeString))

	// ----------------------------------------------------------------
	// perform the change

	err := g.assignMatcher(flags)
	actualResult := fmt.Sprintf("%v", reflect.ValueOf(g.matcher))

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestCompiledGlobAssignMatcherSupportsMatchShortestPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	g := compiledGlob{}
	flags := GlobShortestMatch + GlobAnchorPrefix
	expectedResult := fmt.Sprintf("%v", reflect.ValueOf(g.matchShortestPrefix))

	// ----------------------------------------------------------------
	// perform the change

	err := g.assignMatcher(flags)
	actualResult := fmt.Sprintf("%v", reflect.ValueOf(g.matcher))

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestCompiledGlobAssignMatcherSupportsMatchLongestPrefix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	g := compiledGlob{}
	flags := GlobLongestMatch + GlobAnchorPrefix
	expectedResult := fmt.Sprintf("%v", reflect.ValueOf(g.matchLongestPrefix))

	// ----------------------------------------------------------------
	// perform the change

	err := g.assignMatcher(flags)
	actualResult := fmt.Sprintf("%v", reflect.ValueOf(g.matcher))

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestCompiledGlobAssignMatcherSupportsMatchShortestSuffix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	g := compiledGlob{}
	flags := GlobShortestMatch + GlobAnchorSuffix
	expectedResult := fmt.Sprintf("%v", reflect.ValueOf(g.matchShortestSuffix))

	// ----------------------------------------------------------------
	// perform the change

	err := g.assignMatcher(flags)
	actualResult := fmt.Sprintf("%v", reflect.ValueOf(g.matcher))

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestCompiledGlobAssignMatcherSupportsMatchLongestSuffix(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	g := compiledGlob{}
	flags := GlobLongestMatch + GlobAnchorSuffix
	expectedResult := fmt.Sprintf("%v", reflect.ValueOf(g.matchLongestSuffix))

	// ----------------------------------------------------------------
	// perform the change

	err := g.assignMatcher(flags)
	actualResult := fmt.Sprintf("%v", reflect.ValueOf(g.matcher))

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestCompiledGlobAssignMatcherReturnsErrorOtherwise(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	g := compiledGlob{}
	flags := GlobAnchorPrefix + GlobAnchorSuffix + GlobLongestMatch + 1

	// ----------------------------------------------------------------
	// perform the change

	err := g.assignMatcher(flags)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, err)
}
