package argo

type Iterator[T any] interface {
	HasNext() bool

	Next() T
}

func ForEachRemaining[T any](it Iterator[T], fn func(T)) {
	for it.HasNext() {
		fn(it.Next())
	}
}

func Map[T, R any](it Iterator[T], fn func(T) R) Iterator[R] {
	return mappingIterator[T, R]{it, fn}
}

type mappingIterator[T, R any] struct {
	root Iterator[T]
	mpFn func(T) R
}

func (m mappingIterator[T, R]) HasNext() bool {
	return m.root.HasNext()
}

func (m mappingIterator[T, R]) Next() R {
	return m.mpFn(m.root.Next())
}
