package arg

import (
	"github.com/foxcapades/argonaut/pkg/argo"
)

type ArgumentSpecImpl[T any] struct {
	name    string
	hasName bool

	description    string
	hasDescription bool

	summary    string
	hasSummary bool
}

func (a *ArgumentSpecImpl[T]) WithName(name string) argo.TypedArgumentSpecBuilder[T] {
	a.name = name
	a.hasName = true
	return a
}

func (a *ArgumentSpecImpl[T]) WithDescription(desc string) argo.TypedArgumentSpecBuilder[T] {
	a.description = desc
	a.hasDescription = true
	return a
}

func (a ArgumentSpecImpl[T]) WithSummary(summary string) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithRawDefault(value string) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithDefault(value T) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithPreValidator(validator argo.PreParseArgumentValidator) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithPostValidator(validator argo.PostParseArgumentValidator) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithValueConsumer(consumer argo.ArgumentValueConsumer[T]) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithBinding(binding *T) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) WithUnmarshaler(unmarshaler argo.ValueUnmarshaler[T]) argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) Required() argo.TypedArgumentSpecBuilder[T] {
	// TODO implement me
	panic("implement me")
}

func (a ArgumentSpecImpl[T]) Build() (argo.Argument, error) {
	// TODO implement me
	panic("implement me")
}
