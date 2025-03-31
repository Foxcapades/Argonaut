package utils

func IfElse[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}

	return ifFalse
}

func CallIfElse[T any](condition bool, ifTrue, ifFalse func() T) T {
	if condition {
		return ifTrue()
	}

	return ifFalse()
}
