package cereal

import "github.com/foxcapades/argonaut/pkg/argo"

func NewMagicWrapper[T any](magic argo.MagicUnmarshaler) Deserializer[T] {
	return magicWrapper[T]{magic}
}

type magicWrapper[T any] struct {
	magic argo.MagicUnmarshaler
}

func (m magicWrapper[T]) Deserialize(raw string) (out T, err error) {
	err = m.magic.Unmarshal(raw, &out)
	return
}
