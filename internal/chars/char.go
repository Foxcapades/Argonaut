package chars

const (
	CharTab    = '\t'
	CharLF     = '\n'
	CharCR     = '\r'
	CharSpace  = ' '
	CharDash   = '-'
	CharEquals = '='
)

// IsWhitespace tests if the given character is a whitespace character.
//
// Whitespace characters include spaces, tabs, line feeds, and carriage returns.
func IsWhitespace(c byte) bool {
	return c == CharSpace || c == CharLF || c == CharTab || c == CharCR
}

// IsAlphanumeric tests whether the given character is an alphanumeric character
// meaning it is a standard ASCII alphabet character or digit.
func IsAlphanumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}

// IsAlpha tests whether the given character is a standard ASCII alphabet
// character.
func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// IsNumeric tests whether the given character is a standard ASCII digit
// character.
func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

// IsFlagStringSafe tests whether the given character is safe to appear in a
// flag identifier string.
func IsFlagStringSafe(c byte) bool {
	return IsAlphanumeric(c) || c == '-' || c == '_'
}
