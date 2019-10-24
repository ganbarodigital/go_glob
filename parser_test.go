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

type parserTestDataStruct struct {
	input          string
	expectedResult []parsedPattern
}

func TestParsePattern(t *testing.T) {
	t.Parallel()

	testDataSet := []parserTestDataStruct{
		{
			input: "1234567890",
			expectedResult: []parsedPattern{
				{
					pattern:     "1234567890",
					patternType: patternTypeStatic,
				},
			},
		},
		{
			input: "1*0",
			expectedResult: []parsedPattern{
				{
					pattern:     "1",
					patternType: patternTypeStatic,
				},
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
				{
					pattern:     "0",
					patternType: patternTypeStatic,
				},
			},
		},
		{
			input: "*0*",
			expectedResult: []parsedPattern{
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
				{
					pattern:     "0",
					patternType: patternTypeStatic,
				},
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
			},
		},
		{
			input: "?0*",
			expectedResult: []parsedPattern{
				{
					pattern:     "?",
					patternType: patternTypeSingleMatch,
				},
				{
					pattern:     "0",
					patternType: patternTypeStatic,
				},
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
			},
		},
		{
			input: "*0?",
			expectedResult: []parsedPattern{
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
				{
					pattern:     "0",
					patternType: patternTypeStatic,
				},
				{
					pattern:     "?",
					patternType: patternTypeSingleMatch,
				},
			},
		},
		{
			input: "\\?0*",
			expectedResult: []parsedPattern{
				{
					pattern:     "\\?0",
					patternType: patternTypeStatic,
				},
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
			},
		},
		{
			input: "\\?0\\*",
			expectedResult: []parsedPattern{
				{
					pattern:     "\\?0\\*",
					patternType: patternTypeStatic,
				},
			},
		},
		{
			input: "*.go",
			expectedResult: []parsedPattern{
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
				{
					pattern:     "\\.go",
					patternType: patternTypeStatic,
				},
			},
		},
		{
			input: "*+go",
			expectedResult: []parsedPattern{
				{
					pattern:     "*",
					patternType: patternTypeMultiMatch,
				},
				{
					pattern:     "\\+go",
					patternType: patternTypeStatic,
				},
			},
		},
	}

	for _, testData := range testDataSet {
		// ----------------------------------------------------------------
		// setup your test

		// ----------------------------------------------------------------
		// perform the change

		actualResult := parsePattern(testData.input)

		// ----------------------------------------------------------------
		// test the results

		assert.Equal(t, testData.expectedResult, actualResult)
	}
}
