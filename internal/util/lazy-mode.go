package util

func IfElse[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}

	return ifFalse
}
