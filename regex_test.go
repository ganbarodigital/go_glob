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

type buildRegexTestDataStruct struct {
	input          string
	flags          int
	expectedResult string
}

func TestBuildRegex(t *testing.T) {
	t.Parallel()

	testDataSet := []buildRegexTestDataStruct{
		// static, no flags
		{
			input:          "1234567890",
			expectedResult: "1234567890",
		},
		// static, match from start of string
		{
			input:          "1234567890",
			flags:          GlobAnchorPrefix,
			expectedResult: "^1234567890",
		},
		// static, match against end of string
		{
			input:          "1234567890",
			flags:          GlobAnchorSuffix,
			expectedResult: "1234567890$",
		},
		// static, match whole string
		{
			input:          "1234567890",
			flags:          GlobAnchorSuffix + GlobAnchorPrefix,
			expectedResult: "^1234567890$",
		},
		// single char wildcard, no flags
		{
			input:          "123?5",
			expectedResult: "123.5",
		},
		// single char wildcard, match from start of string
		{
			input:          "123?5",
			flags:          GlobAnchorPrefix,
			expectedResult: "^123.5",
		},
		// single char wildcard, match against end of string
		{
			input:          "123?5",
			flags:          GlobAnchorSuffix,
			expectedResult: "123.5$",
		},
		// single char wildcard, match whole string
		{
			input:          "123?5",
			flags:          GlobAnchorPrefix + GlobAnchorSuffix,
			expectedResult: "^123.5$",
		},
		// two single char wildcards, no flags
		{
			input:          "123??5",
			expectedResult: "123..5",
		},
		// multi char wildcard, no flags
		{
			input:          "123*5",
			expectedResult: "123.*?5",
		},
		// multi char wildcard, match from start of string
		{
			input:          "123*5",
			flags:          GlobAnchorPrefix,
			expectedResult: "^123.*?5",
		},
		// multi char wildcard, match against end of string
		{
			input:          "123*5",
			flags:          GlobAnchorSuffix,
			expectedResult: "123.*?5$",
		},
		// multi char wildcard, match whole string
		{
			input:          "123*5",
			flags:          GlobAnchorPrefix + GlobAnchorSuffix,
			expectedResult: "^123.*?5$",
		},
		// multi char wildcard, shortest match flag
		{
			input:          "123*5",
			flags:          GlobShortestMatch,
			expectedResult: "123.*?5",
		},
		// multi char wildcard, match from start of string, shortest match
		{
			input:          "123*5",
			flags:          GlobAnchorPrefix + GlobShortestMatch,
			expectedResult: "^123.*?5",
		},
		// multi char wildcard, match against end of string, shortest match
		{
			input:          "123*5",
			flags:          GlobAnchorSuffix + GlobShortestMatch,
			expectedResult: "123.*?5$",
		},
		// multi char wildcard, match whole string, shortest match
		{
			input:          "123*5",
			flags:          GlobAnchorPrefix + GlobAnchorSuffix + GlobShortestMatch,
			expectedResult: "^123.*?5$",
		},
		// multi char wildcard, longest match flag
		{
			input:          "123*5",
			flags:          GlobLongestMatch,
			expectedResult: "123.*5",
		},
		// multi char wildcard, match from start of string, longest match
		{
			input:          "123*5",
			flags:          GlobAnchorPrefix + GlobLongestMatch,
			expectedResult: "^123.*5",
		},
		// multi char wildcard, match against end of string, longest match
		{
			input:          "123*5",
			flags:          GlobAnchorSuffix + GlobLongestMatch,
			expectedResult: "123.*5$",
		},
		// multi char wildcard, match whole string, longest match
		{
			input:          "123*5",
			flags:          GlobAnchorPrefix + GlobAnchorSuffix + GlobLongestMatch,
			expectedResult: "^123.*5$",
		},
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		parsedPattern := parsePattern(testData.input)

		// ----------------------------------------------------------------
		// perform the change

		actualResult := buildRegex(parsedPattern, testData.flags)

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, testData.expectedResult, actualResult, testData)
	}
}
