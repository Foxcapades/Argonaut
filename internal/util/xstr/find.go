package xstr

func IndexOfByte(value string, target byte, from int) int {
	for i := from; i < len(value); i++ {
		if value[i] == target {
			return i
		}
	}

	return -1
}

func IndexOfAnyWithin(value string, targets string, max int) int {
	l := min(len(value), max)

	for i := 0; i < l; i++ {
		for j := range targets {
			if value[i] == targets[j] {
				return i
			}
		}
	}

	return -1
}
