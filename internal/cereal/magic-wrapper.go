package cereal

import "github.com/foxcapades/argonaut/pkg/argo"

func NewMagicWrapper[T any](magic argo.MagicUnmarshaler) Deserializer[T] {
	return magicWrapper[T]{magic}
}

type magicWrapper[T any] struct {
	magic argo.MagicUnmarshaler
}

func (m magicWrapper[T]) Deserialize(raw string, prev *T) (out T, err error) {
	// TODO: how does this work with slices and such???
	err = m.magic.Unmarshal(raw, &out)
	return
}
