package argo

const (
	charTab    = '\t'
	charLF     = '\n'
	charCR     = '\r'
	charSpace  = ' '
	charDash   = '-'
	charEquals = '='
)

func isWhitespace(c byte) bool {
	return c == charSpace || c == charLF || c == charTab || c == charCR
}

func isAlphanumeric(c byte) bool {
	return isAlpha(c) || isNumeric(c)
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func isFlagStringSafe(c byte) bool {
	return isAlphanumeric(c) || c == '-' || c == '_'
}

const (
	strDash       = "-"
	strEquals     = "="
	strEmpty      = ""
	strDoubleDash = "--"
)

func nextWhitespace(s string) int {
	for i := range s {
		if isWhitespace(s[i]) {
			return i
		}
	}

	return -1
}

func nextEquals(s string) int {
	for i := range s {
		if s[i] == charEquals {
			return i
		}
	}

	return -1
}

func isBlank(s string) bool {
	if len(s) == 0 {
		return true
	}

	for i := range s {
		if !isWhitespace(s[i]) {
			return false
		}
	}

	return true
}

const defaultGroupName = "174b9e831dec431181e31ede822bb3b5"