package xstr

func IndexOfByte(value string, target byte, from int) int {
	for i := from; i < len(value); i++ {
		if value[i] == target {
			return i
		}
	}

	return -1
}
