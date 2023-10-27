package chars

import "errors"

const (
	CharTab    = '\t'
	CharLF     = '\n'
	CharCR     = '\r'
	CharSpace  = ' '
	CharDash   = '-'
	CharEquals = '='
)

func IsWhitespace(c byte) bool {
	return c == CharSpace || c == CharLF || c == CharTab || c == CharCR
}

func IsAlphanumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}

func IsOctal(c byte) bool {
	return c >= '0' && c <= '7'
}

func IsOctalString(s string) bool {
	for i := 0; i < len(s); i++ {
		if !IsOctal(s[i]) {
			return false
		}
	}
	return true
}

func IsHex(c byte) bool {
	return IsNumeric(c) ||
		(c >= 'A' && c <= 'F') ||
		(c >= 'a' && c <= 'f')
}

func IsHexString(s string) bool {
	for i := range s {
		if !IsHex(s[i]) {
			return false
		}
	}

	return true
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsFlagStringSafe(c byte) bool {
	return IsAlphanumeric(c) || c == '-' || c == '_'
}

func ValidateCommandNodeName(name string) error {
	if len(name) == 0 {
		return errors.New("command names must not be blank")
	}

	if !(IsAlphanumeric(name[0]) || name[0] == '_') {
		return errors.New("command names must begin with an alphanumeric character or an underscore")
	}

	for i := 1; i < len(name); i++ {
		if !IsFlagStringSafe(name[i]) {
			return errors.New("command names may only contain alphanumeric characters, dashes, and/or underscores")
		}
	}

	return nil
}

const (
	StrDash       = "-"
	StrEquals     = "="
	StrEmpty      = ""
	StrDoubleDash = "--"
	StrLF         = "\n"
)

func NextWhitespace(s string) int {
	for i := range s {
		if IsWhitespace(s[i]) {
			return i
		}
	}

	return -1
}

func NextEquals(s string) int {
	for i := range s {
		if s[i] == CharEquals {
			return i
		}
	}

	return -1
}

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

const DefaultGroupName = "174b9e831dec431181e31ede822bb3b5"
