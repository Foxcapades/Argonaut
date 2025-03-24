package utils

import "bufio"

func Write2Strings(buf *bufio.Writer, a, b string) error {
	if _, err := buf.WriteString(a); err != nil {
		return err
	}
	if _, err := buf.WriteString(b); err != nil {
		return err
	}
	return nil
}
