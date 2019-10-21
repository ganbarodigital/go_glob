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

type testDataStruct struct {
	input           string
	pattern         string
	flags           int
	expectedResult  string
	expectedSuccess bool
}

func TestPrefixMatchesStaticStrings(t *testing.T) {
	t.Parallel()

	testDataSet := []testDataStruct{
		{
			input:           "0123456789",
			pattern:         "012345",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// input does not start with static pattern
		{
			input:           "0123456789",
			pattern:         "12345",
			expectedResult:  "",
			expectedSuccess: false,
		},
		// input shorter than static pattern
		{
			input:           "012345",
			pattern:         "0123456789",
			expectedResult:  "",
			expectedSuccess: false,
		},
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		// ----------------------------------------------------------------
		// perform the change

		actualLen, actualSuccess := MatchPrefix(testData.input, testData.pattern, testData.flags)
		actualResult := ""
		if actualSuccess {
			actualResult = testData.input[:actualLen]
		}

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, testData.expectedSuccess, actualSuccess)
		assert.Equal(t, testData.expectedResult, actualResult)
	}
}

func TestPrefixMatchesSingleWildCards(t *testing.T) {
	t.Parallel()

	testDataSet := []testDataStruct{
		{
			input:           "0123456789",
			pattern:         "0?2345",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// multiple single wildcards
		{
			input:           "0123456789",
			pattern:         "0?23?5",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// ALL single wildcards
		{
			input:           "0123456789",
			pattern:         "??????",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// input does not start with pattern
		{
			input:           "0123456789",
			pattern:         "1?345",
			expectedResult:  "",
			expectedSuccess: false,
		},
		// input shorter than pattern
		{
			input:           "012345",
			pattern:         "0?23456789",
			expectedResult:  "",
			expectedSuccess: false,
		},
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		// ----------------------------------------------------------------
		// perform the change

		actualLen, actualSuccess := MatchPrefix(testData.input, testData.pattern, testData.flags)
		actualResult := ""
		if actualSuccess {
			actualResult = testData.input[:actualLen]
		}

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, testData.expectedSuccess, actualSuccess)
		assert.Equal(t, testData.expectedResult, actualResult)
	}
}

func TestPrefixMatchesVariableLengthWildCards(t *testing.T) {
	t.Parallel()

	testDataSet := []testDataStruct{
		// variable-length wildcard, bounded
		{
			input:           "0123456789",
			pattern:         "0*5",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// variable-length wildcard, no prefix
		{
			input:           "0123456789",
			pattern:         "*5",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// variable length wildcard, no suffix
		{
			input:           "0123456789",
			pattern:         "012*",
			expectedResult:  "0123456789",
			expectedSuccess: true,
		},
		// variable length wildcard, match all
		{
			input:           "01234567890",
			pattern:         "*",
			expectedResult:  "01234567890",
			expectedSuccess: true,
		},
		// variable length wildcard, wildcard matches nothing
		{
			input:           "01234567890",
			pattern:         "012*345",
			expectedResult:  "012345",
			expectedSuccess: true,
		},
		// multiple variable length wildcards
		{
			input:           "0123456789",
			pattern:         "0*2*4*6*8*",
			expectedResult:  "0123456789",
			expectedSuccess: true,
		},
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		// ----------------------------------------------------------------
		// perform the change

		actualLen, actualSuccess := MatchPrefix(testData.input, testData.pattern, testData.flags)
		actualResult := ""
		if actualSuccess {
			actualResult = testData.input[:actualLen]
		}

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, testData.expectedSuccess, actualSuccess, testData)
		assert.Equal(t, testData.expectedResult, actualResult, testData)
	}
}
