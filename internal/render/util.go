package render

import (
	"fmt"
	"strings"

	"github.com/Foxcapades/Argonaut/internal/chars"
)

const maxWidth uint8 = 100

func WritePadded(str string, ln int, out *strings.Builder) {
	its := ln - len(str)
	out.WriteString(str)
	for i := 0; i < its; i++ {
		out.WriteByte(chars.CharSpace)
	}
}

// IsBreakChar returns true if the given character is one
// that can be followed immediately by a line break
func IsBreakChar(b byte) bool {
	// TODO: '-' needs to be handled differently than spaces
	//       spaces are removed from the output, however the
	//       dash should be maintained.
	return b == chars.CharSpace || b == chars.CharTab
}

func BreakFmt(str, prefix string, width uint8, out *strings.Builder) {
	str = strings.ReplaceAll(strings.ReplaceAll(str, "\r\n", "\n"), "\r", "\n")

	size := int(width) - len(prefix)
	stln := len(str)

	if size < 1 {
		panic(fmt.Errorf("cannot break string into lengths of %d", size))
	}

	if stln <= size {
		out.WriteString(prefix)
		out.WriteString(str)
		return
	}

	lastSplit := 0
	lastBreak := 0
	for i := 0; i < len(str); i++ {
		b := str[i]

		if i-lastSplit >= size {
			if lastSplit > 0 {
				out.WriteByte(chars.CharLF)
				out.WriteString(prefix)
			}

			if IsBreakChar(b) || b == chars.CharLF {
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
					out.WriteByte(chars.CharDash)
					lastSplit = i - 1
				}
			} else {
				out.WriteString(str[lastSplit:lastBreak])
				lastSplit = lastBreak + 1
			}
		}

		if IsBreakChar(b) {
			lastBreak = i
		} else if b == chars.CharLF {
			// Don't spit out a new line if we hit \n within the
			// first line we are reading.
			if lastSplit > 0 {
				out.WriteByte(chars.CharLF)
				out.WriteString(prefix)
			}
			out.WriteString(str[lastSplit:i])
			lastBreak = i
			lastSplit = i + 1

		}
	}

	if lastSplit < len(str) {
		out.WriteByte(chars.CharLF)
		out.WriteString(prefix)
		out.WriteString(str[lastSplit:])
	}
}
