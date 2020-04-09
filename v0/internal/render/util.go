package render

import (
	"fmt"
	"strings"
)

func writePadded(str string, ln int, out *strings.Builder) {
	its := ln - len(str)
	out.WriteString(str)
	for i := 0; i < its; i++ {
		out.WriteByte(sngSpace)
	}
}

func isBreakChar(b byte) bool {
	return b == ' ' || b == '-'
}

func BreakFmt(str string, offset, width int, out *strings.Builder) {
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

			if isBreakChar(b) {
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

		if isBreakChar(b) {
			lastBreak = i
		}
	}

	if lastSplit < len(str) {
		out.WriteByte('\n')
		out.Write(buf)
		out.WriteString(str[lastSplit:])
	}
}
