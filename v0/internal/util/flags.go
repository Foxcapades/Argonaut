package util

func IsValidShortFlag(b byte) bool {
	return IsLowerLetter(b) || IsUpperLetter(b) || IsDecimal(b)
}

func IsValidLongFlag(s string) bool {
	asB := []byte(s)

	if !IsValidShortFlag(asB[0]) {
		return false
	}

	for i := 1; i < len(asB); i++ {
		b := asB[i]
		if !IsLowerLetter(b) && b != '-' && !IsDecimal(b) &&
			b != '_' && !IsUpperLetter(b) {

			return false
		}
	}

	return true
}
