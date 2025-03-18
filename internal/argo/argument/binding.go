package argument

import "github.com/foxcapades/argonaut/v3/pkg/argo"

type BoundErrorFn[T any] = func(value T) error

func FunctionBinding[T any](fn BoundErrorFn[T]) argo.Binding {
	return binding{argo.BindingTypeErrorFunc, fn}
}

type BoundSimpleFn[T any] = func(value T)

func SimpleFunctionBinding[T any](fn BoundSimpleFn[T]) argo.Binding {
	return binding{argo.BindingTypeSimpleFunc, fn}
}

func UnmarshalerBinding(un argo.ValueUnmarshaler) argo.Binding {
	return binding{argo.BindingTypeUnmarshaler, un}
}

func PointerBinding(ptr any) argo.Binding {
	return binding{argo.BindingTypePointer, ptr}
}

type binding struct {
	kind argo.BindingType
	ref  any
}

func (b binding) Type() argo.BindingType {
	return b.kind
}

func (b binding) BoundTo() any {
	return b.ref
}

func newArgumentBindingError(root error, builder argo.ArgumentBuilder) argo.ArgumentBindingError {
	return &argumentBindingError{root, builder}
}

type argumentBindingError struct {
	root    error
	builder argo.ArgumentBuilder
}

func (a argumentBindingError) Unwrap() error {
	return a.root
}

func (a argumentBindingError) Error() string {
	return "ArgumentBindingError: " + a.root.Error()
}

func (a argumentBindingError) Builder() argo.ArgumentBuilder {
	return a.builder
}
