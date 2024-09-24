package cereal

import "github.com/foxcapades/argonaut/pkg/argo"

func NewCustomWrapper[T any](custom argo.ValueUnmarshaler[T]) Deserializer[T] {
	return customWrapper[T]{custom}
}

type customWrapper[T any] struct {
	custom argo.ValueUnmarshaler[T]
}

func (c customWrapper[T]) Deserialize(raw string) (T, error) {
	return c.custom.Unmarshal(raw)
}
