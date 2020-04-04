package argo

func isDecimal(b byte) bool {
	return b > '/' && b < ':'
}

func isUpperLetter(b byte) bool {
	return b > '@' && b < '['
}

func isLowerLetter(b byte) bool {
	return b > '`' && b < '{'
}
