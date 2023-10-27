package argo

type StringScannerFactory = func(input string) Scanner[string]

// Scanner defines an iterator over a given input that breaks the input into
// segments through some internal mechanism.
//
// For example, a scanner may break an input reader into lines of text, or split
// a string on a specific delimiter.
type Scanner[T any] interface {
	// HasNext indicates whether there is at least one more segment available.
	HasNext() bool

	// Next returns the next segment.
	Next() T
}

// DelimitedSliceScanner returns a new scanner over the given string, breaking
// it into substrings on every delimiter character.
//
// Examples:
//     // Comma separated values:
//     DelimitedSliceScanner("hello,world", ",")
//     // Comma or semicolon separated values.
//     DelimitedSliceScanner("goodbye,cruel;world", ",;")
//
// Parameters:
//   1. input      = Input string that will be scanned.
//   2. delimiters = Set of delimiter characters.  If this string is empty, the
//                   scanner will return the whole input string on the first
//                   call to Next.
func DelimitedSliceScanner(input, delimiters string) Scanner[string] {
	return &delimitedSliceScanner{input: input, delimiters: delimiters}
}

type delimitedSliceScanner struct{ input, delimiters string }

func (d delimitedSliceScanner) HasNext() bool {
	return len(d.input) > 0
}

func (d *delimitedSliceScanner) Next() string {
	if !d.HasNext() {
		panic("no such element")
	}

	idx := d.findNextToken()
	if idx == -1 {
		out := d.input
		d.input = ""
		return out
	}

	out := d.input[:idx]
	d.input = d.input[idx+1:]
	return out
}

func (d delimitedSliceScanner) findNextToken() int {
	for i := 0; i < len(d.input); i++ {
		for j := 0; j < len(d.delimiters); j++ {
			if d.input[i] == d.delimiters[j] {
				return i
			}
		}
	}

	return -1
}
