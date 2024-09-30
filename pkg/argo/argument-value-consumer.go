package argo

type ArgumentValueConsumer[T any] interface {
	Accept(value T) error
}

type ArgumentValueConsumerFn[T any] func(value T) error

func (a ArgumentValueConsumerFn[T]) Accept(value T) error {
	return a(value)
}

func SimpleArgumentValueConsumerFn[T any](fn func(value T)) ArgumentValueConsumer[T] {
	return ArgumentValueConsumerFn[T](func(value T) error {
		fn(value)
		return nil
	})
}
