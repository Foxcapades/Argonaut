package argo

import (
	"fmt"
	"strings"
)

const helpTextMaxWidth uint8 = 100

const (
	flagDivider    = " | "
	paragraphBreak = "\n\n"
)

var headerPadding = [...]string{
	"",
	"  ",
	"    ",
	"      ",
}

var subLinePadding = [...]string{
	"  ",
	"    ",
	"      ",
	"        ",
}

var descriptionPadding = [...]string{
	"         ",
	"           ",
	"             ",
	"               ",
}

func isBreakChar(b byte) bool {
	// TODO: '-' needs to be handled differently than spaces
	//       spaces are removed from the output, however the
	//       dash should be maintained.
	return b == charSpace || b == charTab
}

func breakFmt(str, prefix string, width uint8, out *strings.Builder) {
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
				out.WriteByte(charLF)
				out.WriteString(prefix)
			}

			if isBreakChar(b) || b == charLF {
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
					out.WriteByte(charDash)
					lastSplit = i - 1
				}
			} else {
				out.WriteString(str[lastSplit:lastBreak])
				lastSplit = lastBreak + 1
			}
		}

		if isBreakChar(b) {
			lastBreak = i
		} else if b == charLF {
			// Don't spit out a new line if we hit \n within the
			// first line we are reading.
			if lastSplit > 0 {
				out.WriteByte(charLF)
				out.WriteString(prefix)
			}
			out.WriteString(str[lastSplit:i])
			lastBreak = i
			lastSplit = i + 1

		}
	}

	if lastSplit < len(str) {
		out.WriteByte(charLF)
		out.WriteString(prefix)
		out.WriteString(str[lastSplit:])
	}
}
