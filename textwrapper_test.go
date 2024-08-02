// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package textwrapper

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		input          string
		limit          int
		expectedOutput []string
	}{
		{
			name:  "Simple case",
			input: "This is a sample text that will be wrapped based on the given limit.",
			limit: 20,
			expectedOutput: []string{
				"This is a sample text",
				"that will be wrapped",
				"based on the given",
				"limit.",
			},
		},
		{
			name:  "Single long word",
			input: "Supercalifragilisticexpialidocious",
			limit: 10,
			expectedOutput: []string{
				"Supercalif",
				"ragilistic",
				"expialidoc",
				"ious",
			},
		},
		{
			name:  "Multiple short words",
			input: "Go is fun",
			limit: 5,
			expectedOutput: []string{
				"Go is",
				"fun",
			},
		},
		{
			name:  "Exact limit",
			input: "Golang",
			limit: 6,
			expectedOutput: []string{
				"Golang",
			},
		},
		{
			name:           "Empty input",
			input:          "",
			limit:          10,
			expectedOutput: []string{},
		},
		{
			name:  "Wide characters",
			input: "こんにちは世界",
			limit: 10,
			expectedOutput: []string{
				"こんにちは",
				"世界",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output := Wrap(tc.input, tc.limit)
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}

func TestGetSplitPos(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		word           string
		limit          int
		expectedOutput int
	}{
		{
			name:           "Word within limit",
			word:           "Short",
			limit:          10,
			expectedOutput: utf8.RuneCountInString("Short"),
		},
		{
			name:           "Word exceeding limit",
			word:           "Supercalifragilisticexpialidocious",
			limit:          10,
			expectedOutput: 10,
		},
		{
			name:           "Exact limit",
			word:           "Golang",
			limit:          6,
			expectedOutput: utf8.RuneCountInString("Golang"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output := getSplitPos(tc.word, tc.limit)
			require.Equal(t, tc.expectedOutput, output)
		})
	}
}

func BenchmarkWrap(b *testing.B) {
	testCases := []struct {
		name  string
		input string
		limit int
	}{
		{
			name:  "Simple case",
			input: "This is a sample text that will be wrapped based on the given limit.",
			limit: 20,
		},
		{
			name:  "Single long word",
			input: "Supercalifragilisticexpialidocious",
			limit: 10,
		},
		{
			name:  "Multiple short words",
			input: "Go is fun",
			limit: 5,
		},
		{
			name:  "Exact limit",
			input: "Golang",
			limit: 6,
		},
		{
			name:  "Empty input",
			input: "",
			limit: 10,
		},
		{
			name:  "Wide characters",
			input: "こんにちは世界",
			limit: 10,
		},
	}

	for _, tc := range testCases {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Wrap(tc.input, tc.limit)
			}
		})
	}
}
