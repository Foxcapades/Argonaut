package util

import (
	"fmt"
	A "github.com/Foxcapades/Argonaut/v0/pkg/argo"
	"strings"
)

func RenderHelp(hf A.Flag) string {
	var out strings.Builder

	com := hf.Parent().Parent()
	fgs := com.FlagGroups()
	pos := com.Arguments()

	out.WriteString("  ")
	out.WriteString(com.Name())

	if len(fgs) > 0 {
		out.WriteString(" [OPTIONS]")
	}

	for i, arg := range pos {
		out.WriteByte(' ')
		name := nameForPosArg(i, arg)
		if arg.Required() {
			out.WriteString(name)
		} else {
			out.WriteByte('[')
			out.WriteString(name)
			out.WriteByte(']')
		}
	}

	if com.HasDescription() {
		out.WriteString("\n\n  ")
		BreakFmt(com.Description(), 2, 80, &out)
	}

	for i, fg := range fgs {
		if fg.HasFlags() {
			out.WriteString("\n\n")
			if fg.HasName() {
				out.WriteString(fg.Name())
			} else {
				out.WriteString(fmt.Sprintf("Flag Group %d", i + 1))
			}
			out.WriteString("\n\n")
			renderFg(i, fg, &out)
		}
	}

	return out.String()
}

func renderFg(i int, fg A.FlagGroup, out *strings.Builder) {
	reqFlags := make(map[string]string)
	optFlags := make(map[string]string, len(fg.Flags()))
	maxLn := 0

	for _, flag := range fg.Flags() {
		name := flag.String()
		nln  := len(name)

		if nln > maxLn {
			maxLn = nln
		}

		if flag.Required() {
			reqFlags[name] = "REQUIRED  " + flag.Description()
		} else {
			optFlags[name] = flag.Description()
		}
	}

	maxLn += 6

	if len(reqFlags) > 0 {
		out.WriteString("  Required Flag(s)\n\n")
		for key, val := range reqFlags {
			out.WriteString("    ")
			out.WriteString(key)
			out.WriteString("  ")
			BreakFmt(val, maxLn, 80, out)
			out.WriteByte('\n')
		}
	}

	if len(optFlags) > 0 {
		out.WriteString("  Optional Flag(s)\n\n")
		for key, val := range optFlags {
			out.WriteString("    ")
			out.WriteString(key)
			out.WriteString("  ")
			BreakFmt(val, maxLn, 80, out)
			out.WriteByte('\n')
		}
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

		if i - lastSplit >= size {
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
					out.WriteString(str[lastSplit:i-1])
					out.WriteByte('-')
					lastSplit = i-1
				}
			} else {
				out.WriteString(str[lastSplit:lastBreak])
				lastSplit = lastBreak + 1
			}
		}

		if isBreakChar(b) {lastBreak = i}
	}

	if lastSplit < len(str) {
		out.WriteByte('\n')
		out.Write(buf)
		out.WriteString(str[lastSplit:])
	}
}


func nameForPosArg(i int, arg A.Argument) string {
	if arg.HasName() {
		return arg.Name()
	} else {
		return fmt.Sprintf("arg%d", i)
	}
}