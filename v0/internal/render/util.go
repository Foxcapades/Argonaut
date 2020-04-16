package render

import (
	"fmt"
	"strings"
)

func WritePadded(str string, ln int, out *strings.Builder) {
	its := ln - len(str)
	out.WriteString(str)
	for i := 0; i < its; i++ {
		out.WriteByte(sngSpace)
	}
}

// IsBreakChar returns true if the given character is one
// that can be followed immediately by a line break
func IsBreakChar(b byte) bool {
	// TODO: '-' needs to be handled differently than spaces
	//       spaces are removed from the output, however the
	//       dash should be maintained.
	return b == ' ' || b == '\t' /*|| b == '-'*/
}

func BreakFmt(str string, offset, width int, out *strings.Builder) {
	str = strings.ReplaceAll(strings.ReplaceAll(str, "\r\n", "\n"), "\r", "\n")

	size := width - offset
	stln := len(str)

	if size < 1 {
		panic(fmt.Errorf("cannot break string into lengths of %d", size))
	}

	if stln <= size {
		out.WriteString(str)
		return
	}

	buf := make([]byte, offset)
	for i := range buf {
		buf[i] = ' '
	}

	lastSplit := 0
	lastBreak := 0
	for i := range str {
		b := str[i]

		if i-lastSplit >= size {
			if lastSplit > 0 {
				out.WriteByte('\n')
				out.Write(buf)
			}

			if IsBreakChar(b) || b == '\n' {
				out.WriteString(str[lastSplit:i])
				lastBreak = i
				lastSplit = i + 1
				continue
			}

			// Really long single word
			if lastSplit >= lastBreak {
				if size == 1 {
					out.WriteString(str[lastSplit:i])
					lastSplit = i
				} else {
					out.WriteString(str[lastSplit : i-1])
					out.WriteByte('-')
					lastSplit = i - 1
				}
			} else {
				out.WriteString(str[lastSplit:lastBreak])
				lastSplit = lastBreak + 1
			}
		}

		if IsBreakChar(b) {
			lastBreak = i
		} else if b == '\n' {
			out.WriteString(str[lastSplit : i+1])
			out.Write(buf)
			lastBreak = i
			lastSplit = i + 1
		}
	}

	if lastSplit < len(str) {
		out.WriteByte('\n')
		out.Write(buf)
		out.WriteString(str[lastSplit:])
	}
}
