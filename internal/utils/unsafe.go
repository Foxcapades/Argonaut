package utils

func Cast[T any](value any) T {
	return value.(T)
}
