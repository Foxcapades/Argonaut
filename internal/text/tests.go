package text

func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}

	for i := range s {
		if !IsWhitespace(s[i]) {
			return false
		}
	}

	return true
}

func IsOctalString(s string) bool {
	for i := 0; i < len(s); i++ {
		if !IsOctal(s[i]) {
			return false
		}
	}
	return true
}

func IsHexString(s string) bool {
	for i := range s {
		if !IsHex(s[i]) {
			return false
		}
	}

	return true
}

func IsWhitespace(c byte) bool {
	return c == SpaceByte || c == LineFeedByte || c == TabByte || c == CarriageReturnByte
}

func IsWord(c byte) bool {
	return IsAlpha(c) || IsNumeric(c) || c == UnderscoreByte
}

func IsAlphanumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}

func IsOctal(c byte) bool {
	return c >= '0' && c <= '7'
}

func IsHex(c byte) bool {
	return IsNumeric(c) ||
		(c >= 'A' && c <= 'F') ||
		(c >= 'a' && c <= 'f')
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}
