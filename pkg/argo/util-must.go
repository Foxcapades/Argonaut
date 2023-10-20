package argo

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustReturn[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
