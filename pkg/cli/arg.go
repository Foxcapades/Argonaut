package cli

import (
	"github.com/foxcapades/argonaut/internal/argument"
	"github.com/foxcapades/argonaut/pkg/argo"
)

func Argument[T any]() argo.TypedArgumentSpecBuilder[T] {
	return argument.NewBuilder[T]()
}

func ArgBinding[T any](pointer *T) argo.TypedArgumentSpecBuilder[T] {
	return argument.NewBuilder[T]().WithBinding(pointer)
}
