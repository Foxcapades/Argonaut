package argo

import "iter"

type StringScannerFactory = func(input string) iter.Seq[string]

// DelimitedSliceScanner returns a new scanner over the given string, breaking
// it into substrings on every delimiter character.
//
// Examples:
//
//	// Comma separated values:
//	DelimitedSliceScanner("hello,world", ",")
//	// Comma or semicolon separated values.
//	DelimitedSliceScanner("goodbye,cruel;world", ",;")
//
// Parameters:
//  1. input      = Input string that will be scanned.
//  2. delimiters = Set of delimiter characters.  If this string is empty, the
//     scanner will return the whole input string on the first
//     call to Next.
func DelimitedSliceScanner(input, delimiters string) iter.Seq[string] {
	return func(yield func(string) bool) {
		i := 0
		for {
			if p := findNextToken(i, input, delimiters); p > -1 {
				yield(input[i:p])
				i = p
			} else {
				yield(input[i:])
				break
			}
		}
	}
}

func findNextToken(from int, text, delim string) int {
	for i := from; i < len(text); i++ {
		for j := 0; j < len(delim); j++ {
			if text[i] == delim[j] {
				return i
			}
		}
	}

	return -1
}
