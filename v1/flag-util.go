package argo

func isValidShortFlag(b byte) bool {
	return isLowerLetter(b) || isUpperLetter(b) || isDecimal(b)
}

func isValidLongFlag(s string) bool {
	asB := []byte(s)

	if !isValidShortFlag(asB[1]) {
		return false
	}

	for i := 1; i < len(asB); i++ {
		b := asB[i]
		if !isLowerLetter(b) && b != '-' && !isDecimal(b) &&
			b != '_' && !isUpperLetter(b) {

			return false
		}
	}

	return true
}
