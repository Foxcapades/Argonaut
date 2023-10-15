package chars

const (
	StrDash   = "-"
	StrEquals = "="
	StrEmpty  = ""
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
