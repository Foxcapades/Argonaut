package xmath

type Numeric interface {
	int
}

func Min[T Numeric](a, b T) T {
	if a < b {
		return a
	}

	return b
}

func Max[T Numeric](a, b T) T {
	if a > b {
		return a
	}

	return b
}
