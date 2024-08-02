// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

// Package textwrapper provides functions for wrapping text into lines
// with a specified maximum width. This package is designed to handle
// text wrapping efficiently, ensuring words are not split across lines
// unless they exceed the width limit.
package textwrapper

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/text/width"
)

// Wrap splits the input string into lines with a maximum width limit.
//
// This function is designed to handle input strings efficiently, but it is
// optimized for typical use cases. Based on benchmark tests, it can comfortably
// handle input strings up to several megabytes in size without significant
// performance degradation. For extremely large inputs (tens of megabytes or more),
// performance may vary based on the available system resources.
//
// If processing very large inputs is required frequently, consider additional
// optimizations or using a streaming approach to manage memory usage more effectively.
func Wrap(s string, limit int) []string {
	if s == "" {
		return []string{}
	}

	// Estimate initial capacity for the strings.Builder to reduce allocations
	estimatedCapacity := len(s) / limit
	var lines []string
	lineBuilder := strings.Builder{}
	lineBuilder.Grow(estimatedCapacity * (limit + 1)) // +1 for potential spaces
	var lineWidth int

	words := strings.Fields(s)
	for _, word := range words {
		wordWidth := getWordWidth(word)
		if lineWidth+wordWidth > limit {
			// If the current line is not empty, add it to the lines slice
			if lineWidth > 0 {
				lines = append(lines, lineBuilder.String())
				lineBuilder.Reset()
				lineWidth = 0
			}
			// If the word itself is longer than the limit, split the word
			for wordWidth > limit {
				splitPos := getSplitPos(word, limit)
				lines = append(lines, word[:splitPos])
				word = word[splitPos:]
				wordWidth = getWordWidth(word)
			}
		}
		// If the line is not empty, add a space before the next word
		if lineWidth > 0 {
			lineBuilder.WriteRune(' ')
			lineWidth++
		}
		lineBuilder.WriteString(word)
		lineWidth += wordWidth
	}

	// Add any remaining content in the builder to the lines
	if lineBuilder.Len() > 0 {
		lines = append(lines, lineBuilder.String())
	}

	return lines
}

// getWordWidth calculates the display width of a word.
func getWordWidth(word string) int {
	width := 0
	for _, r := range word {
		width += runeWidth(r)
	}
	return width
}

// runeWidth calculates the display width of a rune.
func runeWidth(r rune) int {
	k := width.LookupRune(r).Kind()
	if k == width.EastAsianFullwidth || k == width.EastAsianWide {
		return 2
	}
	return 1
}

// getSplitPos determines the position to split the word to fit within the limit.
func getSplitPos(word string, limit int) int {
	width := 0
	for i, r := range word {
		width += runeWidth(r)
		if width > limit {
			return i
		}
	}
	return utf8.RuneCountInString(word)
}
