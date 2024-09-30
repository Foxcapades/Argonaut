package argo

type ValueUnmarshaler[T any] interface {
	Unmarshal(rawInput string) (T, error)
}

type ValueUnmarshalerFn[T any] func(rawInput string) (T, error)

func (v ValueUnmarshalerFn[T]) Unmarshal(rawInput string) (T, error) {
	return v(rawInput)
}
