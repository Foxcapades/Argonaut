package argo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var helpTextMaxWidth = 100

func init() {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		return
	}

	width, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}

	helpTextMaxWidth = min(width-20, 100)
}

func pad(size int, out *bufio.Writer) error {
	for i := 0; i < size; i++ {
		if err := out.WriteByte(charSpace); err != nil {
			return err
		}
	}
	return nil
}

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
	"    ",
	"      ",
	"        ",
	"          ",
}

func isBreakChar(b byte) bool {
	// TODO: '-' needs to be handled differently than spaces
	//       spaces are removed from the output, however the
	//       dash should be maintained.
	return b == charSpace || b == charTab
}

func breakFmt(str, prefix string, width int, out *bufio.Writer) error {
	str = strings.ReplaceAll(strings.ReplaceAll(str, "\r\n", "\n"), "\r", "\n")

	size := width - len(prefix)
	// stln := len(str)

	if size < 1 {
		panic(fmt.Errorf("cannot break string into lengths of %d", size))
	}

	// if _, err := out.WriteString(prefix); err != nil {
	// 	return err
	// }

	// if stln <= size {
	// 	if _, err := out.WriteString(str); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	lastSplit := 0
	lastBreak := 0
	for i := 0; i < len(str); i++ {
		b := str[i]

		if i-lastSplit >= size {
			if lastSplit > 0 {
				if err := out.WriteByte(charLF); err != nil {
					return err
				}
				if _, err := out.WriteString(prefix); err != nil {
					return err
				}
			}

			if isBreakChar(b) || b == charLF {
				if _, err := out.WriteString(str[lastSplit:i]); err != nil {
					return err
				}
				lastBreak = i
				lastSplit = i + 1
				continue
			}

			// Really long single word
			if lastSplit >= lastBreak {
				if size == 1 {
					if _, err := out.WriteString(str[lastSplit:i]); err != nil {
						return err
					}
					lastSplit = i
				} else {
					if _, err := out.WriteString(str[lastSplit : i-1]); err != nil {
						return err
					}
					if err := out.WriteByte(charDash); err != nil {
						return err
					}
					lastSplit = i - 1
				}
			} else {
				if _, err := out.WriteString(str[lastSplit:lastBreak]); err != nil {
					return err
				}
				lastSplit = lastBreak + 1
			}
		}

		if isBreakChar(b) {
			lastBreak = i
		} else if b == charLF {
			// Don't spit out a new line if we hit \n within the
			// first line we are reading.
			if lastSplit > 0 {
				if err := out.WriteByte(charLF); err != nil {
					return err
				}
			}
			if _, err := out.WriteString(prefix); err != nil {
				return err
			}
			if _, err := out.WriteString(str[lastSplit:i]); err != nil {
				return err
			}
			lastBreak = i
			lastSplit = i + 1
		}
	}

	if lastSplit < len(str) {
		if lastSplit > 0 {
			if err := out.WriteByte(charLF); err != nil {
				return err
			}
		}
		if _, err := out.WriteString(prefix); err != nil {
			return err
		}
		if _, err := out.WriteString(str[lastSplit:]); err != nil {
			return err
		}
	}

	return nil
}
