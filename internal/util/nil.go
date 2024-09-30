package util

func RequireNonNil[T any](value T) T {
	if value == nil {
		panic("given value must not be nil")
	}

	return value
}
