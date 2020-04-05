package util

func IsDecimal(b byte) bool {
	return b > '/' && b < ':'
}

func IsUpperLetter(b byte) bool {
	return b > '@' && b < '['
}

func IsLowerLetter(b byte) bool {
	return b > '`' && b < '{'
}
