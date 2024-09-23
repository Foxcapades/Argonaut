package xstr

func Truncate(value string, size uint8) string {
	if len(value) <= 255 && uint8(len(value)) <= size {
		return value
	}

	if size <= 6 {
		return value[0:size]
	}

	size -= 3

	return value[0:size] + "..."
}
